package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/internal/middleware"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/pkg/kit"
	"go.uber.org/zap"
)

type RoomRouter struct {
	roomService *service.RoomService
	logger      *zap.SugaredLogger
}

func NewRoomRouter(roomService *service.RoomService, logger *zap.SugaredLogger) *RoomRouter {
	return &RoomRouter{roomService: roomService, logger: logger}
}

// Create 创建房间，当前用户成为房主。
func (r *RoomRouter) Create(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	room, err := r.roomService.Create(userId)
	if err != nil {
		return nil, wrapErr(err)
	}
	return room, nil
}

type joinRequest struct {
	Code string `json:"code"`
}

// Join 通过房间码加入房间。
func (r *RoomRouter) Join(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req joinRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Code == "" {
		return nil, kit.NewInvalidArgumentError()
	}
	room, err := r.roomService.Join(userId, req.Code)
	if err != nil {
		return nil, wrapErr(err)
	}
	return room, nil
}

type roomIdRequest struct {
	RoomId int64 `json:"room_id"`
}

// Leave 离开房间。
func (r *RoomRouter) Leave(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req roomIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RoomId == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	if err := r.roomService.Leave(userId, req.RoomId); err != nil {
		return nil, wrapErr(err)
	}
	return gin.H{"ok": true}, nil
}
