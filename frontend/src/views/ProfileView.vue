<template>
  <div class="profile">
    <!-- 导航栏 -->
    <header class="header">
      <div class="container">
        <div class="header-content">
          <h1 class="logo">上海大学二手商城</h1>
          <nav class="nav">
            <router-link to="/" class="nav-link">首页</router-link>
            <router-link to="/products" class="nav-link">商品</router-link>
            <div v-if="authStore.isLoggedIn" class="user-menu">
              <el-dropdown @command="handleUserCommand">
                <span class="user-name">
                  {{ authStore.user?.nickname }}
                  <el-icon class="el-icon--right">
                    <arrow-down />
                  </el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile" disabled>个人资料</el-dropdown-item>
                    <el-dropdown-item command="create-product" v-if="authStore.user">发布商品</el-dropdown-item>
                    <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </nav>
        </div>
      </div>
    </header>

    <div class="container">
      <div class="profile-content">
        <!-- 左侧菜单 -->
        <div class="profile-sidebar">
          <div class="menu-item" :class="{ active: activeTab === 'info' }" @click="activeTab = 'info'">
            <el-icon><User /></el-icon>
            基本信息
          </div>
          <div class="menu-item" :class="{ active: activeTab === 'email' }" @click="activeTab = 'email'">
            <el-icon><Message /></el-icon>
            邮箱绑定
          </div>
          <div class="menu-item" :class="{ active: activeTab === 'password' }" @click="activeTab = 'password'">
            <el-icon><Lock /></el-icon>
            修改密码
          </div>
          <div class="menu-item" :class="{ active: activeTab === 'avatar' }" @click="activeTab = 'avatar'">
            <el-icon><Picture /></el-icon>
            头像设置
          </div>
        </div>

        <!-- 右侧内容 -->
        <div class="profile-main">
          <!-- 基本信息 -->
          <div v-if="activeTab === 'info'" class="profile-section">
            <h2 class="section-title">基本信息</h2>
            <el-form
              ref="infoFormRef"
              :model="infoForm"
              :rules="infoRules"
              label-width="100px"
              class="profile-form"
            >
              <el-form-item label="用户名">
                <el-input :value="authStore.user?.username" disabled />
              </el-form-item>

              <el-form-item label="昵称" prop="nickname">
                <el-input
                  v-model="infoForm.nickname"
                  placeholder="请输入昵称"
                  :maxlength="10"
                  show-word-limit
                />
              </el-form-item>

              <el-form-item label="邮箱">
                <el-input :value="authStore.user?.email" disabled />
              </el-form-item>

              <el-form-item label="注册时间">
                <el-input :value="formatDate(authStore.user?.createat)" disabled />
              </el-form-item>

              <el-form-item label="账户余额">
                <el-input :value="`¥${authStore.user?.money}`" disabled />
              </el-form-item>

              <el-form-item>
                <el-button
                  type="primary"
                  :loading="updateLoading"
                  @click="updateUserInfo"
                >
                  保存修改
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 邮箱绑定 -->
          <div v-if="activeTab === 'email'" class="profile-section">
            <h2 class="section-title">邮箱绑定</h2>

            <div v-if="authStore.user?.email" class="email-status">
              <el-alert
                title="邮箱已绑定"
                :description="authStore.user.email"
                type="success"
                show-icon
                :closable="false"
              />
            </div>

            <div v-else class="email-bind">
              <p class="email-tip">绑定邮箱可以增强账户安全性，并用于找回密码。</p>

              <el-form
                ref="emailFormRef"
                :model="emailForm"
                :rules="emailRules"
                label-width="100px"
                class="profile-form"
              >
                <el-form-item label="邮箱地址" prop="email">
                  <el-input
                    v-model="emailForm.email"
                    placeholder="请输入邮箱地址"
                    type="email"
                  />
                </el-form-item>

                <el-form-item label="操作类型" prop="operation_type">
                  <el-radio-group v-model="emailForm.operation_type">
                    <el-radio :value="1">绑定邮箱</el-radio>
                    <el-radio :value="2">解绑邮箱</el-radio>
                  </el-radio-group>
                </el-form-item>

                <el-form-item v-if="emailForm.operation_type === 1" label="验证码" prop="code">
                  <div class="code-input-group">
                    <el-input
                      v-model="emailForm.code"
                      placeholder="请输入验证码"
                      style="flex: 1; margin-right: 10px;"
                    />
                    <el-button
                      type="primary"
                      :loading="sendCodeLoading"
                      :disabled="codeButtonDisabled"
                      @click="sendVerificationCode"
                    >
                      {{ codeButtonText }}
                    </el-button>
                  </div>
                </el-form-item>

                <el-form-item>
                  <el-button
                    type="primary"
                    :loading="emailLoading"
                    @click="handleEmailAction"
                  >
                    {{ emailForm.operation_type === 1 ? '绑定邮箱' : '解绑邮箱' }}
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </div>

          <!-- 修改密码 -->
          <div v-if="activeTab === 'password'" class="profile-section">
            <h2 class="section-title">修改密码</h2>
            <el-form
              ref="passwordFormRef"
              :model="passwordForm"
              :rules="passwordRules"
              label-width="120px"
              class="profile-form"
            >
              <el-form-item label="当前密码" prop="currentPassword">
                <el-input
                  v-model="passwordForm.currentPassword"
                  type="password"
                  placeholder="请输入当前密码"
                  show-password
                />
              </el-form-item>

              <el-form-item label="新密码" prop="newPassword">
                <el-input
                  v-model="passwordForm.newPassword"
                  type="password"
                  placeholder="请输入新密码"
                  show-password
                />
              </el-form-item>

              <el-form-item label="确认新密码" prop="confirmPassword">
                <el-input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  placeholder="请再次输入新密码"
                  show-password
                />
              </el-form-item>

              <el-form-item>
                <el-button
                  type="primary"
                  :loading="passwordLoading"
                  @click="changePassword"
                >
                  修改密码
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 头像设置 -->
          <div v-if="activeTab === 'avatar'" class="profile-section">
            <h2 class="section-title">头像设置</h2>

            <div class="avatar-section">
              <div class="current-avatar">
                <img :src="authStore.user?.avater" :alt="authStore.user?.nickname" />
              </div>

              <el-upload
                ref="uploadRef"
                class="avatar-uploader"
                action="/api/v1/avatar"
                :headers="uploadHeaders"
                :show-file-list="false"
                :on-success="handleAvatarSuccess"
                :on-error="handleAvatarError"
                :before-upload="beforeAvatarUpload"
              >
                <el-button type="primary" :loading="avatarLoading">
                  选择新头像
                </el-button>
              </el-upload>

              <p class="avatar-tip">支持 JPG、PNG 格式，文件大小不超过 2MB</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElForm } from 'element-plus'
