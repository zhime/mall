// 登录表单
export interface LoginForm {
  username: string
  password: string
}

// 用户信息
export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  role: string
  status: number
}

// 用户列表项
export interface UserListItem {
  id: number
  username: string
  phone: string
  nickname: string
  status: number
  created_at: string
}

// 用户列表查询参数
export interface UserListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: number
}