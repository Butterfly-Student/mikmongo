import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
    listHotspotProfiles,
    createHotspotProfile,
    updateHotspotProfile,
    deleteHotspotProfile,
    listHotspotUsers,
    createHotspotUser,
    updateHotspotUser,
    deleteHotspotUser,
    listHotspotActive,
    listHotspotHosts,
    listHotspotServers,
} from '@/api/mikrotik/hotspot'
import type { AddHotspotProfileRequest, AddHotspotUserRequest } from '@/lib/schemas/mikrotik'

// ── Hotspot Profiles ──────────────────────────────────────────────────

export function useHotspotProfiles(routerId: string | null) {
    return useQuery({
        queryKey: ['hotspot-profiles', routerId],
        queryFn: () => {
            if (!routerId) throw new Error('No router selected')
            return listHotspotProfiles(routerId)
        },
        enabled: !!routerId,
    })
}

export function useCreateHotspotProfile() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, data }: { routerId: string; data: AddHotspotProfileRequest }) =>
            createHotspotProfile(routerId, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-profiles', variables.routerId] })
            toast.success('Hotspot profile created successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
                'Failed to create hotspot profile'
            toast.error(message)
        },
    })
}

export function useUpdateHotspotProfile() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({
            routerId,
            id,
            data,
        }: {
            routerId: string
            id: string
            data: Partial<AddHotspotProfileRequest>
        }) => updateHotspotProfile(routerId, id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-profiles', variables.routerId] })
            toast.success('Hotspot profile updated successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
                'Failed to update hotspot profile'
            toast.error(message)
        },
    })
}

export function useDeleteHotspotProfile() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, id }: { routerId: string; id: string }) =>
            deleteHotspotProfile(routerId, id),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-profiles', variables.routerId] })
            toast.success('Hotspot profile deleted successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
                'Failed to delete hotspot profile'
            toast.error(message)
        },
    })
}

// ── Hotspot Users ─────────────────────────────────────────────────────

export function useHotspotUsers(routerId: string | null) {
    return useQuery({
        queryKey: ['hotspot-users', routerId],
        queryFn: () => {
            if (!routerId) throw new Error('No router selected')
            return listHotspotUsers(routerId)
        },
        enabled: !!routerId,
    })
}

export function useCreateHotspotUser() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, data }: { routerId: string; data: AddHotspotUserRequest }) =>
            createHotspotUser(routerId, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-users', variables.routerId] })
            toast.success('Hotspot user created successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
                'Failed to create hotspot user'
            toast.error(message)
        },
    })
}

export function useUpdateHotspotUser() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({
            routerId,
            id,
            data,
        }: {
            routerId: string
            id: string
            data: Partial<AddHotspotUserRequest>
        }) => updateHotspotUser(routerId, id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-users', variables.routerId] })
            toast.success('Hotspot user updated successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
                'Failed to update hotspot user'
            toast.error(message)
        },
    })
}

export function useDeleteHotspotUser() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ routerId, id }: { routerId: string; id: string }) =>
            deleteHotspotUser(routerId, id),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['hotspot-users', variables.routerId] })
            queryClient.invalidateQueries({ queryKey: ['hotspot-active', variables.routerId] })
            toast.success('Hotspot user removed successfully')
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
                'Failed to remove hotspot user'
            toast.error(message)
        },
    })
}

// ── Hotspot Active / Hosts / Servers ──────────────────────────────────

export function useHotspotActive(routerId: string | null) {
    return useQuery({
        queryKey: ['hotspot-active', routerId],
        queryFn: () => {
            if (!routerId) throw new Error('No router selected')
            return listHotspotActive(routerId)
        },
        enabled: !!routerId,
        refetchInterval: 30000,
    })
}

export function useHotspotHosts(routerId: string | null) {
    return useQuery({
        queryKey: ['hotspot-hosts', routerId],
        queryFn: () => {
            if (!routerId) throw new Error('No router selected')
            return listHotspotHosts(routerId)
        },
        enabled: !!routerId,
    })
}

export function useHotspotServers(routerId: string | null) {
    return useQuery({
        queryKey: ['hotspot-servers', routerId],
        queryFn: () => {
            if (!routerId) throw new Error('No router selected')
            return listHotspotServers(routerId)
        },
        enabled: !!routerId,
    })
}
