<template>
  <div class="home-page">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <van-search
        v-model="searchKeyword"
        placeholder="搜索商品"
        background="#ee0a24"
        @click="handleSearch"
        @search="handleSearch"
      />
    </div>

    <!-- 轮播图 -->
    <van-swipe class="banner-swipe" :autoplay="3000" indicator-color="white">
      <van-swipe-item v-for="banner in banners" :key="banner.id" @click="handleBannerClick(banner)">
        <img :src="banner.image" :alt="banner.title" />
      </van-swipe-item>
    </van-swipe>

    <!-- 分类导航 -->
    <div class="category-nav">
      <div
        v-for="category in categories.slice(0, 8)"
        :key="category.id"
        class="category-item"
        @click="handleCategoryClick(category)"
      >
        <img :src="category.icon" :alt="category.name" />
        <span>{{ category.name }}</span>
      </div>
    </div>

    <!-- 营销活动 -->
    <div class="marketing-section">
      <van-notice-bar
        color="#ee0a24"
        background="#fff2f1"
        left-icon="volume-o"
        text="限时秒杀活动正在进行中，快来抢购吧！"
      />
    </div>

    <!-- 商品推荐 -->
    <div class="recommend-section">
      <div class="section-header">
        <h3>为你推荐</h3>
      </div>
      
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-list
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="onLoad"
        >
          <div class="product-grid">
            <div
              v-for="product in productList"
              :key="product.id"
              class="product-item"
              @click="handleProductClick(product)"
            >
              <div class="product-image">
                <img v-lazy="product.images[0]" :alt="product.name" />
                <div v-if="product.is_new" class="product-tag tag-new">新品</div>
                <div v-if="product.is_featured" class="product-tag tag-hot">热销</div>
              </div>
              <div class="product-info">
                <h4 class="product-name">{{ product.name }}</h4>
                <div class="product-price">
                  <span class="current-price">{{ formatPrice(product.price) }}</span>
                  <span v-if="product.original_price > product.price" class="original-price">
                    {{ formatPrice(product.original_price) }}
                  </span>
                </div>
                <div class="product-sales">已售{{ product.sales }}件</div>
              </div>
            </div>
          </div>
        </van-list>
      </van-pull-refresh>
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
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { formatPrice } from '@/utils'
import { useCartStore } from '@/store/cart'
import {
  getBanners,
  getCategoryList,
  getRecommendProducts,
  type Product,
  type Category
} from '@/api/product'

interface Banner {
  id: number
  title: string
  image: string
  link: string
  sort: number
}

const router = useRouter()
const cartStore = useCartStore()

const activeTab = ref(0)
const searchKeyword = ref('')
const refreshing = ref(false)
const loading = ref(false)
const finished = ref(false)

const banners = ref<Banner[]>([])
const categories = ref<Category[]>([])
const productList = ref<Product[]>([])
const page = ref(1)

const cartCount = computed(() => cartStore.totalCount)

// 搜索
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push(`/search?keyword=${encodeURIComponent(searchKeyword.value.trim())}`)
  }
}

// 轮播图点击
const handleBannerClick = (banner: Banner) => {
  if (banner.link) {
    if (banner.link.startsWith('http')) {
      window.open(banner.link)
    } else {
      router.push(banner.link)
    }
  }
}

// 分类点击
const handleCategoryClick = (category: Category) => {
  router.push(`/category?id=${category.id}`)
}

// 商品点击
const handleProductClick = (product: Product) => {
  router.push(`/product/${product.id}`)
}

// 加载轮播图
const loadBanners = async () => {
  try {
    const { data } = await getBanners()
    banners.value = data
  } catch (error) {
    console.error('加载轮播图失败:', error)
  }
}

// 加载分类
const loadCategories = async () => {
  try {
    const { data } = await getCategoryList()
    categories.value = data.filter((cat: Category) => cat.level === 1)
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 加载推荐商品
const loadRecommendProducts = async (isRefresh = false) => {
  try {
    const { data } = await getRecommendProducts()
    if (isRefresh) {
      productList.value = data.list
      page.value = 2
    } else {
      productList.value.push(...data.list)
      page.value++
    }
    
    if (data.list.length < 20) {
      finished.value = true
    }
  } catch (error) {
    console.error('加载推荐商品失败:', error)
    showToast('加载失败，请重试')
  }
}

// 下拉刷新
const onRefresh = async () => {
  try {
    refreshing.value = true
    finished.value = false
    await loadRecommendProducts(true)
    showToast('刷新成功')
  } finally {
    refreshing.value = false
  }
}

// 上拉加载
const onLoad = async () => {
  if (finished.value) {
    loading.value = false
    return
  }
  
  try {
    loading.value = true
    await loadRecommendProducts()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadBanners()
  loadCategories()
  loadRecommendProducts(true)
})
</script>

<style scoped lang="scss">
.home-page {
  background-color: #f8f8f8;
  min-height: 100vh;
  padding-bottom: 50px;

  .search-bar {
    position: sticky;
    top: 0;
    z-index: 10;
    background: #ee0a24;
    padding: 8px 16px;
  }

  .banner-swipe {
    height: 200px;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .category-nav {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    padding: 20px 16px;
    background: white;
    margin-bottom: 10px;

    .category-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;

      img {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        margin-bottom: 8px;
      }

      span {
        font-size: 12px;
        color: #333;
      }
    }
  }

  .marketing-section {
    margin-bottom: 10px;
  }

  .recommend-section {
    background: white;

    .section-header {
      padding: 16px;
      border-bottom: 1px solid #f0f0f0;

      h3 {
        margin: 0;
        font-size: 16px;
        font-weight: 600;
        color: #333;
      }
    }

    .product-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1px;
      background: #f0f0f0;

      .product-item {
        background: white;
        padding: 12px;

        .product-image {
          position: relative;
          width: 100%;
          height: 140px;
          margin-bottom: 8px;

          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            border-radius: 4px;
          }

          .product-tag {
            position: absolute;
            top: 6px;
            left: 6px;
            padding: 2px 6px;
            border-radius: 2px;
            font-size: 10px;
            color: white;

            &.tag-new {
              background: #ff976a;
            }

            &.tag-hot {
              background: #ee0a24;
            }
          }
        }

        .product-info {
          .product-name {
            font-size: 14px;
            font-weight: normal;
            line-height: 1.4;
            margin: 0 0 8px 0;
            height: 40px;
            overflow: hidden;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
          }

          .product-price {
            margin-bottom: 4px;

            .current-price {
              font-size: 16px;
              font-weight: 600;
              color: #ee0a24;

              &::before {
                content: '¥';
                font-size: 12px;
              }
            }

            .original-price {
              font-size: 12px;
              color: #999;
              text-decoration: line-through;
              margin-left: 4px;

              &::before {
                content: '¥';
              }
            }
          }

          .product-sales {
            font-size: 12px;
            color: #999;
          }
        }
      }
    }
  }
}
</style>