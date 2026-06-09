package service

import (
	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/model"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (s *UserService) Detail(id int64) (*model.User, error) {
	return s.userDao.FindById(id)
}

// GetOrCreateByOpenid 按 openid 查找用户，不存在则创建（昵称首次随机生成）。
func (s *UserService) GetOrCreateByOpenid(openid, nickname, avatar string) (*model.User, error) {
	user, err := s.userDao.FindByOpenid(openid)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	user = &model.User{Openid: openid, Nickname: nickname, Avatar: avatar}
	if err := s.userDao.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
