import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  PppProfileListResponseSchema,
  PppProfileDetailResponseSchema,
  PppSecretListResponseSchema,
  PppSecretDetailResponseSchema,
  PppActiveListResponseSchema,
  type PppProfile,
  type AddPppProfileRequest,
  type PppSecret,
  type AddPppSecretRequest,
  type PppActive,
} from '@/lib/schemas/mikrotik'

// ── PPP Profiles ──────────────────────────────────────────────────────

export async function listPppProfiles(routerId: string): Promise<PppProfile[]> {
  const response = await adminClient.get(`/routers/${routerId}/ppp/profiles`)
  const parsed = PppProfileListResponseSchema.parse(response.data)
  return parsed.data
}

export async function createPppProfile(
  routerId: string,
  data: AddPppProfileRequest
): Promise<PppProfile> {
  const response = await adminClient.post(`/routers/${routerId}/ppp/profiles`, data)
  const parsed = PppProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getPppProfile(routerId: string, nameOrId: string): Promise<PppProfile> {
  const response = await adminClient.get(`/routers/${routerId}/ppp/profiles/${nameOrId}`)
  const parsed = PppProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updatePppProfile(
  routerId: string,
  id: string,
  data: Partial<AddPppProfileRequest>
): Promise<PppProfile> {
  const response = await adminClient.put(`/routers/${routerId}/ppp/profiles/${id}`, data)
  const parsed = PppProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deletePppProfile(routerId: string, id: string): Promise<string> {
  const response = await adminClient.delete(`/routers/${routerId}/ppp/profiles/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

// ── PPP Secrets ───────────────────────────────────────────────────────

export async function listPppSecrets(routerId: string, profile?: string): Promise<PppSecret[]> {
  const params = profile ? { profile } : undefined
  const response = await adminClient.get(`/routers/${routerId}/ppp/secrets`, { params })
  const parsed = PppSecretListResponseSchema.parse(response.data)
  return parsed.data
}

export async function createPppSecret(
  routerId: string,
  data: AddPppSecretRequest
): Promise<PppSecret> {
  const response = await adminClient.post(`/routers/${routerId}/ppp/secrets`, data)
  const parsed = PppSecretDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getPppSecret(routerId: string, nameOrId: string): Promise<PppSecret> {
  const response = await adminClient.get(`/routers/${routerId}/ppp/secrets/${nameOrId}`)
  const parsed = PppSecretDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updatePppSecret(
  routerId: string,
  id: string,
  data: Partial<AddPppSecretRequest>
): Promise<PppSecret> {
  const response = await adminClient.put(`/routers/${routerId}/ppp/secrets/${id}`, data)
  const parsed = PppSecretDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deletePppSecret(routerId: string, id: string): Promise<string> {
  const response = await adminClient.delete(`/routers/${routerId}/ppp/secrets/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

// ── PPP Active ────────────────────────────────────────────────────────

export async function listPppActive(routerId: string): Promise<PppActive[]> {
  const response = await adminClient.get(`/routers/${routerId}/ppp/active`)
  const parsed = PppActiveListResponseSchema.parse(response.data)
  return parsed.data
}
