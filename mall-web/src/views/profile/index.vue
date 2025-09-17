<template>
  <div class="profile-page">
    <!-- 用户信息头部 -->
    <div class="profile-header">
      <div class="user-info" @click="handleEditProfile">
        <van-image
          class="avatar"
          :src="userStore.avatar || defaultAvatar"
          :alt="userStore.nickname"
          round
        />
        <div class="user-details">
          <div class="nickname">{{ userStore.nickname }}</div>
          <div class="phone">{{ userStore.userInfo?.phone }}</div>
        </div>
        <van-icon name="arrow" />
      </div>
    </div>

    <!-- 功能菜单 -->
    <div class="menu-section">
      <van-cell-group inset>
        <van-cell
          title="我的订单"
          icon="orders-o"
          is-link
          @click="handleOrderList"
        />
        <van-cell
          title="收货地址"
          icon="location-o"
          is-link
          @click="handleAddressList"
        />
        <van-cell
          title="我的收藏"
          icon="star-o"
          is-link
          @click="handleFavorites"
        />
        <van-cell
          title="优惠券"
          icon="coupon-o"
          is-link
          @click="handleCoupons"
        />
      </van-cell-group>
    </div>

    <div class="menu-section">
      <van-cell-group inset>
        <van-cell
          title="客服咨询"
          icon="service-o"
          is-link
          @click="handleCustomerService"
        />
        <van-cell
          title="关于我们"
          icon="info-o"
          is-link
          @click="handleAbout"
        />
        <van-cell
          title="设置"
          icon="setting-o"
          is-link
          @click="handleSettings"
        />
      </van-cell-group>
    </div>

    <!-- 退出登录 -->
    <div class="logout-section">
      <van-button
        type="danger"
        size="large"
        block
        @click="handleLogout"
      >
        退出登录
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
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import { getImageUrl } from '@/utils'

const router = useRouter()
const userStore = useUserStore()
const cartStore = useCartStore()

const activeTab = ref(3)

const cartCount = computed(() => cartStore.totalCount)
const defaultAvatar = getImageUrl('default-avatar.png')

// 编辑个人资料
const handleEditProfile = () => {
  router.push('/profile/edit')
}

// 我的订单
const handleOrderList = () => {
  router.push('/order')
}

// 收货地址
const handleAddressList = () => {
  router.push('/address')
}

// 我的收藏
const handleFavorites = () => {
  showToast('功能开发中...')
}

// 优惠券
const handleCoupons = () => {
  showToast('功能开发中...')
}

// 客服咨询
const handleCustomerService = () => {
  showToast('功能开发中...')
}

// 关于我们
const handleAbout = () => {
  showToast('功能开发中...')
}

// 设置
const handleSettings = () => {
  showToast('功能开发中...')
}

// 退出登录
const handleLogout = async () => {
  try {
    await showConfirmDialog({
      title: '退出登录',
      message: '确定要退出登录吗？'
    })
    
    userStore.logout()
    cartStore.clear()
    showToast('已退出登录')
    router.replace('/login')
  } catch (error) {
    // 用户取消
  }
}
</script>

<style scoped lang="scss">
.profile-page {
  background-color: #f8f8f8;
  min-height: 100vh;
  padding-bottom: 50px;

  .profile-header {
    background: linear-gradient(135deg, #ee0a24 0%, #ff6034 100%);
    padding: 20px 16px 30px;

    .user-info {
      display: flex;
      align-items: center;
      color: white;
      cursor: pointer;

      .avatar {
        width: 60px;
        height: 60px;
        margin-right: 16px;
      }

      .user-details {
        flex: 1;

        .nickname {
          font-size: 18px;
          font-weight: 500;
          margin-bottom: 4px;
        }

        .phone {
          font-size: 14px;
          opacity: 0.8;
        }
      }

      .van-icon {
        font-size: 16px;
        opacity: 0.8;
      }
    }
  }

  .menu-section {
    margin: 16px 0;
  }

  .logout-section {
    margin: 40px 16px 16px;
  }
}
</style>