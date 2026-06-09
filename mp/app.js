const config = require('./config')

App({
  globalData: {
    userInfo: null,
    token: null,
  },

  onLaunch() {
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
    }
    this.login()
  },

  login() {
    wx.login({
      success: (res) => {
        if (!res.code) return
        wx.request({
          url: `${config.baseUrl}/api/auth/login`,
          method: 'POST',
          data: { code: res.code },
          success: (result) => {
            if (result.data && result.data.token) {
              this.globalData.token = result.data.token
              wx.setStorageSync('token', result.data.token)
              if (result.data.userInfo) {
                this.globalData.userInfo = result.data.userInfo
              }
            }
          },
        })
      },
    })
  },

  request(url, method = 'GET', data = {}) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: `${config.baseUrl}${url}`,
        method,
        data,
        header: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${this.globalData.token || ''}`,
        },
        success: (res) => {
          if (res.statusCode === 401) {
            wx.removeStorageSync('token')
            this.globalData.token = null
            this.login()
            reject(new Error('Unauthorized'))
            return
          }
          resolve(res.data)
        },
        fail: reject,
      })
    })
  },
})
