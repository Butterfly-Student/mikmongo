import { adminClient } from '@/lib/axios/admin-client'
import {
  InvoiceListResponseSchema,
  InvoiceDetailResponseSchema,
  type InvoiceResponse,
} from '@/lib/schemas/billing'
import { MessageResponseSchema, ApiResponseSchema } from '@/lib/schemas/auth'
import { invoiceResponseSchema } from '@/lib/schemas/billing'
import { z } from 'zod'

const OverdueInvoiceListSchema = ApiResponseSchema(z.array(invoiceResponseSchema))

export async function listInvoices(
  limit?: number,
  offset?: number
): Promise<{ invoices: InvoiceResponse[]; meta: { total: number; limit: number; offset: number } }> {
  const params: Record<string, number> = {}
  if (limit !== undefined) params.limit = limit
  if (offset !== undefined) params.offset = offset
  const response = await adminClient.get('/invoices', { params })
  const parsed = InvoiceListResponseSchema.parse(response.data)
  const meta = parsed.meta ?? { total: parsed.data.length, limit: limit ?? 50, offset: offset ?? 0 }
  return { invoices: parsed.data, meta }
}

export async function listOverdueInvoices(): Promise<InvoiceResponse[]> {
  const response = await adminClient.get('/invoices/overdue')
  const parsed = OverdueInvoiceListSchema.parse(response.data)
  return parsed.data
}

export async function getInvoice(id: string): Promise<InvoiceResponse> {
  const response = await adminClient.get(`/invoices/${id}`)
  const parsed = InvoiceDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deleteInvoice(id: string): Promise<string> {
  const response = await adminClient.delete(`/invoices/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function triggerMonthlyBilling(): Promise<string> {
  const response = await adminClient.post('/invoices/trigger-monthly')
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}
