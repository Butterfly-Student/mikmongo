import { z } from "zod"

// ─── Generic API envelope ────────────────────────────────────────────────────
export const ApiResponseSchema = <T extends z.ZodTypeAny>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: dataSchema,
    error: z.string().optional().nullable(),
    meta: z
      .object({
        total: z.number(),
        limit: z.number(),
        offset: z.number(),
      })
      .optional()
      .nullable(),
  })

// ─── Admin login / refresh ────────────────────────────────────────────────────
export const AdminUserSchema = z.object({
  id: z.string(),
  email: z.string().email(),
  role: z.enum(["superadmin", "admin", "teknisi"]),
  full_name: z.string(),
})

export const TokenDataSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string(),
  expires_in: z.number(),
})

export const AdminLoginResponseSchema = ApiResponseSchema(
  TokenDataSchema.extend({
    user: AdminUserSchema,
  })
)

export const AdminRefreshResponseSchema = ApiResponseSchema(TokenDataSchema)

// ─── Agent login / refresh ────────────────────────────────────────────────────
export const AgentUserSchema = z.object({
  id: z.string(),
  email: z.string().email(),
  full_name: z.string(),
})

export const AgentLoginResponseSchema = ApiResponseSchema(
  TokenDataSchema.extend({
    user: AgentUserSchema,
  })
)

// ─── Customer login / refresh ─────────────────────────────────────────────────
export const CustomerUserSchema = z.object({
  id: z.string(),
  email: z.string().email(),
  full_name: z.string(),
})

export const CustomerLoginResponseSchema = ApiResponseSchema(
  TokenDataSchema.extend({
    user: CustomerUserSchema,
  })
)

// ─── Login form input ─────────────────────────────────────────────────────────
export const LoginFormSchema = z.object({
  email: z.string().email("Email tidak valid"),
  password: z.string().min(6, "Password minimal 6 karakter"),
})

export type LoginFormValues = z.infer<typeof LoginFormSchema>
export type AdminLoginResponse = z.infer<typeof AdminLoginResponseSchema>
export type AgentLoginResponse = z.infer<typeof AgentLoginResponseSchema>
export type CustomerLoginResponse = z.infer<typeof CustomerLoginResponseSchema>
