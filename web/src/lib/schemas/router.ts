import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'

export const MetaSchema = z.object({
    total: z.number(),
    limit: z.number(),
    offset: z.number(),
})

export const RouterResponseSchema = z.object({
    id: z.string(),
    name: z.string(),
    address: z.string(),
    area: z.string().optional().nullable(),
    api_port: z.number().optional(),
    rest_port: z.number().optional(),
    username: z.string(),
    use_ssl: z.boolean(),
    is_master: z.boolean(),
    is_active: z.boolean(),
    status: z.enum(['online', 'offline', 'unknown']),
    last_seen_at: z.string().nullish(),
    notes: z.string().nullish(),
    created_at: z.string(),
    updated_at: z.string(),
})

export const RouterListResponseSchema = ApiResponseSchema(z.array(RouterResponseSchema)).extend({
    meta: MetaSchema.optional(),
})

export const SelectedRouterResponseSchema = ApiResponseSchema(RouterResponseSchema)

export type RouterResponse = z.infer<typeof RouterResponseSchema>
export type RouterListData = z.infer<typeof RouterListResponseSchema>
