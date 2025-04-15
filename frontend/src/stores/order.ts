import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useApi } from '@/composables/useApi'
import type { Order, OrderStatus, CreateOrderRequest, UpdateOrderStatusRequest } from '@/types/order'

export const useOrderStore = defineStore('order', () => {
  const api = useApi()
  const orders = ref<Order[]>([])
  const currentOrder = ref<Order | null>(null)
  const totalOrders = ref(0)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 获取订单列表
  const fetchOrders = async (params: {
    page: number
    pageSize: number
    search?: string
    status?: OrderStatus
    startDate?: Date
    endDate?: Date
  }) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/orders', { params })
      orders.value = response.data.items
      totalOrders.value = response.data.total
    } catch (err) {
      error.value = '获取订单列表失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 获取热门订单
  const fetchHotOrders = async (limit: number = 6) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/orders/hot', { params: { limit } })
      return response.data
    } catch (err) {
      error.value = '获取热门订单失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 获取订单详情
  const getOrderById = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get(`/orders/${id}`)
      currentOrder.value = response.data
      return response.data
    } catch (err) {
      error.value = '获取订单详情失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 创建订单
  const createOrder = async (order: CreateOrderRequest) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.post('/orders', order)
      return response.data
    } catch (err) {
      error.value = '创建订单失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 更新订单状态
  const updateOrderStatus = async (id: string, data: UpdateOrderStatusRequest) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.patch(`/orders/${id}/status`, data)
      if (currentOrder.value?.id === id) {
        currentOrder.value = response.data
      }
      return response.data
    } catch (err) {
      error.value = '更新订单状态失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 删除订单
  const deleteOrder = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      await api.delete(`/orders/${id}`)
      orders.value = orders.value.filter(order => order.id !== id)
      if (currentOrder.value?.id === id) {
        currentOrder.value = null
      }
    } catch (err) {
      error.value = '删除订单失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 上传文件
  const uploadFile = async (orderId: string, file: File) => {
    loading.value = true
    error.value = null
    try {
      const formData = new FormData()
      formData.append('file', file)
      const response = await api.post(`/orders/${orderId}/files`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      return response.data
    } catch (err) {
      error.value = '上传文件失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 下载文件
  const downloadFile = async (orderId: string, fileId: string) => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get(`/orders/${orderId}/files/${fileId}`, {
        responseType: 'blob'
      })
      return response.data
    } catch (err) {
      error.value = '下载文件失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 删除文件
  const deleteFile = async (orderId: string, fileId: string) => {
    loading.value = true
    error.value = null
    try {
      await api.delete(`/orders/${orderId}/files/${fileId}`)
    } catch (err) {
      error.value = '删除文件失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    orders,
    currentOrder,
    totalOrders,
    loading,
    error,
    fetchOrders,
    fetchHotOrders,
    getOrderById,
    createOrder,
    updateOrderStatus,
    deleteOrder,
    uploadFile,
    downloadFile,
    deleteFile
  }
}) 