import { ArrowDown, User, Message, Lock, Picture } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { updateUser, sendEmail, validEmail, uploadAvatar } from '@/services/user'

const router = useRouter()
const authStore = useAuthStore()

const activeTab = ref('info')
const infoFormRef = ref<InstanceType<typeof ElForm>>()
const emailFormRef = ref<InstanceType<typeof ElForm>>()
const passwordFormRef = ref<InstanceType<typeof ElForm>>()

const updateLoading = ref(false)
const emailLoading = ref(false)
const passwordLoading = ref(false)
const avatarLoading = ref(false)
const sendCodeLoading = ref(false)
const codeButtonDisabled = ref(false)
const countdown = ref(0)

const infoForm = reactive({
  nickname: authStore.user?.nickname || ''
})

const emailForm = reactive({
  email: '',
  operation_type: 1,
  code: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const infoRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 10, message: '昵称长度应为2-10个字符', trigger: 'blur' }
  ]
}

const emailRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  operation_type: [
    { required: true, message: '请选择操作类型', trigger: 'change' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' }
  ]
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
    { min: 8, max: 16, message: '密码长度应为8-16个字符', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, max: 16, message: '密码长度应为8-16个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const uploadHeaders = computed(() => ({
  Authorization: authStore.token
}))

const codeButtonText = computed(() => {
  return countdown.value > 0 ? `${countdown.value}s` : '发送验证码'
})

const handleUserCommand = (command: string) => {
  switch (command) {
    case 'profile':
      break
    case 'create-product':
      router.push('/product/create')
      break
    case 'logout':
      authStore.logout()
      router.push('/')
      break
  }
}

const updateUserInfo = async () => {
  if (!infoFormRef.value) return

  try {
    await infoFormRef.value.validate()
    updateLoading.value = true

    const response = await updateUser({ nickname: infoForm.nickname })
    if (response.status === 200) {
      authStore.updateUser(response.data)
      ElMessage.success('信息更新成功')
    } else {
      ElMessage.error(response.msg || '更新失败')
    }
  } catch (error) {
    console.error('Update validation failed:', error)
  } finally {
    updateLoading.value = false
  }
}

const sendVerificationCode = async () => {
  if (!emailForm.email) {
    ElMessage.warning('请先输入邮箱地址')
    return
  }

  sendCodeLoading.value = true
  try {
    const response = await sendEmail({
      email: emailForm.email,
      password: '', // 发送验证码时不需要密码
      operation_type: emailForm.operation_type
    })

    if (response.status === 200) {
      ElMessage.success('验证码发送成功')
      startCountdown()
    } else {
      ElMessage.error(response.msg || '验证码发送失败')
    }
  } catch (error: any) {
    ElMessage.error('验证码发送失败')
  } finally {
    sendCodeLoading.value = false
  }
}

