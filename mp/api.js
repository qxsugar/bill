// 业务接口封装：统一走 app.request，路径对齐后端 /api/v1/*.* 命名。
const app = getApp()

const api = {
  // 用户
  userDetail: () => app.request('/api/v1/user.detail', 'GET'),
  userUpdate: (data) => app.request('/api/v1/user.update', 'POST', data),
  presetAvatars: () => app.request('/api/v1/user.presetAvatars', 'GET'),

  // 房间
  roomCreate: () => app.request('/api/v1/room.create', 'POST'),
  roomJoin: (code) => app.request('/api/v1/room.join', 'POST', { code }),
  roomLeave: (roomId) => app.request('/api/v1/room.leave', 'POST', { room_id: roomId }),
  roomDetail: (roomId) => app.request(`/api/v1/room.detail?room_id=${roomId}`, 'GET'),
  roomSettle: (roomId) => app.request('/api/v1/room.settle', 'POST', { room_id: roomId }),
  roomLogs: (roomId, limit = 50, offset = 0) =>
    app.request(`/api/v1/room.logs?room_id=${roomId}&limit=${limit}&offset=${offset}`, 'GET'),

  // 交易
  expense: (roomId, items) => app.request('/api/v1/transaction.expense', 'POST', { room_id: roomId, items }),
  revoke: (txId) => app.request('/api/v1/transaction.revoke', 'POST', { tx_id: txId }),
  thank: (txId) => app.request('/api/v1/transaction.thank', 'POST', { tx_id: txId }),

  // 记牌器
  cardDetail: () => app.request('/api/v1/card.detail', 'GET'),
  cardAdjust: (rank, delta) => app.request('/api/v1/card.adjust', 'POST', { rank, delta }),
  cardReset: () => app.request('/api/v1/card.reset', 'POST'),
  cardSetDeck: (deckCount) => app.request('/api/v1/card.setDeck', 'POST', { deck_count: deckCount }),
}

module.exports = api
