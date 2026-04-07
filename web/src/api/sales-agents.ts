import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  SalesAgentListResponseSchema,
  SalesAgentDetailResponseSchema,
  SalesProfilePriceListResponseSchema,
  SalesProfilePriceDetailResponseSchema,
  type SalesAgentResponse,
  type CreateSalesAgentRequest,
  type UpdateSalesAgentRequest,
  type SalesProfilePriceResponse,
  type UpsertProfilePriceRequest,
} from '@/lib/schemas/sales-agent'

export async function listSalesAgents(
  limit?: number,
  offset?: number
): Promise<{ agents: SalesAgentResponse[]; meta?: { total: number; limit: number; offset: number } }> {
  const params: Record<string, number> = {}
  if (limit !== undefined) params.limit = limit
  if (offset !== undefined) params.offset = offset
  const response = await adminClient.get('/sales-agents', { params })
  const parsed = SalesAgentListResponseSchema.parse(response.data)
  const meta = parsed.meta ?? { total: parsed.data.length, limit: limit ?? 50, offset: offset ?? 0 }
  return { agents: parsed.data, meta }
}

export async function createSalesAgent(data: CreateSalesAgentRequest): Promise<SalesAgentResponse> {
  const response = await adminClient.post('/sales-agents', data)
  const parsed = SalesAgentDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getSalesAgent(id: string): Promise<SalesAgentResponse> {
  const response = await adminClient.get(`/sales-agents/${id}`)
  const parsed = SalesAgentDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updateSalesAgent(
  id: string,
  data: UpdateSalesAgentRequest
): Promise<SalesAgentResponse> {
  const response = await adminClient.put(`/sales-agents/${id}`, data)
  const parsed = SalesAgentDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deleteSalesAgent(id: string): Promise<string> {
  const response = await adminClient.delete(`/sales-agents/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function getAgentProfilePrices(id: string): Promise<SalesProfilePriceResponse[]> {
  const response = await adminClient.get(`/sales-agents/${id}/profile-prices`)
  const parsed = SalesProfilePriceListResponseSchema.parse(response.data)
  return parsed.data
}

export async function upsertAgentProfilePrice(
  id: string,
  profile: string,
  data: UpsertProfilePriceRequest
): Promise<SalesProfilePriceResponse> {
  const response = await adminClient.put(`/sales-agents/${id}/profile-prices/${profile}`, data)
  const parsed = SalesProfilePriceDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function listAgentInvoicesBySalesAgent(
  agentId: string,
  limit?: number,
  offset?: number
): Promise<unknown[]> {
  const params: Record<string, number> = {}
  if (limit !== undefined) params.limit = limit
  if (offset !== undefined) params.offset = offset
  const response = await adminClient.get(`/sales-agents/${agentId}/invoices`, { params })
  return response.data?.data ?? []
}

export async function generateAgentInvoice(
  agentId: string,
  data: { period_start: string; period_end: string }
): Promise<unknown> {
  const response = await adminClient.post(`/sales-agents/${agentId}/invoices/generate`, data)
  return response.data?.data
}
