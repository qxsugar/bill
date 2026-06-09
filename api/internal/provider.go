package internal

import (
	"github.com/google/wire"
	"github.com/qxsugar/bill/api/internal/dao"
	"github.com/qxsugar/bill/api/internal/database"
	"github.com/qxsugar/bill/api/internal/logger"
	"github.com/qxsugar/bill/api/internal/router"
	"github.com/qxsugar/bill/api/internal/service"
)

var MiscProviderSet = wire.NewSet(
	logger.NewLogger,
	database.NewDatabase,
)

var DaoProviderSet = wire.NewSet(
	dao.NewUserDao,
	dao.NewRoomDao,
	dao.NewRoomMemberDao,
	dao.NewTransactionDao,
	dao.NewRoomLogDao,
	dao.NewCardTrackerDao,
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	service.NewAuthService,
	service.NewRoomService,
	service.NewTransactionService,
	service.NewCardTrackerService,
)

var RouterProviderSet = wire.NewSet(
	router.NewUserRouter,
	router.NewAuthRouter,
	router.NewRoomRouter,
	router.NewTransactionRouter,
	router.NewCardTrackerRouter,
)
