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

	userDao := dao.NewUserDao(db)
	roomDao := dao.NewRoomDao(db)
	memberDao := dao.NewRoomMemberDao(db)
	transactionDao := dao.NewTransactionDao(db)
	logDao := dao.NewRoomLogDao(db)
	_ = transactionDao

	userService := service.NewUserService(userDao)
	authService := service.NewAuthService(userService)
	roomService := service.NewRoomService(roomDao, memberDao, logDao, userDao)

	userRouter := router.NewUserRouter(userService, log)
	authRouter := router.NewAuthRouter(authService, log)
	roomRouter := router.NewRoomRouter(roomService, log)

	application := NewApplication(log, db, authService, userRouter, authRouter, roomRouter)
	cleanup := func() {
		dbCleanup()
		loggerCleanup()
	}
	return application, cleanup, nil
}
