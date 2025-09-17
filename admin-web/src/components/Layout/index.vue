<template>
  <div class="layout">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside :width="sidebarWidth">
        <div class="sidebar">
          <div class="logo">
            <h2>Mall管理后台</h2>
          </div>
          <el-menu
            :default-active="activeMenu"
            :collapse="isCollapse"
            :unique-opened="false"
            class="sidebar-menu"
            router
          >
            <template v-for="route in menuRoutes" :key="route.path">
              <el-sub-menu
                v-if="route.children && route.children.length > 1"
                :index="route.path"
              >
                <template #title>
                  <el-icon v-if="route.meta?.icon">
                    <component :is="route.meta.icon" />
                  </el-icon>
                  <span>{{ route.meta?.title }}</span>
                </template>
                <el-menu-item
                  v-for="child in route.children"
                  :key="child.path"
                  :index="child.path === 'dashboard' ? '/' + child.path : route.path + '/' + child.path"
                >
                  {{ child.meta?.title }}
                </el-menu-item>
              </el-sub-menu>
              <el-menu-item
                v-else
                :index="route.children?.[0]?.path === 'dashboard' ? '/dashboard' : route.path"
              >
                <el-icon v-if="route.meta?.icon">
                  <component :is="route.meta.icon" />
                </el-icon>
                <template #title>{{ route.meta?.title || route.children?.[0]?.meta?.title }}</template>
              </el-menu-item>
            </template>
          </el-menu>
        </div>
      </el-aside>

      <!-- 主内容区 -->
      <el-container>
        <!-- 顶部导航 -->
        <el-header class="header">
          <div class="header-left">
            <el-button
              type="text"
              @click="toggleSidebar"
            >
              <el-icon><Fold v-if="!isCollapse" /><Expand v-else /></el-icon>
            </el-button>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item
                v-for="item in breadcrumbs"
                :key="item.path"
                :to="{ path: item.path }"
              >
                {{ item.title }}
              </el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-avatar :size="32" :src="userInfo?.avatar" />
                <span class="username">{{ userInfo?.nickname || userInfo?.username }}</span>
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                  <el-dropdown-item command="password">修改密码</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <!-- 主内容 -->
        <el-main class="main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessageBox, ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const userInfo = computed(() => userStore.userInfo)

// 侧边栏宽度
const sidebarWidth = computed(() => isCollapse.value ? '64px' : '200px')

// 当前激活的菜单
const activeMenu = computed(() => route.path)

// 菜单路由
const menuRoutes = computed(() => {
  return router.options.routes.filter(route => 
    route.path !== '/login' && route.meta?.requireAuth
  )
})

// 面包屑
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  return matched.map(item => ({
    title: item.meta?.title,
    path: item.path
  }))
})

// 切换侧边栏
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

// 下拉菜单命令处理
const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      // TODO: 跳转到个人资料页面
      break
    case 'password':
      // TODO: 打开修改密码对话框
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        await userStore.logout()
        ElMessage.success('退出成功')
        router.push('/login')
      } catch (error) {
        // 用户取消
      }
      break
  }
}

onMounted(async () => {
  if (userStore.token && !userStore.userInfo) {
    try {
      await userStore.getUserInfo()
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }
})
</script>

<style scoped lang="scss">
.layout {
  height: 100vh;
  
  .el-container {
    height: 100%;
  }
  
  .sidebar {
    height: 100%;
    background: #304156;
    
    .logo {
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #2b2f3a;
      color: white;
      
      h2 {
        font-size: 18px;
        margin: 0;
      }
    }
    
    .sidebar-menu {
      border: none;
      height: calc(100% - 60px);
      width: 100% !important;
      
      :deep(.el-menu-item),
      :deep(.el-sub-menu__title) {
        color: #bfcbd9;
        
        &:hover {
          background-color: #263445 !important;
          color: #409eff;
        }
      }
      
      :deep(.el-menu-item.is-active) {
        background-color: #409eff !important;
        color: white;
      }
    }
  }
  
  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    background: white;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    
    .header-left {
      display: flex;
      align-items: center;
      gap: 16px;
    }
    
    .header-right {
      .user-info {
        display: flex;
        align-items: center;
        gap: 8px;
        cursor: pointer;
        padding: 8px;
        border-radius: 4px;
        
        &:hover {
          background-color: #f5f5f5;
        }
        
        .username {
          font-size: 14px;
        }
      }
    }
  }
  
  .main {
    background: #f0f2f5;
    padding: 20px;
    overflow-y: auto;
  }
}
</style>