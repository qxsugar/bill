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

// UpdateProfile 修改用户昵称/头像，任一为空则保留原值。
// 返回更新后的用户，供前端实时刷新展示。
func (s *UserService) UpdateProfile(id int64, nickname, avatar string) (*model.User, error) {
	user, err := s.userDao.FindById(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if err := s.userDao.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

// PresetAvatarList 返回预设头像列表，供「修改信息-头像选择」展示。
func (s *UserService) PresetAvatarList() []string {
	return PresetAvatars
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
