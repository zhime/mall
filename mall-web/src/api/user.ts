import request from '@/utils/request'

export interface User {
  id: number
  phone: string
  nickname?: string
  avatar?: string
  gender: number
  birthday?: string
  email?: string
  status: number
  created_at: string
  updated_at: string
}

export interface LoginParams {
  phone: string
  code: string
  type: 'sms' | 'wechat'
}

export interface RegisterParams {
  phone: string
  code: string
  password?: string
  nickname?: string
}

export interface ProfileParams {
  nickname?: string
  avatar?: string
  gender?: number
  birthday?: string
  email?: string
}

// 发送短信验证码
export function sendSmsCode(phone: string, type: 'login' | 'register' = 'login') {
  return request({
    url: '/user/sms/send',
    method: 'post',
    data: { phone, type }
  })
}

// 手机号登录
export function loginByPhone(data: LoginParams) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

// 微信登录
export function loginByWechat(code: string) {
  return request({
    url: '/user/wechat/login',
    method: 'post',
    data: { code }
  })
}

// 注册
export function register(data: RegisterParams) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/user/profile',
    method: 'get'
  })
}

// 更新用户信息
export function updateProfile(data: ProfileParams) {
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

// 上传头像
export function uploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/user/avatar',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}