Component({
  properties: {
    info: { type: Object, value: {} },
    size: { type: String, value: 'medium' },
  },
  computed: {
    initials() {
      const name = this.data.info && this.data.info.nickname
      return name ? name.charAt(0).toUpperCase() : '?'
    },
  },
  lifetimes: {
    attached() {
      const name = (this.data.info && this.data.info.nickname) || ''
      this.setData({ initials: name.charAt(0).toUpperCase() || '?' })
    },
  },
  observers: {
    info(val) {
      const name = (val && val.nickname) || ''
      this.setData({ initials: name.charAt(0).toUpperCase() || '?' })
    },
  },
})
