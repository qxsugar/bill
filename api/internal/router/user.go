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
//
//	@router			/api/v1/user.detail [get]
//	@summary		当前用户信息
//	@description	返回当前登录用户的信息
//	@tags			user
//	@produce		application/json
//	@security		BearerAuth
//	@success		200	{object}	kit.RespBody{resp_data=model.User}	"成功"
//	@failure		401	{object}	kit.RespBody						"未登录"
//	@failure		404	{object}	kit.RespBody						"用户不存在"
//	@failure		500	{object}	kit.RespBody						"服务器内部错误"
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
//
//	@router			/api/v1/user.update [post]
//	@summary		修改用户资料
//	@description	修改当前用户昵称/头像，昵称与头像至少传一个
//	@tags			user
//	@accept			application/json
//	@produce		application/json
//	@security		BearerAuth
//	@param			reqBody	body		updateProfileRequest				true	"资料更新请求"
//	@success		200		{object}	kit.RespBody{resp_data=model.User}	"成功"
//	@failure		400		{object}	kit.RespBody						"请求参数错误"
//	@failure		401		{object}	kit.RespBody						"未登录"
//	@failure		404		{object}	kit.RespBody						"用户不存在"
//	@failure		500		{object}	kit.RespBody						"服务器内部错误"
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
//
//	@router			/api/v1/user.presetAvatars [get]
//	@summary		预设头像列表
//	@description	返回可选的预设头像列表
//	@tags			user
//	@produce		application/json
//	@security		BearerAuth
//	@success		200	{object}	kit.RespBody{resp_data=object{avatars=[]string}}	"成功"
//	@failure		401	{object}	kit.RespBody										"未登录"
func (r *UserRouter) PresetAvatars(ctx *gin.Context) (any, error) {
	return gin.H{"avatars": r.userService.PresetAvatarList()}, nil
}

