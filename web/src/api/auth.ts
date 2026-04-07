import axios from 'axios'
import { adminClient } from '@/lib/axios/admin-client'
import {
  AdminLoginResponseSchema,
  AdminRefreshResponseSchema,
  PortalLoginResponseSchema,
  AgentLoginResponseSchema,
  MessageResponseSchema,
} from '@/lib/schemas/auth'
import type {
  AdminLoginData,
  AdminRefreshData,
  PortalLoginData,
  AgentLoginData,
} from '@/api/types'

// Admin login -- uses plain axios (no token yet)
export async function adminLogin(email: string, password: string): Promise<AdminLoginData> {
  const response = await axios.post('/api/v1/auth/login', { email, password })
  const parsed = AdminLoginResponseSchema.parse(response.data)
  return parsed.data
}

// Admin refresh -- uses plain axios (to avoid interceptor loop)
export async function adminRefreshToken(refreshToken: string): Promise<AdminRefreshData> {
  const response = await axios.post('/api/v1/auth/refresh', { refresh_token: refreshToken })
  const parsed = AdminRefreshResponseSchema.parse(response.data)
  return parsed.data
}

// Admin change password -- uses adminClient (authenticated)
export async function adminChangePassword(oldPassword: string, newPassword: string) {
  const response = await adminClient.post('/auth/change-password', {
    old_password: oldPassword,
    new_password: newPassword,
  })
  return MessageResponseSchema.parse(response.data)
}

// Admin logout -- uses adminClient (authenticated)
export async function adminLogout() {
  const response = await adminClient.post('/auth/logout')
  return MessageResponseSchema.parse(response.data)
}

// Admin get current user -- uses adminClient (authenticated)
export async function adminGetMe() {
  const response = await adminClient.get('/auth/me')
  return response.data
}

// Customer login -- uses plain axios (no token yet)
export async function customerLogin(identifier: string, password: string): Promise<PortalLoginData> {
  const response = await axios.post('/portal/v1/login', { identifier, password })
  const parsed = PortalLoginResponseSchema.parse(response.data)
  return parsed.data
}

// Agent login -- uses plain axios (no token yet)
export async function agentLogin(username: string, password: string): Promise<AgentLoginData> {
  const response = await axios.post('/agent-portal/v1/login', { username, password })
  const parsed = AgentLoginResponseSchema.parse(response.data)
  return parsed.data
}
