import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
    generateVouchers, listVouchers, removeVoucherBatch,
    createMikhmonProfile, updateMikhmonProfile,
    listMikhmonReports, getMikhmonReportSummary,
    setupExpiration, disableExpiration, getExpirationStatus, generateExpirationScript
} from '@/api/mikrotik/mikhmon'
import type {
    GenerateVoucherRequest, CreateMikhmonProfileRequest, UpdateMikhmonProfileRequest, GenerateScriptRequest
} from '@/lib/schemas/mikhmon'

// -- Vouchers
export function useMikhmonVouchers(routerId: string | null, comment: string) {
    return useQuery({
        queryKey: ['mikhmon-vouchers', routerId, comment],
        queryFn: () => {
            if (!routerId) throw new Error("No router selected")
            return listVouchers(routerId, comment)
        },
        enabled: !!routerId,
    })
}

export function useGenerateMikhmonVouchers() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, data }: { routerId: string; data: GenerateVoucherRequest }) => generateVouchers(routerId, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['mikhmon-vouchers', variables.routerId] })
            toast.success("Voucher batch generated successfully")
        },
        onError: (err: unknown) => {
            const message = (err as any)?.response?.data?.error ?? 'Failed to generate vouchers'
            toast.error(message)
        },
    })
}

export function useRemoveMikhmonVoucherBatch() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, comment }: { routerId: string; comment: string }) => removeVoucherBatch(routerId, comment),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['mikhmon-vouchers', variables.routerId] })
            toast.success("Voucher batch deleted successfully")
        },
        onError: (err: unknown) => {
            const message = (err as any)?.response?.data?.error ?? 'Failed to delete voucher batch'
            toast.error(message)
        },
    })
}

// -- Profiles (create/update only — no list endpoint on backend)
export function useCreateMikhmonProfile() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, data }: { routerId: string; data: CreateMikhmonProfileRequest }) => createMikhmonProfile(routerId, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-profiles', variables.routerId] })
            toast.success("Profile created successfully")
        },
        onError: (err: unknown) => {
            const message = (err as any)?.response?.data?.error ?? 'Failed to create profile'
            toast.error(message)
        },
    })
}

export function useUpdateMikhmonProfile() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, id, data }: { routerId: string; id: string; data: UpdateMikhmonProfileRequest }) => updateMikhmonProfile(routerId, id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-profiles', variables.routerId] })
            toast.success("Profile updated successfully")
        },
        onError: (err: unknown) => {
            const message = (err as any)?.response?.data?.error ?? 'Failed to update profile'
            toast.error(message)
        },
    })
}

// -- Reports
export function useMikhmonReports(routerId: string | null) {
    return useQuery({
        queryKey: ['mikhmon-reports', routerId],
        queryFn: () => {
            if (!routerId) throw new Error("No router selected")
            return listMikhmonReports(routerId)
        },
        enabled: !!routerId,
    })
}

export function useMikhmonReportSummary(routerId: string | null) {
    return useQuery({
        queryKey: ['mikhmon-report-summary', routerId],
        queryFn: () => {
            if (!routerId) throw new Error("No router selected")
            return getMikhmonReportSummary(routerId)
        },
        enabled: !!routerId,
    })
}

// -- Expiration Monitor
export function useExpirationStatus(routerId: string | null) {
    return useQuery({
        queryKey: ['mikhmon-expiration-status', routerId],
        queryFn: () => {
            if (!routerId) throw new Error("No router selected")
            return getExpirationStatus(routerId)
        },
        enabled: !!routerId,
        refetchInterval: 30000,
    })
}

export function useSetupExpiration() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, data }: { routerId: string; data: Record<string, unknown> }) => setupExpiration(routerId, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['mikhmon-expiration-status', variables.routerId] })
            toast.success("Expiration monitor setup successfully")
        },
        onError: (err: unknown) => {
            const message = (err as any)?.response?.data?.error ?? 'Failed to setup expiration monitor'
            toast.error(message)
        },
    })
}

export function useDisableExpiration() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId }: { routerId: string }) => disableExpiration(routerId),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['mikhmon-expiration-status', variables.routerId] })
            toast.success("Expiration monitor disabled successfully")
        },
        onError: (err: unknown) => {
            const message = (err as any)?.response?.data?.error ?? 'Failed to disable expiration monitor'
            toast.error(message)
        },
    })
}

export function useGenerateExpirationScript() {
    return useMutation({
        mutationFn: ({ routerId, data }: { routerId: string; data: GenerateScriptRequest }) => generateExpirationScript(routerId, data),
        onError: (err: unknown) => {
            const message = (err as any)?.response?.data?.error ?? 'Failed to generate script'
            toast.error(message)
        },
    })
}
