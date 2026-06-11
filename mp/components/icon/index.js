// 统一图标组件：内联 SVG 线性图标，风格一致（24 viewBox / 圆角端点 / 描边）。
// 用法：<icon name="home" size="44" color="#3B5BDB" />
// 颜色默认 currentColor，可由父级文字色驱动；也可显式传 color。
const PATHS = {
  // 导航 / tabBar
  home: '<path d="M4 11.2 12 4l8 7.2"/><path d="M6 10v9a1 1 0 0 0 1 1h10a1 1 0 0 0 1-1v-9"/><path d="M10 20v-5h4v5"/>',
  'home-fill': '<path d="M3.6 11.5 12 4l8.4 7.5V20a1 1 0 0 1-1 1h-4.4v-5.5h-4V21H4.6a1 1 0 0 1-1-1z"/>',
  user: '<circle cx="12" cy="8" r="3.6"/><path d="M5 20c0-3.6 3.1-6 7-6s7 2.4 7 6"/>',
  'user-fill': '<circle cx="12" cy="7.8" r="4"/><path d="M4.5 20.5c0-3.9 3.4-6.5 7.5-6.5s7.5 2.6 7.5 6.5z"/>',

  // 动作
  plus: '<path d="M12 5v14M5 12h14"/>',
  enter: '<path d="M14 4h4a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-4"/><path d="M10 8l4 4-4 4"/><path d="M14 12H4"/>',
  cards: '<rect x="3" y="6" width="11" height="14" rx="2" transform="rotate(-10 8.5 13)"/><path d="M11 5.5l6 1.6a2 2 0 0 1 1.4 2.5l-2.4 9"/>',
  expense: '<circle cx="12" cy="12" r="8"/><path d="M12 7v10M9.4 9.2c0-1 1.2-1.7 2.6-1.7s2.6.7 2.6 1.7-1.2 1.6-2.6 1.6-2.6.7-2.6 1.7 1.2 1.7 2.6 1.7 2.6-.7 2.6-1.7"/>',
  log: '<rect x="4" y="3.5" width="16" height="17" rx="2"/><path d="M8 8h8M8 12h8M8 16h5"/>',
  edit: '<path d="M14.5 5.5l4 4L8 20H4v-4z"/><path d="M13 7l4 4"/>',
  settle: '<rect x="5" y="3" width="14" height="18" rx="2"/><path d="M8 7h8"/><path d="M8.5 11h0M12 11h0M15.5 11h0M8.5 14.5h0M12 14.5h0M15.5 14.5h0M8.5 18h0M12 18h0"/>',
  leave: '<path d="M10 4H6a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h4"/><path d="M16 8l4 4-4 4"/><path d="M20 12H9"/>',

  // 工具
  megaphone: '<path d="M4 10v4a1 1 0 0 0 1 1h2l8 4V5L7 9H5a1 1 0 0 0-1 1z"/><path d="M18 9a4 4 0 0 1 0 6"/>',
  settings: '<circle cx="12" cy="12" r="3"/><path d="M12 3v2.5M12 18.5V21M21 12h-2.5M5.5 12H3M18.4 5.6l-1.8 1.8M7.4 16.6l-1.8 1.8M18.4 18.4l-1.8-1.8M7.4 7.4 5.6 5.6"/>',
  reset: '<path d="M4 12a8 8 0 1 1 2.5 5.8"/><path d="M4 19v-5h5"/>',
  copy: '<rect x="9" y="9" width="11" height="11" rx="2"/><path d="M5 15H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1v1"/>',
  share: '<circle cx="6.5" cy="12" r="2.5"/><circle cx="17.5" cy="6" r="2.5"/><circle cx="17.5" cy="18" r="2.5"/><path d="M8.7 10.8l6.6-3.6M8.7 13.2l6.6 3.6"/>',
  check: '<path d="M5 12.5l4.5 4.5L19 7"/>',
  close: '<path d="M6 6l12 12M18 6 6 18"/>',
  revoke: '<path d="M4 7v5h5"/><path d="M5.5 12a7.5 7.5 0 1 1 1.8 4.9"/>',
  thank: '<path d="M12 20s-7-4.6-7-9.3A4.2 4.2 0 0 1 12 7.8 4.2 4.2 0 0 1 19 10.7C19 15.4 12 20 12 20z"/>',
}

function buildSvg(name, color, stroke) {
  const inner = PATHS[name] || ''
  const filled = name.endsWith('-fill') || name === 'thank' || name === 'home-fill'
  const fill = filled ? color : 'none'
  const strokeColor = filled ? 'none' : color
  return (
    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" ` +
    `fill="${fill}" stroke="${strokeColor}" stroke-width="${stroke}" ` +
    `stroke-linecap="round" stroke-linejoin="round">${inner}</svg>`
  )
}

Component({
  properties: {
    name: { type: String, value: '' },
    size: { type: null, value: 40 },        // rpx
    color: { type: String, value: '#1F2430' },
    stroke: { type: null, value: 2 },
  },
  data: { uri: '' },
  observers: {
    'name, color, stroke'() { this.render() },
  },
  lifetimes: {
    attached() { this.render() },
  },
  methods: {
    render() {
      const svg = buildSvg(this.data.name, this.data.color, this.data.stroke)
      const uri = 'data:image/svg+xml,' + encodeURIComponent(svg)
      this.setData({ uri })
    },
  },
})
