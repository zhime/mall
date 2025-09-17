<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1>控制台</h1>
      <p>欢迎使用Mall管理后台</p>
    </div>
    
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon user">
              <el-icon><User /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">用户总数</div>
              <div class="stats-value">{{ stats.userCount }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon product">
              <el-icon><Goods /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">商品总数</div>
              <div class="stats-value">{{ stats.productCount }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon order">
              <el-icon><ShoppingBag /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">订单总数</div>
              <div class="stats-value">{{ stats.orderCount }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon revenue">
              <el-icon><Money /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">总收入</div>
              <div class="stats-value">¥{{ stats.revenue }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>销售趋势</span>
            </div>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder">
              销售趋势图表区域
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>热门商品</span>
            </div>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder">
              热门商品图表区域
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 最近订单 -->
    <el-card class="recent-orders">
      <template #header>
        <div class="card-header">
          <span>最近订单</span>
          <el-button type="primary" size="small" @click="$router.push('/order/list')">
            查看全部
          </el-button>
        </div>
      </template>
      
      <el-table :data="recentOrders" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="total_amount" label="金额" width="100">
          <template #default="{ row }">
            ¥{{ row.total_amount }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

// 统计数据
const stats = ref({
  userCount: 1234,
  productCount: 567,
  orderCount: 890,
  revenue: '123,456.78'
})

// 最近订单
const recentOrders = ref([
  {
    order_no: 'ORD20240101001',
    username: '张三',
    total_amount: 299.00,
    status: 1,
    created_at: '2024-01-01 10:30:00'
  },
  {
    order_no: 'ORD20240101002',
    username: '李四',
    total_amount: 199.00,
    status: 2,
    created_at: '2024-01-01 11:15:00'
  },
  {
    order_no: 'ORD20240101003',
    username: '王五',
    total_amount: 399.00,
    status: 3,
    created_at: '2024-01-01 12:00:00'
  }
])

// 获取订单状态类型
const getOrderStatusType = (status: number) => {
  const typeMap: Record<number, string> = {
    1: 'warning',
    2: 'info',
    3: 'primary',
    4: 'success',
    5: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取订单状态文本
const getOrderStatusText = (status: number) => {
  const textMap: Record<number, string> = {
    1: '待付款',
    2: '待发货',
    3: '已发货',
    4: '已完成',
    5: '已取消'
  }
  return textMap[status] || '未知'
}

onMounted(() => {
  // 这里可以调用API获取实际数据
})
</script>

<style scoped lang="scss">
.dashboard {
  .dashboard-header {
    margin-bottom: 20px;
    
    h1 {
      margin: 0 0 8px 0;
      color: #303133;
      font-size: 24px;
      font-weight: 500;
    }
    
    p {
      margin: 0;
      color: #909399;
      font-size: 14px;
    }
  }
  
  .stats-row {
    margin-bottom: 20px;
    
    .stats-card {
      .stats-content {
        display: flex;
        align-items: center;
        
        .stats-icon {
          width: 60px;
          height: 60px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 16px;
          
          .el-icon {
            font-size: 24px;
            color: white;
          }
          
          &.user {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          }
          
          &.product {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
          }
          
          &.order {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
          }
          
          &.revenue {
            background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
          }
        }
        
        .stats-info {
          .stats-title {
            font-size: 14px;
            color: #909399;
            margin-bottom: 8px;
          }
          
          .stats-value {
            font-size: 24px;
            font-weight: 600;
            color: #303133;
          }
        }
      }
    }
  }
  
  .chart-row {
    margin-bottom: 20px;
    
    .chart-container {
      height: 300px;
      
      .chart-placeholder {
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f5f7fa;
        color: #909399;
        font-size: 16px;
      }
    }
  }
  
  .recent-orders {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style>