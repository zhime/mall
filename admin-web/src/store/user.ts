import { defineStore } from 'pinia'
import { login, getUserInfo, logout } from '@/api/auth'
import type { LoginForm, UserInfo } from '@/types/user'

interface UserState {
  token: string
  userInfo: UserInfo | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem('admin_token') || '',
    userInfo: null
  }),

  getters: {
    isLoggedIn: (state) => !!state.token
  },

  actions: {
    // 登录
    async login(loginForm: LoginForm) {
      try {
        const { data } = await login(loginForm)
        this.token = data.token
        localStorage.setItem('admin_token', data.token)
        return Promise.resolve(data)
      } catch (error) {
        return Promise.reject(error)
      }
    },

    // 获取用户信息
    async getUserInfo() {
      try {
        const { data } = await getUserInfo()
        this.userInfo = data
        return Promise.resolve(data)
      } catch (error) {
        this.logout()
        return Promise.reject(error)
      }
    },

    // 登出
    async logout() {
      try {
        await logout()
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        this.token = ''
        this.userInfo = null
        localStorage.removeItem('admin_token')
      }
    },

    // 设置token
    setToken(token: string) {
      this.token = token
      localStorage.setItem('admin_token', token)
    }
  }
})