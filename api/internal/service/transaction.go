package service

import (
	"fmt"

	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/model"
	"github.com/qxsugar/pkg/kit"
	"gorm.io/gorm"
)

type TransactionService struct {
	db             *gorm.DB
	roomDao        *dao.RoomDao
	memberDao      *dao.RoomMemberDao
	transactionDao *dao.TransactionDao
	logDao         *dao.RoomLogDao
	userDao        *dao.UserDao
	broadcaster    Broadcaster
}

func NewTransactionService(
	db *gorm.DB,
	roomDao *dao.RoomDao,
	memberDao *dao.RoomMemberDao,
	transactionDao *dao.TransactionDao,
	logDao *dao.RoomLogDao,
	userDao *dao.UserDao,
) *TransactionService {
	return &TransactionService{
		db:             db,
		roomDao:        roomDao,
		memberDao:      memberDao,
		transactionDao: transactionDao,
		logDao:         logDao,
		userDao:        userDao,
		broadcaster:    noopBroadcaster{},
	}
}

func (s *TransactionService) SetBroadcaster(b Broadcaster) {
	if b != nil {
		s.broadcaster = b
	}
}

// ExpenseItem 一笔支出明细：向某个收入方支出的金额。
type ExpenseItem struct {
	ToUserId int64
	Amount   float64
}

// Expense 记录支出：当前用户（from）向一个或多个收入方（to）发送积分。
// 支持三种前端模式（单笔/均分/统一），后端统一接收为明细列表。
// from 余额减少、各 to 余额增加，单事务内完成并写日志。
func (s *TransactionService) Expense(fromUserId, roomId int64, items []ExpenseItem) error {
	room, err := s.roomDao.FindById(roomId)
	if err != nil {
		return err
	}
	if room == nil {
		return kit.NewNotFoundError().WithInfo("房间不存在")
	}
	if room.Status != model.RoomStatusActive {
		return kit.NewFailedPreconditionError().WithInfo("房间已结算，无法支出")
	}

	// 房间里没有其他用户时禁止支出
	activeCount, err := s.memberDao.CountActiveByRoomId(roomId)
	if err != nil {
		return err
	}
	if activeCount < 2 {
		return kit.NewFailedPreconditionError().WithInfo("人数不足，无法支出")
	}

	if len(items) == 0 {
		return kit.NewInvalidArgumentError().WithInfo("请选择支出对象")
	}
	for _, it := range items {
		if it.Amount <= 0 {
			return kit.NewInvalidArgumentError().WithInfo("支出金额必须大于 0")
		}
		if it.ToUserId == fromUserId {
			return kit.NewInvalidArgumentError().WithInfo("不能向自己支出")
		}
	}

	fromUser, _ := s.userDao.FindById(fromUserId)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		fromMember, err := txFindMember(tx, roomId, fromUserId)
		if err != nil {
			return err
		}
		if fromMember == nil || fromMember.LeftAt != nil {
			return kit.NewFailedPreconditionError().WithInfo("你不在房间内")
		}

		for _, it := range items {
			toMember, err := txFindMember(tx, roomId, it.ToUserId)
			if err != nil {
				return err
			}
			if toMember == nil || toMember.LeftAt != nil {
				return kit.NewFailedPreconditionError().WithInfo("收入方不在房间内")
			}

			// 写交易记录
			if err := tx.Create(&model.Transaction{
				RoomId:     roomId,
				FromUserId: fromUserId,
				ToUserId:   it.ToUserId,
				Amount:     it.Amount,
				Status:     model.TransactionStatusValid,
			}).Error; err != nil {
				return err
			}

			// 更新余额：from 减、to 加
			fromMember.Balance -= it.Amount
			toMember.Balance += it.Amount
			if err := tx.Save(toMember).Error; err != nil {
				return err
			}

			toUser, _ := s.userDao.FindById(it.ToUserId)
			content := fmt.Sprintf("%s 向 %s 发送 %g¥", nameOf(fromUser), nameOf(toUser), it.Amount)
			if err := tx.Create(&model.RoomLog{
				RoomId: roomId, UserId: &fromUserId, LogType: model.LogTypeTransfer, Content: content,
			}).Error; err != nil {
				return err
			}
		}
		return tx.Save(fromMember).Error
	})
	if err != nil {
		return err
	}
	s.broadcaster.BroadcastRoomUpdated(roomId)
	return nil
}

