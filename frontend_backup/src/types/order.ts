export enum OrderStatus {
  PENDING = 'PENDING',
  PROCESSING = 'PROCESSING',
  COMPLETED = 'COMPLETED',
  CANCELLED = 'CANCELLED'
}

export interface Order {
  id: string
  title: string
  description: string
  status: OrderStatus
  createdAt: string
  updatedAt: string
  designerId: string
  factoryId?: string
  fabricId?: string
  quantity: number
  price: number
  deadline: string
  files?: File[]
}

export interface File {
  id: string
  name: string
  url: string
  type: string
  size: number
  createdAt: string
}

export interface OrderFilters {
  searchQuery?: string
  status?: OrderStatus
  startDate?: string
  endDate?: string
}

export interface CreateOrderRequest {
  productName: string
  designerId: number
  customerId: number
  productId: number
  notes?: string
}

export interface UpdateOrderStatusRequest {
  status: OrderStatus
}

export interface UpdateOrderNotesRequest {
  notes: string
} 