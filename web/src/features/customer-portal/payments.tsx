import {
  useReactTable,
  getCoreRowModel,
  getPaginationRowModel,
  flexRender,
  type ColumnDef,
} from '@tanstack/react-table'
import { format } from 'date-fns'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { usePortalPayments } from '@/hooks/use-portal-billing'
import { paymentStatuses } from '@/features/billing/payments/data/schema'
import type { PaymentResponse } from '@/lib/schemas/billing'

function formatRp(amount: number): string {
  return `Rp ${amount.toLocaleString('id-ID')}`
}

function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1).replace(/_/g, ' ')
}

function PaymentStatusBadge({ status }: { status: PaymentResponse['status'] }) {
  const info = paymentStatuses.find((s) => s.value === status)
  return (
    <Badge className={info?.className ?? ''} variant='outline'>
      {info?.label ?? status}
    </Badge>
  )
}

function DataTablePagination({
  table,
}: {
  table: ReturnType<typeof useReactTable<PaymentResponse>>
}) {
  return (
    <div className='flex items-center justify-between px-2 py-2'>
      <div className='text-sm text-muted-foreground'>
        Halaman {table.getState().pagination.pageIndex + 1} dari{' '}
        {table.getPageCount()}
      </div>
      <div className='flex gap-2'>
        <Button
          variant='outline'
          size='sm'
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          Sebelumnya
        </Button>
        <Button
          variant='outline'
          size='sm'
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          Selanjutnya
        </Button>
      </div>
    </div>
  )
}

export default function CustomerPaymentsPage() {
  const { data: payments, isLoading } = usePortalPayments()

  const columns: ColumnDef<PaymentResponse>[] = [
    {
      accessorKey: 'payment_number',
      header: 'No. Referensi',
      cell: ({ row }) => (
        <span className='font-mono text-xs'>{row.original.payment_number}</span>
      ),
    },
    {
      accessorKey: 'payment_method',
      header: 'Metode',
      cell: ({ row }) => capitalize(row.original.payment_method),
    },
    {
      accessorKey: 'amount',
      header: () => <div className='text-right'>Jumlah</div>,
      cell: ({ row }) => (
        <div className='text-right'>{formatRp(row.original.amount)}</div>
      ),
    },
    {
      accessorKey: 'payment_date',
      header: 'Tanggal',
      cell: ({ row }) =>
        format(new Date(row.original.payment_date), 'dd/MM/yyyy'),
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => <PaymentStatusBadge status={row.original.status} />,
    },
  ]

  const table = useReactTable({
    data: payments ?? [],
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    initialState: { pagination: { pageSize: 10 } },
  })

  return (
    <div className='space-y-6'>
      <div>
        <h1 className='text-2xl font-semibold tracking-tight'>
          Riwayat Pembayaran Saya
        </h1>
        <p className='text-sm text-muted-foreground'>
          Lihat semua riwayat pembayaran Anda
        </p>
      </div>

      <div className='rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  {columns.map((_, j) => (
                    <TableCell key={j}>
                      <Skeleton className='h-4 w-full' />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : table.getRowModel().rows.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='text-center py-8 text-muted-foreground'
                >
                  Belum ada riwayat pembayaran.
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination table={table} />
    </div>
  )
}
