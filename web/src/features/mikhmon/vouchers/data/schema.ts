import { z } from 'zod'

export const generateVoucherSchema = z.object({
    quantity: z.coerce.number().int().min(1, 'Min 1').max(500, 'Max 500'),
    profile: z.string().min(1, 'Profile is required'),
    mode: z.enum(['vc', 'up']),
    char_set: z.string().min(1),
    name_length: z.coerce.number().int().min(3).max(12).optional(),
    prefix: z.string().optional(),
    server: z.string().optional(),
    time_limit: z.string().optional(),
    data_limit: z.string().optional(),
    comment: z.string().optional(),
})

export type GenerateVoucherForm = z.infer<typeof generateVoucherSchema>
