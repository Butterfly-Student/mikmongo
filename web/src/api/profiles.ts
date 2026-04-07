import { adminClient } from '@/lib/axios/admin-client'
import { ApiResponseSchema } from '@/lib/schemas/auth'
import { z } from 'zod'
import { profileSchema, type Profile, type CreateProfile } from '@/features/profiles/data/schema'

const ProfileListResponseSchema = ApiResponseSchema(z.array(profileSchema)).extend({
    meta: z.object({
        total: z.number(),
        limit: z.number(),
        offset: z.number(),
    })
})

const SingleProfileResponseSchema = ApiResponseSchema(profileSchema)

export async function listProfiles(routerId: string, limit?: number, offset?: number): Promise<{ profiles: Profile[]; meta: { total: number; limit: number; offset: number } }> {
    const params: Record<string, number> = {}
    if (limit !== undefined) params.limit = limit
    if (offset !== undefined) params.offset = offset
    const response = await adminClient.get(`/routers/${routerId}/bandwidth-profiles`, { params })
    const parsed = ProfileListResponseSchema.parse(response.data)
    return { profiles: parsed.data, meta: parsed.meta }
}

export async function createProfile(routerId: string, data: CreateProfile): Promise<Profile> {
    const response = await adminClient.post(`/routers/${routerId}/bandwidth-profiles`, data)
    const parsed = SingleProfileResponseSchema.parse(response.data)
    return parsed.data
}

export async function deleteProfile(routerId: string, id: string): Promise<void> {
    await adminClient.delete(`/routers/${routerId}/bandwidth-profiles/${id}`)
}

export async function updateProfile(routerId: string, id: string, data: Partial<CreateProfile>): Promise<Profile> {
    const response = await adminClient.put(`/routers/${routerId}/bandwidth-profiles/${id}`, data)
    const parsed = SingleProfileResponseSchema.parse(response.data)
    return parsed.data
}
