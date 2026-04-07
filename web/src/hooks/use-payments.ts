import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listPayments,
  confirmPayment,
  rejectPayment,
  refundPayment,
  initiateGatewayPayment,
} from '@/api/payment'
import { toast } from 'sonner'

export function usePayments(limit?: number, offset?: number) {
  return useQuery({
    queryKey: ['payments', limit, offset],
    queryFn: () => listPayments(limit, offset),
    staleTime: 2 * 60 * 1000,
  })
}

export function useConfirmPayment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => confirmPayment(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] })
      toast.success('Pembayaran berhasil dikonfirmasi')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal mengkonfirmasi pembayaran'
      toast.error(message)
    },
  })
}

export function useRejectPayment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) =>
      rejectPayment(id, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] })
      toast.success('Pembayaran ditolak')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal menolak pembayaran'
      toast.error(message)
    },
  })
}

export function useRefundPayment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      id,
      amount,
      reason,
    }: {
      id: string
      amount: number
      reason: string
    }) => refundPayment(id, { amount, reason }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] })
      toast.success('Refund berhasil diproses')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal memproses refund'
      toast.error(message)
    },
  })
}

export function useInitiateGatewayPayment() {
  return useMutation({
    mutationFn: ({ id, gateway }: { id: string; gateway: string }) =>
      initiateGatewayPayment(id, gateway),
    onSuccess: (data) => {
      window.open(data.payment_url, '_blank')
      toast.success('Halaman pembayaran dibuka di tab baru')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal memuat halaman pembayaran'
      toast.error(message)
    },
  })
}
