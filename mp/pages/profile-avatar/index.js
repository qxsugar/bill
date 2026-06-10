const app = getApp()
const api = require('../../api')

Page({
  data: { avatars: [], selected: '' },

  onLoad() { this.loadAvatars() },

  async loadAvatars() {
    try {
      const res = await api.presetAvatars()
      const current = (app.globalData.userInfo && app.globalData.userInfo.avatar) || ''
      this.setData({ avatars: res.avatars || [], selected: current })
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  selectAvatar(e) {
    this.setData({ selected: e.currentTarget.dataset.emoji })
  },

  async confirm() {
    if (!this.data.selected) return wx.showToast({ title: '请选择头像', icon: 'none' })
    try {
      const user = await api.userUpdate({ avatar: this.data.selected })
      app.globalData.userInfo = user
      wx.setStorageSync('userInfo', user)
      wx.navigateBack()
    } catch (e) {
      wx.showToast({ title: e.message || '保存失败', icon: 'none' })
    }
  },
})
