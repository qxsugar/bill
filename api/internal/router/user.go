package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/internal/middleware"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/pkg/kit"
	"go.uber.org/zap"
)

type UserRouter struct {
	userService *service.UserService
	logger      *zap.SugaredLogger
}

func NewUserRouter(userService *service.UserService, logger *zap.SugaredLogger) *UserRouter {
	return &UserRouter{userService: userService, logger: logger}
}

// Detail 返回当前登录用户的信息。
func (r *UserRouter) Detail(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	user, err := r.userService.Detail(userId)
	if err != nil {
		return nil, kit.NewInternalError().WithErr(err)
	}
	if user == nil {
		return nil, kit.NewNotFoundError()
	}
	return user, nil
}

type updateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// Update 修改当前用户昵称/头像。
func (r *UserRouter) Update(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req updateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, kit.NewInvalidArgumentError()
	}
	if req.Nickname == "" && req.Avatar == "" {
		return nil, kit.NewInvalidArgumentError()
	}
	user, err := r.userService.UpdateProfile(userId, req.Nickname, req.Avatar)
	if err != nil {
		return nil, kit.NewInternalError().WithErr(err)
	}
	if user == nil {
		return nil, kit.NewNotFoundError()
	}
	return user, nil
}

// PresetAvatars 返回预设头像列表。
func (r *UserRouter) PresetAvatars(ctx *gin.Context) (any, error) {
	return gin.H{"avatars": r.userService.PresetAvatarList()}, nil
}

