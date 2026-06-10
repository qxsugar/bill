const app = getApp()
const api = require('../../api')

Page({
  data: { userInfo: {}, saving: false },

  onLoad() { this.loadUser() },
  onShow() { this.loadUser() },

  async loadUser() {
    try {
      const user = await api.userDetail()
      this.setData({ userInfo: user })
      app.globalData.userInfo = user
    } catch (e) {}
  },

  changeAvatar() {
    wx.navigateTo({ url: '/pages/profile-avatar/index' })
  },

  onNicknameInput(e) {
    this.setData({ 'userInfo.nickname': e.detail.value })
  },

  async saveProfile() {
    this.setData({ saving: true })
    try {
      const user = await api.userUpdate({
        nickname: this.data.userInfo.nickname,
        avatar: this.data.userInfo.avatar,
      })
      app.globalData.userInfo = user
      wx.setStorageSync('userInfo', user)
      wx.showToast({ title: '保存成功' })
    } catch (e) {
      wx.showToast({ title: e.message || '保存失败', icon: 'none' })
    } finally {
      this.setData({ saving: false })
    }
  },
})
