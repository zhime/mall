import request from '@/utils/request'
import type { UserListQuery, UserListItem } from '@/types/user'

// 获取用户列表
export const getUserList = (params: UserListQuery) => {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

// 获取用户详情
export const getUserDetail = (id: number) => {
  return request({
    url: `/admin/users/${id}`,
    method: 'get'
  })
}

// 更新用户状态
export const updateUserStatus = (id: number, status: number) => {
  return request({
    url: `/admin/users/${id}/status`,
    method: 'put',
    data: { status }
  })
}

// 删除用户
export const deleteUser = (id: number) => {
  return request({
    url: `/admin/users/${id}`,
    method: 'delete'
  })
}

// 重置用户密码
export const resetUserPassword = (id: number, password: string) => {
  return request({
    url: `/admin/users/${id}/password`,
    method: 'put',
    data: { password }
  })
}