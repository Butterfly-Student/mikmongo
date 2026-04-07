import { adminClient } from '@/lib/axios/admin-client'
import { RawCommandResponseSchema } from '@/lib/schemas/mikrotik'

export async function runRawCommand(
  routerId: string,
  args: string[]
): Promise<Record<string, unknown>[]> {
  const response = await adminClient.post(`/routers/${routerId}/raw/run`, { args })
  const parsed = RawCommandResponseSchema.parse(response.data)
  return parsed.data
}
