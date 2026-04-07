import { customerClient } from '@/lib/axios/customer-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import { CustomerDetailResponseSchema, type CustomerResponse } from '@/lib/schemas/customer'
import {
  PaymentDetailResponseSchema,
  type PaymentResponse,
} from '@/lib/schemas/billing'

export async function getCustomerProfile(): Promise<CustomerResponse> {
  const response = await customerClient.get('/profile')
  const parsed = CustomerDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function changeCustomerPassword(password: string): Promise<string> {
  const response = await customerClient.put('/profile/password', { password })
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function getPortalPayment(id: string): Promise<PaymentResponse> {
  const response = await customerClient.get(`/payments/${id}`)
  const parsed = PaymentDetailResponseSchema.parse(response.data)
  return parsed.data
}
