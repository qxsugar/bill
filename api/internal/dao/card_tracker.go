package dao

import (
	"errors"

	"github.com/qxsugar/bill/api/internal/model"
	"gorm.io/gorm"
)

type CardTrackerDao struct{ db *gorm.DB }

func NewCardTrackerDao(db *gorm.DB) *CardTrackerDao { return &CardTrackerDao{db: db} }

func (d *CardTrackerDao) FindByUserId(userId int64) (*model.CardTracker, error) {
	var item model.CardTracker
	err := d.db.Where("user_id = ?", userId).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *CardTrackerDao) Create(item *model.CardTracker) error { return d.db.Create(item).Error }
func (d *CardTrackerDao) Save(item *model.CardTracker) error   { return d.db.Save(item).Error }
