<template>
  <div class="product-detail">
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
                    <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                    <el-dropdown-item command="create-product" v-if="authStore.user">发布商品</el-dropdown-item>
                    <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div v-else class="auth-links">
              <router-link to="/login" class="auth-link">登录</router-link>
              <router-link to="/register" class="auth-link register">注册</router-link>
            </div>
          </nav>
        </div>
      </div>
    </header>

    <div class="container">
      <div v-if="loading" class="loading">
        <el-icon class="is-loading">
          <Loading />
        </el-icon>
      </div>

      <div v-else-if="!product" class="error">
        <p>商品不存在</p>
        <router-link to="/products" class="back-link">返回商品列表</router-link>
      </div>

      <div v-else class="product-content">
        <div class="product-main">
          <div class="product-images">
            <div class="main-image">
              <img :src="product.img_path" :alt="product.name" />
            </div>
            <div v-if="productImages.length > 0" class="image-gallery">
              <img
                v-for="(image, index) in productImages"
                :key="index"
                :src="image"
                :alt="`商品图片 ${index + 1}`"
                @click="changeMainImage(image)"
                :class="{ active: image === mainImage }"
              />
            </div>
          </div>

          <div class="product-info">
            <h1 class="product-name">{{ product.name }}</h1>
            <p class="product-title">{{ product.title }}</p>

            <div class="product-price">
              <span class="current-price">¥{{ product.discount_price }}</span>
              <span v-if="product.discount_price !== product.price" class="original-price">
                ¥{{ product.price }}
              </span>
            </div>

            <div class="product-meta">
              <div class="meta-item">
                <span class="label">库存:</span>
                <span class="value" :class="{ 'out-of-stock': product.num === 0 }">
                  {{ product.num }}
                </span>
              </div>
              <div class="meta-item">
                <span class="label">状态:</span>
                <span class="value" :class="{ 'not-on-sale': !product.on_sale }">
                  {{ product.on_sale ? '在售' : '已下架' }}
                </span>
              </div>
            </div>

            <div class="seller-info">
              <img :src="product.boss_avatar" :alt="product.boss_name" class="seller-avatar" />
              <div class="seller-details">
                <span class="seller-name">{{ product.boss_name }}</span>
                <span class="seller-label">卖家</span>
              </div>
            </div>

            <div class="product-description">
              <h3>商品描述</h3>
              <p>{{ product.info }}</p>
            </div>

            <!-- 编辑按钮（仅商品拥有者可见） -->
            <div v-if="authStore.user && authStore.user.id === product.boss_id" class="product-actions">
              <el-button type="primary" @click="editProduct">编辑商品</el-button>
              <el-button type="danger" @click="deleteProduct">删除商品</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Loading } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getProductDetail, deleteProduct as deleteProductApi, getProductImages, type Product } from '@/services/product'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const product = ref<Product | null>(null)
const productImages = ref<string[]>([])
const mainImage = ref('')
const loading = ref(true)

const handleUserCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
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

const loadProduct = async () => {
  const productId = route.params.id as string
  if (!productId) {
    ElMessage.error('商品ID不存在')
    router.push('/products')
    return
  }

  try {
    const response = await getProductDetail(productId)
    if (response.status === 200) {
      product.value = response.data
      mainImage.value = product.value.img_path

      // 加载商品图片
      try {
        const imagesResponse = await getProductImages(productId)
        if (imagesResponse.status === 200) {
          productImages.value = imagesResponse.data || []
        }
      } catch (error) {
        console.error('加载商品图片失败:', error)
      }
    } else {
      ElMessage.error(response.msg || '获取商品详情失败')
    }
  } catch (error: any) {
    ElMessage.error('获取商品详情失败')
  } finally {
    loading.value = false
  }
}

const changeMainImage = (imageUrl: string) => {
  mainImage.value = imageUrl
}

const editProduct = () => {
  router.push(`/product/edit/${product.value?.id}`)
}

const deleteProduct = async () => {
  if (!product.value) return

  try {
    await ElMessageBox.confirm('确定要删除这个商品吗？此操作不可恢复。', '确认删除', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const loading = ElLoading.service({
      lock: true,
      text: '正在删除商品...',
      background: 'rgba(0, 0, 0, 0.7)'
    })

    try {
      const response = await deleteProductApi(product.value.id.toString())
      loading.close()

      if (response.status === 200) {
        ElMessage.success('商品删除成功')
        // 延迟跳转，让用户看到成功消息
        setTimeout(() => {
          router.push('/products')
        }, 1500)
      } else {
        ElMessage.error(response.msg || '删除失败')
      }
    } catch (apiError: any) {
      loading.close()

      // 处理不同的错误类型
      if (apiError.response?.status === 403) {
        ElMessage.error('您没有权限删除此商品')
      } else if (apiError.response?.status === 404) {
        ElMessage.error('商品不存在或已被删除')
        // 如果商品不存在，跳转到商品列表页
        setTimeout(() => {
          router.push('/products')
        }, 2000)
      } else if (apiError.response?.status === 401) {
        ElMessage.error('请先登录')
        router.push('/login')
      } else {
        ElMessage.error('删除失败，请稍后重试')
      }
    }
  } catch (error) {
    // 用户取消删除操作
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

onMounted(() => {
  loadProduct()
})
</script>

<style scoped>
.product-detail {
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

.auth-links {
  display: flex;
  gap: 1rem;
}

.auth-link {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s;
}

.auth-link:not(.register) {
  color: #667eea;
  border: 1px solid #667eea;
}

.auth-link:not(.register):hover {
  background: #667eea;
  color: white;
}

.auth-link.register {
  background: #667eea;
  color: white;
}

.auth-link.register:hover {
  background: #5a67d8;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.loading {
  text-align: center;
  padding: 3rem;
}

.error {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.back-link {
  color: #667eea;
  text-decoration: none;
  margin-top: 1rem;
  display: inline-block;
}

.back-link:hover {
  text-decoration: underline;
}

.product-content {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.product-main {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
}

.product-images {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.main-image {
  width: 100%;
  height: 400px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e1e5e9;
}

.main-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-gallery {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.image-gallery img {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.3s;
}

.image-gallery img:hover,
.image-gallery img.active {
  border-color: #667eea;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.product-name {
  font-size: 2rem;
  font-weight: bold;
  color: #333;
  margin: 0;
}

.product-title {
  font-size: 1.1rem;
  color: #666;
  margin: 0;
}

.product-price {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.current-price {
  font-size: 2rem;
  font-weight: bold;
  color: #e74c3c;
}

.original-price {
  font-size: 1.5rem;
  color: #999;
  text-decoration: line-through;
}

.product-meta {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.label {
  font-weight: 500;
  color: #666;
}

.value {
  color: #333;
}

.value.out-of-stock {
  color: #e74c3c;
}

.value.not-on-sale {
  color: #f39c12;
}

.seller-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.seller-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.seller-details {
  display: flex;
  flex-direction: column;
}

.seller-name {
  font-weight: 500;
  color: #333;
}

.seller-label {
  font-size: 0.9rem;
  color: #666;
}

.product-description {
  padding-top: 1rem;
  border-top: 1px solid #e1e5e9;
}

.product-description h3 {
  font-size: 1.2rem;
  font-weight: bold;
  color: #333;
  margin-bottom: 1rem;
}

.product-description p {
  color: #666;
  line-height: 1.6;
}

.product-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .product-main {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .product-name {
    font-size: 1.5rem;
  }

  .current-price {
    font-size: 1.5rem;
  }

  .original-price {
    font-size: 1.2rem;
  }
}
</style>
