const app = getApp()

Page({
  data: { items: [] },
  onLoad(options) {},
  onShow() {},
  confirmSettle() {
    wx.navigateTo({ url: '/pages/settle-done/index' })
  },
})
