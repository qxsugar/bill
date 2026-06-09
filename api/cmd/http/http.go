package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/internal/middleware"
	"github.com/qxsugar/bill/api/internal/router"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/pkg/kit"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	g           *gin.Engine
	logger      *zap.SugaredLogger
	db          *gorm.DB
	authService *service.AuthService
	userRouter  *router.UserRouter
	authRouter  *router.AuthRouter
}

func NewApplication(
	logger *zap.SugaredLogger,
	db *gorm.DB,
	authService *service.AuthService,
	userRouter *router.UserRouter,
	authRouter *router.AuthRouter,
) *Application {
	return &Application{
		g:           gin.New(),
		logger:      logger.Named("gateway"),
		db:          db,
		authService: authService,
		userRouter:  userRouter,
		authRouter:  authRouter,
	}
}

func (app *Application) Start() {
	svr := http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port")),
		Handler: app.g,
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Errorf("server listen error: %v", err)
			os.Exit(1)
		}
	}()

	app.registerInfra()
	app.registerApi()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		app.logger.Errorf("server forced to shutdown: %v", err)
		os.Exit(1)
	}
	app.logger.Info("server exited gracefully")
}

func (app *Application) registerInfra() {
	app.g.Use(requestid.New())
	app.g.Use(middleware.Cors())
	app.g.Use(middleware.AccessLogger(), gin.Recovery(), middleware.Recovery())
}

func (app *Application) registerApi() {
	// 健康检查
	app.g.GET("/health", func(ctx *gin.Context) {
		sqlDB, err := app.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	app.g.GET("/ping", kit.TranslateFunc(func(ctx *gin.Context) (any, error) { return "pong", nil }))

	api := app.g.Group("/api/v1")
	{
		// 公开接口（无需登录）
		api.POST("/auth.login", kit.TranslateFunc(app.authRouter.Login))

		// 需要登录的接口
		authed := api.Group("")
		authed.Use(middleware.Auth(app.authService))
		{
			authed.GET("/user.detail", kit.TranslateFunc(app.userRouter.Detail))
		}
	}
}
