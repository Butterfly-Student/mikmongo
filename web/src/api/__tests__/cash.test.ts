import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  listCashEntries, createCashEntry, getCashEntry, updateCashEntry, deleteCashEntry,
  approveCashEntry, rejectCashEntry,
  listPettyCashFunds, createPettyCashFund, getPettyCashFund, topUpPettyCashFund,
} from '@/api/cash'

const mc = vi.mocked(adminClient)
beforeEach(() => vi.clearAllMocks())

const entry = { id: '550e8400-e29b-41d4-a716-446655440001', entry_number: 'CE-001', type: 'income' as const, source: 'invoice', amount: 100000, description: 'Payment', payment_method: 'cash', entry_date: '2024-01-01', status: 'approved', created_at: '2024-01-01', updated_at: '2024-01-01' }
const fund = { id: '550e8400-e29b-41d4-a716-446655440002', fund_name: 'Kas Kecil', initial_balance: 500000, current_balance: 400000, custodian_id: '550e8400-e29b-41d4-a716-446655440003', status: 'active', created_at: '2024-01-01', updated_at: '2024-01-01' }

describe('listCashEntries', () => {
  it('GET /cash-entries', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [entry] } })
    const result = await listCashEntries()
    expect(mc.get).toHaveBeenCalledWith('/cash-entries', expect.any(Object))
    expect(result.entries).toHaveLength(1)
  })
})

describe('getCashEntry', () => {
  it('GET /cash-entries/:id', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: entry } })
    const result = await getCashEntry('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.get).toHaveBeenCalledWith('/cash-entries/550e8400-e29b-41d4-a716-446655440001')
    expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440001')
  })
})

describe('deleteCashEntry', () => {
  it('DELETE /cash-entries/:id', async () => {
    mc.delete.mockResolvedValueOnce({ data: {} })
    await deleteCashEntry('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.delete).toHaveBeenCalledWith('/cash-entries/550e8400-e29b-41d4-a716-446655440001')
  })
})

describe('approveCashEntry', () => {
  it('POST /cash-entries/:id/approve', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: entry } })
    await approveCashEntry('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.post).toHaveBeenCalledWith('/cash-entries/550e8400-e29b-41d4-a716-446655440001/approve')
  })
})

describe('getPettyCashFund', () => {
  it('GET /petty-cash/:id', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: fund } })
    const result = await getPettyCashFund('550e8400-e29b-41d4-a716-446655440002')
    expect(mc.get).toHaveBeenCalledWith('/petty-cash/550e8400-e29b-41d4-a716-446655440002')
    expect(result.fund_name).toBe('Kas Kecil')
  })
})

describe('topUpPettyCashFund', () => {
  it('POST /petty-cash/:id/topup', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { ...fund, current_balance: 500000 } } })
    const result = await topUpPettyCashFund('550e8400-e29b-41d4-a716-446655440002', 100000)
    expect(mc.post).toHaveBeenCalledWith('/petty-cash/550e8400-e29b-41d4-a716-446655440002/topup', { amount: 100000 })
    expect(result.current_balance).toBe(500000)
  })
})
