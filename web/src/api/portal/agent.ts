import { agentClient } from '@/lib/axios/agent-client'
import { MessageResponseSchema } from '@/lib/schemas/auth'
import {
  SalesAgentDetailResponseSchema,
  type SalesAgentResponse,
} from '@/lib/schemas/sales-agent'

export async function getAgentProfile(): Promise<SalesAgentResponse> {
  const response = await agentClient.get('/profile')
  const parsed = SalesAgentDetailResponseSchema.parse(response.data)
  return parsed.data
}

export async function changeAgentPassword(password: string): Promise<string> {
  const response = await agentClient.put('/profile/password', { password })
  const parsed = MessageResponseSchema.parse(response.data)
  return parsed.data.message
}
