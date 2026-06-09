package dao

import (
	"errors"

	"github.com/qxsugar/bill/api/internal/model"
	"gorm.io/gorm"
)

type RoomDao struct{ db *gorm.DB }

func NewRoomDao(db *gorm.DB) *RoomDao { return &RoomDao{db: db} }

func (d *RoomDao) FindById(id int64) (*model.Room, error) {
	var item model.Room
	err := d.db.Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *RoomDao) FindByCode(code string) (*model.Room, error) {
	var item model.Room
	err := d.db.Where("code = ?", code).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *RoomDao) Create(item *model.Room) error { return d.db.Create(item).Error }
func (d *RoomDao) Save(item *model.Room) error   { return d.db.Save(item).Error }