const startCountdown = () => {
  countdown.value = 60
  codeButtonDisabled.value = true

  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      codeButtonDisabled.value = false
      clearInterval(timer)
    }
  }, 1000)
}

const handleEmailAction = async () => {
  if (!emailFormRef.value) return

  try {
    await emailFormRef.value.validate()
    emailLoading.value = true

    if (emailForm.operation_type === 1) {
      // 绑定邮箱验证
      const response = await validEmail({
        email: emailForm.email,
        operation_type: emailForm.operation_type,
        code: emailForm.code
      })

      if (response.status === 200) {
        ElMessage.success('邮箱绑定成功')
        // 重新获取用户信息
        location.reload()
      } else {
        ElMessage.error(response.msg || '邮箱绑定失败')
      }
    } else {
      // 解绑邮箱
      const response = await sendEmail({
        email: emailForm.email,
        password: '', // 解绑时也不需要密码
        operation_type: emailForm.operation_type
      })

      if (response.status === 200) {
        ElMessage.success('邮箱解绑成功')
        location.reload()
      } else {
        ElMessage.error(response.msg || '邮箱解绑失败')
      }
    }
  } catch (error) {
    console.error('Email action validation failed:', error)
  } finally {
    emailLoading.value = false
  }
}

const changePassword = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    passwordLoading.value = true

    // 使用发送邮箱功能来修改密码
    const response = await sendEmail({
      email: authStore.user?.email || '',
      password: passwordForm.currentPassword,
      operation_type: 3 // 改密码
    })

    if (response.status === 200) {
      ElMessage.success('密码修改成功，请重新登录')
      authStore.logout()
      router.push('/login')
    } else {
      ElMessage.error(response.msg || '密码修改失败')
    }
  } catch (error) {
    console.error('Password change validation failed:', error)
  } finally {
    passwordLoading.value = false
  }
}

const handleAvatarSuccess = (response: any) => {
  if (response.status === 200) {
    authStore.updateUser(response.data)
    ElMessage.success('头像上传成功')
  } else {
    ElMessage.error(response.msg || '头像上传失败')
  }
  avatarLoading.value = false
}

const handleAvatarError = () => {
  ElMessage.error('头像上传失败')
  avatarLoading.value = false
}

const beforeAvatarUpload = (file: File) => {
  const isValidType = ['image/jpeg', 'image/jpg', 'image/png'].includes(file.type)
  const isValidSize = file.size / 1024 / 1024 < 2

  if (!isValidType) {
    ElMessage.error('只能上传 JPG/PNG 格式的图片')
    return false
  }

  if (!isValidSize) {
    ElMessage.error('图片大小不能超过 2MB')
    return false
  }

  avatarLoading.value = true
  return true
}

const formatDate = (timestamp?: number) => {
  if (!timestamp) return ''
  return new Date(timestamp * 1000).toLocaleString()
}
</script>

<style scoped>
.profile {
  min-height: 100vh;
  background: #f8f9fa;
}

.header {
  background: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
}

.logo {
  font-size: 24px;
  font-weight: bold;
  color: #667eea;
  margin: 0;
}

.nav {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.nav-link {
  color: #333;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: #667eea;
}

.user-menu {
  cursor: pointer;
}

.user-name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #333;
  font-weight: 500;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.profile-content {
  display: grid;
  grid-template-columns: 250px 1fr;
  gap: 2rem;
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.profile-sidebar {
  border-right: 1px solid #e1e5e9;
  padding-right: 2rem;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  margin-bottom: 0.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  color: #666;
}

.menu-item:hover,
.menu-item.active {
  background: #f0f2f5;
  color: #667eea;
}

.profile-main {
  padding-left: 2rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: bold;
  color: #333;
  margin-bottom: 2rem;
}

.profile-form {
  max-width: 500px;
}

.email-status {
  margin-bottom: 2rem;
}

.email-tip {
  color: #666;
  margin-bottom: 1rem;
}

.code-input-group {
  display: flex;
  align-items: center;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.current-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #e1e5e9;
}

.current-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-uploader {
  text-align: center;
}

.avatar-tip {
  color: #999;
  font-size: 0.9rem;
  text-align: center;
}

@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .profile-sidebar {
    border-right: none;
    border-bottom: 1px solid #e1e5e9;
    padding-right: 0;
    padding-bottom: 1rem;
  }

  .profile-main {
    padding-left: 0;
    padding-top: 2rem;
  }
}
</style>
