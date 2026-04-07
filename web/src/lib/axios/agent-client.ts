import axios, { type AxiosError } from 'axios'
import { useAuthStore } from '@/stores/auth-store'

export const agentClient = axios.create({
  baseURL: '/agent-portal/v1',
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor: attach Bearer token
agentClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().agentToken
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor: no refresh (portal uses single token)
agentClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().agentClearAuth()
      window.location.href = '/agent/login'
    }
    return Promise.reject(error)
  }
)
