<template>
  <div class="product-detail-page">
    <van-nav-bar
      title="商品详情"
      left-arrow
      fixed
      placeholder
      @click-left="handleBack"
    />

    <div v-if="product" class="product-detail">
      <!-- 商品图片 -->
      <van-swipe class="product-swipe" :autoplay="0" indicator-color="white">
        <van-swipe-item v-for="(image, index) in product.images" :key="index">
          <img :src="image" :alt="product.name" @click="previewImages(index)" />
        </van-swipe-item>
      </van-swipe>

      <!-- 商品信息 -->
      <div class="product-info">
        <h1 class="product-title">{{ product.name }}</h1>
        <div class="product-price">
          <span class="current-price">{{ formatPrice(product.price) }}</span>
          <span v-if="product.original_price > product.price" class="original-price">
            {{ formatPrice(product.original_price) }}
          </span>
        </div>
        <div class="product-stats">
          <span>销量 {{ product.sales }}</span>
          <span>库存 {{ product.stock }}</span>
        </div>
      </div>

      <!-- 商品详情 -->
      <div class="product-description">
        <h3>商品详情</h3>
        <div class="description-content" v-html="product.description"></div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div class="bottom-actions">
      <div class="action-left">
        <div class="action-item" @click="handleCart">
          <van-icon name="shopping-cart-o" />
          <span>购物车</span>
          <van-badge v-if="cartCount" :content="cartCount" />
        </div>
      </div>
      <div class="action-right">
        <van-button type="warning" size="large" @click="handleAddToCart">
          加入购物车
        </van-button>
        <van-button type="danger" size="large" @click="handleBuyNow">
          立即购买
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showImagePreview, showToast } from 'vant'
import { formatPrice } from '@/utils'
import { useCartStore } from '@/store/cart'
import { useUserStore } from '@/store/user'
import { getProduct, type Product } from '@/api/product'

const router = useRouter()
const route = useRoute()
const cartStore = useCartStore()
const userStore = useUserStore()

const product = ref<Product | null>(null)
const quantity = ref(1)

const productId = computed(() => Number(route.params.id))
const cartCount = computed(() => cartStore.totalCount)

const loadProduct = async () => {
  try {
    const { data } = await getProduct(productId.value)
    product.value = data
  } catch (error) {
    showToast('商品不存在')
    handleBack()
  }
}

const handleBack = () => router.back()

const previewImages = (index: number) => {
  if (product.value) {
    showImagePreview({
      images: product.value.images,
      startPosition: index
    })
  }
}

const handleCart = () => router.push('/cart')

const handleAddToCart = () => {
  if (!product.value) return
  
  const cartItem = {
    id: Date.now(),
    product_id: product.value.id,
    name: product.value.name,
    image: product.value.images[0],
    price: product.value.price,
    sku: product.value.sku,
    quantity: quantity.value,
    stock: product.value.stock
  }
  
  cartStore.addToCart(cartItem)
  showToast('已加入购物车')
}

const handleBuyNow = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  handleAddToCart()
  router.push('/checkout')
}

onMounted(() => {
  loadProduct()
})
</script>

<style scoped lang="scss">
.product-detail-page {
  background-color: #f8f8f8;
  min-height: 100vh;
  padding-bottom: 60px;

  .product-swipe {
    height: 375px;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .product-info {
    background: white;
    padding: 16px;
    margin-bottom: 10px;

    .product-title {
      font-size: 16px;
      font-weight: 500;
      line-height: 1.4;
      margin: 0 0 12px 0;
    }

    .product-price {
      margin-bottom: 12px;

      .current-price {
        font-size: 20px;
        font-weight: 600;
        color: #ee0a24;

        &::before {
          content: '¥';
          font-size: 14px;
        }
      }

      .original-price {
        font-size: 14px;
        color: #999;
        text-decoration: line-through;
        margin-left: 8px;

        &::before {
          content: '¥';
        }
      }
    }

    .product-stats {
      display: flex;
      gap: 16px;
      font-size: 13px;
      color: #666;
    }
  }

  .product-description {
    background: white;
    margin-bottom: 10px;
    padding: 16px;

    h3 {
      margin: 0 0 16px 0;
      font-size: 14px;
      font-weight: 500;
    }

    .description-content {
      line-height: 1.6;
      color: #666;
    }
  }

  .bottom-actions {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    display: flex;
    background: white;
    border-top: 1px solid #f0f0f0;
    padding: 8px;

    .action-left {
      display: flex;

      .action-item {
        position: relative;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 8px 16px;
        cursor: pointer;

        .van-icon {
          font-size: 20px;
          margin-bottom: 2px;
        }

        span {
          font-size: 10px;
          color: #666;
        }
      }
    }

    .action-right {
      flex: 1;
      display: flex;
      gap: 8px;
      padding-left: 8px;

      .van-button {
        flex: 1;
      }
    }
  }
}
</style>