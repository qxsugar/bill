package dao

import (
	"github.com/qxsugar/bill/api/internal/model"
	"gorm.io/gorm"
)

type TransactionDao struct{ db *gorm.DB }

func NewTransactionDao(db *gorm.DB) *TransactionDao { return &TransactionDao{db: db} }

func (d *TransactionDao) FindById(id int64) (*model.Transaction, error) {
	var item model.Transaction
	err := d.db.Where("id = ?", id).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *TransactionDao) ListByRoomId(roomId int64) ([]*model.Transaction, error) {
	var list []*model.Transaction
	return list, d.db.Where("room_id = ?", roomId).Order("id desc").Find(&list).Error
}

func (d *TransactionDao) Create(item *model.Transaction) error { return d.db.Create(item).Error }
func (d *TransactionDao) Save(item *model.Transaction) error   { return d.db.Save(item).Error }
