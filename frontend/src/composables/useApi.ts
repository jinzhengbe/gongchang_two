import axios from 'axios'
import { ref } from 'vue'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export function useApi() {
  const loading = ref(false)
  const error = ref<string | null>(null)

  const get = async <T>(url: string, params?: any) => {
    try {
      loading.value = true
      error.value = null
      const response = await api.get<T>(url, { params })
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const post = async <T>(url: string, data?: any) => {
    try {
      loading.value = true
      error.value = null
      const response = await api.post<T>(url, data)
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const put = async <T>(url: string, data?: any) => {
    try {
      loading.value = true
      error.value = null
      const response = await api.put<T>(url, data)
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  const del = async <T>(url: string) => {
    try {
      loading.value = true
      error.value = null
      const response = await api.delete<T>(url)
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    get,
    post,
    put,
    delete: del
  }
} 