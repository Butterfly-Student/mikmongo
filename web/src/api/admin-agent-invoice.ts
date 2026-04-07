import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  AgentInvoiceListResponseSchema,
  AgentInvoiceDetailResponseSchema,
  type AgentInvoiceResponse,
} from '@/lib/schemas/billing'

export async function listAdminAgentInvoices(params?: {
  agent_id?: string
  router_id?: string
  status?: string
  billing_cycle?: string
  billing_year?: number
  billing_month?: number
  billing_week?: number
  limit?: number
  offset?: number
}): Promise<AgentInvoiceResponse[]> {
  const response = await adminClient.get('/agent-invoices', { params })
  const parsed = AgentInvoiceListResponseSchema.parse(response.data)
  return parsed.data
}

export async function getAdminAgentInvoice(id: string): Promise<AgentInvoiceResponse> {
  const response = await adminClient.get(`/agent-invoices/${id}`)
  const parsed = AgentInvoiceDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function payAdminAgentInvoice(
  id: string,
  paidAmount: number
): Promise<AgentInvoiceResponse> {
  const response = await adminClient.put(`/agent-invoices/${id}/pay`, { paid_amount: paidAmount })
  const parsed = AgentInvoiceDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function cancelAdminAgentInvoice(id: string): Promise<string> {
  const response = await adminClient.put(`/agent-invoices/${id}/cancel`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function processScheduledAgentInvoices(): Promise<string> {
  const response = await adminClient.post('/agent-invoices/process')
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}
