import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'
import { MetaSchema } from '@/lib/schemas/router'

export const CustomerResponseSchema = z.object({
    id: z.string().uuid(),
    customer_code: z.string(),
    full_name: z.string(),
    username: z.string().nullish(),
    email: z.string().nullish(),
    phone: z.string(),
    id_card_number: z.string().nullish(),
    address: z.string().nullish(),
    latitude: z.number().nullish(),
    longitude: z.number().nullish(),
    is_active: z.boolean(),
    notes: z.string().nullish(),
    tags: z.array(z.string()).nullish().transform((v) => v ?? []),
    created_at: z.string(),
    updated_at: z.string(),
})

export const CustomerListResponseSchema = ApiResponseSchema(
    z.array(CustomerResponseSchema)
).extend({
    meta: MetaSchema,
})

export const CustomerDetailResponseSchema = ApiResponseSchema(CustomerResponseSchema)

export type CustomerResponse = z.infer<typeof CustomerResponseSchema>
export type CustomerListData = z.infer<typeof CustomerListResponseSchema>

// ── Registration ──

export const RegistrationResponseSchema = z.object({
    id: z.string().uuid(),
    full_name: z.string(),
    email: z.string().nullish(),
    phone: z.string(),
    address: z.string().nullish(),
    latitude: z.number().nullish(),
    longitude: z.number().nullish(),
    notes: z.string().nullish(),
    bandwidth_profile_id: z.string().uuid().nullish(),
    status: z.enum(['pending', 'approved', 'rejected']),
    rejection_reason: z.string().nullish(),
    approved_by: z.string().uuid().nullish(),
    approved_at: z.string().nullish(),
    customer_id: z.string().uuid().nullish(),
    created_at: z.string(),
    updated_at: z.string(),
})

export const RegistrationListResponseSchema = ApiResponseSchema(
    z.array(RegistrationResponseSchema)
).extend({
    meta: MetaSchema,
})

export const RegistrationDetailResponseSchema = ApiResponseSchema(RegistrationResponseSchema)

export type RegistrationResponse = z.infer<typeof RegistrationResponseSchema>
export type RegistrationListData = z.infer<typeof RegistrationListResponseSchema>
