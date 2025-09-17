// pages/index/index.js
const app = getApp()

Page({
  data: {
    banners: [],
    categories: [],
    products: [],
    page: 1,
    loading: false,
    hasMore: true
  },

  onLoad: function (options) {
    this.loadBanners()
    this.loadCategories()
    this.loadProducts()
  },

  onShow: function () {
    // 页面显示
  },

  onReachBottom: function () {
    // 上拉加载更多
    if (this.data.hasMore && !this.data.loading) {
      this.loadProducts()
    }
  },

  onPullDownRefresh: function () {
    // 下拉刷新
    this.setData({
      products: [],
      page: 1,
      hasMore: true
    })
    this.loadProducts()
    wx.stopPullDownRefresh()
  },

  // 加载轮播图
  loadBanners: function() {
    wx.request({
      url: `${app.globalData.apiBase}/api/banners`,
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({
            banners: res.data.data
          })
        }
      }
    })
  },

  // 加载分类
  loadCategories: function() {
    wx.request({
      url: `${app.globalData.apiBase}/api/categories`,
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({
            categories: res.data.data.slice(0, 8)
          })
        }
      }
    })
  },

  // 加载商品
  loadProducts: function() {
    if (this.data.loading) return
    
    this.setData({ loading: true })
    
    wx.request({
      url: `${app.globalData.apiBase}/api/products/recommend`,
      data: {
        page: this.data.page,
        page_size: 20
      },
      success: (res) => {
        if (res.data.code === 200) {
          const newProducts = res.data.data.list
          this.setData({
            products: this.data.page === 1 ? newProducts : [...this.data.products, ...newProducts],
            page: this.data.page + 1,
            hasMore: newProducts.length >= 20
          })
        }
      },
      complete: () => {
        this.setData({ loading: false })
      }
    })
  },

  // 搜索
  goToSearch: function() {
    wx.navigateTo({
      url: '/pages/search/search'
    })
  },

  // 轮播图点击
  onBannerTap: function(e) {
    const item = e.currentTarget.dataset.item
    if (item.link) {
      if (item.link.startsWith('/pages/')) {
        wx.navigateTo({
          url: item.link
        })
      } else if (item.link.startsWith('http')) {
        wx.navigateTo({
          url: `/pages/webview/webview?url=${encodeURIComponent(item.link)}`
        })
      }
    }
  },

  // 分类点击
  onCategoryTap: function(e) {
    const categoryId = e.currentTarget.dataset.id
    wx.switchTab({
      url: `/pages/category/category?id=${categoryId}`
    })
  },

  // 商品点击
  onProductTap: function(e) {
    const productId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/product/product?id=${productId}`
    })
  }
})