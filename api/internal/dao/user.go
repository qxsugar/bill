package dao

import (
	"errors"

	"github.com/qxsugar/bill/api/internal/model"
	"gorm.io/gorm"
)

type UserDao struct{ db *gorm.DB }

func NewUserDao(db *gorm.DB) *UserDao { return &UserDao{db: db} }

func (d *UserDao) FindById(id int64) (*model.User, error) {
	var item model.User
	err := d.db.Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *UserDao) FindByOpenid(openid string) (*model.User, error) {
	var item model.User
	err := d.db.Where("openid = ?", openid).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *UserDao) ListByIds(ids []int64) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}
	var list []*model.User
	return list, d.db.Where("id in ?", ids).Find(&list).Error
}

func (d *UserDao) Save(item *model.User) error   { return d.db.Save(item).Error }
func (d *UserDao) Create(item *model.User) error { return d.db.Create(item).Error }
