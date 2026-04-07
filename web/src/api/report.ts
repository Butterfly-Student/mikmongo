import { adminClient } from '@/lib/axios/admin-client'
import {
  ReportSummaryResponseSchema,
  SubscriptionReportResponseSchema,
  CashFlowReportResponseSchema,
  CashBalanceReportResponseSchema,
  ReconciliationReportResponseSchema,
} from '@/lib/schemas/report'
import type { ReportSummary } from '@/lib/schemas/report'

export async function getReportSummary(from?: string, to?: string): Promise<ReportSummary> {
  const params: Record<string, string> = {}
  if (from) params.from = from
  if (to) params.to = to
  const response = await adminClient.get('/reports/summary', { params })
  const parsed = ReportSummaryResponseSchema.parse(response.data)
  return parsed.data
}

export async function getSubscriptionReport(params?: {
  from?: string
  to?: string
  limit?: number
  offset?: number
}): Promise<unknown[]> {
  const response = await adminClient.get('/reports/subscriptions', { params })
  const parsed = SubscriptionReportResponseSchema.parse(response.data)
  return parsed.data
}

export async function getCashFlowReport(from: string, to: string): Promise<Record<string, unknown>> {
  const response = await adminClient.get('/reports/cash-flow', { params: { from, to } })
  const parsed = CashFlowReportResponseSchema.parse(response.data)
  return parsed.data
}

export async function getCashBalanceReport(): Promise<Record<string, unknown>> {
  const response = await adminClient.get('/reports/cash-balance')
  const parsed = CashBalanceReportResponseSchema.parse(response.data)
  return parsed.data
}

export async function getReconciliationReport(from: string, to: string): Promise<Record<string, unknown>> {
  const response = await adminClient.get('/reports/reconciliation', { params: { from, to } })
  const parsed = ReconciliationReportResponseSchema.parse(response.data)
  return parsed.data
}
