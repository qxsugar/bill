package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/model"
	"github.com/qxsugar/pkg/kit"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// TaskAutoSettle 房间超时自动结算任务的类型标识。
	TaskAutoSettle = "room:auto_settle"
	// roomCodeTTL 房间码在 Redis 中的占位时长，与房间自动结算时限一致。
	roomCodeTTL = 24 * time.Hour
	// roomCodeKeyPrefix 房间码占位键前缀。
	roomCodeKeyPrefix = "bill:room:code:"
)

// autoSettlePayload room:auto_settle 任务载荷。
type autoSettlePayload struct {
	RoomId int64 `json:"room_id"`
}

// Broadcaster 由 websocket 层实现，房间状态变化时通知房间内所有连接。
// 用接口解耦，避免 service 依赖 ws 包形成环。
type Broadcaster interface {
	BroadcastRoomUpdated(roomId int64)
	BroadcastSettled(roomId int64)
}

// noopBroadcaster 默认空实现，ws 层就绪后通过 SetBroadcaster 注入真实实现。
type noopBroadcaster struct{}

func (noopBroadcaster) BroadcastRoomUpdated(int64) {}
func (noopBroadcaster) BroadcastSettled(int64)     {}

type RoomService struct {
	db             *gorm.DB
	rdb            *redis.Client
	asynqClient    *asynq.Client
	roomDao        *dao.RoomDao
	memberDao      *dao.RoomMemberDao
	logDao         *dao.RoomLogDao
	userDao        *dao.UserDao
	transactionDao *dao.TransactionDao
	broadcaster    Broadcaster
}

func NewRoomService(
	db *gorm.DB,
	rdb *redis.Client,
	asynqClient *asynq.Client,
	roomDao *dao.RoomDao,
	memberDao *dao.RoomMemberDao,
	logDao *dao.RoomLogDao,
	userDao *dao.UserDao,
	transactionDao *dao.TransactionDao,
) *RoomService {
	return &RoomService{
		db:             db,
		rdb:            rdb,
		asynqClient:    asynqClient,
		roomDao:        roomDao,
		memberDao:      memberDao,
		logDao:         logDao,
		userDao:        userDao,
		transactionDao: transactionDao,
		broadcaster:    noopBroadcaster{},
	}
}

// SetBroadcaster 注入 websocket 广播实现。
func (s *RoomService) SetBroadcaster(b Broadcaster) {
	if b != nil {
		s.broadcaster = b
	}
}

// Create 创建房间：生成唯一房间码，创建者作为房主自动加入，并写入创建日志。
func (s *RoomService) Create(userId int64) (*model.Room, error) {
	code, err := s.genUniqueCode()
	if err != nil {
		return nil, err
	}

	room := &model.Room{Code: code, OwnerId: userId, Status: model.RoomStatusActive}
	if err := s.roomDao.Create(room); err != nil {
		return nil, err
	}

	// 房主入场
	if err := s.memberDao.Create(&model.RoomMember{RoomId: room.Id, UserId: userId, Balance: 0}); err != nil {
		return nil, err
	}

	// 系统日志：房间已创建
	s.addLog(room.Id, nil, model.LogTypeJoin, fmt.Sprintf("房间已创建，房间代码为 %s", code))
	// 房主加入日志
	if u, _ := s.userDao.FindById(userId); u != nil {
		s.addLog(room.Id, &userId, model.LogTypeJoin, fmt.Sprintf("%s 加入房间", u.Nickname))
	}

	// 投递 24 小时后自动结算任务，超时未结算时由 worker 兜底结算。
	s.enqueueAutoSettle(room.Id)

	return room, nil
}

// Join 通过房间码加入房间。已结算的房间不可加入；已在场则幂等返回。
func (s *RoomService) Join(userId int64, code string) (*model.Room, error) {
	room, err := s.roomDao.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, kit.NewNotFoundError().WithInfo("房间不存在")
	}
	if room.Status == model.RoomStatusSettled {
		return nil, kit.NewFailedPreconditionError().WithInfo("房间已结算，无法加入")
	}

	member, err := s.memberDao.FindByRoomAndUser(room.Id, userId)
	if err != nil {
		return nil, err
	}
	if member != nil && member.LeftAt == nil {
		// 已在房间内，幂等
		return room, nil
	}
	if member != nil {
		// 曾离开过，重新入场：清空离开时间
		member.LeftAt = nil
		if err := s.memberDao.Save(member); err != nil {
			return nil, err
		}
	} else {
		if err := s.memberDao.Create(&model.RoomMember{RoomId: room.Id, UserId: userId, Balance: 0}); err != nil {
			return nil, err
		}
	}

	if u, _ := s.userDao.FindById(userId); u != nil {
		s.addLog(room.Id, &userId, model.LogTypeJoin, fmt.Sprintf("%s 加入房间", u.Nickname))
	}
	s.broadcaster.BroadcastRoomUpdated(room.Id)
	return room, nil
}

// Leave 离开房间。房主不可主动离开（需先结算），非房主直接退出。
func (s *RoomService) Leave(userId, roomId int64) error {
	room, err := s.roomDao.FindById(roomId)
	if err != nil {
		return err
	}
	if room == nil {
		return kit.NewNotFoundError().WithInfo("房间不存在")
	}
	if room.OwnerId == userId && room.Status == model.RoomStatusActive {
		return kit.NewFailedPreconditionError().WithInfo("房主不能主动离开，请点击结算后离开")
	}

	member, err := s.memberDao.FindByRoomAndUser(roomId, userId)
	if err != nil {
		return err
	}
	if member == nil || member.LeftAt == nil {
		if member != nil {
			now := kit.TimeStamp{Time: nowTime()}
			member.LeftAt = &now
			if err := s.memberDao.Save(member); err != nil {
				return err
			}
		}
	}

	if u, _ := s.userDao.FindById(userId); u != nil {
		s.addLog(roomId, &userId, model.LogTypeLeave, fmt.Sprintf("%s 离开房间", u.Nickname))
	}
	s.broadcaster.BroadcastRoomUpdated(roomId)
	return nil
}

