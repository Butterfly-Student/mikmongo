import { useQuery } from '@tanstack/react-query'
import { getReportSummary } from '@/api/report'

export function useReportSummary(from?: string, to?: string) {
    return useQuery({
        queryKey: ['report-summary', from, to],
        queryFn: () => getReportSummary(from, to),
        staleTime: 5 * 60 * 1000,
    })
}
