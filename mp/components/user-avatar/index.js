Component({
  properties: {
    info: { type: Object, value: {} },
    size: { type: String, value: 'medium' },
  },
  data: {
    initials: '?',
    isEmoji: false,
    isImage: false,
    avatar: '',
  },
  lifetimes: {
    attached() {
      this.refresh(this.data.info)
    },
  },
  observers: {
    info(val) {
      this.refresh(val)
    },
  },
  methods: {
    // 头像可能是：图片 URL、emoji 文本，或为空时取昵称首字。
    refresh(info) {
      const avatar = (info && info.avatar) || ''
      const name = (info && info.nickname) || ''
      const isImage = /^https?:\/\//.test(avatar) || /^wxfile:\/\//.test(avatar) || avatar.startsWith('/')
      const isEmoji = !!avatar && !isImage
      this.setData({
        isImage,
        isEmoji,
        avatar,
        initials: name.charAt(0).toUpperCase() || '?',
      })
    },
  },
})
