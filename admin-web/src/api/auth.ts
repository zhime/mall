import request from '@/utils/request'
import type { LoginForm, UserInfo } from '@/types/user'

// 登录
export const login = (data: LoginForm) => {
  return request({
    url: '/auth/login/password',
    method: 'post',
    data
  })
}

// 获取用户信息
export const getUserInfo = () => {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

// 登出
export const logout = () => {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}