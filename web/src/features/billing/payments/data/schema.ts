import type { PaymentResponse } from '@/lib/schemas/billing'

export type { PaymentResponse }

export const paymentStatuses = [
  {
    value: 'pending',
    label: 'Menunggu Konfirmasi',
    className: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300',
  },
  {
    value: 'confirmed',
    label: 'Dikonfirmasi',
    className: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
  },
  {
    value: 'rejected',
    label: 'Ditolak',
    className: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
  },
  {
    value: 'refunded',
    label: 'Dikembalikan',
    className: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
  },
] as const

export const paymentMethods = [
  { value: 'cash', label: 'Cash' },
  { value: 'bank_transfer', label: 'Transfer Bank' },
  { value: 'e-wallet', label: 'E-Wallet' },
  { value: 'credit_card', label: 'Kartu Kredit' },
  { value: 'debit_card', label: 'Kartu Debit' },
  { value: 'qris', label: 'QRIS' },
  { value: 'gateway', label: 'Gateway' },
] as const
