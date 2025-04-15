export interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface CarouselItem {
  id: number
  title: string
  description: string
  imageUrl: string
}

export interface Factory {
  id: number
  name: string
  description: string
  location: string
  rating: number
  imageUrl: string
}

export interface Fabric {
  id: number
  name: string
  description: string
  price: number
  imageUrl: string
  supplier: string
}

export * from './order' 