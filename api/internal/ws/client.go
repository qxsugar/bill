package ws

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qxsugar/bill/api/internal/middleware"
	"github.com/qxsugar/bill/api/internal/service"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
	sendBuffer     = 16
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 小程序与 H5 客户端跨域，统一放行（鉴权已由 token 完成）。
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client 代表一个房间内的 WebSocket 连接。
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	roomId int64
	userId int64
}

// Handler 返回 gin 处理函数：校验 token + room_id，升级连接并接入 Hub。
// 鉴权复用 middleware.Auth：握手前已注入 user_id。
func Handler(hub *Hub, authService *service.AuthService) gin.HandlerFunc {
	authMw := middleware.Auth(authService)
	return func(ctx *gin.Context) {
		// 先走鉴权中间件（支持 query 上的 token），失败会中断。
		authMw(ctx)
		if ctx.IsAborted() {
			return
		}
		userId := middleware.CurrentUserId(ctx)
		roomId, err := strconv.ParseInt(ctx.Query("room_id"), 10, 64)
		if err != nil || roomId == 0 {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			hub.logger.Warnf("ws upgrade failed: %v", err)
			return
		}

		c := &Client{hub: hub, conn: conn, send: make(chan []byte, sendBuffer), roomId: roomId, userId: userId}
		hub.register(c)
		go c.writePump()
		go c.readPump()
	}
}

// readPump 处理心跳并在连接断开时注销。客户端不需要发业务消息，读到即丢弃。
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister(c)
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

// writePump 把 send 通道里的事件写给客户端，并定期发送 ping 维持连接。
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
