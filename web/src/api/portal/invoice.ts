import { customerClient } from '@/lib/axios/customer-client'
import { InvoiceDetailResponseSchema, type InvoiceResponse } from '@/lib/schemas/billing'
import { ApiResponseSchema } from '@/lib/schemas/auth'
import { invoiceResponseSchema } from '@/lib/schemas/billing'
import { z } from 'zod'

const PortalInvoiceListResponseSchema = ApiResponseSchema(z.array(invoiceResponseSchema))

export async function listPortalInvoices(): Promise<InvoiceResponse[]> {
  const response = await customerClient.get('/invoices')
  const parsed = PortalInvoiceListResponseSchema.parse(response.data)
  return parsed.data
}

export async function getPortalInvoice(id: string): Promise<InvoiceResponse> {
  const response = await customerClient.get(`/invoices/${id}`)
  const parsed = InvoiceDetailResponseSchema.parse(response.data)
  return parsed.data
}
