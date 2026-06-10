package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/internal/model"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/pkg/kit"
	"go.uber.org/zap"
)

type AuthRouter struct {
	authService *service.AuthService
	logger      *zap.SugaredLogger
}

func NewAuthRouter(authService *service.AuthService, logger *zap.SugaredLogger) *AuthRouter {
	return &AuthRouter{authService: authService, logger: logger}
}

type loginRequest struct {
	Code string `json:"code"`
}

type loginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Login 用微信 code 换取登录态，返回 JWT 与用户信息。
//
//	@router			/api/v1/auth.login [post]
//	@summary		微信登录
//	@description	用微信 code 换取登录态，返回 JWT 与用户信息
//	@tags			auth
//	@accept			application/json
//	@produce		application/json
//	@param			reqBody	body		loginRequest							true	"登录请求"
//	@success		200		{object}	kit.RespBody{resp_data=loginResponse}	"登录成功"
//	@failure		400		{object}	kit.RespBody							"请求参数错误"
//	@failure		500		{object}	kit.RespBody							"服务器内部错误"
func (r *AuthRouter) Login(ctx *gin.Context) (any, error) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Code == "" {
		return nil, kit.NewInvalidArgumentError()
	}
	token, user, err := r.authService.LoginByCode(req.Code)
	if err != nil {
		return nil, kit.NewInternalError().WithErr(err)
	}
	return loginResponse{Token: token, User: user}, nil
}
