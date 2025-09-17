import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { useUserStore } from '@/store/user'
import router from '@/router'

// 请求和响应的类型定义
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    const userStore = useUserStore()
    
    // 添加 token
    if (userStore.token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${userStore.token}`
      }
    }
    
    // 显示加载提示（除了某些不需要loading的请求）
    if (!config.hideLoading) {
      showLoadingToast({
        message: '加载中...',
        forbidClick: true,
        duration: 0
      })
    }
    
    return config
  },
  (error) => {
    closeToast()
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    closeToast()
    
    const { code, message, data } = response.data
    
    if (code === 200) {
      return { data, message }
    } else if (code === 401) {
      // token 过期或无效
      const userStore = useUserStore()
      userStore.logout()
      router.push('/login')
      showToast('登录已过期，请重新登录')
      return Promise.reject(new Error(message))
    } else {
      showToast(message || '请求失败')
      return Promise.reject(new Error(message))
    }
  },
  (error) => {
    closeToast()
    
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      router.push('/login')
      showToast('登录已过期，请重新登录')
    } else if (error.code === 'ECONNABORTED') {
      showToast('请求超时，请重试')
    } else if (error.message === 'Network Error') {
      showToast('网络连接失败')
    } else {
      showToast(error.response?.data?.message || '请求失败')
    }
    
    return Promise.reject(error)
  }
)

export default request

// 添加请求配置类型扩展
declare module 'axios' {
  interface AxiosRequestConfig {
    hideLoading?: boolean
  }
}