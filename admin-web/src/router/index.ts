import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Layout from '@/components/Layout/index.vue'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      requireAuth: false
    }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '控制台',
          icon: 'Dashboard',
          requireAuth: true
        }
      }
    ]
  },
  {
    path: '/user',
    component: Layout,
    redirect: '/user/list',
    meta: {
      title: '用户管理',
      icon: 'User',
      requireAuth: true
    },
    children: [
      {
        path: 'list',
        name: 'UserList',
        component: () => import('@/views/user/list.vue'),
        meta: {
          title: '用户列表',
          requireAuth: true
        }
      }
    ]
  },
  {
    path: '/product',
    component: Layout,
    redirect: '/product/category',
    meta: {
      title: '商品管理',
      icon: 'Goods',
      requireAuth: true
    },
    children: [
      {
        path: 'category',
        name: 'CategoryList',
        component: () => import('@/views/product/category/index.vue'),
        meta: {
          title: '分类管理',
          requireAuth: true
        }
      },
      {
        path: 'list',
        name: 'ProductList',
        component: () => import('@/views/product/list/index.vue'),
        meta: {
          title: '商品列表',
          requireAuth: true
        }
      },
      {
        path: 'add',
        name: 'ProductAdd',
        component: () => import('@/views/product/form/index.vue'),
        meta: {
          title: '添加商品',
          requireAuth: true
        }
      },
      {
        path: 'edit/:id',
        name: 'ProductEdit',
        component: () => import('@/views/product/form/index.vue'),
        meta: {
          title: '编辑商品',
          requireAuth: true
        }
      }
    ]
  },
  {
    path: '/order',
    component: Layout,
    redirect: '/order/list',
    meta: {
      title: '订单管理',
      icon: 'ShoppingBag',
      requireAuth: true
    },
    children: [
      {
        path: 'list',
        name: 'OrderList',
        component: () => import('@/views/order/list/index.vue'),
        meta: {
          title: '订单列表',
          requireAuth: true
        }
      },
      {
        path: 'detail/:id',
        name: 'OrderDetail',
        component: () => import('@/views/order/detail/index.vue'),
        meta: {
          title: '订单详情',
          requireAuth: true
        }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    redirect: '/system/admin',
    meta: {
      title: '系统管理',
      icon: 'Setting',
      requireAuth: true
    },
    children: [
      {
        path: 'admin',
        name: 'AdminList',
        component: () => import('@/views/system/admin.vue'),
        meta: {
          title: '管理员管理',
          requireAuth: true
        }
      },
      {
        path: 'settings',
        name: 'SystemSettings',
        component: () => import('@/views/system/settings.vue'),
        meta: {
          title: '系统设置',
          requireAuth: true
        }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - Mall管理后台`
  }
  
  // 检查是否需要登录
  if (to.meta?.requireAuth && !userStore.token) {
    next('/login')
  } else if (to.path === '/login' && userStore.token) {
    next('/')
  } else {
    next()
  }
})

export default router