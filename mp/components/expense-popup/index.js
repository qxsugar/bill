Component({
  properties: {
    visible: { type: Boolean, value: false },
  },
  data: { amount: '', remark: '' },
  methods: {
    close() { this.triggerEvent('close') },
    onAmountInput(e) { this.setData({ amount: e.detail.value }) },
    onRemarkInput(e) { this.setData({ remark: e.detail.value }) },
    submit() {
      const amount = parseFloat(this.data.amount)
      if (!amount) return wx.showToast({ title: '请输入金额', icon: 'none' })
      this.triggerEvent('submit', { amount, remark: this.data.remark })
      this.setData({ amount: '', remark: '' })
    },
  },
})
