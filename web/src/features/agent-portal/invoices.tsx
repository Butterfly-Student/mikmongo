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
import { useAgentPortalInvoices, useAgentRequestPayment } from '@/hooks/use-portal-billing'
import type { AgentInvoiceResponse } from '@/lib/schemas/billing'

function formatRp(amount: number): string {
  return `Rp ${amount.toLocaleString('id-ID')}`
}

function formatPeriod(start: string, end: string): string {
  return `${format(new Date(start), 'dd/MM/yyyy')} - ${format(new Date(end), 'dd/MM/yyyy')}`
}

const agentInvoiceStatusInfo: Record<
  AgentInvoiceResponse['status'],
  { label: string; className: string }
> = {
  draft: {
    label: 'Draft',
    className: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
  },
  unpaid: {
    label: 'Belum Lunas',
    className:
      'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300',
  },
  paid: {
    label: 'Lunas',
    className:
      'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
  },
  cancelled: {
    label: 'Dibatalkan',
    className: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
  },
}

function AgentInvoiceStatusBadge({ status }: { status: AgentInvoiceResponse['status'] }) {
  const info = agentInvoiceStatusInfo[status]
  return (
    <Badge className={info?.className ?? ''} variant='outline'>
      {info?.label ?? status}
    </Badge>
  )
}

function DataTablePagination({
  table,
}: {
  table: ReturnType<typeof useReactTable<AgentInvoiceResponse>>
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

function RequestPaymentButton({ invoice }: { invoice: AgentInvoiceResponse }) {
  const { mutate, isPending } = useAgentRequestPayment()
  if (invoice.status !== 'unpaid') return null
  return (
    <Button
      variant='outline'
      size='sm'
      onClick={() => mutate({ id: invoice.id, data: {} })}
      disabled={isPending}
    >
      Ajukan Pembayaran
    </Button>
  )
}

export default function AgentInvoicesPage() {
  const { data: invoices, isLoading } = useAgentPortalInvoices()

  const columns: ColumnDef<AgentInvoiceResponse>[] = [
    {
      accessorKey: 'invoice_number',
      header: 'No. Tagihan',
      cell: ({ row }) => (
        <span className='font-mono text-xs'>{row.original.invoice_number}</span>
      ),
    },
    {
      id: 'period',
      header: 'Periode',
      cell: ({ row }) =>
        formatPeriod(row.original.period_start, row.original.period_end),
    },
    {
      id: 'voucher_count',
      header: 'Voucher',
      cell: ({ row }) => row.original.voucher_count,
    },
    {
      accessorKey: 'total_amount',
      header: () => <div className='text-right'>Total</div>,
      cell: ({ row }) => (
        <div className='text-right'>{formatRp(row.original.total_amount)}</div>
      ),
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => (
        <AgentInvoiceStatusBadge status={row.original.status} />
      ),
    },
    {
      id: 'actions',
      header: 'Aksi',
      cell: ({ row }) => <RequestPaymentButton invoice={row.original} />,
    },
  ]

  const table = useReactTable({
    data: invoices ?? [],
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    initialState: { pagination: { pageSize: 10 } },
  })

  return (
    <div className='space-y-6'>
      <div>
        <h1 className='text-2xl font-semibold tracking-tight'>Tagihan Klien</h1>
        <p className='text-sm text-muted-foreground'>
          Ajukan permintaan pembayaran untuk tagihan pelanggan
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
                  Tidak ada tagihan klien ditemukan.
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
