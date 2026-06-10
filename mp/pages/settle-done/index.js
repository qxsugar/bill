const app = getApp()
const api = require('../../api')

Page({
  data: { roomId: 0, room: null, members: [] },

  onLoad(options) {
    this.setData({ roomId: parseInt(options.room_id, 10) })
    this.loadData()
  },

  async loadData() {
    try {
      const snap = await api.roomDetail(this.data.roomId)
      this.setData({ room: snap.room, members: snap.members || [] })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  backHome() {
    wx.reLaunch({ url: '/pages/index/index' })
  },

  goLog() {
    wx.navigateTo({ url: `/pages/room-log/index?room_id=${this.data.roomId}` })
  },

  onShareAppMessage() {
    return { title: '打牌记账积分汇总', path: '/pages/index/index' }
  },
})
