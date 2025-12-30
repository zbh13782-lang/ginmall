<template>
  <div class="products-page">
    <!-- 导航栏 -->
    <header class="header">
      <div class="container">
        <div class="header-content">
          <h1 class="logo">上海大学二手商城</h1>
          <nav class="nav">
            <router-link to="/" class="nav-link">首页</router-link>
            <router-link to="/products" class="nav-link router-link-active">商品</router-link>
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
      <!-- 搜索和筛选 -->
      <div class="filters-section">
        <div class="search-box">
          <el-input
            v-model="searchQuery"
            placeholder="搜索商品..."
            size="large"
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-button type="primary" @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>

        <div class="filters">
          <el-select v-model="selectedCategory" placeholder="选择分类" @change="loadProducts">
            <el-option label="全部" value="" />
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </div>
      </div>

      <!-- 商品列表 -->
      <div v-if="loading" class="loading">
        <el-icon class="is-loading">
          <Loading />
        </el-icon>
      </div>

      <div v-else-if="products.length === 0" class="no-products">
        <p>{{ searchQuery ? '没有找到相关商品' : '暂无商品' }}</p>
      </div>

      <div v-else class="products-grid">
        <div
          v-for="product in products"
          :key="product.id"
          class="product-card"
          @click="goToProductDetail(product.id)"
        >
          <div class="product-image">
            <img :src="product.img_path" :alt="product.name" />
          </div>
          <div class="product-info">
            <h4 class="product-name">{{ product.name }}</h4>
            <p class="product-title">{{ product.title }}</p>
            <div class="product-price">
              <span class="current-price">¥{{ product.discount_price }}</span>
              <span v-if="product.discount_price !== product.price" class="original-price">
                ¥{{ product.price }}
              </span>
            </div>
            <div class="product-meta">
              <span class="seller">{{ product.boss_name }}</span>
              <span class="stock" :class="{ 'out-of-stock': product.num === 0 }">
                库存: {{ product.num }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[12, 24, 36]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowDown, Search, Loading } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getProducts, searchProducts, getCategories, type Product, type Category } from '@/services/product'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const searchQuery = ref('')
const selectedCategory = ref('')
const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const isSearchMode = ref(false)

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

const handleSearch = async () => {
  if (!searchQuery.value.trim()) {
    isSearchMode.value = false
    loadProducts()
    return
  }

  isSearchMode.value = true
  loading.value = true

  try {
    const response = await searchProducts({ info: searchQuery.value.trim() })
    if (response.status === 200) {
      products.value = response.data.item
      total.value = response.data.total
    }
  } catch (error: any) {
    ElMessage.error('搜索失败')
  } finally {
    loading.value = false
  }
}

const loadProducts = async () => {
  if (isSearchMode.value) return

  loading.value = true
  try {
    const params = {
      pageNum: currentPage.value,
      pageSize: pageSize.value
    }

    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }

    const response = await getProducts(params)
    if (response.status === 200) {
      products.value = response.data.item
      total.value = response.data.total
    }
  } catch (error: any) {
    ElMessage.error('获取商品失败')
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const response = await getCategories()
    if (response.status === 200) {
      categories.value = response.data
    }
  } catch (error: any) {
    console.error('获取分类失败:', error)
  }
}

const goToProductDetail = (productId: number) => {
  router.push(`/product/${productId}`)
}

const handleSizeChange = (newSize: number) => {
  pageSize.value = newSize
  if (isSearchMode.value) {
    handleSearch()
  } else {
    loadProducts()
  }
}

const handleCurrentChange = (newPage: number) => {
  currentPage.value = newPage
  if (isSearchMode.value) {
    handleSearch()
  } else {
    loadProducts()
  }
}

// 监听路由查询参数
watch(
  () => route.query.search,
  (newSearch) => {
    if (newSearch) {
      searchQuery.value = newSearch as string
      handleSearch()
    }
  },
  { immediate: true }
)

onMounted(() => {
  loadCategories()
  if (!route.query.search) {
    loadProducts()
  }
})
</script>

<style scoped>
.products-page {
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
  padding: 0 1rem;
}

.filters-section {
  padding: 2rem 0;
  display: flex;
  gap: 2rem;
  align-items: center;
}

.search-box {
  flex: 1;
  max-width: 400px;
}

.filters {
  display: flex;
  gap: 1rem;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;
}

.product-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.product-image {
  height: 200px;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.product-info {
  padding: 1.5rem;
}

.product-name {
  font-size: 1.2rem;
  font-weight: bold;
  color: #333;
  margin-bottom: 0.5rem;
}

.product-title {
  color: #666;
  margin-bottom: 1rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-price {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.current-price {
  font-size: 1.3rem;
  font-weight: bold;
  color: #e74c3c;
}

.original-price {
  color: #999;
  text-decoration: line-through;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
  color: #666;
}

.seller {
  font-weight: 500;
}

.stock {
  color: #27ae60;
}

.stock.out-of-stock {
  color: #e74c3c;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-bottom: 3rem;
}

.loading {
  text-align: center;
  padding: 3rem;
}

.no-products {
  text-align: center;
  padding: 3rem;
  color: #666;
}
</style>
