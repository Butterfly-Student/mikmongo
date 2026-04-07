import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  getReportSummary,
  getSubscriptionReport,
  getCashFlowReport,
  getCashBalanceReport,
  getReconciliationReport,
} from '@/api/report'

const mc = vi.mocked(adminClient)
beforeEach(() => vi.clearAllMocks())

const summary = { period_start: '2024-01-01', period_end: '2024-01-31', total_revenue: 500000, total_invoiced: 600000, total_invoices: 10, paid_invoices: 8, unpaid_invoices: 2, overdue_invoices: 1, total_payments: 8, total_customers: 50, active_customers: 45, new_customers: 5, active_subscriptions: 45, isolated_subscriptions: 0, suspended_subscriptions: 2 }

describe('getReportSummary', () => {
  it('GET /reports/summary', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: summary } })
    const result = await getReportSummary('2024-01-01', '2024-01-31')
    expect(mc.get).toHaveBeenCalledWith('/reports/summary', expect.any(Object))
    expect(result.total_revenue).toBe(500000)
  })
})

describe('getSubscriptionReport', () => {
  it('GET /reports/subscriptions', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [] } })
    await getSubscriptionReport()
    expect(mc.get).toHaveBeenCalledWith('/reports/subscriptions', expect.any(Object))
  })
})

describe('getCashFlowReport', () => {
  it('GET /reports/cash-flow with from/to', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: { income: 1000000, expense: 300000 } } })
    const result = await getCashFlowReport('2024-01-01', '2024-01-31')
    expect(mc.get).toHaveBeenCalledWith('/reports/cash-flow', { params: { from: '2024-01-01', to: '2024-01-31' } })
    expect(result).toHaveProperty('income')
  })
})

describe('getCashBalanceReport', () => {
  it('GET /reports/cash-balance', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: { balance: 700000 } } })
    await getCashBalanceReport()
    expect(mc.get).toHaveBeenCalledWith('/reports/cash-balance')
  })
})

describe('getReconciliationReport', () => {
  it('GET /reports/reconciliation', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: {} } })
    await getReconciliationReport('2024-01-01', '2024-01-31')
    expect(mc.get).toHaveBeenCalledWith('/reports/reconciliation', { params: { from: '2024-01-01', to: '2024-01-31' } })
  })
})
