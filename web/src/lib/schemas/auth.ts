import { z } from 'zod'

// Generic API response wrapper
export const ApiResponseSchema = <T extends z.ZodTypeAny>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: dataSchema,
  })

// Error response
export const ErrorResponseSchema = z.object({
  success: z.literal(false),
  error: z.string(),
})

// Admin login
export const AdminUserSchema = z.object({
  id: z.string(),
  full_name: z.string(),
  email: z.string().email(),
  phone: z.string(),
  role: z.enum(['superadmin', 'admin', 'cs', 'billing', 'technician', 'readonly']),
  is_active: z.boolean(),
  last_login: z.string().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
})

export const AdminLoginDataSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string(),
  user: AdminUserSchema,
})

export const AdminLoginResponseSchema = ApiResponseSchema(AdminLoginDataSchema)

// Admin refresh -- NOTE: returns `token` NOT `access_token`
export const AdminRefreshDataSchema = z.object({
  token: z.string(),
  refresh_token: z.string(),
})

export const AdminRefreshResponseSchema = ApiResponseSchema(AdminRefreshDataSchema)

// Admin form schemas
export const LoginFormSchema = z.object({
  email: z.string().min(1, 'Email wajib diisi').email('Email tidak valid'),
  password: z.string().min(1, 'Password wajib diisi'),
})

// Change password
export const ChangePasswordSchema = z.object({
  old_password: z.string().min(1, 'Password lama wajib diisi'),
  new_password: z.string().min(8, 'Password baru minimal 8 karakter'),
  confirm_password: z.string().min(1, 'Konfirmasi password wajib diisi'),
}).refine((data) => data.new_password === data.confirm_password, {
  message: 'Password baru dan konfirmasi tidak cocok',
  path: ['confirm_password'],
})

// Message response
export const MessageResponseSchema = z.object({
  success: z.boolean(),
  data: z.object({ message: z.string() }),
})

// Customer portal login
export const CustomerUserSchema = z.object({
  id: z.string(),
  customer_code: z.string(),
  full_name: z.string(),
  username: z.string(),
  email: z.string().email(),
  phone: z.string(),
})

export const PortalLoginDataSchema = z.object({
  token: z.string(),
  customer: CustomerUserSchema,
})

export const PortalLoginResponseSchema = ApiResponseSchema(PortalLoginDataSchema)

export const PortalLoginFormSchema = z.object({
  identifier: z.string().min(1, 'Identifier wajib diisi'),
  password: z.string().min(1, 'Password wajib diisi'),
})

// Agent portal login
export const AgentUserSchema = z.object({
  id: z.string(),
  name: z.string(),
  phone: z.string(),
  username: z.string(),
  status: z.enum(['active', 'inactive']),
})

export const AgentLoginDataSchema = z.object({
  token: z.string(),
  agent: AgentUserSchema,
})

export const AgentLoginResponseSchema = ApiResponseSchema(AgentLoginDataSchema)

export const AgentLoginFormSchema = z.object({
  username: z.string().min(1, 'Username wajib diisi'),
  password: z.string().min(1, 'Password wajib diisi'),
})
