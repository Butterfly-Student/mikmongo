import { z } from 'zod'
import { ApiResponseSchema, MessageResponseSchema as _msg } from '@/lib/schemas/auth'
import { MetaSchema } from '@/lib/schemas/router'

export const SalesAgentResponseSchema = z.object({
  id: z.string().uuid(),
  router_id: z.string().uuid(),
  name: z.string(),
  phone: z.string().nullish(),
  username: z.string(),
  status: z.enum(['active', 'inactive']),
  voucher_mode: z.enum(['mix', 'num', 'alp']).nullish(),
  voucher_length: z.number().int().nullish(),
  voucher_type: z.enum(['upp', 'up']).nullish(),
  bill_discount: z.number().nullish(),
  billing_cycle: z.enum(['weekly', 'monthly']).nullish(),
  billing_day: z.number().int().nullish(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type SalesAgentResponse = z.infer<typeof SalesAgentResponseSchema>

export const CreateSalesAgentRequestSchema = z.object({
  router_id: z.string().uuid(),
  name: z.string().min(1),
  username: z.string().min(1),
  password: z.string().min(6),
  phone: z.string().optional(),
  status: z.enum(['active', 'inactive']).optional(),
  voucher_mode: z.enum(['mix', 'num', 'alp']).optional(),
  voucher_length: z.number().int().optional(),
  voucher_type: z.enum(['upp', 'up']).optional(),
  bill_discount: z.number().optional(),
  billing_cycle: z.enum(['weekly', 'monthly']).optional(),
  billing_day: z.number().int().min(1).max(31).optional(),
})

export type CreateSalesAgentRequest = z.infer<typeof CreateSalesAgentRequestSchema>

export const UpdateSalesAgentRequestSchema = z.object({
  name: z.string().optional(),
  phone: z.string().optional(),
  password: z.string().optional(),
  status: z.enum(['active', 'inactive']).optional(),
  voucher_mode: z.enum(['mix', 'num', 'alp']).optional(),
  voucher_length: z.number().int().optional(),
  voucher_type: z.enum(['upp', 'up']).optional(),
  bill_discount: z.number().optional(),
  billing_cycle: z.enum(['weekly', 'monthly']).optional(),
  billing_day: z.number().int().min(1).max(31).optional(),
})

export type UpdateSalesAgentRequest = z.infer<typeof UpdateSalesAgentRequestSchema>

export const SalesProfilePriceResponseSchema = z.object({
  id: z.string().uuid(),
  sales_agent_id: z.string().uuid(),
  profile_name: z.string(),
  base_price: z.number(),
  selling_price: z.number(),
  voucher_length: z.number().int().nullish(),
  is_active: z.boolean(),
  created_at: z.string(),
})

export type SalesProfilePriceResponse = z.infer<typeof SalesProfilePriceResponseSchema>

export const UpsertProfilePriceRequestSchema = z.object({
  base_price: z.number().optional(),
  selling_price: z.number().optional(),
  voucher_length: z.number().int().optional(),
  is_active: z.boolean().optional(),
})

export type UpsertProfilePriceRequest = z.infer<typeof UpsertProfilePriceRequestSchema>

// Response wrappers
export const SalesAgentListResponseSchema = ApiResponseSchema(z.array(SalesAgentResponseSchema)).extend({
  meta: MetaSchema.optional(),
})

export const SalesAgentDetailResponseSchema = ApiResponseSchema(SalesAgentResponseSchema)

export const SalesProfilePriceListResponseSchema = ApiResponseSchema(z.array(SalesProfilePriceResponseSchema))

export const SalesProfilePriceDetailResponseSchema = ApiResponseSchema(SalesProfilePriceResponseSchema)

// Generate agent invoice
export const GenerateAgentInvoiceRequestSchema = z.object({
  period_start: z.string(),
  period_end: z.string(),
})

export type GenerateAgentInvoiceRequest = z.infer<typeof GenerateAgentInvoiceRequestSchema>

export const MarkPaidRequestSchema = z.object({
  paid_amount: z.number().min(0),
})
