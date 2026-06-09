const app = getApp()

const RANKS = ['A','2','3','4','5','6','7','8','9','10','J','Q','K']
const SUITS = ['♠','♥','♦','♣']

Page({
  data: {
    suits: SUITS.map(name => ({ name, cards: RANKS, played: [] })),
  },
  onLoad(options) {},
  onShow() {},
  toggleCard(e) {
    const { suit, card } = e.currentTarget.dataset
    const suits = this.data.suits.map(s => {
      if (s.name !== suit) return s
      const played = s.played.includes(card)
        ? s.played.filter(c => c !== card)
        : [...s.played, card]
      return { ...s, played }
    })
    this.setData({ suits })
  },
  resetCards() {
    this.setData({
      suits: SUITS.map(name => ({ name, cards: RANKS, played: [] })),
    })
  },
})
