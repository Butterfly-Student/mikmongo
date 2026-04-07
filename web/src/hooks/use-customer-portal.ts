import { useQuery } from '@tanstack/react-query'
import { listPortalSubscriptions } from '@/api/portal/subscription'

export function usePortalSubscriptions() {
    return useQuery({
        queryKey: ['portal-subscriptions'],
        queryFn: () => listPortalSubscriptions(),
        staleTime: 2 * 60 * 1000,
    })
}
