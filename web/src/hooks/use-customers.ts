import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
    listCustomers,
    createCustomer,
    activateCustomerAccount,
    deactivateCustomerAccount,
    deleteCustomer,
    updateCustomer,
    listRegistrations,
    approveRegistration,
    rejectRegistration,
} from '@/api/customer'
import { toast } from 'sonner'

export function useCustomers(limit?: number, offset?: number) {
    return useQuery({
        queryKey: ['customers', limit, offset],
        queryFn: () => listCustomers(limit, offset),
        staleTime: 2 * 60 * 1000,
    })
}

export function useCreateCustomer() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (data: Record<string, unknown>) => createCustomer(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['customers'] })
            toast.success('Customer created successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to create customer'
            toast.error(message)
        },
    })
}

export function useActivateCustomer() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => activateCustomerAccount(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['customers'] })
            toast.success('Customer account activated')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to activate account'
            toast.error(message)
        },
    })
}

export function useDeactivateCustomer() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => deactivateCustomerAccount(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['customers'] })
            toast.success('Customer account deactivated')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to deactivate account'
            toast.error(message)
        },
    })
}

export function useDeleteCustomer() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => deleteCustomer(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['customers'] })
            toast.success('Customer deleted successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to delete customer'
            toast.error(message)
        },
    })
}

export function useUpdateCustomer() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ id, data }: { id: string; data: Record<string, unknown> }) => updateCustomer(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['customers'] })
            toast.success('Customer updated successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to update customer'
            toast.error(message)
        },
    })
}

// ── Registrations ──

export function useRegistrations(limit?: number, offset?: number) {
    return useQuery({
        queryKey: ['registrations', limit, offset],
        queryFn: () => listRegistrations(limit, offset),
        staleTime: 2 * 60 * 1000,
    })
}

export function useApproveRegistration() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({
            id,
            data,
        }: {
            id: string
            data: { router_id: string; profile_id?: string }
        }) => approveRegistration(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['registrations'] })
            queryClient.invalidateQueries({ queryKey: ['customers'] })
            toast.success('Registration approved')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to approve registration'
            toast.error(message)
        },
    })
}

export function useRejectRegistration() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({
            id,
            data,
        }: {
            id: string
            data: { reason: string }
        }) => rejectRegistration(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['registrations'] })
            toast.success('Registration rejected')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to reject registration'
            toast.error(message)
        },
    })
}
