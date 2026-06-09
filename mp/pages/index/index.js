const app = getApp()

Page({
  data: {},
  onLoad(options) {},
  onShow() {},
  createRoom() {
    wx.navigateTo({ url: '/pages/room/index' })
  },
  joinRoom() {
    wx.showModal({
      title: '加入房间',
      editable: true,
      placeholderText: '请输入房间码',
      success(res) {
        if (res.confirm && res.content) {
          wx.navigateTo({ url: `/pages/room/index?code=${res.content}` })
        }
      },
    })
  },
})
