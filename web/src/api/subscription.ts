import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
    SubscriptionListResponseSchema,
    SubscriptionDetailResponseSchema,
} from '@/lib/schemas/subscription'
import type { SubscriptionResponse } from '@/lib/schemas/subscription'

export async function getSubscription(
    routerId: string,
    id: string
): Promise<SubscriptionResponse> {
    const response = await adminClient.get(`/routers/${routerId}/subscriptions/${id}`)
    const parsed = SubscriptionDetailResponseSchema.parse(response.data)
    return parsed.data
}

export async function updateSubscription(
    routerId: string,
    id: string,
    data: Record<string, unknown>
): Promise<SubscriptionResponse> {
    const response = await adminClient.put(`/routers/${routerId}/subscriptions/${id}`, data)
    const parsed = SubscriptionDetailResponseSchema.parse(response.data)
    return parsed.data
}

export async function listSubscriptions(
    routerId: string,
    limit?: number,
    offset?: number
): Promise<{ subscriptions: SubscriptionResponse[]; meta: { total: number; limit: number; offset: number } }> {
    const params: Record<string, number> = {}
    if (limit !== undefined) params.limit = limit
    if (offset !== undefined) params.offset = offset
    const response = await adminClient.get(`/routers/${routerId}/subscriptions`, { params })
    const parsed = SubscriptionListResponseSchema.parse(response.data)
    return { subscriptions: parsed.data, meta: parsed.meta }
}

export async function createSubscription(
    routerId: string,
    data: Record<string, unknown>
): Promise<SubscriptionResponse> {
    const response = await adminClient.post(`/routers/${routerId}/subscriptions`, data)
    const parsed = SubscriptionDetailResponseSchema.parse(response.data)
    return parsed.data
}

export async function activateSubscription(routerId: string, id: string): Promise<void> {
    const response = await adminClient.post(`/routers/${routerId}/subscriptions/${id}/activate`)
    MessageResponseSchema.parse(response.data)
}

export async function suspendSubscription(
    routerId: string,
    id: string,
    reason?: string
): Promise<void> {
    const body: Record<string, string> = {}
    if (reason) body.reason = reason
    const response = await adminClient.post(`/routers/${routerId}/subscriptions/${id}/suspend`, body)
    MessageResponseSchema.parse(response.data)
}

export async function isolateSubscription(
    routerId: string,
    id: string,
    reason?: string
): Promise<void> {
    const body: Record<string, string> = {}
    if (reason) body.reason = reason
    const response = await adminClient.post(`/routers/${routerId}/subscriptions/${id}/isolate`, body)
    MessageResponseSchema.parse(response.data)
}

export async function restoreSubscription(routerId: string, id: string): Promise<void> {
    const response = await adminClient.post(`/routers/${routerId}/subscriptions/${id}/restore`)
    MessageResponseSchema.parse(response.data)
}

export async function terminateSubscription(routerId: string, id: string): Promise<void> {
    const response = await adminClient.post(`/routers/${routerId}/subscriptions/${id}/terminate`)
    MessageResponseSchema.parse(response.data)
}

export async function deleteSubscription(routerId: string, id: string): Promise<void> {
    const response = await adminClient.delete(`/routers/${routerId}/subscriptions/${id}`)
    MessageResponseSchema.parse(response.data)
}
