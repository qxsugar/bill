const app = getApp()

Page({
  data: { userInfo: {} },
  onLoad(options) {
    this.setData({ userInfo: app.globalData.userInfo || {} })
  },
  onShow() {},
  changeAvatar() {
    wx.navigateTo({ url: '/pages/profile-avatar/index' })
  },
  onNicknameInput(e) {
    this.setData({ 'userInfo.nickname': e.detail.value })
  },
  saveProfile() {
    app.globalData.userInfo = this.data.userInfo
    wx.showToast({ title: '保存成功' })
  },
})
