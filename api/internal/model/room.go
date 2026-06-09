package model

import "github.com/qxsugar/pkg/kit"

// 房间状态
const (
	RoomStatusActive  = 0 // 活跃
	RoomStatusSettled = 1 // 已结算
)

// Room 房间表
type Room struct {
	Id        int64          `json:"id" gorm:"column:id"`
	Code      string         `json:"code" gorm:"column:code"`
	OwnerId   int64          `json:"owner_id" gorm:"column:owner_id"`
	Status    int            `json:"status" gorm:"column:status;default:0"`
	SettledAt *kit.TimeStamp `json:"settled_at" gorm:"column:settled_at"`
	CreatedAt kit.TimeStamp  `json:"created_at" gorm:"column:created_at"`
}

func (Room) TableName() string { return "bill_rooms" }

// RoomMember 房间成员表
type RoomMember struct {
	Id       int64          `json:"id" gorm:"column:id"`
	RoomId   int64          `json:"room_id" gorm:"column:room_id"`
	UserId   int64          `json:"user_id" gorm:"column:user_id"`
	Balance  float64        `json:"balance" gorm:"column:balance"`
	JoinedAt kit.TimeStamp  `json:"joined_at" gorm:"column:joined_at"`
	LeftAt   *kit.TimeStamp `json:"left_at" gorm:"column:left_at"`
}

func (RoomMember) TableName() string { return "bill_room_members" }

// RoomLog 房间日志表
const (
	LogTypeJoin     = "join"
	LogTypeLeave    = "leave"
	LogTypeTransfer = "transfer"
	LogTypeRevoke   = "revoke"
	LogTypeThanks   = "thanks"
	LogTypeSettle   = "settle"
)

type RoomLog struct {
	Id        int64         `json:"id" gorm:"column:id"`
	RoomId    int64         `json:"room_id" gorm:"column:room_id"`
	UserId    *int64        `json:"user_id" gorm:"column:user_id"`
	Content   string        `json:"content" gorm:"column:content"`
	LogType   string        `json:"log_type" gorm:"column:log_type"`
	CreatedAt kit.TimeStamp `json:"created_at" gorm:"column:created_at"`
}

func (RoomLog) TableName() string { return "bill_room_logs" }
