// 支出弹窗：支持两种来源
//  - 点击单个头像：targets 为该用户，单笔金额
//  - 点击「支出」：targets 为多个用户，支持「均分」(总额÷人数) 与「统一」(每人相同) 两种模式
// 提交时统一 triggerEvent('submit', { items: [{ to_user_id, amount }] })。
Component({
  properties: {
    visible: { type: Boolean, value: false },
    // 可支出对象：[{ user_id, nickname, avatar }]
    targets: { type: Array, value: [] },
    // 'single' 单个对象；'batch' 多对象（均分/统一）
    mode: { type: String, value: 'single' },
  },
  data: {
    amount: '',          // single / batch-统一 共用；batch-均分时为总额
    batchMode: 'equal',  // equal=均分, same=统一
  },
  methods: {
    close() {
      this.reset()
      this.triggerEvent('close')
    },
    reset() {
      this.setData({ amount: '', batchMode: 'equal' })
    },
    onAmountInput(e) {
      this.setData({ amount: e.detail.value })
    },
    switchBatchMode(e) {
      this.setData({ batchMode: e.currentTarget.dataset.mode })
    },
    submit() {
      const amount = parseFloat(this.data.amount)
      if (!amount || amount <= 0) {
        return wx.showToast({ title: '请输入金额', icon: 'none' })
      }
      const targets = this.data.targets || []
      if (targets.length === 0) {
        return wx.showToast({ title: '没有可支出对象', icon: 'none' })
      }

      let items = []
      if (this.data.mode === 'single') {
        items = [{ to_user_id: targets[0].user_id, amount }]
      } else if (this.data.batchMode === 'equal') {
        // 均分：总额平摊到每人，保留两位小数
        const per = Math.round((amount / targets.length) * 100) / 100
        items = targets.map((t) => ({ to_user_id: t.user_id, amount: per }))
      } else {
        // 统一：每人相同金额
        items = targets.map((t) => ({ to_user_id: t.user_id, amount }))
      }

      this.triggerEvent('submit', { items })
      this.reset()
    },
  },
})
