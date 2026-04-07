import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'
import { MetaSchema } from '@/lib/schemas/router'

export const SubscriptionResponseSchema = z.object({
    id: z.string().uuid(),
    customer_id: z.string().uuid(),
    plan_id: z.string().uuid().nullish(),
    router_id: z.string().uuid(),
    username: z.string(),
    static_ip: z.string().nullish(),
    gateway: z.string().nullish(),
    status: z.enum(['pending', 'active', 'suspended', 'isolated', 'expired', 'terminated']),
    activated_at: z.string().nullish(),
    expiry_date: z.string().nullish(),
    billing_day: z.number().int().nullish(),
    auto_isolate: z.boolean().nullish(),
    grace_period_days: z.number().int().nullish(),
    suspend_reason: z.string().nullish(),
    notes: z.string().nullish(),
    created_at: z.string(),
    updated_at: z.string(),
    mikrotik: z.object({
        service: z.string().nullable(),
        profile: z.string().nullable(),
        local_address: z.string().nullable(),
        remote_address: z.string().nullable(),
    }).nullable(),
})

export const SubscriptionListResponseSchema = ApiResponseSchema(
    z.array(SubscriptionResponseSchema)
).extend({
    meta: MetaSchema,
})

export const SubscriptionDetailResponseSchema = ApiResponseSchema(SubscriptionResponseSchema)

export type SubscriptionResponse = z.infer<typeof SubscriptionResponseSchema>
export type SubscriptionListData = z.infer<typeof SubscriptionListResponseSchema>
