import { adminClient } from '@/lib/axios/admin-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  VoucherBatchApiResponseSchema,
  VoucherListResponseSchema,
  MikhmonProfileDetailResponseSchema,
  ScriptApiResponseSchema,
  MikhmonReportListResponseSchema,
  MikhmonReportDetailResponseSchema,
  MikhmonReportSummaryResponseSchema,
  ExpireStatusApiResponseSchema,
  type VoucherBatchResponse,
  type VoucherResponse,
  type GenerateVoucherRequest,
  type MikhmonProfileResponse,
  type CreateMikhmonProfileRequest,
  type UpdateMikhmonProfileRequest,
  type GenerateScriptRequest,
  type ScriptResponse,
  type MikhmonReportResponse,
  type CreateReportRequest,
  type MikhmonReportSummary,
  type ExpireStatusResponse,
} from '@/lib/schemas/mikhmon'

// ── Vouchers ──────────────────────────────────────────────────────────

export async function generateVouchers(
  routerId: string,
  data: GenerateVoucherRequest
): Promise<VoucherBatchResponse> {
  const response = await adminClient.post(`/routers/${routerId}/mikhmon/vouchers/generate`, data)
  const parsed = VoucherBatchApiResponseSchema.parse(response.data)
  return parsed.data
}

export async function listVouchers(
  routerId: string,
  comment?: string,
): Promise<VoucherResponse[]> {
  const params: Record<string, string> = {}
  if (comment) params.comment = comment
  
  const response = await adminClient.get(`/routers/${routerId}/mikhmon/vouchers`, {
    params,
  })
  const parsed = VoucherListResponseSchema.parse(response.data)
  return parsed.data
}

export async function removeVoucherBatch(routerId: string, comment: string): Promise<string> {
  const response = await adminClient.delete(`/routers/${routerId}/mikhmon/vouchers`, {
    params: { comment },
  })
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

// ── Mikhmon Profiles ──────────────────────────────────────────────────

export async function createMikhmonProfile(
  routerId: string,
  data: CreateMikhmonProfileRequest
): Promise<MikhmonProfileResponse> {
  const response = await adminClient.post(`/routers/${routerId}/mikhmon/profiles`, data)
  const parsed = MikhmonProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function updateMikhmonProfile(
  routerId: string,
  id: string,
  data: UpdateMikhmonProfileRequest
): Promise<MikhmonProfileResponse> {
  const response = await adminClient.put(`/routers/${routerId}/mikhmon/profiles/${id}`, data)
  const parsed = MikhmonProfileDetailResponseSchema.parse(response.data)
  return parsed.data
}


// ── Mikhmon Reports ───────────────────────────────────────────────────

export async function listMikhmonReports(routerId: string): Promise<MikhmonReportResponse[]> {
  const response = await adminClient.get(`/routers/${routerId}/mikhmon/reports`)
  const parsed = MikhmonReportListResponseSchema.parse(response.data)
  return parsed.data
}

export async function createMikhmonReport(
  routerId: string,
  data: CreateReportRequest
): Promise<MikhmonReportResponse> {
  const response = await adminClient.post(`/routers/${routerId}/mikhmon/reports`, data)
  const parsed = MikhmonReportDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function getMikhmonReportSummary(routerId: string): Promise<MikhmonReportSummary> {
  const response = await adminClient.get(`/routers/${routerId}/mikhmon/reports/summary`)
  const parsed = MikhmonReportSummaryResponseSchema.parse(response.data)
  return parsed.data
}

// ── Expiration ────────────────────────────────────────────────────────

export async function setupExpiration(
  routerId: string,
  data: Record<string, unknown>
): Promise<string> {
  const response = await adminClient.post(`/routers/${routerId}/mikhmon/expire/setup`, data)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function disableExpiration(routerId: string): Promise<string> {
  const response = await adminClient.post(`/routers/${routerId}/mikhmon/expire/disable`)
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}

export async function getExpirationStatus(routerId: string): Promise<ExpireStatusResponse> {
  const response = await adminClient.get(`/routers/${routerId}/mikhmon/expire/status`)
  const parsed = ExpireStatusApiResponseSchema.parse(response.data)
  return parsed.data
}

export async function generateExpirationScript(
  routerId: string,
  data: GenerateScriptRequest
): Promise<ScriptResponse> {
  const response = await adminClient.post(
    `/routers/${routerId}/mikhmon/expire/generate-script`,
    data
  )
  const parsed = ScriptApiResponseSchema.parse(response.data)
  return parsed.data
}
