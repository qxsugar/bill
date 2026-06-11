const app = getApp()
const api = require('../../api')

const RANKS = ['BJ', 'SJ', 'A', 'K', 'Q', 'J', '10', '9', '8', '7', '6', '5', '4', '3', '2']
const DISPLAY = { BJ: '大王', SJ: '小王', A: 'A', K: 'K', Q: 'Q', J: 'J', '10': '10', '9': '9', '8': '8', '7': '7', '6': '6', '5': '5', '4': '4', '3': '3', '2': '2' }

Page({
  data: {
    cards: [],
    deckCount: 1,
    remaining: 0,
    tapRank: null,
    deckOptions: ['1 副牌', '2 副牌', '3 副牌', '4 副牌', '5 副牌', '6 副牌', '7 副牌', '8 副牌', '9 副牌', '10 副牌'],
    deckIndex: 0,
  },

  onLoad() { this.loadData() },
  onShow() { this.loadData() },

  async loadData() {
    try {
      const t = await api.cardDetail()
      this.renderState(t)
    } catch (e) {
      wx.showToast({ title: e.message || '加载失败', icon: 'none' })
    }
  },

  renderState(t) {
    const counts = t.counts || {}
    const deckCount = t.deck_count || 1
    let remaining = 0
    const cards = RANKS.map((r) => {
      const count = counts[r] !== undefined ? counts[r] : 0
      remaining += count
      const max = (r === 'BJ' || r === 'SJ') ? deckCount : 4 * deckCount
      return { rank: r, display: DISPLAY[r] || r, count, max }
    })
    this.setData({ cards, deckCount, remaining, deckIndex: deckCount - 1 })
  },

  // 单击扣减，快速双击归还（250ms 内第二次视为双击）
  onCardTap(e) {
    const rank = e.currentTarget.dataset.rank
    if (this._tapTimer && this._tapRank === rank) {
      clearTimeout(this._tapTimer)
      this._tapTimer = null
      this._tapRank = null
      this.adjust(rank, 1)
    } else {
      if (this._tapTimer) clearTimeout(this._tapTimer)
      this._tapRank = rank
      this._tapTimer = setTimeout(() => {
        this._tapTimer = null
        this._tapRank = null
        this.adjust(rank, -1)
      }, 250)
    }
  },

  async adjust(rank, delta) {
    try {
      const t = await api.cardAdjust(rank, delta)
      this.renderState(t)
    } catch (e) {
      wx.showToast({ title: e.message || '操作失败', icon: 'none' })
    }
  },

  async resetCards() {
    try {
      const t = await api.cardReset()
      this.renderState(t)
    } catch (e) {
      wx.showToast({ title: e.message || '重置失败', icon: 'none' })
    }
  },

  // 设置牌副数：系统滚动选择器选 1-10
  async onDeckChange(e) {
    const n = parseInt(e.detail.value, 10) + 1
    if (n === this.data.deckCount) return
    try {
      const t = await api.cardSetDeck(n)
      this.renderState(t)
    } catch (e) {
      wx.showToast({ title: e.message || '设置失败', icon: 'none' })
    }
  },
})
