import { adminClient } from '@/lib/axios/admin-client'
import {
  HotspotSaleListResponseSchema,
  HotspotSaleDetailResponseSchema,
  type HotspotSaleResponse,
} from '@/lib/schemas/billing'

export async function listHotspotSales(params?: {
  router_id?: string
  agent_id?: string
  profile?: string
  batch_code?: string
  date_from?: string
  date_to?: string
  limit?: number
  offset?: number
}): Promise<{ sales: HotspotSaleResponse[]; meta?: { total: number; limit: number; offset: number } }> {
  const response = await adminClient.get('/hotspot-sales', { params })
  const parsed = HotspotSaleListResponseSchema.parse(response.data)
  const meta = parsed.meta
  return { sales: parsed.data, meta }
}

export async function createHotspotSale(
  data: Record<string, unknown>
): Promise<HotspotSaleResponse> {
  const response = await adminClient.post('/hotspot-sales', data)
  const parsed = HotspotSaleDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function listRouterHotspotSales(
  routerId: string,
  params?: {
    profile?: string
    batch_code?: string
    agent_id?: string
    limit?: number
    offset?: number
  }
): Promise<HotspotSaleResponse[]> {
  const response = await adminClient.get(`/routers/${routerId}/hotspot-sales`, { params })
  const parsed = HotspotSaleListResponseSchema.parse(response.data)
  return parsed.data
}
