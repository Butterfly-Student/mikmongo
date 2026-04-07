import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  generateVouchers, listVouchers, removeVoucherBatch,
  createMikhmonProfile, generateProfileScript,
  listMikhmonReports, getMikhmonReportSummary,
  getExpirationStatus, disableExpiration,
} from '@/api/mikrotik/mikhmon'

const mc = vi.mocked(adminClient)
const RID = 'router-1'
beforeEach(() => vi.clearAllMocks())

describe('generateVouchers', () => {
  it('POST /routers/:id/mikhmon/vouchers/generate', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { vouchers: [], quantity: 5, profile: 'basic' } } })
    const result = await generateVouchers(RID, { quantity: 5, profile: 'basic', mode: 'vc', char_set: 'alphanum' })
    expect(mc.post).toHaveBeenCalledWith(`/routers/${RID}/mikhmon/vouchers/generate`, expect.any(Object))
    expect(result.profile).toBe('basic')
  })
})

describe('listVouchers', () => {
  it('GET /routers/:id/mikhmon/vouchers?comment=batch1', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listVouchers(RID, 'batch1')
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/mikhmon/vouchers`, { params: { comment: 'batch1' } })
  })
})

describe('removeVoucherBatch', () => {
  it('DELETE /routers/:id/mikhmon/vouchers?comment=batch1', async () => {
    mc.delete.mockResolvedValueOnce({ data: { success: true, data: { message: 'removed' } } })
    const msg = await removeVoucherBatch(RID, 'batch1')
    expect(mc.delete).toHaveBeenCalledWith(`/routers/${RID}/mikhmon/vouchers`, { params: { comment: 'batch1' } })
    expect(msg).toBe('removed')
  })
})

describe('generateProfileScript', () => {
  it('POST /routers/:id/mikhmon/profiles/generate-script', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { name: 'script', content: '...' } } })
    const result = await generateProfileScript(RID, { mode: 'vc', profile_name: 'basic' })
    expect(mc.post).toHaveBeenCalledWith(`/routers/${RID}/mikhmon/profiles/generate-script`, expect.any(Object))
    expect(result.name).toBe('script')
  })
})

describe('getMikhmonReportSummary', () => {
  it('GET /routers/:id/mikhmon/reports/summary', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: { totalCount: 10, totalRevenue: 500000 } } })
    const result = await getMikhmonReportSummary(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/mikhmon/reports/summary`)
    expect(result.totalCount).toBe(10)
  })
})

describe('getExpirationStatus', () => {
  it('GET /routers/:id/mikhmon/expire/status', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: { enabled: true } } })
    const result = await getExpirationStatus(RID)
    expect(result.enabled).toBe(true)
  })
})

describe('disableExpiration', () => {
  it('POST /routers/:id/mikhmon/expire/disable', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { message: 'disabled' } } })
    const msg = await disableExpiration(RID)
    expect(msg).toBe('disabled')
  })
})
