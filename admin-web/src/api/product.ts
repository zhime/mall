import request from '@/utils/request'

export interface Category {
  id: number
  name: string
  parent_id: number
  level: number
  sort: number
  icon?: string
  description?: string
  status: number
  children?: Category[]
  created_at: string
  updated_at: string
}

export interface Product {
  id: number
  name: string
  category_id: number
  category_name?: string
  sku: string
  price: number
  original_price: number
  stock: number
  sales: number
  images: string[]
  description: string
  specs: ProductSpec[]
  attributes: ProductAttribute[]
  status: number
  is_featured: boolean
  is_new: boolean
  weight: number
  created_at: string
  updated_at: string
}

export interface ProductSpec {
  name: string
  value: string
}

export interface ProductAttribute {
  name: string
  value: string
}

export interface ProductQuery {
  page?: number
  page_size?: number
  keyword?: string
  category_id?: number
  status?: number
  is_featured?: boolean
  is_new?: boolean
  price_min?: number
  price_max?: number
}

export interface CategoryQuery {
  parent_id?: number
  level?: number
  status?: number
}

// 分类相关API
export function getCategoryList(params?: CategoryQuery) {
  return request({
    url: '/api/admin/categories',
    method: 'get',
    params
  })
}

export function getCategoryTree() {
  return request({
    url: '/api/admin/categories/tree',
    method: 'get'
  })
}

export function createCategory(data: Partial<Category>) {
  return request({
    url: '/api/admin/categories',
    method: 'post',
    data
  })
}

export function updateCategory(id: number, data: Partial<Category>) {
  return request({
    url: `/api/admin/categories/${id}`,
    method: 'put',
    data
  })
}

export function deleteCategory(id: number) {
  return request({
    url: `/api/admin/categories/${id}`,
    method: 'delete'
  })
}

export function updateCategoryStatus(id: number, status: number) {
  return request({
    url: `/api/admin/categories/${id}/status`,
    method: 'patch',
    data: { status }
  })
}

// 商品相关API
export function getProductList(params?: ProductQuery) {
  return request({
    url: '/api/admin/products',
    method: 'get',
    params
  })
}

export function getProduct(id: number) {
  return request({
    url: `/api/admin/products/${id}`,
    method: 'get'
  })
}

export function createProduct(data: Partial<Product>) {
  return request({
    url: '/api/admin/products',
    method: 'post',
    data
  })
}

export function updateProduct(id: number, data: Partial<Product>) {
  return request({
    url: `/api/admin/products/${id}`,
    method: 'put',
    data
  })
}

export function deleteProduct(id: number) {
  return request({
    url: `/api/admin/products/${id}`,
    method: 'delete'
  })
}

export function updateProductStatus(id: number, status: number) {
  return request({
    url: `/api/admin/products/${id}/status`,
    method: 'patch',
    data: { status }
  })
}

export function batchUpdateProductStatus(ids: number[], status: number) {
  return request({
    url: '/api/admin/products/batch/status',
    method: 'patch',
    data: { ids, status }
  })
}

export function batchDeleteProducts(ids: number[]) {
  return request({
    url: '/api/admin/products/batch',
    method: 'delete',
    data: { ids }
  })
}

// 图片上传
export function uploadImage(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/api/admin/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}