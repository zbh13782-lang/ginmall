import api from './api'

// 商品相关接口
export interface Product {
  id: number
  name: string
  category_id: number
  title: string
  info: string
  img_path: string
  price: string
  discount_price: string
  view: number
  created_at: number
  num: number
  on_sale: boolean
  boss_id: number
  boss_name: string
  boss_avatar: string
}

export interface Category {
  id: number
  name: string
}

export interface ProductListResponse {
  status: number
  data: {
    item: Product[]
    total: number
  }
  msg: string
}

export interface CreateProductData {
  name: string
  category_id: number
  title: string
  info: string
  price: string
  discount_price: string
  num: number
  on_sale: boolean
}

// 获取商品列表
export const getProducts = (params?: any) => {
  return api.get<ProductListResponse>('/products', { params })
}

// 获取商品详情
export const getProductDetail = (id: string) => {
  return api.get<{ status: number; data: Product; msg: string }>(`/product/${id}`)
}

// 搜索商品
export const searchProducts = (data: { info: string }) => {
  return api.post<ProductListResponse>('/products', data)
}

// 创建商品
export const createProduct = (data: CreateProductData, files: File[]) => {
  const formData = new FormData()

  // 添加商品数据
  Object.keys(data).forEach(key => {
    formData.append(key, String(data[key as keyof CreateProductData]))
  })

  // 添加文件
  files.forEach(file => {
    formData.append('file', file)
  })

  return api.post('/product', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 更新商品
export const updateProduct = (id: string, data: Partial<CreateProductData>) => {
  return api.put(`/product/${id}`, data)
}

// 删除商品
export const deleteProduct = (id: string) => {
  return api.delete(`/product/${id}`)
}

// 获取商品图片
export const getProductImages = (id: string) => {
  return api.get(`/imgs/${id}`)
}

// 获取分类列表
export const getCategories = () => {
  return api.get<{ status: number; data: Category[]; msg: string }>('/categories')
}
