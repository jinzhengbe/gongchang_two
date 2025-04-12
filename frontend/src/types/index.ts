// 轮播图项目类型
export interface CarouselItem {
  id: number
  title: string
  description: string
  image: string
  type: 'order' | 'factory' | 'fabric'
}

// 订单类型
export interface Order {
  id: number
  title: string
  status: 'PENDING' | 'PROCESSING' | 'COMPLETED' | 'CANCELLED'
  createTime: string
  description?: string
}

// 工厂类型
export interface Factory {
  id: number
  name: string
  type: 'premium' | 'standard'
  monthlyCapacity: number
  description?: string
  location?: string
  rating?: number
}

// 面料类型
export interface Fabric {
  id: number
  name: string
  type: 'cotton' | 'silk' | 'wool' | 'synthetic'
  description?: string
  price?: number
  composition?: string
  weight?: number
  width?: number
  minimumOrder?: number
}

// 通用响应类型
export interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

// 分页响应类型
export interface PaginatedResponse<T> {
  total: number
  page: number
  pageSize: number
  items: T[]
}

// 用户类型
export interface User {
  id: number
  username: string
  role: 'designer' | 'factory' | 'supplier'
  avatar?: string
  email?: string
  phone?: string
  company?: string
}

// 错误类型
export interface ApiError {
  code: number
  message: string
  details?: string
} 