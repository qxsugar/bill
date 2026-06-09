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
    const cached = wx.getStorageSync('userInfo')
    if (cached) {
      this.globalData.userInfo = cached
    }
    this.login()
  },

  // 微信登录：code 换 token，存储登录态与用户信息。
  // 返回 Promise，便于页面在登录完成后再发起业务请求。
  login() {
    if (this._loginPromise) return this._loginPromise
    this._loginPromise = new Promise((resolve, reject) => {
      wx.login({
        success: (res) => {
          if (!res.code) {
            reject(new Error('wx.login 未返回 code'))
            return
          }
          wx.request({
            url: `${config.baseUrl}/api/v1/auth.login`,
            method: 'POST',
            data: { code: res.code },
            header: { 'Content-Type': 'application/json' },
            success: (result) => {
              const body = result.data || {}
              if (body.succeeded && body.resp_data) {
                const { token, user } = body.resp_data
                this.globalData.token = token
                this.globalData.userInfo = user
                wx.setStorageSync('token', token)
                wx.setStorageSync('userInfo', user)
                resolve(user)
              } else {
                reject(new Error(body.info || '登录失败'))
              }
            },
            fail: reject,
          })
        },
        fail: reject,
      })
    }).finally(() => {
      this._loginPromise = null
    })
    return this._loginPromise
  },

  // 确保已登录：有 token 直接返回，否则触发登录。
  ensureLogin() {
    if (this.globalData.token) return Promise.resolve(this.globalData.userInfo)
    return this.login()
  },

  // 统一请求：自动带 token，解包 RespBody，业务失败 reject 并带 info。
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
          const body = res.data || {}
          // token 失效：重新登录后由调用方重试
          if (body.code === 40100 || res.statusCode === 401) {
            wx.removeStorageSync('token')
            this.globalData.token = null
            this.login()
            reject(new Error('登录态失效，请重试'))
            return
          }
          if (body.succeeded) {
            resolve(body.resp_data)
          } else {
            const err = new Error(body.info || '请求失败')
            err.code = body.code
            reject(err)
          }
        },
        fail: (err) => reject(err),
      })
    })
  },
})
