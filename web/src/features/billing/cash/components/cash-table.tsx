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
import type { CashEntryResponse } from '@/lib/schemas/billing'
import { cashColumns } from '../data/columns'
import { cashEntryStatuses, cashEntryTypes } from '../data/schema'

interface CashTableProps {
  data: CashEntryResponse[]
  onApprove: (id: string) => void
  onReject: (entry: CashEntryResponse) => void
}

export function CashTable({ data, onApprove, onReject }: CashTableProps) {
  const table = useReactTable({
    data,
    columns: cashColumns,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    globalFilterFn: 'auto',
    meta: {
      onApprove,
      onReject,
    },
  })

  const statusFilterOptions = cashEntryStatuses.map((s) => ({
    label: s.label,
    value: s.value,
  }))

  const typeFilterOptions = cashEntryTypes.map((t) => ({
    label: t.label,
    value: t.value,
  }))

  return (
    <div className='space-y-4'>
      <DataTableToolbar
        table={table}
        searchPlaceholder='Cari deskripsi atau sumber...'
        filters={[
          {
            columnId: 'status',
            title: 'Status',
            options: statusFilterOptions,
          },
          {
            columnId: 'type',
            title: 'Tipe',
            options: typeFilterOptions,
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
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={cashColumns.length} className='h-24 text-center'>
                  <div className='text-sm text-muted-foreground'>
                    Tidak ada entri kas ditemukan.
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
