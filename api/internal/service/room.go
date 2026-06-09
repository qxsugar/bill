package service

import (
	"fmt"

	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/model"
	"github.com/qxsugar/pkg/kit"
	"gorm.io/gorm"
)

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
	roomDao        *dao.RoomDao
	memberDao      *dao.RoomMemberDao
	logDao         *dao.RoomLogDao
	userDao        *dao.UserDao
	transactionDao *dao.TransactionDao
	broadcaster    Broadcaster
}

func NewRoomService(
	db *gorm.DB,
	roomDao *dao.RoomDao,
	memberDao *dao.RoomMemberDao,
	logDao *dao.RoomLogDao,
	userDao *dao.UserDao,
	transactionDao *dao.TransactionDao,
) *RoomService {
	return &RoomService{
		db:             db,
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

	now := nowTime()
	err = s.db.Transaction(func(tx *gorm.DB) error {
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
		owner, _ := s.userDao.FindById(userId)
		return tx.Create(&model.RoomLog{
			RoomId: roomId, UserId: &userId, LogType: model.LogTypeSettle,
			Content: fmt.Sprintf("%s 结算了房间", nameOf(owner)),
		}).Error
	})
	if err != nil {
		return err
	}
	s.broadcaster.BroadcastSettled(roomId)
	return nil
}

// genUniqueCode 按设计规则生成唯一房间码：4 位（优先两两连号）→ 5 位兜底。
func (s *RoomService) genUniqueCode() (string, error) {
	exists := func(code string) bool {
		r, err := s.roomDao.FindByCode(code)
		return err == nil && r != nil
	}
	if code := gen4DigitPreferred(exists); code != "" {
		return code, nil
	}
	if code := gen5Digit(exists); code != "" {
		return code, nil
	}
	return "", kit.NewResourceExhaustedError().WithInfo("房间码已耗尽，请稍后再试")
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
