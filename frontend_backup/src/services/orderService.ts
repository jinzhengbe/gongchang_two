import { ref } from 'vue'
import { format } from 'date-fns'
import { Order, OrderStatus, OrderFilters } from '@/types/order'
import api from './api'

interface OrderStatistics {
  totalOrders: number
  activeOrders: number
  completedOrders: number
  statusCounts: Record<string, number>
  trendData: Array<{
    date: string
    count: number
  }>
}

interface RecentOrder {
  id: number
  orderNumber: string
  productName: string
  status: string
  createdAt: string
}

// 订单服务类
export class OrderService {
  private orders = ref<Order[]>([])
  private filters = ref<OrderFilters>({})
  private baseUrl = '/api'

  // 获取所有订单状态
  getStatuses(): OrderStatus[] {
    return ['待处理', '进行中', '已完成']
  }

  // 获取状态样式类
  getStatusClass(status: OrderStatus): string {
    switch (status) {
      case '待处理':
        return 'pending'
      case '进行中':
        return 'processing'
      case '已完成':
        return 'completed'
      default:
        return ''
    }
  }

  // 格式化日期
  formatDate(date: string): string {
    return format(new Date(date), 'yyyy-MM-dd HH:mm')
  }

  // 设置搜索查询
  setSearchQuery(query: string) {
    this.filters.value.searchQuery = query
  }

  // 设置状态过滤
  setStatusFilter(status: OrderStatus | '') {
    this.filters.value.status = status || undefined
  }

  setDateRange(startDate: string, endDate: string) {
    this.filters.value.startDate = startDate
    this.filters.value.endDate = endDate
  }

  // 获取过滤后的订单列表
  getFilteredOrders(): Order[] {
    let filtered = [...this.orders.value]

    if (this.filters.value.searchQuery) {
      const query = this.filters.value.searchQuery.toLowerCase()
      filtered = filtered.filter(order => 
        order.orderNumber.toLowerCase().includes(query) ||
        order.customerName.toLowerCase().includes(query) ||
        order.productName.toLowerCase().includes(query)
      )
    }

    if (this.filters.value.status) {
      filtered = filtered.filter(order => order.status === this.filters.value.status)
    }

    if (this.filters.value.startDate && this.filters.value.endDate) {
      filtered = filtered.filter(order => {
        const orderDate = new Date(order.createdAt)
        const start = new Date(this.filters.value.startDate!)
        const end = new Date(this.filters.value.endDate!)
        return orderDate >= start && orderDate <= end
      })
    }

    return filtered
  }

  // 获取订单列表
  async fetchOrders() {
    try {
      const response = await api.get('/orders', {
        params: this.filters.value
      })
      this.orders.value = response.data
    } catch (error) {
      console.error('获取订单列表失败:', error)
      throw error
    }
  }

  async getOrderById(id: string): Promise<Order> {
    try {
      const response = await api.get(`/orders/${id}`)
      return response.data
    } catch (error) {
      console.error('获取订单详情失败:', error)
      throw error
    }
  }

  async createOrder(data: FormData): Promise<Order> {
    try {
      const response = await api.post('/orders', data)
      return response.data
    } catch (error) {
      console.error('创建订单失败:', error)
      throw error
    }
  }

  async updateOrderStatus(id: string, status: OrderStatus): Promise<Order> {
    try {
      const response = await api.put(`/orders/${id}/status`, { status })
      return response.data
    } catch (error) {
      console.error('更新订单状态失败:', error)
      throw error
    }
  }

  async updateOrderNotes(id: string, notes: string): Promise<Order> {
    try {
      const response = await api.put(`/orders/${id}/notes`, { notes })
      return response.data
    } catch (error) {
      console.error('更新订单备注失败:', error)
      throw error
    }
  }

  async getOrderStatistics(): Promise<OrderStatistics> {
    const response = await api.get('/orders/statistics')
    return response.data
  }

  async getRecentOrders(limit: number = 5): Promise<RecentOrder[]> {
    const response = await api.get('/orders/recent', {
      params: { limit }
    })
    return response.data
  }

  async searchOrders(query: string): Promise<RecentOrder[]> {
    const response = await api.get('/orders/search', {
      params: { query }
    })
    return response.data
  }
}

export const orderService = new OrderService()

export { OrderStatus } 