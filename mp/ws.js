// 房间 WebSocket：连接 /ws/room，收到 room_updated/settled 事件回调上层刷新。
// 内置断线自动重连，页面 onUnload 时调用 close 释放。
const config = require('./config')

function createRoomSocket(roomId, onEvent) {
  const app = getApp()
  const token = app.globalData.token || ''
  // baseUrl 形如 https://host，需替换为 wss://host
  const wsBase = config.baseUrl.replace(/^http/, 'ws')
  const url = `${wsBase}/ws/room?room_id=${roomId}&token=${token}`

  let socket = null
  let closedByUser = false
  let reconnectTimer = null

  const connect = () => {
    socket = wx.connectSocket({ url })
    socket.onMessage((res) => {
      try {
        const evt = JSON.parse(res.data)
        if (evt && evt.type) onEvent(evt)
      } catch (e) {
        // 忽略非 JSON 消息（如心跳）
      }
    })
    socket.onClose(() => {
      if (closedByUser) return
      // 3 秒后重连
      reconnectTimer = setTimeout(connect, 3000)
    })
    socket.onError(() => {
      // onClose 会随后触发，统一在那里处理重连
    })
  }

  connect()

  return {
    close() {
      closedByUser = true
      if (reconnectTimer) clearTimeout(reconnectTimer)
      if (socket) {
        try {
          socket.close()
        } catch (e) {}
      }
    },
  }
}

module.exports = { createRoomSocket }
