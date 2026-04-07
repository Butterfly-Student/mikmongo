import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listUsers, getUser, createUser, deleteUser } from '@/api/user'
import type { CreateUserRequest } from '@/api/types'

export function useUsers(limit: number, offset: number) {
    return useQuery({
        queryKey: ['users', limit, offset],
        queryFn: () => listUsers(limit, offset),
        staleTime: 2 * 60 * 1000,
    })
}

export function useUser(id: string) {
    return useQuery({
        queryKey: ['users', id],
        queryFn: () => getUser(id),
        enabled: !!id,
    })
}

export function useCreateUser() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (data: CreateUserRequest) => createUser(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['users'] })
        },
    })
}

export function useDeleteUser() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => deleteUser(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['users'] })
        },
    })
}
