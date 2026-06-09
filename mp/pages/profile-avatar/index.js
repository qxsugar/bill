const app = getApp()

Page({
  data: { avatars: [], selected: '' },
  onLoad(options) {},
  onShow() {},
  selectAvatar(e) {
    this.setData({ selected: e.currentTarget.dataset.url })
  },
  uploadCustom() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      success: (res) => {
        const path = res.tempFiles[0].tempFilePath
        this.setData({ selected: path })
      },
    })
  },
  confirm() {
    if (!this.data.selected) return wx.showToast({ title: '请选择头像', icon: 'none' })
    const userInfo = app.globalData.userInfo || {}
    app.globalData.userInfo = { ...userInfo, avatarUrl: this.data.selected }
    wx.navigateBack()
  },
})
