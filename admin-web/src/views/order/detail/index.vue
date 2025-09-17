<template>
  <div class="order-detail-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>订单详情</span>
          <el-button @click="handleBack">返回列表</el-button>
        </div>
      </template>

      <div v-if="order" class="order-detail">
        <!-- 订单基本信息 -->
        <div class="detail-section">
          <h3>订单基本信息</h3>
          <el-descriptions :column="3" border>
            <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
            <el-descriptions-item label="订单状态">
              <el-tag :type="getStatusType(order.status)">
                {{ getStatusText(order.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="支付状态">
              <el-tag :type="getPaymentStatusType(order.payment_status)">
                {{ getPaymentStatusText(order.payment_status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="支付方式">{{ order.payment_method || '未知' }}</el-descriptions-item>
            <el-descriptions-item label="支付时间">{{ order.payment_time || '未支付' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ order.created_at }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 用户信息 -->
        <div class="detail-section">
          <h3>用户信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户姓名">{{ order.user_name }}</el-descriptions-item>
            <el-descriptions-item label="手机号码">{{ order.user_phone }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 收货地址 -->
        <div class="detail-section">
          <h3>收货地址</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="收货人">{{ order.shipping_address.name }}</el-descriptions-item>
            <el-descriptions-item label="联系电话">{{ order.shipping_address.phone }}</el-descriptions-item>
            <el-descriptions-item label="收货地址" :span="2">
              {{ order.shipping_address.province }}{{ order.shipping_address.city }}{{ order.shipping_address.district }}{{ order.shipping_address.address }}
            </el-descriptions-item>
            <el-descriptions-item label="邮政编码">{{ order.shipping_address.postcode || '未填写' }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 商品信息 -->
        <div class="detail-section">
          <h3>商品信息</h3>
          <el-table :data="order.items" border>
            <el-table-column label="商品" min-width="300">
              <template #default="{ row }">
                <div class="product-info">
                  <el-image
                    :src="row.product_image"
                    style="width: 60px; height: 60px; margin-right: 12px"
                    fit="cover"
                  />
                  <div class="product-details">
                    <div class="product-name">{{ row.product_name }}</div>
                    <div class="product-sku">SKU: {{ row.sku }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="price" label="单价" width="120">
              <template #default="{ row }">
                ¥{{ row.price }}
              </template>
            </el-table-column>
            <el-table-column prop="quantity" label="数量" width="80" />
            <el-table-column prop="total_amount" label="小计" width="120">
              <template #default="{ row }">
                ¥{{ row.total_amount }}
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 价格信息 -->
        <div class="detail-section">
          <h3>价格信息</h3>
          <el-descriptions :column="1" border class="price-info">
            <el-descriptions-item label="商品总额">¥{{ order.total_amount }}</el-descriptions-item>
            <el-descriptions-item label="优惠金额">¥{{ order.discount_amount }}</el-descriptions-item>
            <el-descriptions-item label="实付金额">
              <span class="final-amount">¥{{ order.final_amount }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 操作按钮 -->
        <div class="detail-section">
          <h3>订单操作</h3>
          <div class="action-buttons">
            <el-button
              v-if="order.status === 1"
              type="success"
              @click="handleShip"
            >
              <el-icon><Van /></el-icon>
              发货
            </el-button>
            <el-button
              v-if="order.status === 0"
              type="danger"
              @click="handleCancel"
            >
              <el-icon><Close /></el-icon>
              取消订单
            </el-button>
            <el-button
              v-if="order.status === 2"
              type="primary"
              @click="handleConfirm"
            >
              <el-icon><Check /></el-icon>
              确认收货
            </el-button>
            <el-button
              v-if="order.status >= 2"
              type="warning"
              @click="handleRefund"
            >
              <el-icon><RefreshLeft /></el-icon>
              申请退款
            </el-button>
          </div>
        </div>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus'
import { Van, Close, Check, RefreshLeft } from '@element-plus/icons-vue'
import {
  getOrder,
  updateOrderStatus,
  cancelOrder,
  shipOrder,
  confirmOrder,
  refundOrder,
  type Order
} from '@/api/order'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const shipping = ref(false)
const shipDialogVisible = ref(false)
const shipFormRef = ref<FormInstance>()

const orderId = computed(() => route.params.id as string)

const order = ref<Order | null>(null)

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

const loadOrder = async () => {
  try {
    loading.value = true
    const { data } = await getOrder(Number(orderId.value))
    order.value = data
  } catch (error) {
    ElMessage.error('加载订单详情失败')
    handleBack()
  } finally {
    loading.value = false
  }
}

const handleShip = () => {
  shipForm.shipping_company = ''
  shipForm.tracking_no = ''
  shipDialogVisible.value = true
}

const confirmShip = async () => {
  if (!shipFormRef.value || !order.value) return
  
  try {
    await shipFormRef.value.validate()
    shipping.value = true
    
    await shipOrder(order.value.id, shipForm)
    ElMessage.success('发货成功')
    shipDialogVisible.value = false
    loadOrder()
  } catch (error) {
    ElMessage.error('发货失败')
  } finally {
    shipping.value = false
  }
}

const handleCancel = async () => {
  if (!order.value) return
  
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
    
    await cancelOrder(order.value.id, reason)
    ElMessage.success('订单取消成功')
    loadOrder()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleConfirm = async () => {
  if (!order.value) return
  
  try {
    await ElMessageBox.confirm('确认该订单已收货？', '确认收货', { type: 'warning' })
    await confirmOrder(order.value.id)
    ElMessage.success('确认收货成功')
    loadOrder()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleRefund = async () => {
  if (!order.value) return
  
  try {
    const { value } = await ElMessageBox.prompt(
      '请输入退款原因',
      '申请退款',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        inputPlaceholder: '请输入退款原因'
      }
    )
    
    await refundOrder(order.value.id, { amount: order.value.final_amount, reason: value })
    ElMessage.success('退款申请成功')
    loadOrder()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleBack = () => {
  router.push('/order/list')
}

onMounted(() => {
  loadOrder()
})
</script>

<style scoped lang="scss">
.order-detail-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .order-detail {
    .detail-section {
      margin-bottom: 30px;

      h3 {
        margin-bottom: 16px;
        color: #303133;
        font-size: 16px;
        font-weight: 600;
      }
    }

    .product-info {
      display: flex;
      align-items: center;

      .product-details {
        .product-name {
          font-weight: 500;
          margin-bottom: 4px;
        }

        .product-sku {
          font-size: 12px;
          color: #666;
        }
      }
    }

    .price-info {
      max-width: 400px;

      .final-amount {
        color: #f56c6c;
        font-weight: bold;
        font-size: 16px;
      }
    }

    .action-buttons {
      .el-button {
        margin-right: 12px;
        margin-bottom: 12px;
      }
    }
  }
}
</style>