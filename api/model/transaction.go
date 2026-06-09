package model

import "time"

type Transaction struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID     int64      `gorm:"not null;default:0;index" json:"room_id"`
	FromUserID int64      `gorm:"not null;default:0;index" json:"from_user_id"`
	ToUserID   int64      `gorm:"not null;default:0;index" json:"to_user_id"`
	Amount     float64    `gorm:"type:decimal(10,2);not null;default:0.00" json:"amount"`
	Status     int8       `gorm:"not null;default:0" json:"status"` // 0=有效 1=已撤销
	Thanked    int8       `gorm:"not null;default:0" json:"thanked"` // 0=未感谢 1=已感谢
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	RevokedAt  *time.Time `gorm:"default:null" json:"revoked_at"`
}

func (Transaction) TableName() string { return "bill_transactions" }
