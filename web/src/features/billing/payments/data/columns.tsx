import type { ColumnDef } from '@tanstack/react-table'
import { Badge } from '@/components/ui/badge'
import { format } from 'date-fns'
import type { PaymentResponse } from './schema'
import { paymentStatuses, paymentMethods } from './schema'
import { PaymentActionMenu } from '../components/payment-action-menu'

interface PaymentColumnsProps {
  onConfirm: (payment: PaymentResponse) => void
  onReject: (payment: PaymentResponse) => void
  onRefund: (payment: PaymentResponse) => void
  onGateway: (payment: PaymentResponse) => void
}

export function createPaymentColumns({
  onConfirm,
  onReject,
  onRefund,
  onGateway,
}: PaymentColumnsProps): ColumnDef<PaymentResponse>[] {
  return [
    {
      accessorKey: 'payment_number',
      header: 'No. Referensi',
      enableSorting: true,
      cell: ({ row }) => (
        <span className='font-mono text-xs'>{row.original.payment_number}</span>
      ),
    },
    {
      accessorKey: 'transaction_reference',
      header: 'No. Tagihan',
      enableSorting: false,
      cell: ({ row }) => {
        const ref = row.original.transaction_reference
        return (
          <span className='font-mono text-xs text-muted-foreground'>
            {ref ?? '-'}
          </span>
        )
      },
    },
    {
      accessorKey: 'payment_method',
      header: 'Metode',
      enableSorting: false,
      enableColumnFilter: true,
      filterFn: (row, id, value) => {
        return (value as string[]).includes(row.getValue(id))
      },
      cell: ({ row }) => {
        const methodInfo = paymentMethods.find(
          (m) => m.value === row.original.payment_method
        )
        return (
          <span className='text-sm'>{methodInfo?.label ?? row.original.payment_method}</span>
        )
      },
    },
    {
      accessorKey: 'amount',
      header: 'Jumlah',
      enableSorting: true,
      cell: ({ row }) => (
        <span className='text-sm text-right block font-medium'>
          Rp {row.original.amount.toLocaleString('id-ID')}
        </span>
      ),
    },
    {
      accessorKey: 'payment_date',
      header: 'Tanggal',
      enableSorting: true,
      cell: ({ row }) => (
        <span className='text-sm text-muted-foreground'>
          {format(new Date(row.original.payment_date), 'dd/MM/yyyy')}
        </span>
      ),
    },
    {
      accessorKey: 'status',
      header: 'Status',
      enableColumnFilter: true,
      filterFn: (row, id, value) => {
        return (value as string[]).includes(row.getValue(id))
      },
      cell: ({ row }) => {
        const statusInfo = paymentStatuses.find(
          (s) => s.value === row.original.status
        )
        return (
          <Badge className={statusInfo?.className ?? ''} variant='outline'>
            {statusInfo?.label ?? row.original.status}
          </Badge>
        )
      },
    },
    {
      id: 'actions',
      header: 'Aksi',
      cell: ({ row }) => (
        <PaymentActionMenu
          payment={row.original}
          onConfirm={onConfirm}
          onReject={onReject}
          onRefund={onRefund}
          onGateway={onGateway}
        />
      ),
    },
  ]
}
