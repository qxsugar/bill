package model

import "github.com/qxsugar/pkg/kit"

// User 用户表
type User struct {
	Id        int64         `json:"id" gorm:"column:id"`
	Openid    string        `json:"openid" gorm:"column:openid"`
	Nickname  string        `json:"nickname" gorm:"column:nickname"`
	Avatar    string        `json:"avatar" gorm:"column:avatar"`
	CreatedAt kit.TimeStamp `json:"created_at" gorm:"column:created_at"`
	UpdatedAt kit.TimeStamp `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string { return "bill_users" }
