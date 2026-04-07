import { agentClient } from '@/lib/axios/agent-client'
import {
  AgentInvoiceListResponseSchema,
  AgentInvoiceDetailResponseSchema,
  type AgentInvoiceResponse,
} from '@/lib/schemas/billing'

export async function listAgentPortalInvoices(
  limit?: number,
  offset?: number
): Promise<AgentInvoiceResponse[]> {
  const params: Record<string, number> = {}
  if (limit !== undefined) params.limit = limit
  if (offset !== undefined) params.offset = offset
  const response = await agentClient.get('/invoices', { params })
  const parsed = AgentInvoiceListResponseSchema.parse(response.data)
  return parsed.data
}

export async function getAgentPortalInvoice(id: string): Promise<AgentInvoiceResponse> {
  const response = await agentClient.get(`/invoices/${id}`)
  const parsed = AgentInvoiceDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function requestAgentPayment(
  id: string,
  data?: { paid_amount?: number; notes?: string }
): Promise<AgentInvoiceResponse> {
  const response = await agentClient.post(`/invoices/${id}/request-payment`, data ?? {})
  const parsed = AgentInvoiceDetailResponseSchema.parse(response.data)
  return parsed.data
}
