import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn(), delete: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  listInvoices,
  listOverdueInvoices,
  getInvoice,
  deleteInvoice,
  triggerMonthlyBilling,
} from '@/api/invoice'

const mc = vi.mocked(adminClient)
beforeEach(() => vi.clearAllMocks())

const invoice = {
  id: '550e8400-e29b-41d4-a716-446655440001', invoice_number: 'INV-001', customer_id: '550e8400-e29b-41d4-a716-446655440002',
  subscription_id: null, billing_period_start: '2024-01-01', billing_period_end: '2024-01-31',
  billing_month: 1, billing_year: 2024, issue_date: '2024-01-01', due_date: '2024-01-15',
  payment_deadline: '2024-01-15', subtotal: 100000, tax_amount: 0, discount_amount: 0,
  late_fee: 0, total_amount: 100000, paid_amount: 0, balance: 100000,
  status: 'unpaid' as const, payment_date: null, payment_method: null,
  invoice_type: 'recurring' as const, is_auto_generated: true, notes: null,
  created_at: '2024-01-01', updated_at: '2024-01-01',
}

describe('listInvoices', () => {
  it('GET /invoices', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [invoice], meta: { total: 1, limit: 20, offset: 0 } } })
    const result = await listInvoices(20, 0)
    expect(mc.get).toHaveBeenCalledWith('/invoices', expect.any(Object))
    expect(result.invoices).toHaveLength(1)
  })
})

describe('listOverdueInvoices', () => {
  it('GET /invoices/overdue', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [invoice] } })
    const result = await listOverdueInvoices()
    expect(mc.get).toHaveBeenCalledWith('/invoices/overdue')
    expect(result).toHaveLength(1)
  })
})

describe('getInvoice', () => {
  it('GET /invoices/:id', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: invoice } })
    const result = await getInvoice('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.get).toHaveBeenCalledWith('/invoices/550e8400-e29b-41d4-a716-446655440001')
    expect(result.invoice_number).toBe('INV-001')
  })
})

describe('deleteInvoice', () => {
  it('DELETE /invoices/:id', async () => {
    mc.delete.mockResolvedValueOnce({ data: { success: true, data: { message: 'cancelled' } } })
    const msg = await deleteInvoice('550e8400-e29b-41d4-a716-446655440001')
    expect(mc.delete).toHaveBeenCalledWith('/invoices/550e8400-e29b-41d4-a716-446655440001')
    expect(msg).toBe('cancelled')
  })
})

describe('triggerMonthlyBilling', () => {
  it('POST /invoices/trigger-monthly', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: { message: 'triggered' } } })
    const msg = await triggerMonthlyBilling()
    expect(mc.post).toHaveBeenCalledWith('/invoices/trigger-monthly')
    expect(msg).toBe('triggered')
  })
})
