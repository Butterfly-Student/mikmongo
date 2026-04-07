import { adminClient } from '@/lib/axios/admin-client'
import {
  CashEntryListResponseSchema,
  CashEntryDetailResponseSchema,
  PettyCashFundListResponseSchema,
  PettyCashFundDetailResponseSchema,
  type CashEntryResponse,
  type CreateCashEntry,
  type PettyCashFundResponse,
} from '@/lib/schemas/billing'

export async function listCashEntries(params?: {
  type?: string
  source?: string
  status?: string
  date_from?: string
  date_to?: string
  limit?: number
  offset?: number
}): Promise<{ entries: CashEntryResponse[]; meta?: { total: number; limit: number; offset: number } }> {
  const response = await adminClient.get('/cash-entries', { params })
  const parsed = CashEntryListResponseSchema.parse(response.data)
  return { entries: parsed.data, meta: parsed.meta }
}

export async function createCashEntry(data: CreateCashEntry): Promise<CashEntryResponse> {
  const response = await adminClient.post('/cash-entries', data)
  const parsed = CashEntryDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function approveCashEntry(id: string): Promise<CashEntryResponse> {
  const response = await adminClient.post(`/cash-entries/${id}/approve`)
  const parsed = CashEntryDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function rejectCashEntry(id: string, reason: string): Promise<CashEntryResponse> {
  const response = await adminClient.post(`/cash-entries/${id}/reject`, { reason })
  const parsed = CashEntryDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function listPettyCashFunds(): Promise<PettyCashFundResponse[]> {
  const response = await adminClient.get('/petty-cash')
  const parsed = PettyCashFundListResponseSchema.parse(response.data)
  return parsed.data
}

export async function createPettyCashFund(data: {
  fund_name: string
  initial_balance: number
  custodian_id: string
}): Promise<PettyCashFundResponse> {
  const response = await adminClient.post('/petty-cash', data)
  const parsed = PettyCashFundDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function topUpPettyCashFund(
  id: string,
  amount: number
): Promise<PettyCashFundResponse> {
  const response = await adminClient.post(`/petty-cash/${id}/topup`, { amount })
  const parsed = PettyCashFundDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getCashEntry(id: string): Promise<CashEntryResponse> {
  const response = await adminClient.get(`/cash-entries/${id}`)
  const parsed = CashEntryDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updateCashEntry(
  id: string,
  data: Partial<CreateCashEntry>
): Promise<CashEntryResponse> {
  const response = await adminClient.put(`/cash-entries/${id}`, data)
  const parsed = CashEntryDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deleteCashEntry(id: string): Promise<void> {
  await adminClient.delete(`/cash-entries/${id}`)
}

export async function getPettyCashFund(id: string): Promise<PettyCashFundResponse> {
  const response = await adminClient.get(`/petty-cash/${id}`)
  const parsed = PettyCashFundDetailResponseSchema.parse(response.data)
  return parsed.data
}