// Settle 结算房间：仅房主可操作。
// 标记房间已结算、记录结算时间，所有在场成员标记离开，写结算日志并广播 settled。
// 结算后房间不可再操作，所有用户跳转到「结算-已结算」页。
func (s *RoomService) Settle(userId, roomId int64) error {
	room, err := s.roomDao.FindById(roomId)
	if err != nil {
		return err
	}
	if room == nil {
		return kit.NewNotFoundError().WithInfo("房间不存在")
	}
	if room.OwnerId != userId {
		return kit.NewPermissionDeniedError().WithInfo("只有房主可以点击结算")
	}
	if room.Status == model.RoomStatusSettled {
		return kit.NewFailedPreconditionError().WithInfo("房间已结算")
	}

	owner, _ := s.userDao.FindById(userId)
	return s.doSettle(room, &userId, fmt.Sprintf("%s 结算了房间", nameOf(owner)))
}

// doSettle 执行结算的核心逻辑：标记房间已结算、成员离场、写日志并广播。
// 调用方负责校验前置条件（如房主权限、是否已结算）。
// operatorId 为 nil 表示系统自动结算。
func (s *RoomService) doSettle(room *model.Room, operatorId *int64, logContent string) error {
	now := nowTime()
	roomId := room.Id
	err := s.db.Transaction(func(tx *gorm.DB) error {
		settledAt := kit.TimeStamp{Time: now}
		room.Status = model.RoomStatusSettled
		room.SettledAt = &settledAt
		if err := tx.Save(room).Error; err != nil {
			return err
		}
		// 所有在场成员标记离开
		if err := tx.Model(&model.RoomMember{}).
			Where("room_id = ? and left_at is null", roomId).
			Update("left_at", now).Error; err != nil {
			return err
		}
		return tx.Create(&model.RoomLog{
			RoomId: roomId, UserId: operatorId, LogType: model.LogTypeSettle,
			Content: logContent,
		}).Error
	})
	if err != nil {
		return err
	}
	s.broadcaster.BroadcastSettled(roomId)
	return nil
}

// Logs 分页返回房间日志。日志 DAO 按 id 倒序取出，
// 设计要求房间日志页从旧到新展示，故在此反转为正序。
func (s *RoomService) Logs(roomId int64, limit, offset int) ([]*model.RoomLog, int64, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	list, total, err := s.logDao.ListByRoomId(roomId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	// 反转为从旧到新
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list, total, nil
}

// genUniqueCode 按设计规则生成唯一房间码：4 位（优先两两连号）→ 5 位兜底。
// 通过 Redis SETNX 原子占位（TTL 24h），避免并发创建撞码；
// 同时排除数据库中仍存活的同码房间。占位成功返回 true。
func (s *RoomService) genUniqueCode() (string, error) {
	ctx := context.Background()
	claim := func(code string) bool {
		// 先排除库中已有同码房间（占位过期但房间仍在的情况）。
		if r, err := s.roomDao.FindByCode(code); err != nil || r != nil {
			return false
		}
		ok, err := s.rdb.SetNX(ctx, roomCodeKeyPrefix+code, 1, roomCodeTTL).Result()
		return err == nil && ok
	}
	if code := gen4DigitPreferred(claim); code != "" {
		return code, nil
	}
	if code := gen5Digit(claim); code != "" {
		return code, nil
	}
	return "", kit.NewResourceExhaustedError().WithInfo("房间码已耗尽，请稍后再试")
}

// enqueueAutoSettle 投递房间超时自动结算任务，24 小时后触发。失败仅记录日志，不阻断创建。
func (s *RoomService) enqueueAutoSettle(roomId int64) {
	if s.asynqClient == nil {
		return
	}
	payload, err := json.Marshal(autoSettlePayload{RoomId: roomId})
	if err != nil {
		zap.S().Errorf("marshal auto-settle payload failed: room=%d err=%v", roomId, err)
		return
	}
	_, err = s.asynqClient.Enqueue(
		asynq.NewTask(TaskAutoSettle, payload),
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.ProcessIn(roomCodeTTL),
	)
	if err != nil {
		zap.S().Errorf("enqueue auto-settle task failed: room=%d err=%v", roomId, err)
	}
}

// HandleAutoSettle 处理 room:auto_settle 任务：房间到期仍活跃则自动结算。
// 幂等：房间已结算或不存在时直接返回，不视为失败。
func (s *RoomService) HandleAutoSettle(ctx context.Context, task *asynq.Task) error {
	var p autoSettlePayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		// 载荷不可解析，重试也无意义。
		return fmt.Errorf("unmarshal auto-settle payload: %w: %w", err, asynq.SkipRetry)
	}

	room, err := s.roomDao.FindById(p.RoomId)
	if err != nil {
		return err
	}
	if room == nil || room.Status == model.RoomStatusSettled {
		return nil
	}

	if err := s.doSettle(room, nil, "房间超过 24 小时未结算，已自动结算"); err != nil {
		return err
	}
	zap.S().Infof("room auto-settled: room=%d", p.RoomId)
	return nil
}

// addLog 写入房间日志，忽略错误（日志不应阻断主流程）。
func (s *RoomService) addLog(roomId int64, userId *int64, logType, content string) {
	_ = s.logDao.Create(&model.RoomLog{
		RoomId:  roomId,
		UserId:  userId,
		LogType: logType,
		Content: content,
	})
}
