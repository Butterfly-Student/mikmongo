import { adminClient } from '@/lib/axios/admin-client'
import {
  SystemResourceResponseSchema,
  InterfaceStatsListResponseSchema,
  type SystemResource,
  type InterfaceStats,
} from '@/lib/schemas/mikrotik'

export async function getSystemResource(routerId: string): Promise<SystemResource> {
  const response = await adminClient.get(`/routers/${routerId}/monitor/system-resource`)
  const parsed = SystemResourceResponseSchema.parse(response.data)
  return parsed.data
}

export async function getInterfaces(routerId: string): Promise<InterfaceStats[]> {
  const response = await adminClient.get(`/routers/${routerId}/monitor/interfaces`)
  const parsed = InterfaceStatsListResponseSchema.parse(response.data)
  return parsed.data
}
