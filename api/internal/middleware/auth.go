package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qxsugar/bill/api/internal/service"
	"github.com/qxsugar/pkg/kit"
)

// CtxUserId 是写入 gin.Context 的当前用户 id 键。
const CtxUserId = "user_id"

// Auth 解析 Authorization: Bearer <token>，校验后把用户 id 写入上下文。
// 校验失败时直接以统一响应格式返回未认证错误并中断后续处理。
func Auth(authService *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := bearerToken(ctx)
		if token == "" {
			abortUnauthorized(ctx)
			return
		}
		userId, err := authService.ParseToken(token)
		if err != nil || userId == 0 {
			abortUnauthorized(ctx)
			return
		}
		ctx.Set(CtxUserId, userId)
		ctx.Next()
	}
}

// bearerToken 从请求头取出 Bearer token，也兼容 query 上的 token（WebSocket 握手用）。
func bearerToken(ctx *gin.Context) string {
	h := ctx.GetHeader("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	}
	return ctx.Query("token")
}

func abortUnauthorized(ctx *gin.Context) {
	ex := kit.NewUnauthenticatedError()
	ctx.AbortWithStatusJSON(http.StatusOK, kit.RespBody{
		Succeeded: false,
		Code:      ex.Code(),
		Info:      ex.Info(),
	})
}

// CurrentUserId 从上下文取出当前用户 id，未认证返回 0。
func CurrentUserId(ctx *gin.Context) int64 {
	v, ok := ctx.Get(CtxUserId)
	if !ok {
		return 0
	}
	id, _ := v.(int64)
	return id
}
