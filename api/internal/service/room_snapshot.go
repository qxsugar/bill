package service

import (
	"fmt"

	"github.com/qxsugar/bill/api/internal/model"
	"github.com/qxsugar/pkg/kit"
)

// MemberView 房间成员视图：成员积分余额 + 用户展示信息。
type MemberView struct {
	UserId   int64   `json:"user_id"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Balance  float64 `json:"balance"`
	IsOwner  bool    `json:"is_owner"`
	Left     bool    `json:"left"`
}

// MessageView 房间消息视图：一条有效交易渲染成的消息，支持撤销/感谢按钮判定。
type MessageView struct {
	Id         int64   `json:"id"`
	FromUserId int64   `json:"from_user_id"`
	FromName   string  `json:"from_name"`
	ToUserId   int64   `json:"to_user_id"`
	ToName     string  `json:"to_name"`
	Amount     float64 `json:"amount"`
	Thanked    bool    `json:"thanked"`
	CreatedAt  int64   `json:"created_at"`
	Text       string  `json:"text"`
}

// RoomSnapshot 房间快照：房间状态 + 成员列表 + 消息列表。
// 前端在进入房间和收到 ws 推送时拉取这一份完整状态。
type RoomSnapshot struct {
	Room     *model.Room    `json:"room"`
	Members  []*MemberView  `json:"members"`
	Messages []*MessageView `json:"messages"`
}

// Snapshot 组装房间快照。currentUserId 用于前端判断「自己发出的可撤销」等。
func (s *RoomService) Snapshot(roomId int64) (*RoomSnapshot, error) {
	room, err := s.roomDao.FindById(roomId)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, kit.NewNotFoundError().WithInfo("房间不存在")
	}

	members, err := s.memberDao.ListByRoomId(roomId)
	if err != nil {
		return nil, err
	}

	// 收集所有相关用户，批量取用户信息，避免 N+1。
	userIds := make([]int64, 0, len(members))
	for _, m := range members {
		userIds = append(userIds, m.UserId)
	}
	users, err := s.userDao.ListByIds(userIds)
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]*model.User, len(users))
	for _, u := range users {
		userMap[u.Id] = u
	}

	memberViews := make([]*MemberView, 0, len(members))
	for _, m := range members {
		u := userMap[m.UserId]
		mv := &MemberView{
			UserId:  m.UserId,
			Balance: m.Balance,
			IsOwner: m.UserId == room.OwnerId,
			Left:    m.LeftAt != nil,
		}
		if u != nil {
			mv.Nickname = u.Nickname
			mv.Avatar = u.Avatar
		}
		memberViews = append(memberViews, mv)
	}

	txs, err := s.transactionDao.ListValidByRoomId(roomId)
	if err != nil {
		return nil, err
	}
	messages := make([]*MessageView, 0, len(txs))
	for _, tx := range txs {
		from := userMap[tx.FromUserId]
		to := userMap[tx.ToUserId]
		mv := &MessageView{
			Id:         tx.Id,
			FromUserId: tx.FromUserId,
			ToUserId:   tx.ToUserId,
			Amount:     tx.Amount,
			Thanked:    tx.Thanked == model.TransactionThanked,
			CreatedAt:  tx.CreatedAt.Unix(),
		}
		if from != nil {
			mv.FromName = from.Nickname
		}
		if to != nil {
			mv.ToName = to.Nickname
		}
		mv.Text = fmt.Sprintf("%s 向 %s 发送 %g¥", mv.FromName, mv.ToName, tx.Amount)
		messages = append(messages, mv)
	}

	return &RoomSnapshot{Room: room, Members: memberViews, Messages: messages}, nil
}
