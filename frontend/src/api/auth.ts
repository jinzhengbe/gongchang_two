import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export interface LoginCredentials {
  username: string
  password: string
}

export interface RegisterData {
  username: string
  password: string
  email: string
}

export interface AuthResponse {
  token: string
  user: {
    id: number
    username: string
    email: string
  }
}

export const login = async (credentials: LoginCredentials): Promise<AuthResponse> => {
  const response = await axios.post(`${API_URL}/auth/login`, credentials)
  return response.data
}

export const register = async (userData: RegisterData): Promise<AuthResponse> => {
  const response = await axios.post(`${API_URL}/auth/register`, userData)
  return response.data
} 