import { adminClient } from '@/lib/axios/admin-client'
import {
    RouterListResponseSchema,
    SelectedRouterResponseSchema,
} from '@/lib/schemas/router'
import type { RouterResponse } from '@/lib/schemas/router'

export async function listRouters(
    limit?: number,
    offset?: number
): Promise<{ routers: RouterResponse[]; meta: { total: number; limit: number; offset: number } }> {
    const params: Record<string, number> = {}
    if (limit !== undefined) params.limit = limit
    if (offset !== undefined) params.offset = offset
    const response = await adminClient.get('/routers', { params })
    const parsed = RouterListResponseSchema.parse(response.data)
    const meta = parsed.meta ?? { total: parsed.data.length, limit: limit ?? 50, offset: offset ?? 0 }
    return { routers: parsed.data, meta }
}

export async function getSelectedRouter(): Promise<RouterResponse> {
    const response = await adminClient.get('/routers/selected')
    const parsed = SelectedRouterResponseSchema.parse(response.data)
    return parsed.data
}

export async function selectRouter(id: string): Promise<RouterResponse> {
    const response = await adminClient.post(`/routers/select/${id}`)
    const parsed = SelectedRouterResponseSchema.parse(response.data)
    return parsed.data
}

export async function createRouter(data: any): Promise<RouterResponse> {
    const response = await adminClient.post('/routers', data)
    const parsed = SelectedRouterResponseSchema.parse(response.data)
    return parsed.data
}

export async function syncRouter(id: string): Promise<void> {
    await adminClient.post(`/routers/${id}/sync`)
}

export async function testConnection(id: string): Promise<void> {
    await adminClient.post(`/routers/${id}/test-connection`)
}

export async function updateRouter(id: string, data: Record<string, unknown>): Promise<RouterResponse> {
    const response = await adminClient.put(`/routers/${id}`, data)
    const parsed = SelectedRouterResponseSchema.parse(response.data)
    return parsed.data
}

export async function deleteRouter(id: string): Promise<void> {
    await adminClient.delete(`/routers/${id}`)
}

export async function syncAllRouters(): Promise<void> {
    await adminClient.post('/routers/sync-all')
}
