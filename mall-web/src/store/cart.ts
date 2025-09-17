import { defineStore } from 'pinia'
import { storage } from '@/utils'

export interface CartItem {
  id: number
  product_id: number
  name: string
  image: string
  price: number
  sku: string
  quantity: number
  selected: boolean
  stock: number
  specs?: { name: string; value: string }[]
}

interface CartState {
  items: CartItem[]
}

export const useCartStore = defineStore('cart', {
  state: (): CartState => ({
    items: storage.get('cartItems') || []
  }),
  
  getters: {
    // 购物车商品总数
    totalCount: (state) => state.items.reduce((sum, item) => sum + item.quantity, 0),
    
    // 选中商品总数
    selectedCount: (state) => state.items.filter(item => item.selected).reduce((sum, item) => sum + item.quantity, 0),
    
    // 选中商品总价
    selectedTotal: (state) => state.items.filter(item => item.selected).reduce((sum, item) => sum + item.price * item.quantity, 0),
    
    // 是否全选
    isAllSelected: (state) => state.items.length > 0 && state.items.every(item => item.selected),
    
    // 选中的商品
    selectedItems: (state) => state.items.filter(item => item.selected)
  },
  
  actions: {
    // 保存到本地存储
    saveToStorage() {
      storage.set('cartItems', this.items)
    },
    
    // 添加商品到购物车
    addToCart(product: Omit<CartItem, 'selected'>) {
      const existingItem = this.items.find(item => 
        item.product_id === product.product_id && 
        JSON.stringify(item.specs || []) === JSON.stringify(product.specs || [])
      )
      
      if (existingItem) {
        existingItem.quantity += product.quantity
      } else {
        this.items.push({ ...product, selected: true })
      }
      
      this.saveToStorage()
    },
    
    // 更新商品数量
    updateQuantity(id: number, quantity: number) {
      const item = this.items.find(item => item.id === id)
      if (item) {
        item.quantity = Math.max(1, Math.min(quantity, item.stock))
        this.saveToStorage()
      }
    },
    
    // 删除商品
    removeItem(id: number) {
      const index = this.items.findIndex(item => item.id === id)
      if (index > -1) {
        this.items.splice(index, 1)
        this.saveToStorage()
      }
    },
    
    // 切换商品选中状态
    toggleSelect(id: number) {
      const item = this.items.find(item => item.id === id)
      if (item) {
        item.selected = !item.selected
        this.saveToStorage()
      }
    },
    
    // 全选/取消全选
    toggleSelectAll() {
      const allSelected = this.isAllSelected
      this.items.forEach(item => {
        item.selected = !allSelected
      })
      this.saveToStorage()
    },
    
    // 清空购物车
    clear() {
      this.items = []
      this.saveToStorage()
    },
    
    // 删除选中商品
    removeSelected() {
      this.items = this.items.filter(item => !item.selected)
      this.saveToStorage()
    }
  }
})