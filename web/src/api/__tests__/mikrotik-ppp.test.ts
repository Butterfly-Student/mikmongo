import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  listPppProfiles,
  createPppProfile,
  getPppProfile,
  updatePppProfile,
  deletePppProfile,
  listPppSecrets,
  createPppSecret,
  getPppSecret,
  deletePppSecret,
  listPppActive,
} from '@/api/mikrotik/ppp'

const mc = vi.mocked(adminClient)
const RID = 'router-1'
beforeEach(() => vi.clearAllMocks())

const profile = { name: '10M-Profile', localAddress: '10.0.0.1', remoteAddress: 'pool-10m', rateLimit: '10M/10M' }
const secret = { name: 'user01', password: 'secret', profile: '10M-Profile', service: 'pppoe' }
const active = { name: 'user01', address: '10.0.0.100', uptime: '1h', service: 'pppoe' }

describe('listPppProfiles', () => {
  it(`GET /routers/:id/ppp/profiles`, async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [profile] } })
    const result = await listPppProfiles(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/ppp/profiles`)
    expect(result).toHaveLength(1)
  })
})

describe('createPppProfile', () => {
  it('POST /routers/:id/ppp/profiles', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: profile } })
    const result = await createPppProfile(RID, { name: '10M-Profile' })
    expect(mc.post).toHaveBeenCalledWith(`/routers/${RID}/ppp/profiles`, expect.any(Object))
    expect(result.name).toBe('10M-Profile')
  })
})

describe('getPppProfile', () => {
  it('GET /routers/:id/ppp/profiles/:name', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: profile } })
    await getPppProfile(RID, '10M-Profile')
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/ppp/profiles/10M-Profile`)
  })
})

describe('deletePppProfile', () => {
  it('DELETE /routers/:id/ppp/profiles/:id', async () => {
    mc.delete.mockResolvedValueOnce({ data: { success: true, data: { message: 'removed' } } })
    const msg = await deletePppProfile(RID, '*1')
    expect(mc.delete).toHaveBeenCalledWith(`/routers/${RID}/ppp/profiles/*1`)
    expect(msg).toBe('removed')
  })
})

describe('listPppSecrets', () => {
  it('GET /routers/:id/ppp/secrets', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [secret] } })
    const result = await listPppSecrets(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/ppp/secrets`, { params: undefined })
    expect(result).toHaveLength(1)
  })

  it('passes profile filter', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listPppSecrets(RID, '10M-Profile')
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/ppp/secrets`, { params: { profile: '10M-Profile' } })
  })
})

describe('createPppSecret', () => {
  it('POST /routers/:id/ppp/secrets', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: secret } })
    const result = await createPppSecret(RID, { name: 'user01', password: 'secret' })
    expect(result.name).toBe('user01')
  })
})

describe('listPppActive', () => {
  it('GET /routers/:id/ppp/active', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [active] } })
    const result = await listPppActive(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/ppp/active`)
    expect(result).toHaveLength(1)
  })
})
