import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  listPayments,
  getPayment,
  createPayment,
  confirmPayment,
  rejectPayment,
  refundPayment,
  initiateGatewayPayment,
} from '@/api/payment'

const mc = vi.mocked(adminClient)
beforeEach(() => vi.clearAllMocks())

const payment = {
  id: '550e8400-e29b-41d4-a716-446655440001', payment_number: 'PAY-001', customer_id: '550e8400-e29b-41d4-a716-446655440002',
  amount: 100000, allocated_amount: 0, remaining_amount: 100000,
  payment_method: 'cash' as const, payment_date: '2024-01-01',
  bank_name: null, bank_account_number: null, bank_account_name: null,
  transaction_reference: null, ewallet_provider: null, ewallet_number: null,
  gateway_name: null, gateway_trx_id: null, proof_image: null, receipt_number: null,
  status: 'pending' as const, processed_at: null, rejection_reason: null,
  refund_amount: null, refund_date: null, refund_reason: null, notes: null,
  created_at: '2024-01-01', updated_at: '2024-01-01',
}

describe('listPayments', () => {
  it('GET /payments', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [payment], meta: { total: 1, limit: 20, offset: 0 } } })
    const result = await listPayments()
    expect(mc.get).toHaveBeenCalledWith('/payments', expect.any(Object))
    expect(result.payments).toHaveLength(1)
  })
})

describe('getPayment', () => {
  it('GET /payments/:id', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: payment } })
    const result = await getPayment('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.get).toHaveBeenCalledWith('/payments/550e8400-e29b-41d4-a716-446655440001')
    expect(result.payment_number).toBe('PAY-001')
  })
})

describe('createPayment', () => {
  it('POST /payments', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: payment } })
    const result = await createPayment({
      customer_id: '550e8400-e29b-41d4-a716-446655440002', amount: 100000,
      payment_method: 'cash', payment_date: '2024-01-01',
    })
    expect(mc.post).toHaveBeenCalledWith('/payments', expect.any(Object))
    expect(result.id).toBe('550e8400-e29b-41d4-a716-446655440001')
  })
})

describe('confirmPayment', () => {
  it('POST /payments/:id/confirm', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { message: 'confirmed' } } })
    const msg = await confirmPayment('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.post).toHaveBeenCalledWith('/payments/550e8400-e29b-41d4-a716-446655440001/confirm')
    expect(msg).toBe('confirmed')
  })
})

describe('rejectPayment', () => {
  it('POST /payments/:id/reject', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { message: 'rejected' } } })
    await rejectPayment('550e8400-e29b-41d4-a716-446655440001', 'invalid')
    expect(mc.post).toHaveBeenCalledWith('/payments/550e8400-e29b-41d4-a716-446655440001/reject', { reason: 'invalid' })
  })
})

describe('refundPayment', () => {
  it('POST /payments/:id/refund', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { message: 'refunded' } } })
    await refundPayment('550e8400-e29b-41d4-a716-446655440001', { amount: 50000, reason: 'error' })
    expect(mc.post).toHaveBeenCalledWith('/payments/550e8400-e29b-41d4-a716-446655440001/refund', { amount: 50000, reason: 'error' })
  })
})

describe('initiateGatewayPayment', () => {
  it('POST /payments/:id/initiate-gateway', async () => {
    mc.post.mockResolvedValueOnce({
      data: { success: true, data: { payment_url: 'https://pay.me', expires_at: '2024-01-01', gateway_id: 'gw1' } },
    })
    const result = await initiateGatewayPayment('550e8400-e29b-41d4-a716-446655440001', 'midtrans')
    expect(result.payment_url).toBe('https://pay.me')
  })
})
