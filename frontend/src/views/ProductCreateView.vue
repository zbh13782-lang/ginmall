<template>
  <div class="product-create">
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
                    <el-dropdown-item command="create-product" disabled>发布商品</el-dropdown-item>
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
      <div class="create-content">
        <h2 class="page-title">发布商品</h2>

        <el-form
          ref="productFormRef"
          :model="productForm"
          :rules="productRules"
          label-width="120px"
          class="product-form"
        >
          <el-form-item label="商品名称" prop="name">
            <el-input
              v-model="productForm.name"
              placeholder="请输入商品名称"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="商品标题" prop="title">
            <el-input
              v-model="productForm.title"
              placeholder="请输入商品标题"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="商品分类" prop="category_id">
            <el-select v-model="productForm.category_id" placeholder="请选择商品分类">
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="商品价格" prop="price">
            <el-input-number
              v-model="productForm.price"
              :precision="2"
              :min="0"
              :max="999999"
              placeholder="请输入商品价格"
              style="width: 200px;"
            />
          </el-form-item>

          <el-form-item label="折扣价格" prop="discount_price">
            <el-input-number
              v-model="productForm.discount_price"
              :precision="2"
              :min="0"
              :max="999999"
              placeholder="请输入折扣价格"
              style="width: 200px;"
            />
            <span class="form-tip">留空表示无折扣</span>
          </el-form-item>

          <el-form-item label="库存数量" prop="num">
            <el-input-number
              v-model="productForm.num"
              :min="0"
              :max="999999"
              placeholder="请输入库存数量"
              style="width: 200px;"
            />
          </el-form-item>

          <el-form-item label="商品描述" prop="info">
            <el-input
              v-model="productForm.info"
              type="textarea"
              :rows="6"
              placeholder="请输入商品详细描述"
              maxlength="1000"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="是否上架" prop="on_sale">
            <el-radio-group v-model="productForm.on_sale">
              <el-radio :value="true">上架</el-radio>
              <el-radio :value="false">下架</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="商品图片" prop="images">
            <el-upload
              ref="uploadRef"
              v-model:file-list="fileList"
              class="upload-demo"
              action=""
              :auto-upload="false"
              :on-change="handleFileChange"
              :on-remove="handleFileRemove"
              list-type="picture-card"
              :limit="5"
              accept="image/*"
            >
              <el-icon><Plus /></el-icon>
              <template #tip>
                <div class="upload-tip">
                  支持 JPG、PNG 格式，单张图片不超过 5MB，最多上传 5 张图片
                </div>
              </template>
            </el-upload>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="submitLoading"
              @click="submitProduct"
            >
              发布商品
            </el-button>
            <el-button size="large" @click="cancel">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElForm, UploadFile } from 'element-plus'
import { ArrowDown, Plus } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { createProduct, getCategories, type Category } from '@/services/product'

const router = useRouter()
const authStore = useAuthStore()

const productFormRef = ref<InstanceType<typeof ElForm>>()
const uploadRef = ref()

const submitLoading = ref(false)
const fileList = ref<UploadFile[]>([])
const categories = ref<Category[]>([])

const productForm = reactive({
  name: '',
  title: '',
  category_id: undefined as number | undefined,
  price: undefined as number | undefined,
  discount_price: undefined as number | undefined,
  num: undefined as number | undefined,
  info: '',
  on_sale: true,
  images: [] as File[]
})

const productRules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 1, max: 100, message: '商品名称长度应为1-100个字符', trigger: 'blur' }
  ],
  title: [
    { required: true, message: '请输入商品标题', trigger: 'blur' },
    { min: 1, max: 200, message: '商品标题长度应为1-200个字符', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入商品价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格必须大于等于0', trigger: 'blur' }
  ],
  discount_price: [
    {
      type: 'number',
      min: 0,
      message: '折扣价格必须大于等于0',
      trigger: 'blur'
    },
    {
      validator: (rule: any, value: number, callback: Function) => {
        if (value && productForm.price && value > productForm.price) {
          callback(new Error('折扣价格不能高于原价'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  num: [
    { required: true, message: '请输入库存数量', trigger: 'blur' },
    { type: 'number', min: 0, message: '库存数量必须大于等于0', trigger: 'blur' }
  ],
  info: [
    { required: true, message: '请输入商品描述', trigger: 'blur' },
    { min: 10, max: 1000, message: '商品描述长度应为10-1000个字符', trigger: 'blur' }
  ],
  on_sale: [
    { required: true, message: '请选择是否上架', trigger: 'change' }
  ]
}

const handleUserCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'create-product':
      break
    case 'logout':
      authStore.logout()
      router.push('/')
      break
  }
}

const handleFileChange = (file: UploadFile, fileList: UploadFile[]) => {
  // 检查文件类型
  const isValidType = ['image/jpeg', 'image/jpg', 'image/png'].includes(file.raw?.type || '')
  if (!isValidType) {
    ElMessage.error('只能上传 JPG/PNG 格式的图片')
    return false
  }

  // 检查文件大小 (5MB)
  const isValidSize = (file.size || 0) / 1024 / 1024 < 5
  if (!isValidSize) {
    ElMessage.error('图片大小不能超过 5MB')
    return false
  }

  // 更新表单中的图片文件
  productForm.images = fileList.map(f => f.raw as File).filter(Boolean)
}

const handleFileRemove = (file: UploadFile, fileList: UploadFile[]) => {
  productForm.images = fileList.map(f => f.raw as File).filter(Boolean)
}

const submitProduct = async () => {
  if (!productFormRef.value) return

  try {
    await productFormRef.value.validate()

    if (productForm.images.length === 0) {
      ElMessage.error('请至少上传一张商品图片')
      return
    }

    submitLoading.value = true

    // 准备提交数据
    const submitData = {
      name: productForm.name,
      title: productForm.title,
      category_id: productForm.category_id!,
      price: productForm.price!.toString(),
      discount_price: productForm.discount_price ? productForm.discount_price.toString() : productForm.price!.toString(),
      num: productForm.num!,
      info: productForm.info,
      on_sale: productForm.on_sale
    }

    const response = await createProduct(submitData, productForm.images)

    if (response.status === 200) {
      ElMessage.success('商品发布成功')
      router.push('/products')
    } else {
      ElMessage.error(response.msg || '商品发布失败')
    }
  } catch (error) {
    console.error('Product creation validation failed:', error)
  } finally {
    submitLoading.value = false
  }
}

const cancel = () => {
  router.back()
}

const loadCategories = async () => {
  try {
    const response = await getCategories()
    if (response.status === 200) {
      categories.value = response.data
    }
  } catch (error: any) {
    ElMessage.error('获取分类失败')
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.product-create {
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
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.create-content {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  color: #333;
  text-align: center;
  margin-bottom: 2rem;
}

.product-form {
  max-width: 600px;
  margin: 0 auto;
}

.form-tip {
  color: #999;
  font-size: 0.9rem;
  margin-left: 1rem;
}

.upload-demo {
  width: 100%;
}

.upload-tip {
  color: #999;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

:deep(.el-upload--picture-card) {
  width: 100px;
  height: 100px;
}

:deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 100px;
  height: 100px;
}
</style>
