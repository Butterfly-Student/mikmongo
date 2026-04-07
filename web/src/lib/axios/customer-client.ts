import axios, { type AxiosError } from 'axios'
import { useAuthStore } from '@/stores/auth-store'

export const customerClient = axios.create({
  baseURL: '/portal/v1',
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor: attach Bearer token
customerClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().customerToken
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor: no refresh (portal uses single token)
customerClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().customerClearAuth()
      window.location.href = '/customer/login'
    }
    return Promise.reject(error)
  }
)
