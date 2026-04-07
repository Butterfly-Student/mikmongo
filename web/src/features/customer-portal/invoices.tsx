import { useState, useMemo } from 'react'
import {
  useReactTable,
  getCoreRowModel,
  getPaginationRowModel,
  getFilteredRowModel,
  flexRender,
  type ColumnDef,
} from '@tanstack/react-table'
import { format, isPast } from 'date-fns'
import { id as idLocale } from 'date-fns/locale'
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
import { Input } from '@/components/ui/input'
import { usePortalInvoices, usePortalInitiatePayment } from '@/hooks/use-portal-billing'
import { InvoiceDetailSheet } from '@/features/billing/invoices/components/invoice-detail-sheet'
import { invoiceStatuses } from '@/features/billing/invoices/data/schema'
import type { InvoiceResponse } from '@/lib/schemas/billing'

function formatRp(amount: number): string {
  return `Rp ${amount.toLocaleString('id-ID')}`
}

function formatPeriod(start: string, end: string): string {
  const startDate = new Date(start)
  const endDate = new Date(end)
  const startStr = format(startDate, 'MMM yyyy', { locale: idLocale })
  const endStr = format(endDate, 'MMM yyyy', { locale: idLocale })
  return startStr === endStr ? startStr : `${startStr} - ${endStr}`
}

function InvoiceStatusBadge({ status }: { status: InvoiceResponse['status'] }) {
  const info = invoiceStatuses.find((s) => s.value === status)
  return (
    <Badge className={info?.className ?? ''} variant='outline'>
      {info?.label ?? status}
    </Badge>
  )
}

function DataTablePagination({
  table,
}: {
  table: ReturnType<typeof useReactTable<InvoiceResponse>>
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

function PayButton({ invoice }: { invoice: InvoiceResponse }) {
  const { mutate, isPending } = usePortalInitiatePayment()
  if (invoice.status !== 'unpaid' && invoice.status !== 'overdue') return null
  return (
    <Button
      variant='default'
      size='sm'
      onClick={() => mutate(invoice.id)}
      disabled={isPending}
    >
      Bayar Sekarang
    </Button>
  )
}

export default function CustomerInvoicesPage() {
  const { data: invoices, isLoading } = usePortalInvoices()
  const [selectedInvoice, setSelectedInvoice] = useState<InvoiceResponse | null>(null)
  const [sheetOpen, setSheetOpen] = useState(false)
  const [globalFilter, setGlobalFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')

  const filtered = useMemo(() => {
    if (!invoices) return []
    return invoices.filter((inv) => {
      const matchesSearch = globalFilter
        ? inv.invoice_number.toLowerCase().includes(globalFilter.toLowerCase())
        : true
      const matchesStatus = statusFilter !== 'all' ? inv.status === statusFilter : true
      return matchesSearch && matchesStatus
    })
  }, [invoices, globalFilter, statusFilter])

  const columns: ColumnDef<InvoiceResponse>[] = [
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
        formatPeriod(
          row.original.billing_period_start,
          row.original.billing_period_end,
        ),
    },
    {
      accessorKey: 'payment_deadline',
      header: 'Jatuh Tempo',
      cell: ({ row }) => {
        const deadline = row.original.payment_deadline
        const isOverdue =
          isPast(new Date(deadline)) && row.original.status !== 'paid'
        return (
          <span className={isOverdue ? 'text-red-600' : undefined}>
            {format(new Date(deadline), 'dd/MM/yyyy')}
          </span>
        )
      },
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
      cell: ({ row }) => <InvoiceStatusBadge status={row.original.status} />,
    },
    {
      id: 'actions',
      header: 'Aksi',
      cell: ({ row }) => (
        <div className='flex items-center gap-2'>
          <Button
            variant='outline'
            size='sm'
            onClick={() => {
              setSelectedInvoice(row.original)
              setSheetOpen(true)
            }}
          >
            Detail
          </Button>
          <PayButton invoice={row.original} />
        </div>
      ),
    },
  ]

  const table = useReactTable({
    data: filtered,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    initialState: { pagination: { pageSize: 10 } },
  })

  const statusOptions = [
    { value: 'all', label: 'Semua' },
    { value: 'unpaid', label: 'Belum Lunas' },
    { value: 'paid', label: 'Lunas' },
    { value: 'overdue', label: 'Terlambat' },
  ]

  return (
    <div className='space-y-6'>
      <div>
        <h1 className='text-2xl font-semibold tracking-tight'>Tagihan Saya</h1>
        <p className='text-sm text-muted-foreground'>
          Lihat dan bayar tagihan langganan Anda
        </p>
      </div>

      {/* Toolbar */}
      <div className='flex flex-col gap-3 sm:flex-row sm:items-center'>
        <Input
          placeholder='Cari nomor tagihan...'
          value={globalFilter}
          onChange={(e) => setGlobalFilter(e.target.value)}
          className='max-w-xs'
        />
        <div className='flex gap-2 flex-wrap'>
          {statusOptions.map((opt) => (
            <Button
              key={opt.value}
              variant={statusFilter === opt.value ? 'default' : 'outline'}
              size='sm'
              onClick={() => setStatusFilter(opt.value)}
            >
              {opt.label}
            </Button>
          ))}
        </div>
      </div>

      {/* Table */}
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
                <TableCell colSpan={columns.length} className='text-center py-8 text-muted-foreground'>
                  Belum ada tagihan.
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

      <InvoiceDetailSheet
        invoice={selectedInvoice}
        open={sheetOpen}
        onOpenChange={setSheetOpen}
      />
    </div>
  )
}
