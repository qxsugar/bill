package worker

import (
	"os"

	"github.com/hibiken/asynq"
	"github.com/qxsugar/bill/api/internal/task"
	"go.uber.org/zap"
)

// Worker 运行 asynq 服务端，消费异步任务（如房间超时自动结算）。
type Worker struct {
	logger        *zap.SugaredLogger
	asynqRedisOpt asynq.RedisClientOpt
	handler       *task.Handler
}

func NewWorker(
	opt asynq.RedisClientOpt,
	handler *task.Handler,
	logger *zap.SugaredLogger,
) *Worker {
	return &Worker{
		asynqRedisOpt: opt,
		handler:       handler,
		logger:        logger.Named("worker"),
	}
}

func (w *Worker) Start() {
	srv := asynq.NewServer(w.asynqRedisOpt, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	w.handler.Register(mux)

	w.logger.Info("worker started")
	if err := srv.Run(mux); err != nil {
		w.logger.Errorf("failed to run worker: %v", err)
		os.Exit(1)
	}
}
