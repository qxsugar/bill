// Code generated manually to satisfy InitializeWorker.
//go:build !wireinject
// +build !wireinject

package worker

import (
	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/database"
	"github.com/qxsugar/bill/api/internal/logger"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/bill/api/internal/task"
)

func InitializeWorker() (*Worker, func(), error) {
	log, loggerCleanup, err := logger.NewLogger()
	if err != nil {
		return nil, nil, err
	}
	db, dbCleanup, err := database.NewDatabase()
	if err != nil {
		loggerCleanup()
		return nil, nil, err
	}
	rdb, redisCleanup, err := database.NewClient()
	if err != nil {
		dbCleanup()
		loggerCleanup()
		return nil, nil, err
	}

	asynqOpt := database.NewAsynqRedisOpt()
	asynqClient, asynqCleanup := database.NewAsynqClient(asynqOpt)

	userDao := dao.NewUserDao(db)
	roomDao := dao.NewRoomDao(db)
	memberDao := dao.NewRoomMemberDao(db)
	transactionDao := dao.NewTransactionDao(db)
	logDao := dao.NewRoomLogDao(db)

	roomService := service.NewRoomService(db, rdb, asynqClient, roomDao, memberDao, logDao, userDao, transactionDao)

	handler := task.NewHandler(roomService, log)
	w := NewWorker(asynqOpt, handler, log)

	cleanup := func() {
		asynqCleanup()
		redisCleanup()
		dbCleanup()
		loggerCleanup()
	}
	return w, cleanup, nil
}
