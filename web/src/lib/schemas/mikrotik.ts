import { z } from 'zod'
import { ApiResponseSchema } from '@/lib/schemas/auth'

// ── PPP ──────────────────────────────────────────────────────────────

export const PppProfileSchema = z.object({
  '.id': z.string().optional(),
  name: z.string(),
  localAddress: z.string().nullish(),
  remoteAddress: z.string().nullish(),
  rateLimit: z.string().nullish(),
  dnsServer: z.string().nullish(),
  sessionTimeout: z.string().nullish(),
  idleTimeout: z.string().nullish(),
  parentQueue: z.string().nullish(),
  queueType: z.string().nullish(),
  onlyOne: z.boolean().nullish(),
  useCompression: z.boolean().nullish(),
  useEncryption: z.boolean().nullish(),
  changeTCPMSS: z.boolean().nullish(),
  bridge: z.string().nullish(),
  addressList: z.string().nullish(),
  comment: z.string().nullish(),
}).passthrough()

export type PppProfile = z.infer<typeof PppProfileSchema>

export const AddPppProfileRequestSchema = z.object({
  name: z.string().min(1),
  local_address: z.string().optional(),
  remote_address: z.string().optional(),
  rate_limit: z.string().optional(),
  only_one: z.string().optional(),
  comment: z.string().optional(),
})

export type AddPppProfileRequest = z.infer<typeof AddPppProfileRequestSchema>

export const PppSecretSchema = z.object({
  '.id': z.string().optional(),
  name: z.string(),
  service: z.string().nullish(),
  callerID: z.string().nullish(),
  password: z.string().nullish(),
  profile: z.string().nullish(),
  routes: z.string().nullish(),
  localAddress: z.string().nullish(),
  remoteAddress: z.string().nullish(),
  limitBytesIn: z.number().nullish(),
  limitBytesOut: z.number().nullish(),
  lastLoggedOut: z.string().nullish(),
  comment: z.string().nullish(),
  lastCallerID: z.string().nullish(),
  lastDisconnectReason: z.string().nullish(),
  disabled: z.boolean().nullish(),
}).passthrough()

export type PppSecret = z.infer<typeof PppSecretSchema>

export const AddPppSecretRequestSchema = z.object({
  name: z.string().min(1),
  password: z.string().min(1),
  profile: z.string().optional(),
  service: z.string().optional(),
  caller_id: z.string().optional(),
  local_address: z.string().optional(),
  remote_address: z.string().optional(),
  routes: z.string().optional(),
  comment: z.string().optional(),
  disabled: z.boolean().optional(),
})

export type AddPppSecretRequest = z.infer<typeof AddPppSecretRequestSchema>

export const PppActiveSchema = z.object({
  name: z.string(),
  service: z.string().nullish(),
  callerID: z.string().nullish(),
  encoding: z.string().nullish(),
  address: z.string().nullish(),
  uptime: z.string().nullish(),
}).passthrough()

export type PppActive = z.infer<typeof PppActiveSchema>

// ── Hotspot ──────────────────────────────────────────────────────────

export const HotspotProfileSchema = z.object({
  '.id': z.string().optional(),
  name: z.string(),
  addressPool: z.string().nullish(),
  sharedUsers: z.number().nullish(),
  rateLimit: z.string().nullish(),
  parentQueue: z.string().nullish(),
  sessionTimeout: z.string().nullish(),
  idleTimeout: z.string().nullish(),
  keepaliveTimeout: z.string().nullish(),
  statusAutorefresh: z.string().nullish(),
  onLogin: z.string().nullish(),
  onLogout: z.string().nullish(),
  transparentProxy: z.boolean().nullish(),
  advertise: z.boolean().nullish(),
}).passthrough()

export type HotspotProfile = z.infer<typeof HotspotProfileSchema>

export const AddHotspotProfileRequestSchema = z.object({
  name: z.string().min(1),
  shared_users: z.string().optional(),
  rate_limit: z.string().optional(),
  expired_mode: z.string().optional(),
  validity_time: z.string().optional(),
  keepalive_timeout: z.string().optional(),
  status_autorefresh: z.string().optional(),
  on_login: z.string().optional(),
  on_logout: z.string().optional(),
})

export type AddHotspotProfileRequest = z.infer<typeof AddHotspotProfileRequestSchema>

