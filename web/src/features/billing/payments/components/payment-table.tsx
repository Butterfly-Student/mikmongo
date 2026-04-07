import { useState, useMemo } from 'react'
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
import type { PaymentResponse } from '../data/schema'
import { paymentStatuses, paymentMethods } from '../data/schema'
import { createPaymentColumns } from '../data/columns'
import { DateRangeFilter } from './date-range-filter'

interface PaymentTableProps {
  data: PaymentResponse[]
  onConfirm: (payment: PaymentResponse) => void
  onReject: (payment: PaymentResponse) => void
  onRefund: (payment: PaymentResponse) => void
  onGateway: (payment: PaymentResponse) => void
}

export function PaymentTable({
  data,
  onConfirm,
  onReject,
  onRefund,
  onGateway,
}: PaymentTableProps) {
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')

  const columns = useMemo(
    () => createPaymentColumns({ onConfirm, onReject, onRefund, onGateway }),
    [onConfirm, onReject, onRefund, onGateway]
  )

  const filteredData = useMemo(() => {
    if (!dateFrom && !dateTo) return data
    return data.filter((payment) => {
      const paymentDate = payment.payment_date.slice(0, 10)
      if (dateFrom && paymentDate < dateFrom) return false
      if (dateTo && paymentDate > dateTo) return false
      return true
    })
  }, [data, dateFrom, dateTo])

  const table = useReactTable({
    data: filteredData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const statusFilterOptions = paymentStatuses.map((s) => ({
    label: s.label,
    value: s.value,
  }))

  const methodFilterOptions = paymentMethods.map((m) => ({
    label: m.label,
    value: m.value,
  }))

  return (
    <div className='space-y-4'>
      <div className='flex items-center gap-2 flex-wrap'>
        <DataTableToolbar
          table={table}
          searchPlaceholder='Cari nomor referensi atau tagihan...'
          filters={[
            {
              columnId: 'payment_method',
              title: 'Metode',
              options: methodFilterOptions,
            },
            {
              columnId: 'status',
              title: 'Status',
              options: statusFilterOptions,
            },
          ]}
        />
        <DateRangeFilter
          dateFrom={dateFrom}
          dateTo={dateTo}
          onDateFromChange={setDateFrom}
          onDateToChange={setDateTo}
        />
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
                          header.getContext()
                        )}
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
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center'
                >
                  <div className='text-sm text-muted-foreground'>
                    Tidak ada pembayaran ditemukan.
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
