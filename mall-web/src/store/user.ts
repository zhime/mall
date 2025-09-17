import { defineStore } from 'pinia'
import { storage } from '@/utils'
import type { User } from '@/api/user'

interface UserState {
  token: string | null
  userInfo: User | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: storage.get('token'),
    userInfo: storage.get('userInfo')
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    avatar: (state) => state.userInfo?.avatar || '',
    nickname: (state) => state.userInfo?.nickname || state.userInfo?.phone || '用户'
  },
  
  actions: {
    // 设置 token
    setToken(token: string) {
      this.token = token
      storage.set('token', token)
    },
    
    // 设置用户信息
    setUserInfo(userInfo: User) {
      this.userInfo = userInfo
      storage.set('userInfo', userInfo)
    },
    
    // 登录
    login(token: string, userInfo: User) {
      this.setToken(token)
      this.setUserInfo(userInfo)
    },
    
    // 登出
    logout() {
      this.token = null
      this.userInfo = null
      storage.remove('token')
      storage.remove('userInfo')
    },
    
    // 更新用户信息
    updateUserInfo(userInfo: Partial<User>) {
      if (this.userInfo) {
        this.userInfo = { ...this.userInfo, ...userInfo }
        storage.set('userInfo', this.userInfo)
      }
    }
  }
})