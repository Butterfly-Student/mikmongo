import { z } from 'zod'

export const pppSecretInfoSchema = z.object({
    service: z.string().nullable(),
    profile: z.string().nullable(),
    local_address: z.string().nullable(),
    remote_address: z.string().nullable(),
})

export const subscriptionSchema = z.object({
    id: z.string().uuid(),
    customer_id: z.string().uuid(),
    plan_id: z.string().uuid().nullable(),
    router_id: z.string().uuid(),
    username: z.string(),
    static_ip: z.string().nullable(),
    gateway: z.string().nullable(),
    status: z.enum(['pending', 'active', 'suspended', 'isolated', 'expired', 'terminated']),
    activated_at: z.string().datetime().nullable(),
    expiry_date: z.string().nullable(),
    billing_day: z.number().int().nullable(),
    auto_isolate: z.boolean().nullable(),
    grace_period_days: z.number().int().nullable(),
    suspend_reason: z.string().nullable(),
    notes: z.string().nullable(),
    created_at: z.string().datetime(),
    updated_at: z.string().datetime(),
    mikrotik: pppSecretInfoSchema.nullable(),
})

export type Subscription = z.infer<typeof subscriptionSchema>

export const createSubscriptionSchema = z.object({
    customer_id: z.string().uuid('Customer is required'),
    plan_id: z.string().uuid('Bandwidth profile is required'),
    username: z.string().min(1, 'Username is required'),
    password: z.string().optional(),
    static_ip: z.string().optional(),
    gateway: z.string().optional(),
    billing_day: z.coerce.number().int().optional(),
    auto_isolate: z.boolean().default(true),
    grace_period_days: z.coerce.number().int().optional(),
    notes: z.string().optional(),
})

export type CreateSubscription = z.infer<typeof createSubscriptionSchema>
