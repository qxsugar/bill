package task

import (
	"github.com/hibiken/asynq"
	"github.com/qxsugar/bill/api/internal/service"
	"go.uber.org/zap"
)

// Handler 汇聚各类异步任务的处理函数，由 worker 注册到 asynq mux。
type Handler struct {
	logger      *zap.SugaredLogger
	roomService *service.RoomService
}

func NewHandler(roomService *service.RoomService, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger:      logger.Named("[Task]"),
		roomService: roomService,
	}
}

// Register 将所有任务处理函数注册到 mux。
func (h *Handler) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(service.TaskAutoSettle, h.roomService.HandleAutoSettle)
}
