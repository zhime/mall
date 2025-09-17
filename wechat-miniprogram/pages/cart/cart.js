// pages/cart/cart.js
const app = getApp()

Page({
  data: {
    cartItems: [],
    selectAll: false,
    totalPrice: 0,
    selectedCount: 0
  },

  onLoad: function () {
    this.loadCartItems()
  },

  onShow: function () {
    this.loadCartItems()
  },

  loadCartItems: function() {
    if (!app.globalData.token) {
      wx.navigateTo({
        url: '/pages/login/login'
      })
      return
    }

    // 模拟购物车数据
    const cartItems = [
      {
        id: 1,
        product_id: 1,
        name: '商品名称1',
        image: '/images/product1.jpg',
        price: 99.00,
        quantity: 1,
        selected: false,
        stock: 100
      }
    ]
    
    this.setData({ cartItems })
    this.calculateTotal()
  },

  onSelectItem: function(e) {
    const index = e.currentTarget.dataset.index
    const cartItems = this.data.cartItems
    cartItems[index].selected = !cartItems[index].selected
    
    this.setData({ cartItems })
    this.calculateTotal()
  },

  onSelectAll: function() {
    const selectAll = !this.data.selectAll
    const cartItems = this.data.cartItems.map(item => ({
      ...item,
      selected: selectAll
    }))
    
    this.setData({ cartItems, selectAll })
    this.calculateTotal()
  },

  calculateTotal: function() {
    const selectedItems = this.data.cartItems.filter(item => item.selected)
    const totalPrice = selectedItems.reduce((sum, item) => sum + item.price * item.quantity, 0)
    const selectedCount = selectedItems.reduce((sum, item) => sum + item.quantity, 0)
    const selectAll = this.data.cartItems.length > 0 && selectedItems.length === this.data.cartItems.length
    
    this.setData({ totalPrice, selectedCount, selectAll })
  },

  onCheckout: function() {
    if (this.data.selectedCount === 0) {
      wx.showToast({ title: '请选择商品', icon: 'none' })
      return
    }
    
    wx.navigateTo({
      url: '/pages/checkout/checkout'
    })
  }
})