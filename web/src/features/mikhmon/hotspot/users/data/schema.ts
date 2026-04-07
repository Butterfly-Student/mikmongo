import { z } from 'zod'

export const createHotspotUserSchema = z.object({
    name: z.string().min(1, 'Username is required'),
    password: z.string().optional(),
    profile: z.string().optional(),
    server: z.string().optional(),
    mac_address: z.string().optional(),
    limit_uptime: z.string().optional(),
    limit_bytes_total: z.string().optional(),
    comment: z.string().optional(),
    disabled: z.boolean().optional(),
})

export type CreateHotspotUserForm = z.infer<typeof createHotspotUserSchema>
