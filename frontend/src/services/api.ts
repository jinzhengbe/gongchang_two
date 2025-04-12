import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.error('未授权，请重新登录')
          // TODO: 跳转到登录页
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器错误')
          break
        default:
          ElMessage.error(error.response.data?.error || '请求失败')
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export const orderApi = {
  // 获取订单列表
  getOrders: (params?: any) => {
    return api.get('/orders', { params })
  },

  // 获取订单详情
  getOrderById: (id: string) => {
    return api.get(`/orders/${id}`)
  },

  // 创建订单
  createOrder: (data: any) => {
    return api.post('/orders', data)
  },

  // 更新订单状态
  updateOrderStatus: (id: string, status: string) => {
    return api.patch(`/orders/${id}/status`, { status })
  },

  // 更新订单备注
  updateOrderNotes: (id: string, notes: string) => {
    return api.patch(`/orders/${id}/notes`, { notes })
  }
}

export default api 