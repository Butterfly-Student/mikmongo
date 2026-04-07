import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  SimpleQueueListResponseSchema,
  FirewallRuleListResponseSchema,
  NatRuleListResponseSchema,
  AddressListResponseSchema,
  IpPoolListResponseSchema,
  IpPoolDetailResponseSchema,
  IpAddressListResponseSchema,
  type SimpleQueue,
  type FirewallRule,
  type NatRule,
  type AddressListEntry,
  type IpPool,
  type AddIpPoolRequest,
  type IpAddress,
} from '@/lib/schemas/mikrotik'

export async function listQueues(routerId: string): Promise<SimpleQueue[]> {
  const response = await adminClient.get(`/routers/${routerId}/queue/simple`)
  const parsed = SimpleQueueListResponseSchema.parse(response.data)
  return parsed.data
}

export async function listFirewallFilters(routerId: string): Promise<FirewallRule[]> {
  const response = await adminClient.get(`/routers/${routerId}/firewall/filter`)
  const parsed = FirewallRuleListResponseSchema.parse(response.data)
  return parsed.data
}

export async function listFirewallNat(routerId: string): Promise<NatRule[]> {
  const response = await adminClient.get(`/routers/${routerId}/firewall/nat`)
  const parsed = NatRuleListResponseSchema.parse(response.data)
  return parsed.data
}

export async function listFirewallAddressList(routerId: string): Promise<AddressListEntry[]> {
  const response = await adminClient.get(`/routers/${routerId}/firewall/address-list`)
  const parsed = AddressListResponseSchema.parse(response.data)
  return parsed.data
}

// ── IP Pools ──────────────────────────────────────────────────────────

export async function listIpPools(routerId: string): Promise<IpPool[]> {
  const response = await adminClient.get(`/routers/${routerId}/ip/pools`)
  const parsed = IpPoolListResponseSchema.parse(response.data)
  return parsed.data
}

export async function createIpPool(routerId: string, data: AddIpPoolRequest): Promise<IpPool> {
  const response = await adminClient.post(`/routers/${routerId}/ip/pools`, data)
  const parsed = IpPoolDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getIpPool(routerId: string, id: string): Promise<IpPool> {
  const response = await adminClient.get(`/routers/${routerId}/ip/pools/${id}`)
  const parsed = IpPoolDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updateIpPool(
  routerId: string,
  id: string,
  data: Partial<AddIpPoolRequest>
): Promise<IpPool> {
  const response = await adminClient.put(`/routers/${routerId}/ip/pools/${id}`, data)
  const parsed = IpPoolDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function deleteIpPool(routerId: string, id: string): Promise<string> {
  const response = await adminClient.delete(`/routers/${routerId}/ip/pools/${id}`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function listIpAddresses(routerId: string): Promise<IpAddress[]> {
  const response = await adminClient.get(`/routers/${routerId}/ip/addresses`)
  const parsed = IpAddressListResponseSchema.parse(response.data)
  return parsed.data
}
