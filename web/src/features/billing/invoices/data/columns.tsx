import type { ColumnDef } from '@tanstack/react-table'
import { Badge } from '@/components/ui/badge'
import type { InvoiceResponse } from '@/lib/schemas/billing'
import { invoiceStatuses } from './schema'

export const invoiceColumns: ColumnDef<InvoiceResponse>[] = [
  {
    accessorKey: 'invoice_number',
    header: 'No. Tagihan',
    enableSorting: true,
    cell: ({ row }) => (
      <span className='font-mono text-xs'>{row.original.invoice_number}</span>
    ),
  },
  {
    accessorKey: 'customer_id',
    header: 'Pelanggan',
    enableSorting: true,
    cell: ({ row }) => (
      <span className='font-mono text-xs text-muted-foreground'>
        {row.original.customer_id.slice(0, 8)}...
      </span>
    ),
  },
  {
    accessorKey: 'issue_date',
    header: 'Tanggal',
    enableSorting: true,
    cell: ({ row }) => (
      <span className='text-sm text-muted-foreground'>
        {new Date(row.original.issue_date).toLocaleDateString('id-ID')}
      </span>
    ),
  },
  {
    accessorKey: 'payment_deadline',
    header: 'Jatuh Tempo',
    enableSorting: true,
    cell: ({ row }) => {
      const isOverdue =
        new Date(row.original.payment_deadline) < new Date() &&
        row.original.status !== 'paid'
      return (
        <span
          className={`text-sm ${isOverdue ? 'text-red-600 font-medium' : 'text-muted-foreground'}`}
        >
          {new Date(row.original.payment_deadline).toLocaleDateString('id-ID')}
        </span>
      )
    },
  },
  {
    accessorKey: 'total_amount',
    header: 'Total',
    enableSorting: true,
    cell: ({ row }) => (
      <span className='text-sm text-right block font-medium'>
        Rp {row.original.total_amount.toLocaleString('id-ID')}
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
      const statusInfo = invoiceStatuses.find((s) => s.value === row.original.status)
      return (
        <Badge className={statusInfo?.className ?? ''} variant='outline'>
          {statusInfo?.label ?? row.original.status}
        </Badge>
      )
    },
  },
]
