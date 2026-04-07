import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  SettingsListResponseSchema,
  SettingDetailResponseSchema,
  type SystemSettingResponse,
  type UpsertSystemSettingRequest,
} from '@/lib/schemas/settings'

export async function listSettings(): Promise<SystemSettingResponse[]> {
  const response = await adminClient.get('/settings')
  const parsed = SettingsListResponseSchema.parse(response.data)
  return parsed.data
}

export async function upsertSetting(
  data: UpsertSystemSettingRequest
): Promise<SystemSettingResponse> {
  const response = await adminClient.post('/settings', data)
  const parsed = SettingDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getSetting(id: string): Promise<SystemSettingResponse> {
  const response = await adminClient.get(`/settings/${id}`)
  const parsed = SettingDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updateSetting(
  id: string,
  data: UpsertSystemSettingRequest
): Promise<SystemSettingResponse> {
  const response = await adminClient.put(`/settings/${id}`, data)
  const parsed = SettingDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deleteSetting(id: string): Promise<string> {
  const response = await adminClient.delete(`/settings/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}
