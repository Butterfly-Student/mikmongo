import { adminClient } from '@/lib/axios/admin-client'
import {
  PaymentListResponseSchema,
  PaymentDetailResponseSchema,
  GatewayPaymentDetailResponseSchema,
  type PaymentResponse,
  type GatewayPaymentResponse,
  type CreatePaymentRequest,
} from '@/lib/schemas/billing'
import { MessageResponseSchema } from '@/lib/schemas/auth'

export async function listPayments(
  limit?: number,
  offset?: number
): Promise<{ payments: PaymentResponse[]; meta: { total: number; limit: number; offset: number } }> {
  const params: Record<string, number> = {}
  if (limit !== undefined) params.limit = limit
  if (offset !== undefined) params.offset = offset
  const response = await adminClient.get('/payments', { params })
  const parsed = PaymentListResponseSchema.parse(response.data)
  const meta = parsed.meta ?? { total: parsed.data.length, limit: limit ?? 50, offset: offset ?? 0 }
  return { payments: parsed.data, meta }
}

export async function getPayment(id: string): Promise<PaymentResponse> {
  const response = await adminClient.get(`/payments/${id}`)
  const parsed = PaymentDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function createPayment(data: CreatePaymentRequest): Promise<PaymentResponse> {
  const response = await adminClient.post('/payments', data)
  const parsed = PaymentDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function confirmPayment(id: string): Promise<string> {
  const response = await adminClient.post(`/payments/${id}/confirm`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function rejectPayment(id: string, reason: string): Promise<string> {
  const response = await adminClient.post(`/payments/${id}/reject`, { reason })
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function refundPayment(
  id: string,
  data: { amount: number; reason: string }
): Promise<string> {
  const response = await adminClient.post(`/payments/${id}/refund`, data)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function initiateGatewayPayment(
  id: string,
  gateway: string
): Promise<GatewayPaymentResponse> {
  const response = await adminClient.post(`/payments/${id}/initiate-gateway?gateway=${gateway}`)
  const parsed = GatewayPaymentDetailResponseSchema.parse(response.data)
  return parsed.data
}
