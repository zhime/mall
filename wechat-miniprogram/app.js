// app.js
App({
  globalData: {
    userInfo: null,
    token: null,
    apiBase: 'https://api.mall.com'
  },

  onLaunch: function () {
    // 小程序启动
    console.log('Mall小程序启动')
    
    // 检查登录状态
    this.checkLoginStatus()
    
    // 检查更新
    this.checkUpdate()
  },

  onShow: function (options) {
    // 小程序显示
    console.log('小程序显示', options)
  },

  onHide: function () {
    // 小程序隐藏
    console.log('小程序隐藏')
  },

  onError: function (msg) {
    // 错误监听
    console.error('小程序错误:', msg)
  },

  // 检查登录状态
  checkLoginStatus: function() {
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
      this.getUserInfo()
    }
  },

  // 获取用户信息
  getUserInfo: function() {
    wx.request({
      url: `${this.globalData.apiBase}/api/user/profile`,
      header: {
        'Authorization': `Bearer ${this.globalData.token}`
      },
      success: (res) => {
        if (res.data.code === 200) {
          this.globalData.userInfo = res.data.data
        }
      }
    })
  },

  // 检查更新
  checkUpdate: function() {
    const updateManager = wx.getUpdateManager()
    
    updateManager.onCheckForUpdate(function (res) {
      console.log('检查更新:', res.hasUpdate)
    })

    updateManager.onUpdateReady(function () {
      wx.showModal({
        title: '更新提示',
        content: '新版本已经准备好，是否重启应用？',
        success: function (res) {
          if (res.confirm) {
            updateManager.applyUpdate()
          }
        }
      })
    })

    updateManager.onUpdateFailed(function () {
      console.log('新版本下载失败')
    })
  },

  // 登录
  login: function(callback) {
    wx.login({
      success: (res) => {
        if (res.code) {
          wx.request({
            url: `${this.globalData.apiBase}/api/user/wechat/login`,
            method: 'POST',
            data: {
              code: res.code
            },
            success: (result) => {
              if (result.data.code === 200) {
                this.globalData.token = result.data.data.token
                this.globalData.userInfo = result.data.data.user
                wx.setStorageSync('token', result.data.data.token)
                wx.setStorageSync('userInfo', result.data.data.user)
                
                if (callback) callback(true)
              } else {
                if (callback) callback(false, result.data.message)
              }
            },
            fail: () => {
              if (callback) callback(false, '登录失败')
            }
          })
        }
      }
    })
  },

  // 退出登录
  logout: function() {
    this.globalData.token = null
    this.globalData.userInfo = null
    wx.removeStorageSync('token')
    wx.removeStorageSync('userInfo')
  }
})