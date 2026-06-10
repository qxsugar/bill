const app = getApp()
const api = require('../../api')

Page({
  data: { roomId: 0, logs: [], loading: false },

  onLoad(options) {
    this.setData({ roomId: parseInt(options.room_id, 10) })
  },

  onShow() { this.loadLogs() },

  async loadLogs() {
    const { roomId } = this.data
    if (!roomId) return
    this.setData({ loading: true })
    try {
      const res = await api.roomLogs(roomId, 100, 0)
      const logs = (res.list || []).map((l) => ({
        ...l,
        time: formatTs(l.created_at),
      }))
      this.setData({ logs })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    } finally {
      this.setData({ loading: false })
    }
  },
})

function formatTs(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const pad = (n) => String(n).padStart(2, '0')
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}
