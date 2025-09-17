<template>
  <div class="order-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>订单管理</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出订单
          </el-button>
        </div>
      </template>

      <!-- 统计卡片 -->
      <div class="statistics-row">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-number">{{ statistics.total_orders }}</div>
                <div class="stat-label">总订单数</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-number pending">{{ statistics.pending_orders }}</div>
                <div class="stat-label">待处理订单</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-number success">¥{{ statistics.total_amount }}</div>
                <div class="stat-label">总交易额</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-number info">{{ statistics.today_orders }}</div>
                <div class="stat-label">今日订单</div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <div class="search-bar">
        <el-form :model="searchForm" inline>
          <el-form-item label="订单号">
            <el-input
              v-model="searchForm.order_no"
              placeholder="请输入订单号"
              style="width: 200px"
              clearable
            />
          </el-form-item>
          <el-form-item label="用户手机号">
            <el-input
              v-model="searchForm.user_phone"
              placeholder="请输入手机号"
              style="width: 150px"
              clearable
            />
          </el-form-item>
          <el-form-item label="订单状态">
            <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 120px">
              <el-option label="待付款" :value="0" />
              <el-option label="待发货" :value="1" />
              <el-option label="已发货" :value="2" />
              <el-option label="已完成" :value="3" />
              <el-option label="已取消" :value="4" />
            </el-select>
          </el-form-item>
          <el-form-item label="支付状态">
            <el-select v-model="searchForm.payment_status" placeholder="支付状态" clearable style="width: 120px">
              <el-option label="未支付" :value="0" />
              <el-option label="已支付" :value="1" />
              <el-option label="已退款" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="创建时间">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 240px"
              @change="handleDateChange"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="toolbar">
        <el-button
          type="success"
          :disabled="!selectedOrders.length"
          @click="handleBatchStatus(1)"
        >
          批量发货
        </el-button>
        <el-button
          type="warning"
          :disabled="!selectedOrders.length"
          @click="handleBatchStatus(3)"
        >
          批量完成
        </el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="orderList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="order_no" label="订单号" width="200" />
        <el-table-column label="用户信息" width="150">
          <template #default="{ row }">
            <div>
              <div>{{ row.user_name }}</div>
              <div class="user-phone">{{ row.user_phone }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="商品信息" min-width="200">
          <template #default="{ row }">
            <div class="order-items">
              <div
                v-for="item in row.items.slice(0, 2)"
                :key="item.id"
                class="order-item"
              >
                <el-image
                  :src="item.product_image"
                  style="width: 40px; height: 40px; margin-right: 8px"
                  fit="cover"
                />
                <div class="item-info">
                  <div class="item-name">{{ item.product_name }}</div>
                  <div class="item-detail">¥{{ item.price }} × {{ item.quantity }}</div>
                </div>
              </div>
              <div v-if="row.items.length > 2" class="more-items">
                +{{ row.items.length - 2 }}个商品
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="订单金额" width="120">
          <template #default="{ row }">
            <div>
              <div class="final-amount">¥{{ row.final_amount }}</div>
              <div v-if="row.discount_amount > 0" class="discount">
                优惠: ¥{{ row.discount_amount }}
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="订单状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="支付状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getPaymentStatusType(row.payment_status)">
              {{ getPaymentStatusText(row.payment_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button
              v-if="row.status === 1"
              type="success"
              link
              @click="handleShip(row)"
            >
              发货
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="danger"
              link
              @click="handleCancel(row)"
            >
              取消
            </el-button>
            <el-dropdown v-if="row.status >= 2" @command="handleCommand($event, row)">
              <el-button type="primary" link>
                更多<el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="confirm">确认收货</el-dropdown-item>
                  <el-dropdown-item command="refund">申请退款</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 发货对话框 -->
    <el-dialog v-model="shipDialogVisible" title="订单发货" width="500px">
      <el-form ref="shipFormRef" :model="shipForm" :rules="shipFormRules" label-width="100px">
        <el-form-item label="快递公司" prop="shipping_company">
          <el-select v-model="shipForm.shipping_company" placeholder="请选择快递公司" style="width: 100%">
            <el-option label="顺丰速运" value="SF" />
            <el-option label="圆通速递" value="YTO" />
            <el-option label="中通快递" value="ZTO" />
            <el-option label="韵达速递" value="YD" />
            <el-option label="申通快递" value="STO" />
            <el-option label="京东快递" value="JD" />
            <el-option label="邮政EMS" value="EMS" />
          </el-select>
        </el-form-item>
        <el-form-item label="快递单号" prop="tracking_no">
          <el-input v-model="shipForm.tracking_no" placeholder="请输入快递单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="shipping" @click="confirmShip">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus'
import { Download, ArrowDown } from '@element-plus/icons-vue'
import {
  getOrderList,
  updateOrderStatus,
  batchUpdateOrderStatus,
  cancelOrder,
  shipOrder,
  confirmOrder,
  refundOrder,
  getOrderStatistics,
  exportOrders,
  type Order,
  type OrderQuery,
  type OrderStatistics
} from '@/api/order'

const router = useRouter()
const loading = ref(false)
const shipping = ref(false)
const shipDialogVisible = ref(false)
const shipFormRef = ref<FormInstance>()

const searchForm = reactive({
  order_no: '',
  user_phone: '',
  status: undefined as number | undefined,
  payment_status: undefined as number | undefined,
  start_date: '',
  end_date: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const shipForm = reactive({
  shipping_company: '',
  tracking_no: ''
})

const shipFormRules = {
  shipping_company: [
    { required: true, message: '请选择快递公司', trigger: 'change' }
  ],
  tracking_no: [
    { required: true, message: '请输入快递单号', trigger: 'blur' }
  ]
}

const orderList = ref<Order[]>([])
const selectedOrders = ref<Order[]>([])
const statistics = ref<OrderStatistics>({
  total_orders: 0,
  pending_orders: 0,
  processing_orders: 0,
  completed_orders: 0,
  total_amount: 0,
  today_orders: 0,
  today_amount: 0
})
const dateRange = ref<[Date, Date] | null>(null)
const currentOrder = ref<Order | null>(null)

const loadOrders = async () => {
  try {
    loading.value = true
    const params: OrderQuery = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    }
    
    const { data } = await getOrderList(params)
    orderList.value = data.list
    pagination.total = data.total
  } catch (error) {
    ElMessage.error('加载订单列表失败')
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  try {
    const { data } = await getOrderStatistics()
    statistics.value = data
  } catch (error) {
    ElMessage.error('加载统计数据失败')
  }
}

const getStatusType = (status: number) => {
  const types = ['warning', 'info', 'primary', 'success', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status: number) => {
  const texts = ['待付款', '待发货', '已发货', '已完成', '已取消']
  return texts[status] || '未知'
}

const getPaymentStatusType = (status: number) => {
  const types = ['warning', 'success', 'danger']
  return types[status] || 'info'
}

const getPaymentStatusText = (status: number) => {
  const texts = ['未支付', '已支付', '已退款']
  return texts[status] || '未知'
}

const handleDateChange = (dates: [Date, Date] | null) => {
  if (dates) {
    searchForm.start_date = dates[0].toISOString().split('T')[0]
    searchForm.end_date = dates[1].toISOString().split('T')[0]
  } else {
    searchForm.start_date = ''
    searchForm.end_date = ''
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadOrders()
}

const handleReset = () => {
  Object.assign(searchForm, {
    order_no: '',
    user_phone: '',
    status: undefined,
    payment_status: undefined,
    start_date: '',
    end_date: ''
  })
  dateRange.value = null
  pagination.page = 1
  loadOrders()
}

const handleView = (row: Order) => {
  router.push(`/order/detail/${row.id}`)
}

const handleShip = (row: Order) => {
  currentOrder.value = row
  shipForm.shipping_company = ''
  shipForm.tracking_no = ''
  shipDialogVisible.value = true
}

const confirmShip = async () => {
  if (!shipFormRef.value || !currentOrder.value) return
  
  try {
    await shipFormRef.value.validate()
    shipping.value = true
    
    await shipOrder(currentOrder.value.id, shipForm)
    ElMessage.success('发货成功')
    shipDialogVisible.value = false
    loadOrders()
  } catch (error) {
    ElMessage.error('发货失败')
  } finally {
    shipping.value = false
  }
}

const handleCancel = async (row: Order) => {
  try {
    const { value: reason } = await ElMessageBox.prompt(
      '请输入取消原因',
      '取消订单',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        inputPlaceholder: '请输入取消原因'
      }
    )
    
    await cancelOrder(row.id, reason)
    ElMessage.success('订单取消成功')
    loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleCommand = async (command: string, row: Order) => {
  switch (command) {
    case 'confirm':
      try {
        await ElMessageBox.confirm('确认该订单已收货？', '确认收货', { type: 'warning' })
        await confirmOrder(row.id)
        ElMessage.success('确认收货成功')
        loadOrders()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('操作失败')
        }
      }
      break
    case 'refund':
      try {
        const { value } = await ElMessageBox.prompt(
          '请输入退款金额和原因',
          '申请退款',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消'
          }
        )
        
        await refundOrder(row.id, { amount: row.final_amount, reason: value })
        ElMessage.success('退款申请成功')
        loadOrders()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('操作失败')
        }
      }
      break
  }
}

const handleSelectionChange = (selection: Order[]) => {
  selectedOrders.value = selection
}

const handleBatchStatus = async (status: number) => {
  try {
    const statusText = getStatusText(status)
    await ElMessageBox.confirm(
      `确认将选中的 ${selectedOrders.value.length} 个订单状态更新为"${statusText}"吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    
    const ids = selectedOrders.value.map(order => order.id)
    await batchUpdateOrderStatus(ids, status)
    ElMessage.success('批量操作成功')
    loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量操作失败')
    }
  }
}

const handleExport = async () => {
  try {
    await exportOrders(searchForm)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  loadOrders()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadOrders()
}

onMounted(() => {
  loadStatistics()
  loadOrders()
})
</script>

<style scoped lang="scss">
.order-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .statistics-row {
    margin-bottom: 20px;

    .stat-card {
      .stat-content {
        text-align: center;
        
        .stat-number {
          font-size: 24px;
          font-weight: bold;
          margin-bottom: 8px;
          
          &.pending {
            color: #E6A23C;
          }
          
          &.success {
            color: #67C23A;
          }
          
          &.info {
            color: #409EFF;
          }
        }
        
        .stat-label {
          color: #909399;
          font-size: 14px;
        }
      }
    }
  }

  .search-bar {
    margin-bottom: 20px;
  }

  .toolbar {
    margin-bottom: 20px;
  }

  .user-phone {
    font-size: 12px;
    color: #666;
  }

  .order-items {
    .order-item {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      
      .item-info {
        .item-name {
          font-size: 12px;
          margin-bottom: 2px;
        }
        
        .item-detail {
          font-size: 12px;
          color: #666;
        }
      }
    }
    
    .more-items {
      font-size: 12px;
      color: #999;
    }
  }

  .final-amount {
    color: #f56c6c;
    font-weight: 500;
  }

  .discount {
    font-size: 12px;
    color: #67c23a;
  }

  .pagination {
    margin-top: 20px;
    text-align: center;
  }
}
</style>