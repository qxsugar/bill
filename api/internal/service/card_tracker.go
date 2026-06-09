package service

import (
	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/model"
	"github.com/qxsugar/pkg/kit"
)

type CardTrackerService struct {
	dao *dao.CardTrackerDao
}

func NewCardTrackerService(d *dao.CardTrackerDao) *CardTrackerService {
	return &CardTrackerService{dao: d}
}

// Get 返回用户的记牌器状态，不存在则用默认牌（1 副）初始化。
func (s *CardTrackerService) Get(userId int64) (*model.CardTracker, error) {
	t, err := s.dao.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}
	t = &model.CardTracker{UserId: userId, DeckCount: 1, Counts: model.DefaultCounts(1)}
	if err := s.dao.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

// Adjust 调整某个牌面的剩余数量：delta=-1 表示点击扣除，delta=+1 表示双击增加。
// 数量被夹在 [0, 上限] 之间，上限为该牌面默认张数（按牌副数）。
func (s *CardTrackerService) Adjust(userId int64, rank string, delta int) (*model.CardTracker, error) {
	t, err := s.Get(userId)
	if err != nil {
		return nil, err
	}
	def := model.DefaultCounts(t.DeckCount)
	max, ok := def[rank]
	if !ok {
		return nil, kit.NewInvalidArgumentError().WithInfo("无效的牌面")
	}
	v := t.Counts[rank] + delta
	if v < 0 {
		v = 0
	}
	if v > max {
		v = max
	}
	t.Counts[rank] = v
	if err := s.dao.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

// Reset 按当前牌副数重置所有牌面为默认数量。
func (s *CardTrackerService) Reset(userId int64) (*model.CardTracker, error) {
	t, err := s.Get(userId)
	if err != nil {
		return nil, err
	}
	t.Counts = model.DefaultCounts(t.DeckCount)
	if err := s.dao.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

// SetDeckCount 设置牌副数并按新配置重置（对应「设置」页确认后的行为）。
func (s *CardTrackerService) SetDeckCount(userId int64, deckCount int) (*model.CardTracker, error) {
	if deckCount < 1 || deckCount > 10 {
		return nil, kit.NewInvalidArgumentError().WithInfo("牌副数需在 1-10 之间")
	}
	t, err := s.Get(userId)
	if err != nil {
		return nil, err
	}
	t.DeckCount = deckCount
	t.Counts = model.DefaultCounts(deckCount)
	if err := s.dao.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}
