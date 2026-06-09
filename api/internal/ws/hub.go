// Package ws 提供房间内的 WebSocket 实时广播。
// 每个房间维护一组连接，房间状态变化时由 service 层调用 Hub 广播事件，
// 客户端收到事件后重新拉取 room.detail 快照刷新界面。
package ws

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

// 事件类型
const (
	EventRoomUpdated = "room_updated" // 房间状态变化（支出/撤销/感谢/加入/离开）
	EventSettled     = "settled"      // 房间已结算，客户端应跳转到已结算页
)

// Event 推送给客户端的消息体。
type Event struct {
	Type   string `json:"type"`
	RoomId int64  `json:"room_id"`
}

// Hub 管理所有房间的连接集合，并实现 service.Broadcaster 接口。
type Hub struct {
	mu     sync.RWMutex
	rooms  map[int64]map[*Client]struct{}
	logger *zap.SugaredLogger
}

func NewHub(logger *zap.SugaredLogger) *Hub {
	return &Hub{
		rooms:  make(map[int64]map[*Client]struct{}),
		logger: logger.Named("ws"),
	}
}

// register 将连接加入房间。
func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[c.roomId] == nil {
		h.rooms[c.roomId] = make(map[*Client]struct{})
	}
	h.rooms[c.roomId][c] = struct{}{}
}

// unregister 将连接从房间移除，并关闭其发送通道。
func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.rooms[c.roomId]; ok {
		if _, exists := set[c]; exists {
			delete(set, c)
			close(c.send)
		}
		if len(set) == 0 {
			delete(h.rooms, c.roomId)
		}
	}
}

// broadcast 向房间内所有连接发送事件；发送缓冲满则丢弃该连接的本条消息。
func (h *Hub) broadcast(roomId int64, evt Event) {
	payload, err := json.Marshal(evt)
	if err != nil {
		h.logger.Warnf("marshal event failed: %v", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.rooms[roomId] {
		select {
		case c.send <- payload:
		default:
			// 客户端消费过慢，丢弃以避免阻塞广播
		}
	}
}

// BroadcastRoomUpdated 实现 service.Broadcaster。
func (h *Hub) BroadcastRoomUpdated(roomId int64) {
	h.broadcast(roomId, Event{Type: EventRoomUpdated, RoomId: roomId})
}

// BroadcastSettled 实现 service.Broadcaster。
func (h *Hub) BroadcastSettled(roomId int64) {
	h.broadcast(roomId, Event{Type: EventSettled, RoomId: roomId})
}
