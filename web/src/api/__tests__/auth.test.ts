import { vi, describe, it, expect, beforeEach } from 'vitest'
import axios from 'axios'

vi.mock('axios', () => ({
  default: { post: vi.fn(), get: vi.fn() },
}))
vi.mock('@/lib/axios/admin-client', () => ({
  adminClient: { post: vi.fn(), get: vi.fn() },
}))

import { adminClient } from '@/lib/axios/admin-client'
import {
  adminLogin,
  adminRefreshToken,
  adminChangePassword,
  adminLogout,
  adminGetMe,
  customerLogin,
  agentLogin,
} from '@/api/auth'

const mockedAxios = vi.mocked(axios)
const mockedClient = vi.mocked(adminClient)

beforeEach(() => vi.clearAllMocks())

describe('adminLogin', () => {
  it('POST /api/v1/auth/login and returns tokens + user', async () => {
    mockedAxios.post.mockResolvedValueOnce({
      data: {
        success: true,
        data: {
          access_token: 'acc',
          refresh_token: 'ref',
          user: {
            id: 'u1', full_name: 'Admin', email: 'a@b.com', phone: '08',
            role: 'admin', is_active: true, last_login: null,
            created_at: '2024-01-01', updated_at: '2024-01-01',
          },
        },
      },
    })
    const result = await adminLogin('a@b.com', 'pass')
    expect(mockedAxios.post).toHaveBeenCalledWith('/api/v1/auth/login', { email: 'a@b.com', password: 'pass' })
    expect(result.access_token).toBe('acc')
  })
})

describe('adminRefreshToken', () => {
  it('POST /api/v1/auth/refresh and returns new token', async () => {
    mockedAxios.post.mockResolvedValueOnce({
      data: { success: true, data: { token: 'new-acc', refresh_token: 'new-ref' } },
    })
    const result = await adminRefreshToken('old-ref')
    expect(result.token).toBe('new-acc')
  })
})

describe('adminChangePassword', () => {
  it('POST /auth/change-password', async () => {
    mockedClient.post.mockResolvedValueOnce({
      data: { success: true, data: { message: 'ok' } },
    })
    await adminChangePassword('old', 'new123456')
    expect(mockedClient.post).toHaveBeenCalledWith('/auth/change-password', {
      old_password: 'old', new_password: 'new123456',
    })
  })
})

describe('adminLogout', () => {
  it('POST /auth/logout', async () => {
    mockedClient.post.mockResolvedValueOnce({
      data: { success: true, data: { message: 'logged out' } },
    })
    await adminLogout()
    expect(mockedClient.post).toHaveBeenCalledWith('/auth/logout')
  })
})

describe('adminGetMe', () => {
  it('GET /auth/me', async () => {
    mockedClient.get.mockResolvedValueOnce({ data: { success: true, data: {} } })
    await adminGetMe()
    expect(mockedClient.get).toHaveBeenCalledWith('/auth/me')
  })
})

describe('customerLogin', () => {
  it('POST /portal/v1/login', async () => {
    mockedAxios.post.mockResolvedValueOnce({
      data: {
        success: true,
        data: {
          token: 'portal-tok',
          customer: {
            id: 'c1', customer_code: 'C001', full_name: 'Cust',
            username: 'cust1', email: 'c@b.com', phone: '08',
          },
        },
      },
    })
    const result = await customerLogin('cust1', 'pass')
    expect(result.token).toBe('portal-tok')
  })
})

describe('agentLogin', () => {
  it('POST /agent-portal/v1/login', async () => {
    mockedAxios.post.mockResolvedValueOnce({
      data: {
        success: true,
        data: {
          token: 'agent-tok',
          agent: { id: 'a1', name: 'Agent', phone: '08', username: 'agt1', status: 'active' },
        },
      },
    })
    const result = await agentLogin('agt1', 'pass')
    expect(result.token).toBe('agent-tok')
  })
})
