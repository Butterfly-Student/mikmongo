import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import { listSettings, upsertSetting, getSetting, updateSetting, deleteSetting } from '@/api/settings'

const mc = vi.mocked(adminClient)
beforeEach(() => vi.clearAllMocks())

const setting = { id: '550e8400-e29b-41d4-a716-446655440001', group_name: 'general', key_name: 'app_name', value: 'MikMongo', type: 'string' as const, label: 'App Name', description: null, is_encrypted: false, is_public: true, updated_at: '2024-01-01' }

describe('listSettings', () => {
  it('GET /settings', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [setting] } })
    const result = await listSettings()
    expect(mc.get).toHaveBeenCalledWith('/settings')
    expect(result).toHaveLength(1)
  })
})

describe('upsertSetting', () => {
  it('POST /settings', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: setting } })
    const result = await upsertSetting({ key_name: 'app_name', value: 'MikMongo' })
    expect(mc.post).toHaveBeenCalledWith('/settings', expect.any(Object))
    expect(result.key_name).toBe('app_name')
  })
})

describe('getSetting', () => {
  it('GET /settings/:id', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: setting } })
    const result = await getSetting('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.get).toHaveBeenCalledWith('/settings/550e8400-e29b-41d4-a716-446655440001')
    expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440001')
  })
})

describe('deleteSetting', () => {
  it('DELETE /settings/:id', async () => {
    mc.delete.mockResolvedValueOnce({ data: { success: true, data: { message: 'deleted' } } })
    const msg = await deleteSetting('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.delete).toHaveBeenCalledWith('/settings/550e8400-e29b-41d4-a716-446655440001')
    expect(msg).toBe('deleted')
  })
})
