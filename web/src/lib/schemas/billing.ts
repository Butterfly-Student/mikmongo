import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'
import { MetaSchema } from '@/lib/schemas/router'


// Invoice
export const invoiceResponseSchema = z.object({
  id: z.string().uuid(),
  invoice_number: z.string(),
  customer_id: z.string().uuid(),
  subscription_id: z.string().uuid().nullish(),
  billing_period_start: z.string(),
  billing_period_end: z.string(),
  billing_month: z.number().int(),
  billing_year: z.number().int(),
  issue_date: z.string(),
  due_date: z.string(),
  payment_deadline: z.string(),
  subtotal: z.number(),
  tax_amount: z.number(),
  discount_amount: z.number(),
  late_fee: z.number(),
  total_amount: z.number(),
  paid_amount: z.number(),
  balance: z.number(),
  status: z.enum(['draft', 'sent', 'unpaid', 'partial', 'paid', 'overpaid', 'overdue', 'cancelled', 'refunded']),
  payment_date: z.string().nullish(),
  payment_method: z.string().nullish(),
  invoice_type: z.enum(['recurring', 'installation', 'additional', 'refund']),
  is_auto_generated: z.boolean(),
  notes: z.string().nullish(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type InvoiceResponse = z.infer<typeof invoiceResponseSchema>

// Payment
export const paymentResponseSchema = z.object({
  id: z.string().uuid(),
  payment_number: z.string(),
  customer_id: z.string().uuid(),
  amount: z.number(),
  allocated_amount: z.number(),
  remaining_amount: z.number(),
  payment_method: z.enum(['cash', 'bank_transfer', 'e-wallet', 'credit_card', 'debit_card', 'check', 'qris', 'gateway']),
  payment_date: z.string(),
  bank_name: z.string().nullish(),
  bank_account_number: z.string().nullish(),
  bank_account_name: z.string().nullish(),
  transaction_reference: z.string().nullish(),
  ewallet_provider: z.string().nullish(),
  ewallet_number: z.string().nullish(),
  gateway_name: z.string().nullish(),
  gateway_trx_id: z.string().nullish(),
  proof_image: z.string().nullish(),
  receipt_number: z.string().nullish(),
  status: z.enum(['pending', 'confirmed', 'rejected', 'refunded']),
  processed_at: z.string().nullish(),
  rejection_reason: z.string().nullish(),
  refund_amount: z.number().nullish(),
  refund_date: z.string().nullish(),
  refund_reason: z.string().nullish(),
  notes: z.string().nullish(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type PaymentResponse = z.infer<typeof paymentResponseSchema>

// Gateway payment
export const gatewayPaymentResponseSchema = z.object({
  payment_url: z.string(),
  expires_at: z.string(),
  gateway_id: z.string(),
})

export type GatewayPaymentResponse = z.infer<typeof gatewayPaymentResponseSchema>

// Payment form schemas
export const rejectPaymentSchema = z.object({
  reason: z.string().min(1, 'Alasan wajib diisi'),
})

export const refundPaymentSchema = z.object({
  amount: z.number().min(0.01),
  reason: z.string().min(1, 'Alasan wajib diisi'),
})

// Cash entry
export const cashEntryResponseSchema = z.object({
  id: z.string().uuid(),
  entry_number: z.string(),
  type: z.enum(['income', 'expense']),
  source: z.string(),
  amount: z.number(),
  description: z.string(),
  reference_type: z.string().nullish(),
  reference_id: z.string().nullish(),
  payment_method: z.string().nullish(),
  bank_name: z.string().nullish(),
  account_number: z.string().nullish(),
  petty_cash_fund_id: z.string().uuid().nullish(),
  entry_date: z.string(),
  status: z.string(),
  created_by: z.string().uuid().nullish(),
  approved_by: z.string().uuid().nullish(),
  approved_at: z.string().nullish(),
  notes: z.string().nullish(),
  receipt_image: z.string().nullish(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type CashEntryResponse = z.infer<typeof cashEntryResponseSchema>

export const createCashEntrySchema = z.object({
  type: z.enum(['income', 'expense']),
  source: z.enum(['invoice', 'agent_invoice', 'installation', 'penalty', 'other', 'operational', 'upstream', 'purchase', 'salary']),
  amount: z.number().min(0.01, 'Jumlah harus lebih dari 0'),
  description: z.string().min(1, 'Deskripsi wajib diisi'),
  payment_method: z.string().min(1, 'Metode pembayaran wajib diisi'),
  bank_name: z.string().optional(),
  account_number: z.string().optional(),
  petty_cash_fund_id: z.string().uuid().optional(),
  entry_date: z.string().optional(),
  notes: z.string().optional(),
})

export type CreateCashEntry = z.infer<typeof createCashEntrySchema>

// Petty cash fund
export const pettyCashFundResponseSchema = z.object({
  id: z.string().uuid(),
  fund_name: z.string(),
  initial_balance: z.number(),
  current_balance: z.number(),
  custodian_id: z.string().uuid(),
  status: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type PettyCashFundResponse = z.infer<typeof pettyCashFundResponseSchema>

export const topUpFundSchema = z.object({
  amount: z.number().min(0.01, 'Jumlah harus lebih dari 0'),
})

// Agent invoice
export const agentInvoiceResponseSchema = z.object({
  id: z.string().uuid(),
  agent_id: z.string().uuid(),
  router_id: z.string().uuid(),
  invoice_number: z.string(),
  billing_cycle: z.enum(['weekly', 'monthly']),
  period_start: z.string(),
  period_end: z.string(),
  billing_month: z.number().int(),
  billing_week: z.number().int().nullish(),
  billing_year: z.number().int(),
  voucher_count: z.number().int(),
  subtotal: z.number(),
  selling_total: z.number(),
  profit: z.number(),
  discount_amount: z.number(),
  total_amount: z.number(),
  paid_amount: z.number(),
  balance: z.number(),
  status: z.enum(['draft', 'unpaid', 'paid', 'cancelled']),
  notes: z.string().nullish(),
  created_at: z.string(),
  updated_at: z.string(),
})

export type AgentInvoiceResponse = z.infer<typeof agentInvoiceResponseSchema>

export const agentRequestPaymentSchema = z.object({
  paid_amount: z.number().optional(),
  notes: z.string().optional(),
})

// List/detail response schemas
export const InvoiceListResponseSchema = ApiResponseSchema(z.array(invoiceResponseSchema)).extend({
  meta: MetaSchema.optional(),
})

export const InvoiceDetailResponseSchema = ApiResponseSchema(invoiceResponseSchema)

export const PaymentListResponseSchema = ApiResponseSchema(z.array(paymentResponseSchema)).extend({
  meta: MetaSchema.optional(),
})

export const PaymentDetailResponseSchema = ApiResponseSchema(paymentResponseSchema)

export const GatewayPaymentDetailResponseSchema = ApiResponseSchema(gatewayPaymentResponseSchema)

export const CashEntryListResponseSchema = z.object({
  success: z.boolean(),
  data: z.array(cashEntryResponseSchema),
  meta: MetaSchema.optional(),
})

export const CashEntryDetailResponseSchema = ApiResponseSchema(cashEntryResponseSchema)

export const PettyCashFundListResponseSchema = ApiResponseSchema(z.array(pettyCashFundResponseSchema))

export const PettyCashFundDetailResponseSchema = ApiResponseSchema(pettyCashFundResponseSchema)

export const AgentInvoiceListResponseSchema = ApiResponseSchema(z.array(agentInvoiceResponseSchema))

export const AgentInvoiceDetailResponseSchema = ApiResponseSchema(agentInvoiceResponseSchema)

// ── Create payment request ──
export const createPaymentRequestSchema = z.object({
  customer_id: z.string().uuid(),
  amount: z.number().min(0.01),
  payment_method: z.enum(['cash', 'bank_transfer', 'e-wallet', 'credit_card', 'debit_card', 'check', 'qris', 'gateway']),
  payment_date: z.string(),
  bank_name: z.string().optional(),
  bank_account_number: z.string().optional(),
  bank_account_name: z.string().optional(),
  transaction_reference: z.string().optional(),
  ewallet_provider: z.string().optional(),
  ewallet_number: z.string().optional(),
  proof_image: z.string().optional(),
  notes: z.string().optional(),
})

export type CreatePaymentRequest = z.infer<typeof createPaymentRequestSchema>

// ── Hotspot sale ──
export const hotspotSaleResponseSchema = z.object({
  id: z.string().uuid(),
  router_id: z.string().uuid(),
  username: z.string().nullish(),
  profile: z.string().nullish(),
  price: z.number().nullish(),
  selling_price: z.number().nullish(),
  prefix: z.string().nullish(),
  batch_code: z.string().nullish(),
  sales_agent_id: z.string().uuid().nullish(),
  created_at: z.string(),
})

export type HotspotSaleResponse = z.infer<typeof hotspotSaleResponseSchema>

export const HotspotSaleListResponseSchema = ApiResponseSchema(z.array(hotspotSaleResponseSchema)).extend({
  meta: MetaSchema.optional(),
})

export const HotspotSaleDetailResponseSchema = ApiResponseSchema(hotspotSaleResponseSchema)
