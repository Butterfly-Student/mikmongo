import { adminClient } from '@/lib/axios/admin-client'
import { ApiResponseSchema, MessageResponseSchema } from '@/lib/schemas/auth'
import {
    UserListResponseSchema,
    UserResponseSchema,
    CreateUserRequestSchema,
} from '@/lib/schemas/user'
import type { UserResponse, CreateUserRequest } from '@/lib/schemas/user'

export async function listUsers(
    limit: number,
    offset: number
): Promise<{ users: UserResponse[]; meta: { total: number; limit: number; offset: number } }> {
    const response = await adminClient.get('/users', { params: { limit, offset } })
    const parsed = UserListResponseSchema.parse(response.data)
    return { users: parsed.data, meta: parsed.meta }
}

export async function getUser(id: string): Promise<UserResponse> {
    const response = await adminClient.get(`/users/${id}`)
    const parsed = ApiResponseSchema(UserResponseSchema).parse(response.data)
    return parsed.data
}

export async function createUser(data: CreateUserRequest): Promise<UserResponse> {
    // Validate input
    CreateUserRequestSchema.parse(data)
    const response = await adminClient.post('/users', data)
    const parsed = ApiResponseSchema(UserResponseSchema).parse(response.data)
    return parsed.data
}

export async function deleteUser(id: string): Promise<{ message: string }> {
    const response = await adminClient.delete(`/users/${id}`)
    const parsed = MessageResponseSchema.parse(response.data)
    return parsed.data
}
