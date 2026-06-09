package middleware

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AccessLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()
		path := ctx.FullPath()
		if path == "" {
			path = ctx.Request.URL.Path
		}

		logger := zap.S().Named("http").With(
			"request_id", requestid.Get(ctx),
			"method", ctx.Request.Method,
			"path", path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"ip", ctx.ClientIP(),
		)
		switch {
		case status >= 500:
			logger.Error("request completed")
		case status >= 400:
			logger.Warn("request completed")
		default:
			logger.Info("request completed")
		}
	}
}
