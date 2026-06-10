const app = getApp()
const api = require('../../api')

Page({
  data: {
    roomId: 0,
    room: null,
    members: [],
    myUserId: 0,
    isOwner: false,
    settling: false,
  },

  onLoad(options) {
    const roomId = parseInt(options.room_id, 10)
    const myUserId = (app.globalData.userInfo && app.globalData.userInfo.id) || 0
    this.setData({ roomId, myUserId })
    this.refresh()
  },

  async refresh() {
    try {
      const snap = await api.roomDetail(this.data.roomId)
      this.setData({
        room: snap.room,
        members: snap.members || [],
        isOwner: snap.room && snap.room.owner_id === this.data.myUserId,
      })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  async confirmSettle() {
    if (!this.data.isOwner) {
      return wx.showToast({ title: '只有房主可以点击结算', icon: 'none' })
    }
    wx.showModal({
      title: '确认结算',
      content: '结算后不可再操作，所有用户将退出房间',
      success: async (res) => {
        if (!res.confirm) return
        this.setData({ settling: true })
        try {
          await api.roomSettle(this.data.roomId)
          wx.redirectTo({ url: `/pages/settle-done/index?room_id=${this.data.roomId}` })
        } catch (e) {
          wx.showToast({ title: e.message || '结算失败', icon: 'none' })
        } finally {
          this.setData({ settling: false })
        }
      },
    })
  },

  backToRoom() { wx.navigateBack() },
})
