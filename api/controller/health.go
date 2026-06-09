package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/bootstrap"
)

func Health(c *gin.Context) {
	sqlDB, err := bootstrap.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "database unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
