const app = getApp()

Page({
  data: { roomCode: '', qrUrl: '' },
  onLoad(options) {},
  onShow() {},
  copyCode() {
    wx.setClipboardData({ data: this.data.roomCode })
  },
})
