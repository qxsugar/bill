package router

import (
	"strconv"

	"github.com/gin-gonic/gin"
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

func (r *UserRouter) Detail(ctx *gin.Context) (any, error) {
	id, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		return nil, kit.NewInvalidArgumentError()
	}
	user, err := r.userService.Detail(id)
	if err != nil {
		return nil, kit.NewInternalError().WithErr(err)
	}
	if user == nil {
		return nil, kit.NewNotFoundError()
	}
	return user, nil
}
