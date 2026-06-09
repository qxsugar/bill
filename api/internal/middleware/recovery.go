package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/qxsugar/pkg/kit"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			logger := zap.S().Named("recovery")
			if err := recover(); err != nil {
				if body, ok := err.(kit.RespBody); ok {
					ctx.PureJSON(http.StatusOK, body)
				} else {
					logger.With("request_id", requestid.Get(ctx)).Errorf("panic recovered: %v", err)
					logger.Debugf("stacktrace: %s", string(debug.Stack()))
					ctx.PureJSON(http.StatusOK, kit.RespBody{Succeeded: false, Code: kit.ErrUnknown, Info: "服务器内部错误"})
				}
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
