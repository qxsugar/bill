const app = getApp()
const api = require('../../api')
const { createRoomSocket } = require('../../ws')

Page({
  data: {
    roomId: 0,
    room: null,
    members: [],      // 在场成员（含自己）
    others: [],       // 除自己外的在场成员，用于批量支出
    messages: [],
    myUserId: 0,
    isOwner: false,
    // 支出弹窗
    popupVisible: false,
    popupMode: 'single',
    popupTargets: [],
  },

  onLoad(options) {
    const roomId = parseInt(options.room_id, 10)
    if (!roomId) {
      wx.showToast({ title: '房间不存在', icon: 'none' })
      return
    }
    const myUserId = (app.globalData.userInfo && app.globalData.userInfo.id) || 0
    this.setData({ roomId, myUserId })
    this.refresh()
    this.connectSocket()
  },

  onUnload() {
    if (this.socket) this.socket.close()
  },

  // 拉取房间快照并整理视图数据。
  async refresh() {
    try {
      const snap = await api.roomDetail(this.data.roomId)
      const myUserId = this.data.myUserId
      const activeMembers = (snap.members || []).filter((m) => !m.left)
      const others = activeMembers.filter((m) => m.user_id !== myUserId)
      this.setData({
        room: snap.room,
        members: activeMembers,
        others,
        messages: this.decorateMessages(snap.messages || []),
        isOwner: snap.room && snap.room.owner_id === myUserId,
      })
      // 房间已结算：跳转到已结算页
      if (snap.room && snap.room.status === 1) {
        this.goSettleDone()
      }
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  // 给消息打上「可撤销」（自己发出）/「可感谢」（自己收到且未感谢）标记。
  decorateMessages(messages) {
    const myUserId = this.data.myUserId
    return messages.map((m) => ({
      ...m,
      canRevoke: m.from_user_id === myUserId,
      canThank: m.to_user_id === myUserId && !m.thanked,
    }))
  },

  connectSocket() {
    this.socket = createRoomSocket(this.data.roomId, (evt) => {
      if (evt.type === 'settled') {
        this.goSettleDone()
      } else {
        this.refresh()
      }
    })
  },

  // 点击成员头像：向单个对象支出。
  onMemberTap(e) {
    const userId = e.currentTarget.dataset.userId
    if (userId === this.data.myUserId) return
    if (this.data.others.length === 0) {
      return wx.showToast({ title: '人数不足，无法支出', icon: 'none' })
    }
    const target = this.data.members.find((m) => m.user_id === userId)
    this.setData({ popupVisible: true, popupMode: 'single', popupTargets: [target] })
  },

  // 点击「支出」：向所有其他对象批量支出（均分/统一）。
  onExpenseTap() {
    if (this.data.others.length === 0) {
      return wx.showToast({ title: '人数不足，无法支出', icon: 'none' })
    }
    this.setData({ popupVisible: true, popupMode: 'batch', popupTargets: this.data.others })
  },

  closePopup() {
    this.setData({ popupVisible: false })
  },

  async onExpenseSubmit(e) {
    const { items } = e.detail
    this.setData({ popupVisible: false })
    try {
      await api.expense(this.data.roomId, items)
      // 广播会触发 refresh，这里也主动刷新一次以求即时
      this.refresh()
    } catch (err) {
      wx.showToast({ title: err.message || '支出失败', icon: 'none' })
    }
  },

  async onRevoke(e) {
    const txId = e.currentTarget.dataset.id
    try {
      await api.revoke(txId)
      this.refresh()
    } catch (err) {
      wx.showToast({ title: err.message || '撤销失败', icon: 'none' })
    }
  },

  async onThank(e) {
    const txId = e.currentTarget.dataset.id
    try {
      await api.thank(txId)
      this.refresh()
    } catch (err) {
      wx.showToast({ title: err.message || '操作失败', icon: 'none' })
    }
  },

  goInvite() {
    wx.navigateTo({ url: `/pages/room-invite/index?room_id=${this.data.roomId}&code=${this.data.room.code}` })
  },

  goLog() {
    wx.navigateTo({ url: `/pages/room-log/index?room_id=${this.data.roomId}` })
  },

  goCardTracker() {
    wx.navigateTo({ url: '/pages/card-tracker/index' })
  },

  goProfile() {
    wx.navigateTo({ url: '/pages/profile/index' })
  },

  // 结算：仅房主可进入结算页；非房主提示。
  goSettle() {
    if (!this.data.isOwner) {
      return wx.showToast({ title: '只有房主可以点击结算', icon: 'none' })
    }
    wx.navigateTo({ url: `/pages/settle/index?room_id=${this.data.roomId}` })
  },

  // 离开：非房主直接离开；房主提示需先结算。
  async onLeave() {
    if (this.data.isOwner) {
      return wx.showToast({ title: '房主不能主动离开，请点击结算后离开', icon: 'none' })
    }
    try {
      await api.roomLeave(this.data.roomId)
      wx.reLaunch({ url: '/pages/index/index' })
    } catch (err) {
      wx.showToast({ title: err.message || '离开失败', icon: 'none' })
    }
  },

  goSettleDone() {
    wx.redirectTo({ url: `/pages/settle-done/index?room_id=${this.data.roomId}` })
  },
})