export const HotspotUserSchema = z.object({
  '.id': z.string().optional(),
  name: z.string(),
  password: z.string().nullish(),
  profile: z.string().nullish(),
  server: z.string().nullish(),
  address: z.string().nullish(),
  macAddress: z.string().nullish(),
  limitUptime: z.string().nullish(),
  limitBytesTotal: z.number().nullish(),
  comment: z.string().nullish(),
  disabled: z.boolean().nullish(),
  uptime: z.string().nullish(),
  bytesIn: z.number().nullish(),
  bytesOut: z.number().nullish(),
}).passthrough()

export type HotspotUser = z.infer<typeof HotspotUserSchema>

export const AddHotspotUserRequestSchema = z.object({
  name: z.string().min(1),
  password: z.string().optional(),
  profile: z.string().optional(),
  server: z.string().optional(),
  limit_uptime: z.string().optional(),
  limit_bytes_total: z.string().optional(),
  comment: z.string().optional(),
  disabled: z.boolean().optional(),
})

export type AddHotspotUserRequest = z.infer<typeof AddHotspotUserRequestSchema>

export const HotspotActiveSchema = z.object({
  server: z.string().nullish(),
  user: z.string().nullish(),
  domain: z.string().nullish(),
  address: z.string().nullish(),
  macAddress: z.string().nullish(),
  loginBy: z.string().nullish(),
  uptime: z.string().nullish(),
  idleTime: z.string().nullish(),
  sessionTimeLeft: z.string().nullish(),
  bytesIn: z.number().nullish(),
  bytesOut: z.number().nullish(),
}).passthrough()

export type HotspotActive = z.infer<typeof HotspotActiveSchema>

export const HotspotHostSchema = z.object({
  '.id': z.string().optional(),
  macAddress: z.string().nullish(),
  address: z.string().nullish(),
  toAddress: z.string().nullish(),
  server: z.string().nullish(),
  uptime: z.string().nullish(),
  bytesIn: z.number().nullish(),
  bytesOut: z.number().nullish(),
}).passthrough()

export type HotspotHost = z.infer<typeof HotspotHostSchema>

export const HotspotServerSchema = z.object({
  '.id': z.string().optional(),
  name: z.string(),
  interface: z.string().nullish(),
  addressPool: z.string().nullish(),
  profile: z.string().nullish(),
  disabled: z.boolean().nullish(),
}).passthrough()

export type HotspotServer = z.infer<typeof HotspotServerSchema>

// ── Network ──────────────────────────────────────────────────────────

export const SimpleQueueSchema = z.object({
  id: z.string().optional(),
  name: z.string(),
  target: z.string().nullish(),
  dst: z.string().nullish(),
  maxLimit: z.string().nullish(),
  limitAt: z.string().nullish(),
  burstLimit: z.string().nullish(),
  burstThreshold: z.string().nullish(),
  burstTime: z.string().nullish(),
  priority: z.string().nullish(),
  parent: z.string().nullish(),
  comment: z.string().nullish(),
  disabled: z.boolean().nullish(),
  dynamic: z.boolean().nullish(),
}).passthrough()

export type SimpleQueue = z.infer<typeof SimpleQueueSchema>

export const FirewallRuleSchema = z.object({
  id: z.string().optional(),
  chain: z.string().nullish(),
  action: z.string().nullish(),
  protocol: z.string().nullish(),
  srcAddress: z.string().nullish(),
  dstAddress: z.string().nullish(),
  srcPort: z.string().nullish(),
  dstPort: z.string().nullish(),
  inInterface: z.string().nullish(),
  outInterface: z.string().nullish(),
  comment: z.string().nullish(),
  disabled: z.boolean().nullish(),
}).passthrough()

export type FirewallRule = z.infer<typeof FirewallRuleSchema>

export const NatRuleSchema = z.object({
  id: z.string().optional(),
  chain: z.string().nullish(),
  action: z.string().nullish(),
  protocol: z.string().nullish(),
  srcAddress: z.string().nullish(),
  dstAddress: z.string().nullish(),
  srcPort: z.string().nullish(),
  dstPort: z.string().nullish(),
  toAddresses: z.string().nullish(),
  toPorts: z.string().nullish(),
  comment: z.string().nullish(),
  disabled: z.boolean().nullish(),
  bytes: z.number().nullish(),
  packets: z.number().nullish(),
}).passthrough()

export type NatRule = z.infer<typeof NatRuleSchema>

export const AddressListEntrySchema = z.object({
  id: z.string().optional(),
  list: z.string().nullish(),
  address: z.string().nullish(),
  timeout: z.string().nullish(),
  comment: z.string().nullish(),
  disabled: z.boolean().nullish(),
}).passthrough()

export type AddressListEntry = z.infer<typeof AddressListEntrySchema>

