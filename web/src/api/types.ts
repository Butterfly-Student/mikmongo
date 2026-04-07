import type { z } from 'zod'
import type {
  AdminLoginDataSchema,
  AdminRefreshDataSchema,
  PortalLoginDataSchema,
  AgentLoginDataSchema,
  AdminUserSchema,
  CustomerUserSchema,
  AgentUserSchema,
  LoginFormSchema,
  PortalLoginFormSchema,
  AgentLoginFormSchema,
  ChangePasswordSchema,
} from '@/lib/schemas/auth'
import type {
  RouterResponseSchema,
  RouterListResponseSchema,
} from '@/lib/schemas/router'
import type {
  UserResponseSchema,
  CreateUserRequestSchema,
  CreateUserFormSchema,
} from '@/lib/schemas/user'
import type { ReportSummarySchema } from '@/lib/schemas/report'

export type AdminUser = z.infer<typeof AdminUserSchema>
export type CustomerUser = z.infer<typeof CustomerUserSchema>
export type AgentUser = z.infer<typeof AgentUserSchema>

export type AdminLoginData = z.infer<typeof AdminLoginDataSchema>
export type AdminRefreshData = z.infer<typeof AdminRefreshDataSchema>
export type PortalLoginData = z.infer<typeof PortalLoginDataSchema>
export type AgentLoginData = z.infer<typeof AgentLoginDataSchema>

export type LoginFormValues = z.infer<typeof LoginFormSchema>
export type PortalLoginFormValues = z.infer<typeof PortalLoginFormSchema>
export type AgentLoginFormValues = z.infer<typeof AgentLoginFormSchema>
export type ChangePasswordValues = z.infer<typeof ChangePasswordSchema>

export type ApiResponse<T> = {
  success: boolean
  data: T
}

export type Meta = {
  total: number
  limit: number
  offset: number
}

// Router types
export type RouterResponse = z.infer<typeof RouterResponseSchema>
export type RouterListData = z.infer<typeof RouterListResponseSchema>

// User types
export type UserResponse = z.infer<typeof UserResponseSchema>
export type CreateUserRequest = z.infer<typeof CreateUserRequestSchema>
export type CreateUserFormValues = z.infer<typeof CreateUserFormSchema>

// Report types
export type ReportSummary = z.infer<typeof ReportSummarySchema>
