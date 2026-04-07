import type { InvoiceResponse } from '@/lib/schemas/billing'

export type { InvoiceResponse }

export const invoiceStatuses = [
  {
    value: 'paid',
    label: 'Lunas',
    className: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
  },
  {
    value: 'unpaid',
    label: 'Belum Lunas',
    className: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300',
  },
  {
    value: 'overdue',
    label: 'Terlambat',
    className: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
  },
  {
    value: 'partial',
    label: 'Sebagian',
    className: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
  },
  {
    value: 'draft',
    label: 'Draft',
    className: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
  },
  {
    value: 'sent',
    label: 'Terkirim',
    className: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
  },
  {
    value: 'cancelled',
    label: 'Dibatalkan',
    className: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
  },
  {
    value: 'refunded',
    label: 'Dikembalikan',
    className: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
  },
  {
    value: 'overpaid',
    label: 'Lebih Bayar',
    className: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
  },
] as const

export const overdueOptions = [
  { value: 'yes', label: 'Terlambat' },
  { value: 'no', label: 'Belum Terlambat' },
] as const
