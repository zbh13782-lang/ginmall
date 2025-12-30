import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginData, RegisterData } from '@/services/user'
import { userLogin, userRegister } from '@/services/user'
import { ElMessage } from 'element-plus'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)

  // 从localStorage恢复用户状态
  const initUser = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (error) {
        console.error('Failed to parse stored user:', error)
        localStorage.removeItem('user')
      }
    }
  }

  // 是否已登录
  const isLoggedIn = computed(() => !!token.value && !!user.value)

  // 登录
  const login = async (data: LoginData) => {
    try {
      const response = await userLogin(data)
      if (response.status === 200) {
        token.value = response.data.token
        user.value = response.data.user

        // 保存到localStorage
        localStorage.setItem('token', token.value)
        localStorage.setItem('user', JSON.stringify(user.value))

        ElMessage.success('登录成功')
        return true
      } else {
        ElMessage.error(response.msg || '登录失败')
        return false
      }
    } catch (error: any) {
      ElMessage.error(error.response?.data?.msg || '登录失败')
      return false
    }
  }

  // 注册
  const register = async (data: RegisterData) => {
    try {
      const response = await userRegister(data)
      if (response.status === 200) {
        token.value = response.data.token
        user.value = response.data.user

        // 保存到localStorage
        localStorage.setItem('token', token.value)
        localStorage.setItem('user', JSON.stringify(user.value))

        ElMessage.success('注册成功')
        return true
      } else {
        ElMessage.error(response.msg || '注册失败')
        return false
      }
    } catch (error: any) {
      ElMessage.error(error.response?.data?.msg || '注册失败')
      return false
    }
  }

  // 登出
  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    ElMessage.success('已登出')
  }

  // 更新用户信息
  const updateUser = (newUser: User) => {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  return {
    token,
    user,
    isLoggedIn,
    initUser,
    login,
    register,
    logout,
    updateUser
  }
})
