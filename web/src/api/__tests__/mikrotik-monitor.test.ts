import { vi, describe, it, expect, beforeEach } from 'vitest'

vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { get: vi.fn(), post: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import { getSystemResource, getInterfaces } from '@/api/mikrotik/monitor'
import { runRawCommand } from '@/api/mikrotik/raw'

const mc = vi.mocked(adminClient)
const RID = 'router-1'
beforeEach(() => vi.clearAllMocks())

describe('getSystemResource', () => {
  it('GET /routers/:id/monitor/system-resource', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: { cpuLoad: 10, boardName: 'RB750' } } })
    const result = await getSystemResource(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/monitor/system-resource`)
    expect(result.boardName).toBe('RB750')
  })
})

describe('getInterfaces', () => {
  it('GET /routers/:id/monitor/interfaces', async () => {
    mc.get.mockResolvedValueOnce({ data: { success: true, data: [{ name: 'ether1', running: true }] } })
    const result = await getInterfaces(RID)
    expect(mc.get).toHaveBeenCalledWith(`/routers/${RID}/monitor/interfaces`)
    expect(result[0].name).toBe('ether1')
  })
})

describe('runRawCommand', () => {
  it('POST /routers/:id/raw/run', async () => {
    mc.post.mockResolvedValueOnce({ data: { success: true, data: [{ '.id': '*1' }] } })
    const result = await runRawCommand(RID, ['/ip/address/print'])
    expect(mc.post).toHaveBeenCalledWith(`/routers/${RID}/raw/run`, { args: ['/ip/address/print'] })
    expect(result).toHaveLength(1)
  })
})
