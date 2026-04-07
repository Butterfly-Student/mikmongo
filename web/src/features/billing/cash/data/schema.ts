import type { CashEntryResponse, PettyCashFundResponse } from '@/lib/schemas/billing'

export type { CashEntryResponse, PettyCashFundResponse }

export const cashEntryStatuses = [
  {
    value: 'pending',
    label: 'Menunggu',
    className: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300',
  },
  {
    value: 'approved',
    label: 'Disetujui',
    className: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
  },
  {
    value: 'rejected',
    label: 'Ditolak',
    className: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
  },
] as const

export const cashEntryTypes = [
  {
    value: 'income',
    label: 'Masuk',
    className: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
  },
  {
    value: 'expense',
    label: 'Keluar',
    className: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
  },
] as const

export const cashEntrySources = [
  { value: 'invoice', label: 'Tagihan' },
  { value: 'agent_invoice', label: 'Tagihan Agen' },
  { value: 'installation', label: 'Instalasi' },
  { value: 'penalty', label: 'Denda' },
  { value: 'other', label: 'Lainnya' },
  { value: 'operational', label: 'Operasional' },
  { value: 'upstream', label: 'Upstream' },
  { value: 'purchase', label: 'Pembelian' },
  { value: 'salary', label: 'Gaji' },
] as const
