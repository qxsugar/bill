const app = getApp()

Page({
  data: { roomId: 0, code: '' },

  onLoad(options) {
    this.setData({
      roomId: parseInt(options.room_id, 10) || 0,
      code: options.code || '',
    })
  },

  copyCode() {
    wx.setClipboardData({
      data: this.data.code,
      success: () => wx.showToast({ title: '已复制' }),
    })
  },

  onShareAppMessage() {
    const code = this.data.code
    return {
      title: `${app.globalData.userInfo && app.globalData.userInfo.nickname || '朋友'} 邀请你加入房间一起玩耍`,
      path: `/pages/index/index`,
      imageUrl: '',
    }
  },
})
