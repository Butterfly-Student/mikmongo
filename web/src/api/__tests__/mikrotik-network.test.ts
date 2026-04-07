import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  listQueues, listFirewallFilters, listFirewallNat, listFirewallAddressList,
  listIpPools, createIpPool, getIpPool, updateIpPool, deleteIpPool, listIpAddresses,
} from '@/api/mikrotik/network'

const mc = vi.mocked(adminClient)
const RID = 'router-1'
beforeEach(() => vi.clearAllMocks())

const pool = { name: 'pool-10m', ranges: '10.0.0.2-10.0.0.254' }

describe('listQueues', () => {
  it('GET /routers/:id/queue/simple', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listQueues(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/queue/simple`)
  })
})

describe('listFirewallFilters', () => {
  it('GET /routers/:id/firewall/filter', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listFirewallFilters(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/firewall/filter`)
  })
})

describe('listFirewallNat', () => {
  it('GET /routers/:id/firewall/nat', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listFirewallNat(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/firewall/nat`)
  })
})

describe('listFirewallAddressList', () => {
  it('GET /routers/:id/firewall/address-list', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listFirewallAddressList(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/firewall/address-list`)
  })
})

describe('createIpPool', () => {
  it('POST /routers/:id/ip/pools', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: pool } })
    const result = await createIpPool(RID, { name: 'pool-10m', ranges: '10.0.0.2-10.0.0.254' })
    expect(mc.post).toHaveBeenCalledWith(`/routers/${RID}/ip/pools`, expect.any(Object))
    expect(result.name).toBe('pool-10m')
  })
})

describe('deleteIpPool', () => {
  it('DELETE /routers/:id/ip/pools/:id', async () => {
    mc.delete.mockResolvedValueOnce({ data: { success: true, data: { message: 'removed' } } })
    const msg = await deleteIpPool(RID, '*1')
    expect(mc.delete).toHaveBeenCalledWith(`/routers/${RID}/ip/pools/*1`)
    expect(msg).toBe('removed')
  })
})

describe('listIpAddresses', () => {
  it('GET /routers/:id/ip/addresses', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await listIpAddresses(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/ip/addresses`)
  })
})
