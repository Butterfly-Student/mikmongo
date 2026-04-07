import { z } from 'zod'

export const customerSchema = z.object({
    id: z.string().uuid(),
    customer_code: z.string().optional(),
    full_name: z.string(),
    username: z.string().optional().nullable(),
    email: z.string().email().optional().nullable(),
    phone: z.string(),
    id_card_number: z.string().optional().nullable(),
    address: z.string().optional().nullable(),
    latitude: z.number().optional().nullable(),
    longitude: z.number().optional().nullable(),
    is_active: z.boolean(),
    notes: z.string().optional().nullable(),
    tags: z.array(z.string()).default([]),
    created_at: z.string().datetime(),
    updated_at: z.string().datetime(),
})

export type Customer = z.infer<typeof customerSchema>

export const createCustomerSchema = z.object({
    full_name: z.string().min(1, 'Full name is required'),
    phone: z.string().min(1, 'Phone is required'),
    email: z.string().email('Invalid email address').optional().or(z.literal('')),
    address: z.string().optional(),
    latitude: z.number().optional(),
    longitude: z.number().optional(),
    plan_id: z.string().uuid().optional(),
    router_id: z.string().uuid().optional(),
    username: z.string().optional(),
    password: z.string().optional(),
    static_ip: z.string().optional(),
})

export type CreateCustomer = z.infer<typeof createCustomerSchema>

export const registrationSchema = z.object({
    id: z.string().uuid(),
    full_name: z.string(),
    email: z.string().email().optional().nullable(),
    phone: z.string(),
    address: z.string().optional().nullable(),
    latitude: z.number().optional().nullable(),
    longitude: z.number().optional().nullable(),
    notes: z.string().optional().nullable(),
    bandwidth_profile_id: z.string().uuid().optional().nullable(),
    status: z.enum(['pending', 'approved', 'rejected']),
    rejection_reason: z.string().optional().nullable(),
    approved_by: z.string().uuid().optional().nullable(),
    approved_at: z.string().datetime().optional().nullable(),
    customer_id: z.string().uuid().optional().nullable(),
    created_at: z.string().datetime(),
    updated_at: z.string().datetime(),
})

export type Registration = z.infer<typeof registrationSchema>

export const approveRegistrationSchema = z.object({
    router_id: z.string().uuid('Router is required'),
    profile_id: z.string().uuid('Profile is required').optional(),
})

export type ApproveRegistration = z.infer<typeof approveRegistrationSchema>

export const rejectRegistrationSchema = z.object({
    reason: z.string().min(1, 'Rejection reason is required'),
})

export type RejectRegistration = z.infer<typeof rejectRegistrationSchema>
