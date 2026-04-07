import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'
import { MetaSchema } from '@/lib/schemas/router'

export const UserResponseSchema = z.object({
    id: z.string(),
    full_name: z.string(),
    email: z.string().email(),
    phone: z.string().optional(),
    role: z.enum(['superadmin', 'admin', 'cs', 'billing', 'technician', 'readonly']),
    is_active: z.boolean(),
    last_login: z.string().nullish(),
    created_at: z.string(),
    updated_at: z.string(),
})

export const CreateUserRequestSchema = z.object({
    full_name: z.string().min(1, 'Full name is required'),
    email: z
        .string()
        .min(1, 'Email is required')
        .email('Please enter a valid email address'),
    phone: z.string().optional(),
    role: z.enum(['superadmin', 'admin', 'cs', 'billing', 'technician', 'readonly'], {
        required_error: 'Please select a role',
    }),
    password: z.string().min(8, 'Password must be at least 8 characters'),
})

export const CreateUserFormSchema = CreateUserRequestSchema.extend({
    confirm_password: z.string().min(1, 'Please confirm your password'),
}).refine((data) => data.password === data.confirm_password, {
    message: 'Passwords do not match',
    path: ['confirm_password'],
})

export const UserListResponseSchema = ApiResponseSchema(z.array(UserResponseSchema)).extend({
    meta: MetaSchema,
})

export type UserResponse = z.infer<typeof UserResponseSchema>
export type CreateUserRequest = z.infer<typeof CreateUserRequestSchema>
export type CreateUserFormValues = z.infer<typeof CreateUserFormSchema>
