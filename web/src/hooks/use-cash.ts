import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listCashEntries,
  createCashEntry,
  approveCashEntry,
  rejectCashEntry,
  listPettyCashFunds,
  createPettyCashFund,
  topUpPettyCashFund,
} from '@/api/cash'
import type { CreateCashEntry } from '@/lib/schemas/billing'
import { toast } from 'sonner'

export function useCashEntries(params?: {
  type?: string
  status?: string
  limit?: number
  offset?: number
}) {
  return useQuery({
    queryKey: ['cash-entries', params],
    queryFn: () => listCashEntries(params),
    staleTime: 2 * 60 * 1000,
  })
}

export function useCreateCashEntry() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateCashEntry) => createCashEntry(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cash-entries'] })
      toast.success('Entri kas berhasil dibuat')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal membuat entri kas'
      toast.error(message)
    },
  })
}

export function useApproveCashEntry() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => approveCashEntry(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cash-entries'] })
      toast.success('Entri kas disetujui')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal menyetujui entri kas'
      toast.error(message)
    },
  })
}

export function useRejectCashEntry() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) =>
      rejectCashEntry(id, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cash-entries'] })
      toast.success('Entri kas ditolak')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal menolak entri kas'
      toast.error(message)
    },
  })
}

export function usePettyCashFunds() {
  return useQuery({
    queryKey: ['petty-cash'],
    queryFn: () => listPettyCashFunds(),
    staleTime: 5 * 60 * 1000,
  })
}

export function useCreatePettyCashFund() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: { fund_name: string; initial_balance: number; custodian_id: string }) =>
      createPettyCashFund(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['petty-cash'] })
      toast.success('Dana kecil berhasil dibuat')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal membuat dana kecil'
      toast.error(message)
    },
  })
}

export function useTopUpPettyCashFund() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, amount }: { id: string; amount: number }) =>
      topUpPettyCashFund(id, amount),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['petty-cash'] })
      toast.success('Saldo berhasil ditambahkan')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal menambahkan saldo'
      toast.error(message)
    },
  })
}
