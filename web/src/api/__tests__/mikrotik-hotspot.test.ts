import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  listHotspotProfiles,
  createHotspotProfile,
  listHotspotUsers,
  createHotspotUser,
  getHotspotUser,
  deleteHotspotUser,
  listHotspotActive,
  listHotspotHosts,
  listHotspotServers,
} from '@/api/mikrotik/hotspot'

const mc = vi.mocked(adminClient)
const RID = 'router-1'
beforeEach(() => vi.clearAllMocks())

describe('listHotspotProfiles', () => {
  it('GET /routers/:id/hotspot/profiles', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [{ name: 'default' }] } })
    const result = await listHotspotProfiles(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/hotspot/profiles`)
    expect(result[0].name).toBe('default')
  })
})

describe('createHotspotProfile', () => {
  it('POST /routers/:id/hotspot/profiles', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { name: 'gold' } } })
    const result = await createHotspotProfile(RID, { name: 'gold' })
    expect(mc.post).toHaveBeenCalledWith(`/routers/${RID}/hotspot/profiles`, { name: 'gold' })
    expect(result.name).toBe('gold')
  })
})

describe('listHotspotUsers', () => {
  it('GET /routers/:id/hotspot/users', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [{ name: 'user1' }] } })
    const result = await listHotspotUsers(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/hotspot/users`)
    expect(result).toHaveLength(1)
  })
})

describe('createHotspotUser', () => {
  it('POST /routers/:id/hotspot/users', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { name: 'user1', profile: 'default' } } })
    const result = await createHotspotUser(RID, { name: 'user1' })
    expect(result.name).toBe('user1')
  })
})

describe('deleteHotspotUser', () => {
  it('DELETE /routers/:id/hotspot/users/:id', async () => {
    mc.delete.mockResolvedValueOnce({ data: { success: true, data: { message: 'removed' } } })
    const msg = await deleteHotspotUser(RID, '*5')
    expect(mc.delete).toHaveBeenCalledWith(`/routers/${RID}/hotspot/users/*5`)
    expect(msg).toBe('removed')
  })
})

describe('listHotspotActive', () => {
  it('GET /routers/:id/hotspot/active', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listHotspotActive(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/hotspot/active`)
  })
})

describe('listHotspotHosts', () => {
  it('GET /routers/:id/hotspot/hosts', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listHotspotHosts(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/hotspot/hosts`)
  })
})

describe('listHotspotServers', () => {
  it('GET /routers/:id/hotspot/servers', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listHotspotServers(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/hotspot/servers`)
  })
})
