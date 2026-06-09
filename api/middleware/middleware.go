package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return param.TimeStamp.Format(time.RFC3339) + " | " +
			param.Method + " " + param.Path + " | " +
			http.StatusText(param.StatusCode) + "\n"
	})
}

func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
