package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Openid    string    `gorm:"type:varchar(64);uniqueIndex;not null;default:''" json:"openid"`
	Nickname  string    `gorm:"type:varchar(50);not null;default:''" json:"nickname"`
	Avatar    string    `gorm:"type:varchar(255);not null;default:''" json:"avatar"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string { return "bill_users" }
