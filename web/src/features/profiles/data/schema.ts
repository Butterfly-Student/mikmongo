import { z } from "zod";

export const pppProfileInfoSchema = z.object({
    rate_limit: z.string().optional().nullable(),
    local_address: z.string().optional().nullable(),
    remote_address: z.string().optional().nullable(),
    parent_queue: z.string().optional().nullable(),
    queue_type: z.string().optional().nullable(),
    dns_server: z.string().optional().nullable(),
    session_timeout: z.string().optional().nullable(),
    idle_timeout: z.string().optional().nullable(),
});

export const profileSchema = z.object({
    id: z.string().uuid(),
    router_id: z.string().uuid(),
    profile_code: z.string(),
    name: z.string(),
    description: z.string().optional().nullable(),
    download_speed: z.number().int(),
    upload_speed: z.number().int(),
    price_monthly: z.number(),
    tax_rate: z.number().optional().nullable(),
    billing_cycle: z.enum(["daily", "weekly", "monthly", "yearly"]),
    billing_day: z.number().int().optional().nullable(),
    is_active: z.boolean(),
    is_visible: z.boolean(),
    sort_order: z.number().int().optional().nullable(),
    grace_period_days: z.number().int().optional().nullable(),
    isolate_profile_name: z.string().optional().nullable(),
    created_at: z.string().datetime(),
    updated_at: z.string().datetime(),
    mikrotik: pppProfileInfoSchema.optional().nullable(),
});

export type Profile = z.infer<typeof profileSchema>;

export const createProfileSchema = z.object({
    profile_code: z.string().min(1, { message: "Profile code is required" }),
    name: z.string().min(1, { message: "Name is required" }),
    description: z.string().optional(),
    download_speed: z.coerce.number().int().min(1),
    upload_speed: z.coerce.number().int().min(1),
    price_monthly: z.coerce.number().min(0.01),
    tax_rate: z.coerce.number().optional(),
    billing_cycle: z.enum(["daily", "weekly", "monthly", "yearly"]).default("monthly"),
    billing_day: z.coerce.number().int().optional(),
    grace_period_days: z.coerce.number().int().optional(),
    isolate_profile_name: z.string().optional(),
    sort_order: z.coerce.number().int().optional(),
    is_visible: z.boolean().default(true),
    mt_local_address: z.string().optional(),
    mt_remote_address: z.string().optional(),
    mt_parent_queue: z.string().optional(),
    mt_queue_type: z.string().optional(),
    mt_dns_server: z.string().optional(),
    mt_session_timeout: z.string().optional(),
    mt_idle_timeout: z.string().optional(),
});

export type CreateProfile = z.infer<typeof createProfileSchema>;
