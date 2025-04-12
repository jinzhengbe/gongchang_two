import { ElMessage } from 'element-plus'
import type { ApiError } from '@/types'

export const handleApiError = (error: ApiError | Error) => {
  if ('code' in error) {
    const apiError = error as ApiError
    switch (apiError.code) {
      case 401:
        ElMessage.error('请先登录')
        // 跳转到登录页
        window.location.href = '/login'
        break
      case 403:
        ElMessage.error('权限不足')
        break
      case 404:
        ElMessage.error('资源不存在')
        break
      case 500:
        ElMessage.error('服务器错误')
        break
      default:
        ElMessage.error(apiError.message || '请求失败')
    }
  } else {
    ElMessage.error(error.message || '网络错误')
  }
}

export const createErrorHandler = (context: string) => {
  return (error: ApiError | Error) => {
    console.error(`Error in ${context}:`, error)
    handleApiError(error)
  }
}

export const isApiError = (error: any): error is ApiError => {
  return 'code' in error && 'message' in error
}

export const extractErrorMessage = (error: any): string => {
  if (isApiError(error)) {
    return error.message
  }
  if (error instanceof Error) {
    return error.message
  }
  return '未知错误'
} 