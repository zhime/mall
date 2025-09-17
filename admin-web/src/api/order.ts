import request from '@/utils/request'

export interface Order {
  id: number
  order_no: string
  user_id: number
  user_name?: string
  user_phone?: string
  total_amount: number
  discount_amount: number
  final_amount: number
  status: number
  payment_status: number
  payment_method?: string
  payment_time?: string
  shipping_address: ShippingAddress
  items: OrderItem[]
  created_at: string
  updated_at: string
}

export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  product_name: string
  product_image: string
  sku: string
  price: number
  quantity: number
  total_amount: number
}

export interface ShippingAddress {
  name: string
  phone: string
  province: string
  city: string
  district: string
  address: string
  postcode?: string
}

export interface OrderQuery {
  page?: number
  page_size?: number
  order_no?: string
  user_phone?: string
  status?: number
  payment_status?: number
  payment_method?: string
  start_date?: string
  end_date?: string
}

export interface OrderStatistics {
  total_orders: number
  pending_orders: number
  processing_orders: number
  completed_orders: number
  total_amount: number
  today_orders: number
  today_amount: number
}

// 订单列表
export function getOrderList(params?: OrderQuery) {
  return request({
    url: '/api/admin/orders',
    method: 'get',
    params
  })
}

// 订单详情
export function getOrder(id: number) {
  return request({
    url: `/api/admin/orders/${id}`,
    method: 'get'
  })
}

// 更新订单状态
export function updateOrderStatus(id: number, status: number) {
  return request({
    url: `/api/admin/orders/${id}/status`,
    method: 'patch',
    data: { status }
  })
}

// 更新支付状态
export function updatePaymentStatus(id: number, payment_status: number) {
  return request({
    url: `/api/admin/orders/${id}/payment-status`,
    method: 'patch',
    data: { payment_status }
  })
}

// 批量更新订单状态
export function batchUpdateOrderStatus(ids: number[], status: number) {
  return request({
    url: '/api/admin/orders/batch/status',
    method: 'patch',
    data: { ids, status }
  })
}

// 取消订单
export function cancelOrder(id: number, reason?: string) {
  return request({
    url: `/api/admin/orders/${id}/cancel`,
    method: 'patch',
    data: { reason }
  })
}

// 发货
export function shipOrder(id: number, data: { tracking_no: string; shipping_company: string }) {
  return request({
    url: `/api/admin/orders/${id}/ship`,
    method: 'patch',
    data
  })
}

// 确认收货
export function confirmOrder(id: number) {
  return request({
    url: `/api/admin/orders/${id}/confirm`,
    method: 'patch'
  })
}

// 申请退款
export function refundOrder(id: number, data: { amount: number; reason: string }) {
  return request({
    url: `/api/admin/orders/${id}/refund`,
    method: 'post',
    data
  })
}

// 订单统计
export function getOrderStatistics() {
  return request({
    url: '/api/admin/orders/statistics',
    method: 'get'
  })
}

// 导出订单
export function exportOrders(params?: OrderQuery) {
  return request({
    url: '/api/admin/orders/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}