// Revoke 撤销自己发出的有效交易：恢复双方余额，标记撤销，写日志。
func (s *TransactionService) Revoke(userId, txId int64) error {
	t, err := s.transactionDao.FindById(txId)
	if err != nil {
		return err
	}
	if t == nil {
		return kit.NewNotFoundError().WithInfo("记录不存在")
	}
	if t.FromUserId != userId {
		return kit.NewPermissionDeniedError().WithInfo("只能撤销自己发送的记录")
	}
	if t.Status != model.TransactionStatusValid {
		return kit.NewFailedPreconditionError().WithInfo("该记录已撤销")
	}
	room, _ := s.roomDao.FindById(t.RoomId)
	if room != nil && room.Status != model.RoomStatusActive {
		return kit.NewFailedPreconditionError().WithInfo("房间已结算，无法撤销")
	}

	fromUser, _ := s.userDao.FindById(t.FromUserId)
	toUser, _ := s.userDao.FindById(t.ToUserId)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		fromMember, err := txFindMember(tx, t.RoomId, t.FromUserId)
		if err != nil {
			return err
		}
		toMember, err := txFindMember(tx, t.RoomId, t.ToUserId)
		if err != nil {
			return err
		}
		if fromMember != nil {
			fromMember.Balance += t.Amount
			if err := tx.Save(fromMember).Error; err != nil {
				return err
			}
		}
		if toMember != nil {
			toMember.Balance -= t.Amount
			if err := tx.Save(toMember).Error; err != nil {
				return err
			}
		}

		now := kit.TimeStamp{Time: nowTime()}
		t.Status = model.TransactionStatusRevoked
		t.RevokedAt = &now
		if err := tx.Save(t).Error; err != nil {
			return err
		}

		content := fmt.Sprintf("%s 撤销了向 %s 发送的 %g¥", nameOf(fromUser), nameOf(toUser), t.Amount)
		return tx.Create(&model.RoomLog{
			RoomId: t.RoomId, UserId: &userId, LogType: model.LogTypeRevoke, Content: content,
		}).Error
	})
	if err != nil {
		return err
	}
	s.broadcaster.BroadcastRoomUpdated(t.RoomId)
	return nil
}

// Thank 对自己收到的交易发送感谢，仅限收入方且只能感谢一次。
func (s *TransactionService) Thank(userId, txId int64) error {
	t, err := s.transactionDao.FindById(txId)
	if err != nil {
		return err
	}
	if t == nil {
		return kit.NewNotFoundError().WithInfo("记录不存在")
	}
	if t.ToUserId != userId {
		return kit.NewPermissionDeniedError().WithInfo("只能对自己收到的积分表示感谢")
	}
	if t.Status != model.TransactionStatusValid {
		return kit.NewFailedPreconditionError().WithInfo("该记录已撤销")
	}
	if t.Thanked == model.TransactionThanked {
		return kit.NewFailedPreconditionError().WithInfo("已经感谢过了")
	}

	fromUser, _ := s.userDao.FindById(t.FromUserId)
	toUser, _ := s.userDao.FindById(t.ToUserId)

	t.Thanked = model.TransactionThanked
	if err := s.transactionDao.Save(t); err != nil {
		return err
	}
	// 感谢方向：收入方 to 向支出方 from 表达感谢
	content := fmt.Sprintf("%s 向 %s 发送了感谢", nameOf(toUser), nameOf(fromUser))
	_ = s.logDao.Create(&model.RoomLog{
		RoomId: t.RoomId, UserId: &userId, LogType: model.LogTypeThanks, Content: content,
	})
	s.broadcaster.BroadcastRoomUpdated(t.RoomId)
	return nil
}

// txFindMember 在事务内按房间+用户查成员。
func txFindMember(tx *gorm.DB, roomId, userId int64) (*model.RoomMember, error) {
	var m model.RoomMember
	err := tx.Where("room_id = ? and user_id = ?", roomId, userId).First(&m).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func nameOf(u *model.User) string {
	if u == nil {
		return "某玩家"
	}
	return u.Nickname
}
