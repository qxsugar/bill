const app = getApp()
const api = require('../../api')

Page({
  data: {
    notice: '欢迎使用打牌记账，公平记分不伤感情',
  },

  onLoad() {
    app.ensureLogin().catch(() => {})
  },

  // 创建房间：登录后调用后端，成功跳转到房间页。
  async createRoom() {
    try {
      await app.ensureLogin()
      wx.showLoading({ title: '创建中', mask: true })
      const room = await api.roomCreate()
      wx.hideLoading()
      wx.navigateTo({ url: `/pages/room/index?room_id=${room.id}` })
    } catch (e) {
      wx.hideLoading()
      wx.showToast({ title: e.message || '创建失败', icon: 'none' })
    }
  },

  // 加入房间：输入房间码 → 后端校验 → 进入房间页。
  joinRoom() {
    wx.showModal({
      title: '加入房间',
      editable: true,
      placeholderText: '请输入房间代码',
      success: async (res) => {
        if (!res.confirm || !res.content) return
        const code = res.content.trim()
        try {
          await app.ensureLogin()
          wx.showLoading({ title: '加入中', mask: true })
          const room = await api.roomJoin(code)
          wx.hideLoading()
          wx.navigateTo({ url: `/pages/room/index?room_id=${room.id}` })
        } catch (e) {
          wx.hideLoading()
          wx.showToast({ title: e.message || '加入失败', icon: 'none' })
        }
      },
    })
  },

  goCardTracker() {
    wx.navigateTo({ url: '/pages/card-tracker/index' })
  },
})
