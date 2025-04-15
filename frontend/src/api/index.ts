import axios from 'axios'
import type { ApiResponse, PaginatedResponse, CarouselItem, Order, Factory, Fabric } from '@/types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response) {
      const { status, data } = error.response
      switch (status) {
        case 401:
          // 未授权，跳转到登录页
          window.location.href = '/login'
          break
        case 403:
          // 权限不足
          console.error('权限不足')
          break
        case 500:
          // 服务器错误
          console.error('服务器错误')
          break
        default:
          console.error('请求失败')
      }
      return Promise.reject(data)
    }
    return Promise.reject(error)
  }
)

// 首页API
export const homeApi = {
  // 获取轮播图数据
  getCarouselItems(): Promise<ApiResponse<CarouselItem[]>> {
    return api.get('/home/carousel')
  },

  // 获取新品面料
  getNewFabrics(): Promise<ApiResponse<Fabric[]>> {
    return api.get('/fabrics/new')
  },

  // 获取热门面料
  getHotFabrics(): Promise<ApiResponse<Fabric[]>> {
    return api.get('/fabrics/hot')
  },

  // 获取推荐工厂
  getRecommendedFactories(): Promise<ApiResponse<Factory[]>> {
    return api.get('/factories/recommended')
  },

  // 获取最新订单
  getLatestOrders(): Promise<ApiResponse<Order[]>> {
    return api.get('/orders/latest')
  }
}

// 订单API
export const orderApi = {
  // 获取订单列表
  getOrders(params: { page: number; pageSize: number }): Promise<ApiResponse<PaginatedResponse<Order>>> {
    return api.get('/orders', { params })
  },

  // 获取订单详情
  getOrderDetail(id: number): Promise<ApiResponse<Order>> {
    return api.get(`/orders/${id}`)
  },

  // 创建订单
  createOrder(data: Partial<Order>): Promise<ApiResponse<Order>> {
    return api.post('/orders', data)
  },

  // 更新订单
  updateOrder(id: number, data: Partial<Order>): Promise<ApiResponse<Order>> {
    return api.put(`/orders/${id}`, data)
  },

  // 删除订单
  deleteOrder(id: number): Promise<ApiResponse<void>> {
    return api.delete(`/orders/${id}`)
  }
}

// 工厂API
export const factoryApi = {
  // 获取工厂列表
  getFactories(params: { page: number; pageSize: number }): Promise<ApiResponse<PaginatedResponse<Factory>>> {
    return api.get('/factories', { params })
  },

  // 获取工厂详情
  getFactoryDetail(id: number): Promise<ApiResponse<Factory>> {
    return api.get(`/factories/${id}`)
  }
}

// 面料API
export const fabricApi = {
  // 获取面料列表
  getFabrics(params: { page: number; pageSize: number }): Promise<ApiResponse<PaginatedResponse<Fabric>>> {
    return api.get('/fabrics', { params })
  },

  // 获取面料详情
  getFabricDetail(id: number): Promise<ApiResponse<Fabric>> {
    return api.get(`/fabrics/${id}`)
  }
}

export default api 