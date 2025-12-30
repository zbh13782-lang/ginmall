<template>
  <div class="image-test">
    <h1>上海大学二手商城 - 图片显示测试</h1>

    <div class="test-section">
      <h2>测试本地图片</h2>
      <div class="image-container">
        <img src="/static/imgs/avatar/debug_avatar.jpg" alt="本地头像" />
        <p>本地图片: /static/imgs/avatar/debug_avatar.jpg</p>
      </div>
    </div>

    <div class="test-section">
      <h2>测试API返回的图片</h2>
      <div v-for="product in products" :key="product.id" class="product-item">
        <h3>{{ product.name }}</h3>
        <div class="image-container">
          <img :src="product.img_path" :alt="product.name" />
          <p>商品图片: {{ product.img_path }}</p>
        </div>
        <div class="image-container">
          <img :src="product.boss_avatar" :alt="product.boss_name" />
          <p>卖家头像: {{ product.boss_avatar }}</p>
        </div>
      </div>
    </div>

    <div class="test-section">
      <h2>测试七牛云图片</h2>
      <div class="image-container">
        <img src="http://t7b8f3uri.hd-bkt.clouddn.com/FiiADE5apF17Q3EJ0apsJotGu5_x" alt="七牛云图片" />
        <p>七牛云图片: http://t7b8f3uri.hd-bkt.clouddn.com/FiiADE5apF17Q3EJ0apsJotGu5_x</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProducts, type Product } from '@/services/product'

const products = ref<Product[]>([])

const loadProducts = async () => {
  try {
    const response = await getProducts({ pageNum: 1, pageSize: 3 })
    if (response.status === 200) {
      products.value = response.data.item
    }
  } catch (error) {
    console.error('加载商品失败:', error)
  }
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.image-test {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.test-section {
  margin-bottom: 40px;
  border: 1px solid #ddd;
  padding: 20px;
  border-radius: 8px;
}

.test-section h2 {
  color: #333;
  margin-bottom: 20px;
}

.image-container {
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 4px;
}

.image-container img {
  max-width: 200px;
  max-height: 200px;
  border: 1px solid #ccc;
  margin-bottom: 10px;
}

.image-container p {
  font-size: 12px;
  color: #666;
  word-break: break-all;
}

.product-item {
  margin-bottom: 30px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.product-item h3 {
  color: #333;
  margin-bottom: 15px;
}
</style>
