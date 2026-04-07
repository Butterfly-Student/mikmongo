import { z } from "zod";

export const routerSchema = z.object({
    id: z.string().uuid(),
    name: z.string(),
    address: z.string(),
    area: z.string().optional().nullable(),
    api_port: z.number().int(),
    rest_port: z.number().int(),
    username: z.string(),
    use_ssl: z.boolean(),
    is_master: z.boolean(),
    is_active: z.boolean(),
    status: z.enum(["online", "offline", "unknown"]),
    last_seen_at: z.string().datetime().optional().nullable(),
    notes: z.string().optional().nullable(),
    created_at: z.string().datetime(),
    updated_at: z.string().datetime(),
});

export type Router = z.infer<typeof routerSchema>;

export const createRouterSchema = z.object({
    name: z.string().min(1, { message: "Name is required" }),
    address: z.string().min(1, { message: "Address is required" }),
    username: z.string().min(1, { message: "Username is required" }),
    password: z.string().min(1, { message: "Password is required" }),
    area: z.string().optional(),
    api_port: z.coerce.number().int().default(8728),
    rest_port: z.coerce.number().int().default(80),
    use_ssl: z.boolean().default(false),
    is_master: z.boolean().default(false),
    notes: z.string().optional(),
});

export type CreateRouter = z.infer<typeof createRouterSchema>;
