const app = getApp()

Page({
  data: { roomCode: '', members: [] },
  onLoad(options) {
    if (options.code) this.setData({ roomCode: options.code })
  },
  onShow() {},
  goInvite() { wx.navigateTo({ url: '/pages/room-invite/index' }) },
  goLog() { wx.navigateTo({ url: '/pages/room-log/index' }) },
  goSettle() { wx.navigateTo({ url: '/pages/settle/index' }) },
})
