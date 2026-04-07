import { customerClient } from '@/lib/axios/customer-client'
import {
  GatewayPaymentDetailResponseSchema,
  type PaymentResponse,
  type GatewayPaymentResponse,
} from '@/lib/schemas/billing'
import { ApiResponseSchema } from '@/lib/schemas/auth'
import { paymentResponseSchema } from '@/lib/schemas/billing'
import { z } from 'zod'

const PortalPaymentListResponseSchema = ApiResponseSchema(z.array(paymentResponseSchema))

export async function listPortalPayments(): Promise<PaymentResponse[]> {
  const response = await customerClient.get('/payments')
  const parsed = PortalPaymentListResponseSchema.parse(response.data)
  return parsed.data
}

export async function initiatePortalPayment(id: string): Promise<GatewayPaymentResponse> {
  const response = await customerClient.post(`/payments/${id}/pay`)
  const parsed = GatewayPaymentDetailResponseSchema.parse(response.data)
  return parsed.data
}
