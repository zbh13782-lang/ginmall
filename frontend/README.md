# Gin Mall 前端项目

基于 Vue 3 + TypeScript + Element Plus 构建的商城前端应用。

## 功能特性

- ✅ 用户注册和登录
- ✅ 用户信息修改和邮箱绑定
- ✅ 商品浏览、搜索和详情查看
- ✅ Boss用户发布和管理商品
- ✅ 响应式设计，支持移动端
- ✅ 现代化UI界面

## 技术栈

- **Vue 3** - 渐进式JavaScript框架
- **TypeScript** - 类型安全的JavaScript
- **Vue Router** - 官方路由管理器
- **Pinia** - 状态管理库
- **Element Plus** - Vue 3 UI组件库
- **Axios** - HTTP客户端
- **Vite** - 构建工具

## 快速开始

### 环境要求

- Node.js 18+
- npm 或 yarn

### 安装依赖

```bash
cd frontend
npm install
```

### 开发环境运行

```bash
npm run dev
```

项目将在 `http://localhost:5173` 启动。

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
frontend/
├── src/
│   ├── components/     # 公共组件
│   ├── views/         # 页面组件
│   ├── router/        # 路由配置
│   ├── stores/        # Pinia状态管理
│   ├── services/      # API服务
│   ├── types/         # TypeScript类型定义
│   ├── utils/         # 工具函数
│   ├── App.vue        # 根组件
│   └── main.ts        # 入口文件
├── public/            # 静态资源
├── index.html         # HTML模板
├── vite.config.ts     # Vite配置
├── tsconfig.json      # TypeScript配置
└── package.json       # 项目配置
```

## 主要页面

- `/` - 首页（商品列表）
- `/login` - 用户登录
- `/register` - 用户注册
- `/profile` - 用户个人资料
- `/products` - 商品列表页
- `/product/:id` - 商品详情页
- `/product/create` - 发布商品
- `/product/edit/:id` - 编辑商品

## API集成

项目通过代理配置连接后端API：

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:3000', // 后端服务器地址
      changeOrigin: true
    }
  }
}
```

## 开发说明

1. 确保后端服务正在运行（默认端口3000）
2. 前端开发服务器将在5173端口启动
3. 修改后端地址需要在 `vite.config.ts` 中更新代理配置

## 注意事项

- 注册用户时需要提供16位密钥（用于密码加密）
- 商品图片上传支持 JPG/PNG 格式，单张不超过5MB
- Boss用户可以发布和管理自己的商品
- 用户可以绑定邮箱用于增强账户安全性
