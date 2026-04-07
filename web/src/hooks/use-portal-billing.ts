import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listPortalInvoices, getPortalInvoice } from '@/api/portal/invoice'
import { listPortalPayments, initiatePortalPayment } from '@/api/portal/payment'
import {
  listAgentPortalInvoices,
  requestAgentPayment,
} from '@/api/portal/agent-invoice'
import { toast } from 'sonner'

export function usePortalInvoices() {
  return useQuery({
    queryKey: ['portal-invoices'],
    queryFn: () => listPortalInvoices(),
    staleTime: 2 * 60 * 1000,
  })
}

export function usePortalInvoice(id: string | null) {
  return useQuery({
    queryKey: ['portal-invoices', id],
    queryFn: () => getPortalInvoice(id!),
    enabled: !!id,
    staleTime: 2 * 60 * 1000,
  })
}

export function usePortalPayments() {
  return useQuery({
    queryKey: ['portal-payments'],
    queryFn: () => listPortalPayments(),
    staleTime: 2 * 60 * 1000,
  })
}

export function usePortalInitiatePayment() {
  return useMutation({
    mutationFn: (id: string) => initiatePortalPayment(id),
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

export function useAgentPortalInvoices(limit?: number, offset?: number) {
  return useQuery({
    queryKey: ['agent-invoices', limit, offset],
    queryFn: () => listAgentPortalInvoices(limit, offset),
    staleTime: 2 * 60 * 1000,
  })
}

export function useAgentRequestPayment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string
      data?: { paid_amount?: number; notes?: string }
    }) => requestAgentPayment(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['agent-invoices'] })
      toast.success('Permintaan pembayaran berhasil diajukan')
    },
    onError: (err: unknown) => {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data
          ?.error ?? 'Gagal mengajukan pembayaran'
      toast.error(message)
    },
  })
}
