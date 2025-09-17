import request from '@/utils/request'

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
}

export interface ProductQuery {
  page?: number
  page_size?: number
  keyword?: string
  category_id?: number
  price_min?: number
  price_max?: number
  sort?: 'sales' | 'price_asc' | 'price_desc' | 'newest'
  is_featured?: boolean
  is_new?: boolean
}

// 获取商品列表
export function getProductList(params?: ProductQuery) {
  return request({
    url: '/products',
    method: 'get',
    params
  })
}

// 获取商品详情
export function getProduct(id: number) {
  return request({
    url: `/products/${id}`,
    method: 'get'
  })
}

// 获取分类列表
export function getCategoryList() {
  return request({
    url: '/categories',
    method: 'get'
  })
}

// 获取分类树
export function getCategoryTree() {
  return request({
    url: '/categories/tree',
    method: 'get'
  })
}

// 搜索商品
export function searchProducts(keyword: string, params?: Omit<ProductQuery, 'keyword'>) {
  return request({
    url: '/products/search',
    method: 'get',
    params: { keyword, ...params }
  })
}

// 获取热门搜索关键词
export function getHotKeywords() {
  return request({
    url: '/search/hot',
    method: 'get'
  })
}

// 获取搜索建议
export function getSearchSuggestions(keyword: string) {
  return request({
    url: '/search/suggest',
    method: 'get',
    params: { keyword }
  })
}

// 获取首页推荐商品
export function getRecommendProducts() {
  return request({
    url: '/products/recommend',
    method: 'get'
  })
}

// 获取首页轮播图
export function getBanners() {
  return request({
    url: '/banners',
    method: 'get'
  })
}