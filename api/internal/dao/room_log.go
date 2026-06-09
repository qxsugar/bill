package dao

import (
	"github.com/qxsugar/bill/api/internal/model"
	"gorm.io/gorm"
)

type RoomLogDao struct{ db *gorm.DB }

func NewRoomLogDao(db *gorm.DB) *RoomLogDao { return &RoomLogDao{db: db} }

func (d *RoomLogDao) ListByRoomId(roomId int64, limit, offset int) ([]*model.RoomLog, int64, error) {
	query := d.db.Model(&model.RoomLog{}).Where("room_id = ?", roomId)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*model.RoomLog
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (d *RoomLogDao) Create(item *model.RoomLog) error { return d.db.Create(item).Error }
