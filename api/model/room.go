package model

import "time"

type Room struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code       string     `gorm:"type:varchar(10);uniqueIndex;not null;default:''" json:"code"`
	OwnerID    int64      `gorm:"not null;default:0;index" json:"owner_id"`
	Status     int8       `gorm:"not null;default:0" json:"status"` // 0=活跃 1=已结算
	SettledAt  *time.Time `gorm:"default:null" json:"settled_at"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (Room) TableName() string { return "bill_rooms" }

type RoomMember struct {
	ID       int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID   int64      `gorm:"not null;default:0;index" json:"room_id"`
	UserID   int64      `gorm:"not null;default:0;index" json:"user_id"`
	Balance  float64    `gorm:"type:decimal(10,2);not null;default:0.00" json:"balance"`
	JoinedAt time.Time  `gorm:"autoCreateTime" json:"joined_at"`
	LeftAt   *time.Time `gorm:"default:null" json:"left_at"`
}

func (RoomMember) TableName() string { return "bill_room_members" }

type RoomLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID    int64     `gorm:"not null;default:0;index" json:"room_id"`
	UserID    *int64    `gorm:"default:null" json:"user_id"`
	Content   string    `gorm:"type:varchar(255);not null;default:''" json:"content"`
	LogType   string    `gorm:"type:varchar(30);not null;default:''" json:"log_type"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (RoomLog) TableName() string { return "bill_room_logs" }
