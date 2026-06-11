// Code generated manually to satisfy InitializeApplication.
//go:build !wireinject
// +build !wireinject

package http

import (
	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/database"
	"github.com/qxsugar/bill/api/internal/logger"
	"github.com/qxsugar/bill/api/internal/router"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/bill/api/internal/weapp"
	"github.com/qxsugar/bill/api/internal/ws"
)

func InitializeApplication() (*Application, func(), error) {
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

	weappClient := weapp.NewClient(rdb)

	userDao := dao.NewUserDao(db)
	roomDao := dao.NewRoomDao(db)
	memberDao := dao.NewRoomMemberDao(db)
	transactionDao := dao.NewTransactionDao(db)
	logDao := dao.NewRoomLogDao(db)
	cardTrackerDao := dao.NewCardTrackerDao(db)

	userService := service.NewUserService(userDao)
	authService := service.NewAuthService(userService, weappClient)
	roomService := service.NewRoomService(db, roomDao, memberDao, logDao, userDao, transactionDao)
	transactionService := service.NewTransactionService(db, roomDao, memberDao, transactionDao, logDao, userDao)
	cardTrackerService := service.NewCardTrackerService(cardTrackerDao)

	// WebSocket Hub 作为广播器注入房间/交易服务
	hub := ws.NewHub(log)
	roomService.SetBroadcaster(hub)
	transactionService.SetBroadcaster(hub)

	userRouter := router.NewUserRouter(userService, log)
	authRouter := router.NewAuthRouter(authService, log)
	roomRouter := router.NewRoomRouter(roomService, log)
	transactionRouter := router.NewTransactionRouter(transactionService, log)
	cardTrackerRouter := router.NewCardTrackerRouter(cardTrackerService, log)

	application := NewApplication(log, db, authService, userRouter, authRouter, roomRouter, transactionRouter, cardTrackerRouter, hub)
	cleanup := func() {
		redisCleanup()
		dbCleanup()
		loggerCleanup()
	}
	return application, cleanup, nil
}
