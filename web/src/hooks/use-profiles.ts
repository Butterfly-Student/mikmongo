import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listProfiles, createProfile, deleteProfile, updateProfile } from '@/api/profiles'
import { toast } from 'sonner'
import type { CreateProfile } from '@/features/profiles/data/schema'

export function useProfiles(routerId: string | null, limit?: number, offset?: number) {
    return useQuery({
        queryKey: ['profiles', routerId, limit, offset],
        queryFn: () => {
            if (!routerId) throw new Error("No router selected")
            return listProfiles(routerId, limit, offset)
        },
        enabled: !!routerId,
        staleTime: 2 * 60 * 1000,
    })
}

export function useCreateProfile() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ routerId, data }: { routerId: string; data: CreateProfile }) => createProfile(routerId, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['profiles', variables.routerId] })
            toast.success("Bandwidth profile created successfully")
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to create profile'
            toast.error(message)
        },
    })
}

export function useDeleteProfile() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ routerId, id }: { routerId: string; id: string }) => deleteProfile(routerId, id),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['profiles', variables.routerId] })
            toast.success("Bandwidth profile deleted successfully")
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to delete profile'
            toast.error(message)
        },
    })
}

export function useUpdateProfile() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ routerId, id, data }: { routerId: string; id: string; data: Partial<CreateProfile> }) =>
            updateProfile(routerId, id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['profiles', variables.routerId] })
            toast.success("Bandwidth profile updated successfully")
        },
        onError: (err: unknown) => {
            const message =
                (err as { response?: { data?: { error?: string } } })?.response?.data
                    ?.error ?? 'Failed to update profile'
            toast.error(message)
        },
    })
}
