package model

import "github.com/qxsugar/pkg/kit"

// 交易状态
const (
	TransactionStatusValid   = 0 // 有效
	TransactionStatusRevoked = 1 // 已撤销

	TransactionNotThanked = 0 // 未感谢
	TransactionThanked    = 1 // 已感谢
)

// Transaction 交易记录表
type Transaction struct {
	Id         int64          `json:"id" gorm:"column:id"`
	RoomId     int64          `json:"room_id" gorm:"column:room_id"`
	FromUserId int64          `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserId   int64          `json:"to_user_id" gorm:"column:to_user_id"`
	Amount     float64        `json:"amount" gorm:"column:amount"`
	Status     int            `json:"status" gorm:"column:status;default:0"`
	Thanked    int            `json:"thanked" gorm:"column:thanked;default:0"`
	CreatedAt  kit.TimeStamp  `json:"created_at" gorm:"column:created_at"`
	RevokedAt  *kit.TimeStamp `json:"revoked_at" gorm:"column:revoked_at"`
}

func (Transaction) TableName() string { return "bill_transactions" }
