import axios, { AxiosError, InternalAxiosRequestConfig } from "axios"
import {
  getCustomerAccessToken,
  getCustomerRefreshToken,
  customerAuthActions,
} from "@/store"

const customerClient = axios.create({
  baseURL: "/portal/v1",
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
    if (error) reject(error)
    else resolve(token!)
  })
  failedQueue = []
}

customerClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getCustomerAccessToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

customerClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean
    }

    if (error.response?.status !== 401 || original._retry) {
      return Promise.reject(error)
    }

    if (isRefreshing) {
      return new Promise<string>((resolve, reject) => {
        failedQueue.push({ resolve, reject })
      })
        .then((token) => {
          original.headers.Authorization = `Bearer ${token}`
          return customerClient(original)
        })
        .catch((err) => Promise.reject(err))
    }

    original._retry = true
    isRefreshing = true

    try {
      const refreshToken = getCustomerRefreshToken()
      if (!refreshToken) {
        throw new Error("No refresh token available")
      }
      const { data } = await axios.post("/portal/v1/auth/refresh", {
        refresh_token: refreshToken,
      })
      const newAccessToken: string = data.data.access_token
      const newRefreshToken: string = data.data.refresh_token
      customerAuthActions.setTokens(newAccessToken, newRefreshToken)
      processQueue(null, newAccessToken)
      original.headers.Authorization = `Bearer ${newAccessToken}`
      return customerClient(original)
    } catch (refreshError) {
      processQueue(refreshError as AxiosError, null)
      customerAuthActions.clearAuth()
      window.location.href = "/customer/login"
      return Promise.reject(refreshError)
    } finally {
      isRefreshing = false
    }
  }
)

export { customerClient }
