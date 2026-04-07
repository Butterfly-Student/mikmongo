import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
    listSubscriptions,
    createSubscription,
    activateSubscription,
    suspendSubscription,
    isolateSubscription,
    restoreSubscription,
    terminateSubscription,
    deleteSubscription,
} from '@/api/subscription'
import { toast } from 'sonner'

export function useSubscriptions(routerId: string | null, limit?: number, offset?: number) {
    return useQuery({
        queryKey: ['subscriptions', routerId, limit, offset],
        queryFn: () => {
            if (!routerId) throw new Error('No router selected')
            return listSubscriptions(routerId, limit, offset)
        },
        enabled: !!routerId,
        staleTime: 2 * 60 * 1000,
    })
}

export function useCreateSubscription(routerId: string) {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (data: Record<string, unknown>) =>
            createSubscription(routerId, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['subscriptions', routerId] })
            toast.success('Subscription created successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to create subscription'
            toast.error(message)
        },
    })
}

export function useActivateSubscription(routerId: string) {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => activateSubscription(routerId, id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['subscriptions', routerId] })
            toast.success('Subscription activated')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to activate subscription'
            toast.error(message)
        },
    })
}

export function useSuspendSubscription(routerId: string) {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ id, reason }: { id: string; reason?: string }) =>
            suspendSubscription(routerId, id, reason),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['subscriptions', routerId] })
            toast.success('Subscription suspended')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to suspend subscription'
            toast.error(message)
        },
    })
}

export function useIsolateSubscription(routerId: string) {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ id, reason }: { id: string; reason?: string }) =>
            isolateSubscription(routerId, id, reason),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['subscriptions', routerId] })
            toast.success('Subscription isolated')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to isolate subscription'
            toast.error(message)
        },
    })
}

export function useRestoreSubscription(routerId: string) {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => restoreSubscription(routerId, id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['subscriptions', routerId] })
            toast.success('Subscription restored')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to restore subscription'
            toast.error(message)
        },
    })
}

export function useTerminateSubscription(routerId: string) {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => terminateSubscription(routerId, id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['subscriptions', routerId] })
            toast.success('Subscription terminated')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to terminate subscription'
            toast.error(message)
        },
    })
}

export function useDeleteSubscription(routerId: string) {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => deleteSubscription(routerId, id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['subscriptions', routerId] })
            toast.success('Subscription deleted successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to delete subscription'
            toast.error(message)
        },
    })
}
