import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'

// 格式化日期
export const formatDate = (date: string | Date, formatStr: string = 'yyyy-MM-dd'): string => {
  try {
    const dateObj = typeof date === 'string' ? new Date(date) : date
    return format(dateObj, formatStr, { locale: zhCN })
  } catch (error) {
    console.error('Date formatting error:', error)
    return ''
  }
}

// 格式化金额
export const formatCurrency = (amount: number): string => {
  try {
    return new Intl.NumberFormat('zh-CN', {
      style: 'currency',
      currency: 'CNY'
    }).format(amount)
  } catch (error) {
    console.error('Currency formatting error:', error)
    return ''
  }
}

// 格式化数字（添加千位分隔符）
export const formatNumber = (num: number): string => {
  try {
    return new Intl.NumberFormat('zh-CN').format(num)
  } catch (error) {
    console.error('Number formatting error:', error)
    return ''
  }
}

// 文件大小格式化
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

// 防抖函数
export const debounce = <T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timeoutId: NodeJS.Timeout

  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}

// 节流函数
export const throttle = <T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let lastCall = 0

  return (...args: Parameters<T>) => {
    const now = Date.now()
    if (now - lastCall >= delay) {
      fn(...args)
      lastCall = now
    }
  }
}

// 深拷贝
export const deepClone = <T>(obj: T): T => {
  if (obj === null || typeof obj !== 'object') {
    return obj
  }

  if (obj instanceof Date) {
    return new Date(obj.getTime()) as any
  }

  if (obj instanceof Array) {
    return obj.map(item => deepClone(item)) as any
  }

  if (obj instanceof Object) {
    return Object.fromEntries(
      Object.entries(obj).map(([key, value]) => [key, deepClone(value)])
    ) as any
  }

  return obj
}

// 生成随机ID
export const generateId = (): string => {
  return Math.random().toString(36).substr(2, 9)
}

// 检查是否为空对象
export const isEmptyObject = (obj: object): boolean => {
  return Object.keys(obj).length === 0
}

// 检查是否为有效的URL
export const isValidUrl = (url: string): boolean => {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

// 获取文件扩展名
export const getFileExtension = (filename: string): string => {
  return filename.slice((filename.lastIndexOf('.') - 1 >>> 0) + 2)
}

// 检查是否为图片文件
export const isImageFile = (filename: string): boolean => {
  const ext = getFileExtension(filename).toLowerCase()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)
}

// 检查是否为移动设备
export const isMobile = (): boolean => {
  return window.innerWidth <= 768
}

// 获取浏览器语言
export const getBrowserLanguage = (): string => {
  return navigator.language || (navigator as any).userLanguage
}

// 平滑滚动到顶部
export const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

// 复制文本到剪贴板
export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch (error) {
    console.error('Failed to copy text:', error)
    return false
  }
}

// 图片加载错误处理
export const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/images/placeholder.jpg'
  img.alt = '图片加载失败'
} 