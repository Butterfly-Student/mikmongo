import axios, { AxiosError, InternalAxiosRequestConfig } from "axios"
import {
  getAdminAccessToken,
  getAdminRefreshToken,
  adminAuthActions,
} from "@/store"

const adminClient = axios.create({
  baseURL: "/api/v1",
  headers: { "Content-Type": "application/json" },
  timeout: 15000,
})

let isRefreshing = false
let failedQueue: Array<{
  resolve: (token: string) => void
  reject: (err: AxiosError) => void
}> = []

const processQueue = (error: AxiosError | null, token: string | null): void => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error)
    } else {
      resolve(token!)
    }
  })
  failedQueue = []
}

// Attach Bearer token on every request
adminClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAdminAccessToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Handle 401: refresh token, retry, or redirect to login
adminClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean
    }

    if (error.response?.status !== 401 || original._retry) {
      return Promise.reject(error)
    }

    // Another request is already refreshing — queue this one
    if (isRefreshing) {
      return new Promise<string>((resolve, reject) => {
        failedQueue.push({ resolve, reject })
      })
        .then((token) => {
          original.headers.Authorization = `Bearer ${token}`
          return adminClient(original)
        })
        .catch((err) => Promise.reject(err))
    }

    original._retry = true
    isRefreshing = true

    try {
      const refreshToken = getAdminRefreshToken()
      if (!refreshToken) {
        throw new Error("No refresh token available")
      }
      // CRITICAL: use plain axios — NOT adminClient — to avoid interceptor loop
      const { data } = await axios.post("/api/v1/auth/refresh", {
        refresh_token: refreshToken,
      })
      const newAccessToken: string = data.data.access_token
      const newRefreshToken: string = data.data.refresh_token
      adminAuthActions.setTokens(newAccessToken, newRefreshToken)
      processQueue(null, newAccessToken)
      original.headers.Authorization = `Bearer ${newAccessToken}`
      return adminClient(original)
    } catch (refreshError) {
      processQueue(refreshError as AxiosError, null)
      adminAuthActions.clearAuth()
      // Router not available at module scope — use location redirect as fallback
      window.location.href = "/login"
      return Promise.reject(refreshError)
    } finally {
      isRefreshing = false
    }
  }
)

export { adminClient }
