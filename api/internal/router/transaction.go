package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/internal/middleware"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/pkg/kit"
	"go.uber.org/zap"
)

type TransactionRouter struct {
	txService *service.TransactionService
	logger    *zap.SugaredLogger
}

func NewTransactionRouter(txService *service.TransactionService, logger *zap.SugaredLogger) *TransactionRouter {
	return &TransactionRouter{txService: txService, logger: logger}
}

type expenseItemDTO struct {
	ToUserId int64   `json:"to_user_id"`
	Amount   float64 `json:"amount"`
}

type expenseRequest struct {
	RoomId int64            `json:"room_id"`
	Items  []expenseItemDTO `json:"items"`
}

// Expense 记录支出。前端三种模式（单笔/均分/统一）均归一为 items 明细列表提交。
func (r *TransactionRouter) Expense(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req expenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RoomId == 0 || len(req.Items) == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	items := make([]service.ExpenseItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, service.ExpenseItem{ToUserId: it.ToUserId, Amount: it.Amount})
	}
	if err := r.txService.Expense(userId, req.RoomId, items); err != nil {
		return nil, wrapErr(err)
	}
	return gin.H{"ok": true}, nil
}

type txIdRequest struct {
	TxId int64 `json:"tx_id"`
}

// Revoke 撤销自己发出的交易。
func (r *TransactionRouter) Revoke(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req txIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.TxId == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	if err := r.txService.Revoke(userId, req.TxId); err != nil {
		return nil, wrapErr(err)
	}
	return gin.H{"ok": true}, nil
}

// Thank 对收到的交易发送感谢。
func (r *TransactionRouter) Thank(ctx *gin.Context) (any, error) {
	userId := middleware.CurrentUserId(ctx)
	var req txIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.TxId == 0 {
		return nil, kit.NewInvalidArgumentError()
	}
	if err := r.txService.Thank(userId, req.TxId); err != nil {
		return nil, wrapErr(err)
	}
	return gin.H{"ok": true}, nil
}
