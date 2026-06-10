package router

import (
	"strconv"

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
//
//	@router			/api/v1/room.create [post]
//	@summary		创建房间
//	@description	创建房间，当前用户成为房主
//	@tags			room
//	@produce		application/json
//	@security		BearerAuth
//	@success		200	{object}	kit.RespBody{resp_data=model.Room}	"成功"
//	@failure		401	{object}	kit.RespBody						"未登录"
//	@failure		500	{object}	kit.RespBody						"服务器内部错误"
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
//
//	@router			/api/v1/room.join [post]
//	@summary		加入房间
//	@description	通过房间码加入房间
//	@tags			room
//	@accept			application/json
//	@produce		application/json
//	@security		BearerAuth
//	@param			reqBody	body		joinRequest							true	"加入房间请求"
//	@success		200		{object}	kit.RespBody{resp_data=model.Room}	"成功"
//	@failure		400		{object}	kit.RespBody						"请求参数错误"
//	@failure		401		{object}	kit.RespBody						"未登录"
//	@failure		500		{object}	kit.RespBody						"服务器内部错误"
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
//
//	@router			/api/v1/room.leave [post]
//	@summary		离开房间
//	@description	离开房间
//	@tags			room
//	@accept			application/json
//	@produce		application/json
//	@security		BearerAuth
//	@param			reqBody	body		roomIdRequest							true	"房间 ID 请求"
//	@success		200		{object}	kit.RespBody{resp_data=object{ok=bool}}	"成功"
//	@failure		400		{object}	kit.RespBody							"请求参数错误"
//	@failure		401		{object}	kit.RespBody							"未登录"
//	@failure		500		{object}	kit.RespBody							"服务器内部错误"
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

// Detail 返回房间快照（房间 + 成员 + 消息），供进入房间与 ws 推送后拉取。
//
//	@router			/api/v1/room.detail [get]
//	@summary		房间快照
//	@description	返回房间快照（房间 + 成员 + 消息），供进入房间与 ws 推送后拉取
//	@tags			room
//	@produce		application/json
//	@security		BearerAuth
//	@param			room_id	query		int												true	"房间 ID"
//	@success		200		{object}	kit.RespBody{resp_data=service.RoomSnapshot}	"成功"
//	@failure		400		{object}	kit.RespBody									"请求参数错误"
//	@failure		401		{object}	kit.RespBody									"未登录"
//	@failure		404		{object}	kit.RespBody									"房间不存在"
//	@failure		500		{object}	kit.RespBody									"服务器内部错误"
func (r *RoomRouter) Detail(ctx *gin.Context) (any, error) {
	roomId, err := strconv.ParseInt(ctx.Query("room_id"), 10, 64)
	if err != nil || roomId == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	snap, err := r.roomService.Snapshot(roomId)
	if err != nil {
		return nil, wrapErr(err)
	}
	return snap, nil
}

// Settle 结算房间（仅房主）。
//
//	@router			/api/v1/room.settle [post]
//	@summary		结算房间
//	@description	结算房间（仅房主）
//	@tags			room
//	@accept			application/json
//	@produce		application/json
//	@security		BearerAuth
//	@param			reqBody	body		roomIdRequest							true	"房间 ID 请求"
//	@success		200		{object}	kit.RespBody{resp_data=object{ok=bool}}	"成功"
//	@failure		400		{object}	kit.RespBody							"请求参数错误"
//	@failure		401		{object}	kit.RespBody							"未登录"
//	@failure		500		{object}	kit.RespBody							"服务器内部错误"
func (r *RoomRouter) Settle(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req roomIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RoomId == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	if err := r.roomService.Settle(userId, req.RoomId); err != nil {
		return nil, wrapErr(err)
	}
	return gin.H{"ok": true}, nil
}

// Logs 分页返回房间日志（从旧到新）。
//
//	@router			/api/v1/room.logs [get]
//	@summary		房间日志
//	@description	分页返回房间日志（从旧到新）
//	@tags			room
//	@produce		application/json
//	@security		BearerAuth
//	@param			room_id	query		int															true	"房间 ID"
//	@param			limit	query		int															false	"分页大小"	default(50)
//	@param			offset	query		int															false	"分页偏移"	default(0)
//	@success		200		{object}	kit.RespBody{resp_data=kit.PageBody{list=[]model.RoomLog}}	"成功"
//	@failure		400		{object}	kit.RespBody												"请求参数错误"
//	@failure		401		{object}	kit.RespBody												"未登录"
//	@failure		500		{object}	kit.RespBody												"服务器内部错误"
func (r *RoomRouter) Logs(ctx *gin.Context) (any, error) {
	roomId, err := strconv.ParseInt(ctx.Query("room_id"), 10, 64)
	if err != nil || roomId == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	list, total, err := r.roomService.Logs(roomId, limit, offset)
	if err != nil {
		return nil, wrapErr(err)
	}
	return kit.PageBody{Offset: offset, Limit: limit, Total: total, List: list}, nil
}
