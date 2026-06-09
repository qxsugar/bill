package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/controller"
	"github.com/qxsugar/bill/api/middleware"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger(), middleware.Recovery())

	r.GET("/health", controller.Health)

	v1 := r.Group("/api/v1")
	{
		_ = v1 // routes will be added here
	}

	return r
}
