import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'

export const ReportSummarySchema = z.object({
    period_start: z.string(),
    period_end: z.string(),
    total_revenue: z.number(),
    total_invoiced: z.number(),
    total_invoices: z.number(),
    paid_invoices: z.number(),
    unpaid_invoices: z.number(),
    overdue_invoices: z.number(),
    total_payments: z.number(),
    total_customers: z.number(),
    active_customers: z.number(),
    new_customers: z.number(),
    active_subscriptions: z.number(),
    isolated_subscriptions: z.number(),
    suspended_subscriptions: z.number(),
}).passthrough()

export const ReportSummaryResponseSchema = ApiResponseSchema(ReportSummarySchema)

export type ReportSummary = z.infer<typeof ReportSummarySchema>

// Subscription report — returns array of subscription objects
export const SubscriptionReportResponseSchema = ApiResponseSchema(z.array(z.record(z.string(), z.unknown())))

// Cash flow, cash balance, reconciliation — backend returns freeform object
export const CashFlowReportResponseSchema = ApiResponseSchema(z.record(z.string(), z.unknown()))
export const CashBalanceReportResponseSchema = ApiResponseSchema(z.record(z.string(), z.unknown()))
export const ReconciliationReportResponseSchema = ApiResponseSchema(z.record(z.string(), z.unknown()))
