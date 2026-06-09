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
	_ = dao.NewRoomDao(db)
	_ = dao.NewRoomMemberDao(db)
	_ = dao.NewTransactionDao(db)
	_ = dao.NewRoomLogDao(db)

	userService := service.NewUserService(userDao)
	userRouter := router.NewUserRouter(userService, log)

	application := NewApplication(log, db, userRouter)
	cleanup := func() {
		dbCleanup()
		loggerCleanup()
	}
	return application, cleanup, nil
}
