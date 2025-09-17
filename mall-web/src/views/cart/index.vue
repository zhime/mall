<template>
  <div class="cart-page">
    <van-nav-bar title="购物车" fixed placeholder />

    <div v-if="cartStore.items.length > 0" class="cart-content">
      <!-- 购物车商品列表 -->
      <div class="cart-list">
        <van-checkbox-group v-model="checkedItems">
          <van-swipe-cell
            v-for="item in cartStore.items"
            :key="item.id"
            :right-width="65"
          >
            <div class="cart-item">
              <van-checkbox :name="item.id" />
              <div class="item-image">
                <img :src="item.image" :alt="item.name" />
              </div>
              <div class="item-info">
                <h4 class="item-name">{{ item.name }}</h4>
                <div v-if="item.specs?.length" class="item-specs">
                  {{ item.specs.map(spec => `${spec.name}: ${spec.value}`).join(', ') }}
                </div>
                <div class="item-bottom">
                  <div class="item-price">¥{{ formatPrice(item.price) }}</div>
                  <van-stepper
                    v-model="item.quantity"
                    :min="1"
                    :max="item.stock"
                    @change="handleQuantityChange(item.id, $event)"
                  />
                </div>
              </div>
            </div>
            <template #right>
              <van-button
                square
                type="danger"
                text="删除"
                class="delete-button"
                @click="handleDeleteItem(item.id)"
              />
            </template>
          </van-swipe-cell>
        </van-checkbox-group>
      </div>

      <!-- 底部操作栏 -->
      <div class="cart-footer">
        <div class="footer-left">
          <van-checkbox
            v-model="selectAll"
            @change="handleSelectAll"
          >
            全选
          </van-checkbox>
          <span class="total-info">
            共{{ cartStore.selectedCount }}件，合计：
            <span class="total-price">¥{{ formatPrice(cartStore.selectedTotal) }}</span>
          </span>
        </div>
        <van-button
          type="danger"
          size="large"
          :disabled="cartStore.selectedCount === 0"
          @click="handleCheckout"
        >
          结算({{ cartStore.selectedCount }})
        </van-button>
      </div>
    </div>

    <!-- 空购物车 -->
    <div v-else class="empty-cart">
      <div class="empty-icon">🛒</div>
      <div class="empty-text">购物车还是空的</div>
      <van-button type="primary" @click="handleGoShopping">
        去逛逛
      </van-button>
    </div>

    <!-- 底部导航 -->
    <van-tabbar v-model="activeTab" fixed placeholder>
      <van-tabbar-item icon="home-o" to="/">首页</van-tabbar-item>
      <van-tabbar-item icon="apps-o" to="/category">分类</van-tabbar-item>
      <van-tabbar-item icon="shopping-cart-o" to="/cart" :badge="cartCount">购物车</van-tabbar-item>
      <van-tabbar-item icon="user-o" to="/profile">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { formatPrice } from '@/utils'
import { useCartStore } from '@/store/cart'
import { useUserStore } from '@/store/user'

const router = useRouter()
const cartStore = useCartStore()
const userStore = useUserStore()

const activeTab = ref(2)
const checkedItems = ref<number[]>([])

const cartCount = computed(() => cartStore.totalCount)

// 全选状态
const selectAll = computed({
  get: () => cartStore.isAllSelected,
  set: (value: boolean) => {
    if (value) {
      checkedItems.value = cartStore.items.map(item => item.id)
    } else {
      checkedItems.value = []
    }
  }
})

// 更新商品数量
const handleQuantityChange = (id: number, quantity: number) => {
  cartStore.updateQuantity(id, quantity)
}

// 删除商品
const handleDeleteItem = async (id: number) => {
  try {
    await showConfirmDialog({
      title: '确认删除',
      message: '确定要删除这件商品吗？'
    })
    
    cartStore.removeItem(id)
    showToast('删除成功')
  } catch (error) {
    // 用户取消
  }
}

// 全选/取消全选
const handleSelectAll = () => {
  cartStore.toggleSelectAll()
}

// 结算
const handleCheckout = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  if (cartStore.selectedCount === 0) {
    showToast('请选择商品')
    return
  }
  
  router.push('/checkout')
}

// 去逛逛
const handleGoShopping = () => {
  router.push('/')
}
</script>

<style scoped lang="scss">
.cart-page {
  background-color: #f8f8f8;
  min-height: 100vh;
  padding-bottom: 50px;

  .cart-content {
    padding-bottom: 60px;
  }

  .cart-list {
    .cart-item {
      display: flex;
      align-items: center;
      padding: 16px;
      background: white;
      border-bottom: 1px solid #f0f0f0;

      .van-checkbox {
        margin-right: 12px;
      }

      .item-image {
        width: 80px;
        height: 80px;
        margin-right: 12px;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          border-radius: 8px;
        }
      }

      .item-info {
        flex: 1;

        .item-name {
          font-size: 14px;
          font-weight: normal;
          line-height: 1.4;
          margin: 0 0 8px 0;
          overflow: hidden;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }

        .item-specs {
          font-size: 12px;
          color: #666;
          margin-bottom: 8px;
        }

        .item-bottom {
          display: flex;
          justify-content: space-between;
          align-items: center;

          .item-price {
            font-size: 16px;
            font-weight: 600;
            color: #ee0a24;
          }
        }
      }
    }

    .delete-button {
      height: 100%;
    }
  }

  .cart-footer {
    position: fixed;
    bottom: 50px;
    left: 0;
    right: 0;
    display: flex;
    align-items: center;
    padding: 12px 16px;
    background: white;
    border-top: 1px solid #f0f0f0;

    .footer-left {
      flex: 1;
      display: flex;
      align-items: center;

      .van-checkbox {
        margin-right: 12px;
      }

      .total-info {
        font-size: 12px;
        color: #666;

        .total-price {
          font-size: 16px;
          font-weight: 600;
          color: #ee0a24;
        }
      }
    }
  }

  .empty-cart {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px 20px;

    .empty-icon {
      font-size: 80px;
      margin-bottom: 16px;
      opacity: 0.3;
    }

    .empty-text {
      font-size: 16px;
      color: #666;
      margin-bottom: 24px;
    }
  }
}
</style>