import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  listSalesAgents,
  createSalesAgent,
  getSalesAgent,
  updateSalesAgent,
  deleteSalesAgent,
  getAgentProfilePrices,
  upsertAgentProfilePrice,
} from '@/api/sales-agents'

const mc = vi.mocked(adminClient)
beforeEach(() => vi.clearAllMocks())

const agentData = {
  id: '550e8400-e29b-41d4-a716-446655440001', router_id: '550e8400-e29b-41d4-a716-446655440002', name: 'Agent', phone: '08',
  username: 'agt', status: 'active' as const,
  voucher_mode: null, voucher_length: null, voucher_type: null,
  bill_discount: null, billing_cycle: null, billing_day: null,
  created_at: '2024-01-01', updated_at: '2024-01-01',
}

describe('listSalesAgents', () => {
  it('GET /sales-agents and returns agents', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [agentData] } })
    const result = await listSalesAgents()
    expect(mc.get).toHaveBeenCalledWith('/sales-agents', expect.any(Object))
    expect(result.agents).toHaveLength(1)
  })
})

describe('createSalesAgent', () => {
  it('POST /sales-agents', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: agentData } })
    const result = await createSalesAgent({ router_id: '550e8400-e29b-41d4-a716-446655440002', name: 'A', username: 'u', password: 'pass123' })
    expect(mc.post).toHaveBeenCalledWith('/sales-agents', expect.any(Object))
    expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440001')
  })
})

describe('getSalesAgent', () => {
  it('GET /sales-agents/:id', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: agentData } })
    const result = await getSalesAgent('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.get).toHaveBeenCalledWith('/sales-agents/550e8400-e29b-41d4-a716-446655440001')
    expect(result.name).toBe('Agent')
  })
})

describe('updateSalesAgent', () => {
  it('PUT /sales-agents/:id', async () => {
    mc.put.mockResolvedValueOnce({ data: { success: true, data: { ...agentData, name: 'Updated' } } })
    const result = await updateSalesAgent('550e8400-e29b-41d4-a716-446655440001', { name: 'Updated' })
    expect(mc.put).toHaveBeenCalledWith('/sales-agents/550e8400-e29b-41d4-a716-446655440001', { name: 'Updated' })
    expect(result.name).toBe('Updated')
  })
})

describe('deleteSalesAgent', () => {
  it('DELETE /sales-agents/:id', async () => {
    mc.delete.mockResolvedValueOnce({ data: { success: true, data: { message: 'deleted' } } })
    const msg = await deleteSalesAgent('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.delete).toHaveBeenCalledWith('/sales-agents/550e8400-e29b-41d4-a716-446655440001')
    expect(msg).toBe('deleted')
  })
})

describe('getAgentProfilePrices', () => {
  it('GET /sales-agents/:id/profile-prices', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    const result = await getAgentProfilePrices('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.get).toHaveBeenCalledWith('/sales-agents/550e8400-e29b-41d4-a716-446655440001/profile-prices')
    expect(result).toEqual([])
  })
})

describe('upsertAgentProfilePrice', () => {
  it('PUT /sales-agents/:id/profile-prices/:profile', async () => {
    const price = { id: '550e8400-e29b-41d4-a716-446655440003', sales_agent_id: '550e8400-e29b-41d4-a716-446655440001', profile_name: 'basic', base_price: 10000, selling_price: 12000, voucher_length: null, is_active: true, created_at: '2024-01-01' }
    mc.put.mockResolvedValueOnce({ data: { success: true, data: price } })
    const result = await upsertAgentProfilePrice('550e8400-e29b-41d4-a716-446655440001', 'basic', { selling_price: 12000 })
    expect(mc.put).toHaveBeenCalledWith('/sales-agents/550e8400-e29b-41d4-a716-446655440001/profile-prices/basic', { selling_price: 12000 })
    expect(result.profile_name).toBe('basic')
  })
})
