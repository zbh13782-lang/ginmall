<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">登录</h2>
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="0px"
        class="login-form"
      >
        <el-form-item prop="user_name">
          <el-input
            v-model="loginForm.user_name"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <span>还没有账号？</span>
        <router-link to="/register" class="register-link">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElForm, ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { LoginData } from '@/services/user'

const router = useRouter()
const authStore = useAuthStore()

const loginFormRef = ref<InstanceType<typeof ElForm>>()
const loading = ref(false)

const loginForm = reactive<LoginData>({
  user_name: '',
  password: ''
})

const loginRules = {
  user_name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 5, max: 15, message: '用户名长度应为5-15个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, max: 16, message: '密码长度应为8-16个字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    const success = await authStore.login(loginForm)
    if (success) {
      router.push('/')
    }
  } catch (error) {
    console.error('Login validation failed:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.login-title {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
  font-size: 28px;
  font-weight: 600;
}

.login-form {
  margin-bottom: 20px;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
}

.login-footer {
  text-align: center;
  color: #666;
}

.register-link {
  color: #667eea;
  text-decoration: none;
  margin-left: 5px;
  font-weight: 500;
}

.register-link:hover {
  text-decoration: underline;
}
</style>
