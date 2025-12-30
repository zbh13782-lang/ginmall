import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(ElementPlus)

// 初始化用户状态
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore(pinia)
authStore.initUser()

app.mount('#app')