export const IpPoolSchema = z.object({
  id: z.string().optional(),
  name: z.string(),
  ranges: z.string().nullish(),
  nextPool: z.string().nullish(),
  comment: z.string().nullish(),
}).passthrough()

export type IpPool = z.infer<typeof IpPoolSchema>

export const AddIpPoolRequestSchema = z.object({
  name: z.string().min(1),
  ranges: z.string().min(1),
  comment: z.string().optional(),
})

export type AddIpPoolRequest = z.infer<typeof AddIpPoolRequestSchema>

export const IpAddressSchema = z.object({
  id: z.string().optional(),
  address: z.string().nullish(),
  network: z.string().nullish(),
  interface: z.string().nullish(),
  disabled: z.boolean().nullish(),
  comment: z.string().nullish(),
}).passthrough()

export type IpAddress = z.infer<typeof IpAddressSchema>

// ── Monitor ──────────────────────────────────────────────────────────

export const SystemResourceSchema = z.object({
  uptime: z.string().nullish(),
  version: z.string().nullish(),
  buildTime: z.string().nullish(),
  freeMemory: z.number().nullish(),
  totalMemory: z.number().nullish(),
  freeHddSpace: z.number().nullish(),
  totalHddSpace: z.number().nullish(),
  architectureName: z.string().nullish(),
  boardName: z.string().nullish(),
  platform: z.string().nullish(),
  cpu: z.string().nullish(),
  cpuCount: z.number().nullish(),
  cpuFrequency: z.number().nullish(),
  cpuLoad: z.number().nullish(),
}).passthrough()

export type SystemResource = z.infer<typeof SystemResourceSchema>

export const InterfaceStatsSchema = z.object({
  id: z.string().optional(),
  name: z.string(),
  type: z.string().nullish(),
  mtu: z.number().nullish(),
  macAddress: z.string().nullish(),
  running: z.boolean().nullish(),
  disabled: z.boolean().nullish(),
  comment: z.string().nullish(),
}).passthrough()

export type InterfaceStats = z.infer<typeof InterfaceStatsSchema>

// ── Raw command ───────────────────────────────────────────────────────

export const RawCommandRequestSchema = z.object({
  args: z.array(z.string()).min(1).max(20),
})

export type RawCommandRequest = z.infer<typeof RawCommandRequestSchema>

// ── Response wrappers ─────────────────────────────────────────────────

export const PppProfileListResponseSchema = ApiResponseSchema(z.array(PppProfileSchema))
export const PppProfileDetailResponseSchema = ApiResponseSchema(PppProfileSchema)
export const PppSecretListResponseSchema = ApiResponseSchema(z.array(PppSecretSchema))
export const PppSecretDetailResponseSchema = ApiResponseSchema(PppSecretSchema)
export const PppActiveListResponseSchema = ApiResponseSchema(z.array(PppActiveSchema))

export const HotspotProfileListResponseSchema = ApiResponseSchema(z.array(HotspotProfileSchema))
export const HotspotProfileDetailResponseSchema = ApiResponseSchema(HotspotProfileSchema)
export const HotspotUserListResponseSchema = ApiResponseSchema(z.array(HotspotUserSchema))
export const HotspotUserDetailResponseSchema = ApiResponseSchema(HotspotUserSchema)
export const HotspotActiveListResponseSchema = ApiResponseSchema(z.array(HotspotActiveSchema))
export const HotspotHostListResponseSchema = ApiResponseSchema(z.array(HotspotHostSchema))
export const HotspotServerListResponseSchema = ApiResponseSchema(z.array(HotspotServerSchema))

export const SimpleQueueListResponseSchema = ApiResponseSchema(z.array(SimpleQueueSchema))
export const FirewallRuleListResponseSchema = ApiResponseSchema(z.array(FirewallRuleSchema))
export const NatRuleListResponseSchema = ApiResponseSchema(z.array(NatRuleSchema))
export const AddressListResponseSchema = ApiResponseSchema(z.array(AddressListEntrySchema))
export const IpPoolListResponseSchema = ApiResponseSchema(z.array(IpPoolSchema))
export const IpPoolDetailResponseSchema = ApiResponseSchema(IpPoolSchema)
export const IpAddressListResponseSchema = ApiResponseSchema(z.array(IpAddressSchema))

export const SystemResourceResponseSchema = ApiResponseSchema(SystemResourceSchema)
export const InterfaceStatsListResponseSchema = ApiResponseSchema(z.array(InterfaceStatsSchema))

export const RawCommandResponseSchema = ApiResponseSchema(z.array(z.record(z.string(), z.unknown())))
