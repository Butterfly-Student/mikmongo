import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listRouters, selectRouter } from '@/api/router'
import { useRouterStore } from '@/stores/router-store'

export function useRouters() {
    return useQuery({
        queryKey: ['routers'],
        queryFn: () => listRouters(),
        staleTime: 2 * 60 * 1000,
    })
}

export function useSelectRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => selectRouter(id),
        onSuccess: (router) => {
            useRouterStore.getState().setSelectedRouter(router.id, router.name)
            queryClient.invalidateQueries({ queryKey: ['routers'] })
            queryClient.invalidateQueries({ queryKey: ['report-summary'] })
        },
    })
}

import { createRouter, syncRouter, testConnection, updateRouter, deleteRouter, syncAllRouters } from '@/api/router'
import { toast } from 'sonner'
import type { CreateRouter } from '@/features/routers/data/schema'

export function useCreateRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (data: CreateRouter) => createRouter(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['routers'] })
            toast.success("Router created successfully")
        },
        onError: (err: any) => {
            toast.error(err.response?.data?.error || "Failed to create router")
        }
    })
}

export function useSyncRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => syncRouter(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['routers'] })
            toast.success("Router synced successfully")
        },
        onError: (err: any) => {
            toast.error(err.response?.data?.error || "Failed to sync router")
        }
    })
}

export function useTestRouterConnection() {
    return useMutation({
        mutationFn: (id: string) => testConnection(id),
        onSuccess: () => {
            toast.success("Connection test passed")
        },
        onError: (err: unknown) => {
            const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? "Connection test failed"
            toast.error(message)
        }
    })
}

export function useUpdateRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ id, data }: { id: string; data: Record<string, unknown> }) => updateRouter(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['routers'] })
            toast.success("Router updated successfully")
        },
        onError: (err: unknown) => {
            const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? "Failed to update router"
            toast.error(message)
        }
    })
}

export function useDeleteRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => deleteRouter(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['routers'] })
            toast.success("Router deleted successfully")
        },
        onError: (err: unknown) => {
            const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? "Failed to delete router"
            toast.error(message)
        }
    })
}

export function useSyncAllRouters() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: () => syncAllRouters(),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['routers'] })
            toast.success("All routers synced successfully")
        },
        onError: (err: unknown) => {
            const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? "Failed to sync all routers"
            toast.error(message)
        }
    })
}
