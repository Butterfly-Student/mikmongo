import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'

// ── Vouchers ──────────────────────────────────────────────────────────

export const VoucherResponseSchema = z.object({
  id: z.string(),
  name: z.string(),
  password: z.string().nullish(),
  profile: z.string().nullish(),
  server: z.string().nullish(),
  comment: z.string().nullish(),
  code: z.string().nullish(),
  mode: z.string().nullish(),
  date: z.string().nullish(),
}).passthrough()

export type VoucherResponse = z.infer<typeof VoucherResponseSchema>

export const VoucherBatchResponseSchema = z.object({
  code: z.string().nullish(),
  quantity: z.number().int().nullish(),
  profile: z.string().nullish(),
  server: z.string().nullish(),
  time_limit: z.string().nullish(),
  data_limit: z.string().nullish(),
  vouchers: z.array(VoucherResponseSchema),
})

export type VoucherBatchResponse = z.infer<typeof VoucherBatchResponseSchema>

export const GenerateVoucherRequestSchema = z.object({
  quantity: z.number().int().min(1).max(1000),
  profile: z.string().min(1),
  mode: z.enum(['vc', 'up']),
  char_set: z.string().min(1),
  server: z.string().optional(),
  name_length: z.number().int().min(3).max(12).optional(),
  prefix: z.string().optional(),
  time_limit: z.string().optional(),
  data_limit: z.string().optional(),
  comment: z.string().optional(),
})

export type GenerateVoucherRequest = z.infer<typeof GenerateVoucherRequestSchema>

// ── Mikhmon Profiles ──────────────────────────────────────────────────

export const MikhmonProfileResponseSchema = z.object({
  id: z.string().optional(),
  name: z.string(),
  address_pool: z.string().nullish(),
  rate_limit: z.string().nullish(),
  shared_users: z.number().nullish(),
  parent_queue: z.string().nullish(),
  price: z.number().nullish(),
  selling_price: z.number().nullish(),
  validity: z.string().nullish(),
  expire_mode: z.string().nullish(),
  lock_user: z.boolean().nullish(),
  lock_server: z.boolean().nullish(),
  on_login_script: z.string().nullish(),
  created_at: z.string().nullish(),
  updated_at: z.string().nullish(),
}).passthrough()

export type MikhmonProfileResponse = z.infer<typeof MikhmonProfileResponseSchema>

export const CreateMikhmonProfileRequestSchema = z.object({
  name: z.string().min(1),
  address_pool: z.string().optional(),
  rate_limit: z.string().optional(),
  shared_users: z.number().int().optional(),
  parent_queue: z.string().optional(),
  config: z.object({
    name: z.string().optional(),
    address_pool: z.string().optional(),
    rate_limit: z.string().optional(),
    shared_users: z.number().int().optional(),
    parent_queue: z.string().optional(),
    price: z.number().int().optional(),
    selling_price: z.number().int().optional(),
    validity: z.string().optional(),
    expire_mode: z.string().optional(),
    lock_user: z.boolean().optional(),
    lock_server: z.boolean().optional(),
    on_login_script: z.string().optional(),
  }).optional(),
})

export type CreateMikhmonProfileRequest = z.infer<typeof CreateMikhmonProfileRequestSchema>

export const UpdateMikhmonProfileRequestSchema = z.object({
  name: z.string().optional(),
  address_pool: z.string().optional(),
  rate_limit: z.string().optional(),
  shared_users: z.number().int().optional(),
  parent_queue: z.string().optional(),
  price: z.number().int().optional(),
  selling_price: z.number().int().optional(),
  validity: z.string().optional(),
  expire_mode: z.string().optional(),
  lock_user: z.boolean().optional(),
  lock_server: z.boolean().optional(),
})

export type UpdateMikhmonProfileRequest = z.infer<typeof UpdateMikhmonProfileRequestSchema>

export const GenerateScriptRequestSchema = z.object({
  mode: z.string().min(1),
  profile_name: z.string().min(1),
  price: z.number().int().optional(),
  validity: z.string().optional(),
  selling_price: z.number().int().optional(),
  no_exp: z.boolean().optional(),
  lock_user: z.string().optional(),
  lock_server: z.string().optional(),
})

export type GenerateScriptRequest = z.infer<typeof GenerateScriptRequestSchema>

export const ScriptResponseSchema = z.object({
  name: z.string().nullish(),
  content: z.string().nullish(),
})

export type ScriptResponse = z.infer<typeof ScriptResponseSchema>

// ── Mikhmon Reports ───────────────────────────────────────────────────

export const MikhmonReportResponseSchema = z.object({
  id: z.string().optional(),
  user: z.string().nullish(),
  price: z.number().nullish(),
  ip: z.string().nullish(),
  mac: z.string().nullish(),
  validity: z.string().nullish(),
  profile: z.string().nullish(),
  comment: z.string().nullish(),
  created_at: z.string().nullish(),
}).passthrough()

export type MikhmonReportResponse = z.infer<typeof MikhmonReportResponseSchema>

export const CreateReportRequestSchema = z.object({
  user: z.string().min(1),
  price: z.number().int().optional(),
  ip: z.string().optional(),
  mac: z.string().optional(),
  validity: z.string().optional(),
  profile: z.string().optional(),
  comment: z.string().optional(),
})

export type CreateReportRequest = z.infer<typeof CreateReportRequestSchema>

export const MikhmonReportSummarySchema = z.object({
  totalCount: z.number().int().nullish(),
  totalSales: z.number().nullish(),
  totalRevenue: z.number().nullish(),
}).passthrough()

export type MikhmonReportSummary = z.infer<typeof MikhmonReportSummarySchema>

// ── Expiration ────────────────────────────────────────────────────────

export const ExpireStatusResponseSchema = z.object({
  enabled: z.boolean().nullish(),
  last_run: z.string().nullish(),
  next_run: z.string().nullish(),
  user_count: z.number().int().nullish(),
}).passthrough()

export type ExpireStatusResponse = z.infer<typeof ExpireStatusResponseSchema>

// ── Response wrappers ─────────────────────────────────────────────────

export const VoucherBatchApiResponseSchema = ApiResponseSchema(VoucherBatchResponseSchema)
export const VoucherListResponseSchema = ApiResponseSchema(z.array(VoucherResponseSchema))
export const VoucherDetailResponseSchema = ApiResponseSchema(VoucherResponseSchema)

export const MikhmonProfileListResponseSchema = ApiResponseSchema(z.array(MikhmonProfileResponseSchema))
export const MikhmonProfileDetailResponseSchema = ApiResponseSchema(MikhmonProfileResponseSchema)
export const ScriptApiResponseSchema = ApiResponseSchema(ScriptResponseSchema)

export const MikhmonReportListResponseSchema = ApiResponseSchema(z.array(MikhmonReportResponseSchema))
export const MikhmonReportDetailResponseSchema = ApiResponseSchema(MikhmonReportResponseSchema)
export const MikhmonReportSummaryResponseSchema = ApiResponseSchema(MikhmonReportSummarySchema)

export const ExpireStatusApiResponseSchema = ApiResponseSchema(ExpireStatusResponseSchema)
