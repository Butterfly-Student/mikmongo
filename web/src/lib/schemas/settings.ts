import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'

export const SystemSettingResponseSchema = z.object({
  id: z.string().uuid(),
  group_name: z.string().nullish(),
  key_name: z.string(),
  value: z.string().nullish(),
  type: z.enum(['string', 'integer', 'boolean', 'json', 'password']).nullish(),
  label: z.string().nullish(),
  description: z.string().nullish(),
  is_encrypted: z.boolean().nullish(),
  is_public: z.boolean().nullish(),
  updated_at: z.string().nullish(),
})

export type SystemSettingResponse = z.infer<typeof SystemSettingResponseSchema>

export const UpsertSystemSettingRequestSchema = z.object({
  group_name: z.string().optional(),
  key_name: z.string(),
  value: z.string().optional(),
  type: z.enum(['string', 'integer', 'boolean', 'json', 'password']).optional(),
  label: z.string().optional(),
  description: z.string().optional(),
  is_encrypted: z.boolean().optional(),
  is_public: z.boolean().optional(),
})

export type UpsertSystemSettingRequest = z.infer<typeof UpsertSystemSettingRequestSchema>

// Response wrappers
export const SettingsListResponseSchema = ApiResponseSchema(z.array(SystemSettingResponseSchema))

export const SettingDetailResponseSchema = ApiResponseSchema(SystemSettingResponseSchema)
