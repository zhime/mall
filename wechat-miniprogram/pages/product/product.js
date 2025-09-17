// pages/product/product.js
const app = getApp()

Page({
  data: {
    product: null,
    selectedSpecs: {},
    quantity: 1,
    showSpecModal: false
  },

  onLoad: function (options) {
    const productId = options.id
    if (productId) {
      this.loadProduct(productId)
    }
  },

  loadProduct: function(id) {
    wx.showLoading({ title: '加载中...' })
    
    wx.request({
      url: `${app.globalData.apiBase}/api/products/${id}`,
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({
            product: res.data.data
          })
        } else {
          wx.showToast({ title: '商品不存在', icon: 'none' })
          wx.navigateBack()
        }
      },
      complete: () => {
        wx.hideLoading()
      }
    })
  },

  onImageTap: function(e) {
    const current = e.currentTarget.dataset.src
    wx.previewImage({
      current: current,
      urls: this.data.product.images
    })
  },

  onAddToCart: function() {
    if (!app.globalData.token) {
      wx.navigateTo({
        url: '/pages/login/login'
      })
      return
    }

    if (this.data.product.specs && this.data.product.specs.length > 0) {
      this.setData({ showSpecModal: true })
    } else {
      this.addToCart()
    }
  },

  addToCart: function() {
    wx.showToast({ title: '已加入购物车', icon: 'success' })
    this.setData({ showSpecModal: false })
  },

  onBuyNow: function() {
    if (!app.globalData.token) {
      wx.navigateTo({
        url: '/pages/login/login'
      })
      return
    }

    this.addToCart()
    wx.switchTab({
      url: '/pages/cart/cart'
    })
  }
})