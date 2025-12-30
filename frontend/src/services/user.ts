import api from './api'

// 用户相关接口
export interface User {
  id: number
  username: string
  email: string
  nickname: string
  status: string
  avater: string
  money: string
  createat: number
}

export interface LoginData {
  user_name: string
  password: string
}

export interface RegisterData {
  nickname: string
  user_name: string
  password: string
  key: string
}

export interface RegisterFormData extends RegisterData {
  confirmPassword: string
}

export interface UpdateUserData {
  nickname: string
}

export interface SendEmailData {
  email: string
  password: string
  operation_type: number
}

export interface TokenData {
  token: string
  user: User
}

// 用户注册
export const userRegister = (data: RegisterData) => {
  return api.post<TokenData>('/user/register', data)
}

// 用户登录
export const userLogin = (data: LoginData) => {
  return api.post<TokenData>('/user/login', data)
}

// 更新用户信息
export const updateUser = (data: UpdateUserData) => {
  return api.put<User>('/user', data)
}

// 发送邮箱验证
export const sendEmail = (data: SendEmailData) => {
  return api.post('/user/sending-email', data)
}

// 验证邮箱
export const validEmail = (data: any) => {
  return api.post('/user/valid-email', data)
}

// 上传头像
export const uploadAvatar = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post('/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
