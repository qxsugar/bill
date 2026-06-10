package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/internal/middleware"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/pkg/kit"
	"go.uber.org/zap"
)

type CardTrackerRouter struct {
	svc    *service.CardTrackerService
	logger *zap.SugaredLogger
}

func NewCardTrackerRouter(svc *service.CardTrackerService, logger *zap.SugaredLogger) *CardTrackerRouter {
	return &CardTrackerRouter{svc: svc, logger: logger}
}

// Detail 返回当前用户记牌器状态（不存在则初始化）。
//
//	@router			/api/v1/card.detail [get]
//	@summary		记牌器状态
//	@description	返回当前用户记牌器状态（不存在则初始化）
//	@tags			card
//	@produce		application/json
//	@security		BearerAuth
//	@success		200	{object}	kit.RespBody{resp_data=model.CardTracker}	"成功"
//	@failure		401	{object}	kit.RespBody								"未登录"
//	@failure		500	{object}	kit.RespBody								"服务器内部错误"
func (r *CardTrackerRouter) Detail(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	t, err := r.svc.Get(userId)
	if err != nil {
		return nil, wrapErr(err)
	}
	return t, nil
}

type adjustRequest struct {
	Rank  string `json:"rank"`
	Delta int    `json:"delta"`
}

// Adjust 调整某牌面剩余数量：delta=-1 点击扣除，delta=+1 双击增加。
//
//	@router			/api/v1/card.adjust [post]
//	@summary		调整牌面数量
//	@description	调整某牌面剩余数量：delta=-1 点击扣除，delta=+1 双击增加
//	@tags			card
//	@accept			application/json
//	@produce		application/json
//	@security		BearerAuth
//	@param			reqBody	body		adjustRequest								true	"调整请求"
//	@success		200		{object}	kit.RespBody{resp_data=model.CardTracker}	"成功"
//	@failure		400		{object}	kit.RespBody								"请求参数错误"
//	@failure		401		{object}	kit.RespBody								"未登录"
//	@failure		500		{object}	kit.RespBody								"服务器内部错误"
func (r *CardTrackerRouter) Adjust(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req adjustRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Rank == "" || req.Delta == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	// 归一化为 ±1，避免越级调整
	step := 1
	if req.Delta < 0 {
		step = -1
	}
	t, err := r.svc.Adjust(userId, req.Rank, step)
	if err != nil {
		return nil, wrapErr(err)
	}
	return t, nil
}

// Reset 按当前牌副数重置。
//
//	@router			/api/v1/card.reset [post]
//	@summary		重置记牌器
//	@description	按当前牌副数重置
//	@tags			card
//	@produce		application/json
//	@security		BearerAuth
//	@success		200	{object}	kit.RespBody{resp_data=model.CardTracker}	"成功"
//	@failure		401	{object}	kit.RespBody								"未登录"
//	@failure		500	{object}	kit.RespBody								"服务器内部错误"
func (r *CardTrackerRouter) Reset(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	t, err := r.svc.Reset(userId)
	if err != nil {
		return nil, wrapErr(err)
	}
	return t, nil
}

type deckRequest struct {
	DeckCount int `json:"deck_count"`
}

// SetDeck 设置牌副数并重置（设置页确认）。
//
//	@router			/api/v1/card.setDeck [post]
//	@summary		设置牌副数
//	@description	设置牌副数并重置（设置页确认）
//	@tags			card
//	@accept			application/json
//	@produce		application/json
//	@security		BearerAuth
//	@param			reqBody	body		deckRequest									true	"牌副数请求"
//	@success		200		{object}	kit.RespBody{resp_data=model.CardTracker}	"成功"
//	@failure		400		{object}	kit.RespBody								"请求参数错误"
//	@failure		401		{object}	kit.RespBody								"未登录"
//	@failure		500		{object}	kit.RespBody								"服务器内部错误"
func (r *CardTrackerRouter) SetDeck(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req deckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.DeckCount == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	t, err := r.svc.SetDeckCount(userId, req.DeckCount)
	if err != nil {
		return nil, wrapErr(err)
	}
	return t, nil
}
