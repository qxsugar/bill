const app = getApp()

Page({
  data: {},
  onLoad(options) {},
  onShow() {},
  backHome() {
    wx.reLaunch({ url: '/pages/index/index' })
  },
})
