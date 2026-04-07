import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DataTableToolbar } from '@/components/data-table/toolbar'
import { DataTablePagination } from '@/components/data-table/pagination'
import type { InvoiceResponse } from '@/lib/schemas/billing'
import { invoiceColumns } from '../data/columns'
import { invoiceStatuses } from '../data/schema'

interface InvoiceTableProps {
  data: InvoiceResponse[]
  onRowClick: (invoice: InvoiceResponse) => void
}

export function InvoiceTable({ data, onRowClick }: InvoiceTableProps) {
  const table = useReactTable({
    data,
    columns: invoiceColumns,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const statusFilterOptions = invoiceStatuses.map((s) => ({
    label: s.label,
    value: s.value,
  }))

  return (
    <div className='space-y-4'>
      <DataTableToolbar
        table={table}
        searchPlaceholder='Cari nomor tagihan atau pelanggan...'
        filters={[
          {
            columnId: 'status',
            title: 'Status',
            options: statusFilterOptions,
          },
        ]}
      />
      <div className='rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  className='cursor-pointer hover:bg-muted/50'
                  onClick={() => onRowClick(row.original)}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={invoiceColumns.length} className='h-24 text-center'>
                  <div className='text-sm text-muted-foreground'>
                    Tidak ada tagihan ditemukan.
                  </div>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination table={table} />
    </div>
  )
}
