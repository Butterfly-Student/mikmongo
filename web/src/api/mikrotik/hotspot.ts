import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  HotspotProfileListResponseSchema,
  HotspotProfileDetailResponseSchema,
  HotspotUserListResponseSchema,
  HotspotUserDetailResponseSchema,
  HotspotActiveListResponseSchema,
  HotspotHostListResponseSchema,
  HotspotServerListResponseSchema,
  type HotspotProfile,
  type AddHotspotProfileRequest,
  type HotspotUser,
  type AddHotspotUserRequest,
  type HotspotActive,
  type HotspotHost,
  type HotspotServer,
} from '@/lib/schemas/mikrotik'

// ── Hotspot Profiles ──────────────────────────────────────────────────

export async function listHotspotProfiles(routerId: string): Promise<HotspotProfile[]> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot/profiles`)
  const parsed = HotspotProfileListResponseSchema.parse(response.data)
  return parsed.data
}

export async function createHotspotProfile(
  routerId: string,
  data: AddHotspotProfileRequest
): Promise<HotspotProfile> {
  const response = await adminClient.post(`/routers/${routerId}/hotspot/profiles`, data)
  const parsed = HotspotProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getHotspotProfile(
  routerId: string,
  nameOrId: string
): Promise<HotspotProfile> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot/profiles/${nameOrId}`)
  const parsed = HotspotProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updateHotspotProfile(
  routerId: string,
  id: string,
  data: Partial<AddHotspotProfileRequest>
): Promise<HotspotProfile> {
  const response = await adminClient.put(`/routers/${routerId}/hotspot/profiles/${id}`, data)
  const parsed = HotspotProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deleteHotspotProfile(routerId: string, id: string): Promise<string> {
  const response = await adminClient.delete(`/routers/${routerId}/hotspot/profiles/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

// ── Hotspot Users ─────────────────────────────────────────────────────

export async function listHotspotUsers(routerId: string): Promise<HotspotUser[]> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot/users`)
  const parsed = HotspotUserListResponseSchema.parse(response.data)
  return parsed.data
}

export async function createHotspotUser(
  routerId: string,
  data: AddHotspotUserRequest
): Promise<HotspotUser> {
  const response = await adminClient.post(`/routers/${routerId}/hotspot/users`, data)
  const parsed = HotspotUserDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getHotspotUser(routerId: string, nameOrId: string): Promise<HotspotUser> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot/users/${nameOrId}`)
  const parsed = HotspotUserDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updateHotspotUser(
  routerId: string,
  id: string,
  data: Partial<AddHotspotUserRequest>
): Promise<HotspotUser> {
  const response = await adminClient.put(`/routers/${routerId}/hotspot/users/${id}`, data)
  const parsed = HotspotUserDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deleteHotspotUser(routerId: string, id: string): Promise<string> {
  const response = await adminClient.delete(`/routers/${routerId}/hotspot/users/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

// ── Hotspot Active / Hosts / Servers ──────────────────────────────────

export async function listHotspotActive(routerId: string): Promise<HotspotActive[]> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot/active`)
  const parsed = HotspotActiveListResponseSchema.parse(response.data)
  return parsed.data
}

export async function listHotspotHosts(routerId: string): Promise<HotspotHost[]> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot/hosts`)
  const parsed = HotspotHostListResponseSchema.parse(response.data)
  return parsed.data
}

export async function listHotspotServers(routerId: string): Promise<HotspotServer[]> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot/servers`)
  const parsed = HotspotServerListResponseSchema.parse(response.data)
  return parsed.data
}
