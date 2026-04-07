import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listInvoices,
  listOverdueInvoices,
  getInvoice,
  triggerMonthlyBilling,
} from '@/api/invoice'
import { toast } from 'sonner'

export function useInvoices(limit?: number, offset?: number) {
  return useQuery({
    queryKey: ['invoices', limit, offset],
    queryFn: () => listInvoices(limit, offset),
    staleTime: 2 * 60 * 1000,
  })
}

export function useOverdueInvoices() {
  return useQuery({
    queryKey: ['invoices', 'overdue'],
    queryFn: () => listOverdueInvoices(),
    staleTime: 2 * 60 * 1000,
  })
}

export function useInvoice(id: string | null) {
  return useQuery({
    queryKey: ['invoices', id],
    queryFn: () => getInvoice(id!),
    enabled: !!id,
    staleTime: 2 * 60 * 1000,
  })
}

export function useTriggerMonthlyBilling() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => triggerMonthlyBilling(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] })
      toast.success('Tagihan bulanan berhasil dibuat')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal membuat tagihan bulanan'
      toast.error(message)
    },
  })
}
