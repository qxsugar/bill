package dao

import (
	"github.com/qxsugar/bill/api/internal/model"
	"gorm.io/gorm"
)

type RoomMemberDao struct{ db *gorm.DB }

func NewRoomMemberDao(db *gorm.DB) *RoomMemberDao { return &RoomMemberDao{db: db} }

func (d *RoomMemberDao) ListByRoomId(roomId int64) ([]*model.RoomMember, error) {
	var list []*model.RoomMember
	return list, d.db.Where("room_id = ?", roomId).Order("id asc").Find(&list).Error
}

func (d *RoomMemberDao) FindByRoomAndUser(roomId, userId int64) (*model.RoomMember, error) {
	var item model.RoomMember
	err := d.db.Where("room_id = ? and user_id = ?", roomId, userId).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *RoomMemberDao) Create(item *model.RoomMember) error { return d.db.Create(item).Error }
func (d *RoomMemberDao) Save(item *model.RoomMember) error   { return d.db.Save(item).Error }
