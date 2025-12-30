<template>
  <div class="register-container">
    <div class="register-card">
      <h2 class="register-title">注册</h2>
      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-width="0px"
        class="register-form"
      >
        <el-form-item prop="nickname">
          <el-input
            v-model="registerForm.nickname"
            placeholder="昵称"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="user_name">
          <el-input
            v-model="registerForm.user_name"
            placeholder="用户名"
            size="large"
            :prefix-icon="Avatar"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="确认密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item prop="key">
          <el-input
            v-model="registerForm.key"
            placeholder="密钥 (16位)"
            size="large"
            :prefix-icon="Key"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="register-btn"
            :loading="loading"
            @click="handleRegister"
          >
            注册
          </el-button>
        </el-form-item>
      </el-form>

      <div class="register-footer">
        <span>已有账号？</span>
        <router-link to="/login" class="login-link">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElForm, ElMessage } from 'element-plus'
import { User, Avatar, Lock, Key } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { RegisterFormData } from '@/services/user'

const router = useRouter()
const authStore = useAuthStore()

const registerFormRef = ref<InstanceType<typeof ElForm>>()
const loading = ref(false)

const registerForm = reactive<RegisterFormData>({
  nickname: '',
  user_name: '',
  password: '',
  key: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule: any, value: string, callback: Function) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const registerRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 10, message: '昵称长度应为2-10个字符', trigger: 'blur' }
  ],
  user_name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 5, max: 15, message: '用户名长度应为5-15个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, max: 16, message: '密码长度应为8-16个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  key: [
    { required: true, message: '请输入密钥', trigger: 'blur' },
    { len: 16, message: '密钥必须为16位', trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return

  try {
    await registerFormRef.value.validate()
    loading.value = true

    // 移除confirmPassword字段，只发送后端需要的字段
    const { confirmPassword, ...registerData } = registerForm
    const success = await authStore.register(registerData)
    if (success) {
      router.push('/')
    }
  } catch (error) {
    console.error('Register validation failed:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.register-title {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
  font-size: 28px;
  font-weight: 600;
}

.register-form {
  margin-bottom: 20px;
}

.register-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
}

.register-footer {
  text-align: center;
  color: #666;
}

.login-link {
  color: #667eea;
  text-decoration: none;
  margin-left: 5px;
  font-weight: 500;
}

.login-link:hover {
  text-decoration: underline;
}
</style>
