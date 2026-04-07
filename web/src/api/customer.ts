import axios from 'axios'
import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
    CustomerListResponseSchema,
    CustomerDetailResponseSchema,
    RegistrationListResponseSchema,
    RegistrationDetailResponseSchema,
} from '@/lib/schemas/customer'
import type { CustomerResponse, RegistrationResponse } from '@/lib/schemas/customer'

export async function listCustomers(
    limit?: number,
    offset?: number
): Promise<{ customers: CustomerResponse[]; meta: { total: number; limit: number; offset: number } }> {
    const params: Record<string, number> = {}
    if (limit !== undefined) params.limit = limit
    if (offset !== undefined) params.offset = offset
    const response = await adminClient.get('/customers', { params })
    const parsed = CustomerListResponseSchema.parse(response.data)
    return { customers: parsed.data, meta: parsed.meta }
}

export async function getCustomer(id: string): Promise<CustomerResponse> {
    const response = await adminClient.get(`/customers/${id}`)
    const parsed = CustomerDetailResponseSchema.parse(response.data)
    return parsed.data
}

export async function createCustomer(
    data: Record<string, unknown>
): Promise<CustomerResponse> {
    const response = await adminClient.post('/customers', data)
    const parsed = CustomerDetailResponseSchema.parse(response.data)
    return parsed.data
}

export async function updateCustomer(
    id: string,
    data: Record<string, unknown>
): Promise<CustomerResponse> {
    const response = await adminClient.put(`/customers/${id}`, data)
    const parsed = CustomerDetailResponseSchema.parse(response.data)
    return parsed.data
}

export async function deleteCustomer(id: string): Promise<void> {
    const response = await adminClient.delete(`/customers/${id}`)
    MessageResponseSchema.parse(response.data)
}

export async function activateCustomerAccount(id: string): Promise<void> {
    const response = await adminClient.post(`/customers/${id}/activate-account`)
    MessageResponseSchema.parse(response.data)
}

export async function deactivateCustomerAccount(id: string): Promise<void> {
    const response = await adminClient.post(`/customers/${id}/deactivate-account`)
    MessageResponseSchema.parse(response.data)
}

// ── Registrations ──

export async function listRegistrations(
    limit?: number,
    offset?: number
): Promise<{ registrations: RegistrationResponse[]; meta: { total: number; limit: number; offset: number } }> {
    const params: Record<string, number> = {}
    if (limit !== undefined) params.limit = limit
    if (offset !== undefined) params.offset = offset
    const response = await adminClient.get('/registrations', { params })
    const parsed = RegistrationListResponseSchema.parse(response.data)
    return { registrations: parsed.data, meta: parsed.meta }
}

export async function approveRegistration(
    id: string,
    data: { router_id: string; profile_id?: string }
): Promise<void> {
    const response = await adminClient.post(`/registrations/${id}/approve`, data)
    MessageResponseSchema.parse(response.data)
}

export async function rejectRegistration(
    id: string,
    data: { reason: string }
): Promise<void> {
    const response = await adminClient.post(`/registrations/${id}/reject`, data)
    MessageResponseSchema.parse(response.data)
}

export async function getRegistration(id: string): Promise<RegistrationResponse> {
    const response = await adminClient.get(`/registrations/${id}`)
    const parsed = RegistrationDetailResponseSchema.parse(response.data)
    return parsed.data
}

export async function publicRegister(data: Record<string, unknown>): Promise<void> {
    await axios.post('/api/v1/register', data)
}
