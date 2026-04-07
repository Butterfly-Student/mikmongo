import { flexRender, getCoreRowModel, useReactTable, type ColumnDef } from '@tanstack/react-table'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import type { MikhmonReportResponse } from '@/lib/schemas/mikhmon'

interface ReportTableProps {
    columns: ColumnDef<MikhmonReportResponse>[]
    data: MikhmonReportResponse[]
    isLoading: boolean
    pagination: { pageIndex: number; pageSize: number }
    onPaginationChange: (pagination: { pageIndex: number; pageSize: number }) => void
    search: string
    onSearchChange: (search: string) => void
}

export function ReportTable({ columns, data, isLoading, pagination, onPaginationChange, search, onSearchChange }: ReportTableProps) {
    const paginatedData = data.slice(pagination.pageIndex * pagination.pageSize, (pagination.pageIndex + 1) * pagination.pageSize)
    const table = useReactTable({ data: paginatedData, columns, getCoreRowModel: getCoreRowModel() })

    return (
        <div className='space-y-4'>
            <Input
                placeholder='Search by user, profile, IP...'
                value={search}
                onChange={(e) => onSearchChange(e.target.value)}
                className='h-8 w-[250px]'
            />
            <div className='rounded-md border'>
                <Table>
                    <TableHeader>
                        {table.getHeaderGroups().map((hg) => (
                            <TableRow key={hg.id}>
                                {hg.headers.map((h) => (
                                    <TableHead key={h.id}>
                                        {h.isPlaceholder ? null : flexRender(h.column.columnDef.header, h.getContext())}
                                    </TableHead>
                                ))}
                            </TableRow>
                        ))}
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>{columns.map((_, ci) => <TableCell key={ci}><Skeleton className='h-5 w-full' /></TableCell>)}</TableRow>
                            ))
                        ) : table.getRowModel().rows?.length ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow key={row.id}>
                                    {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                                    ))}
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                <TableCell colSpan={columns.length} className='h-24 text-center text-sm text-muted-foreground'>
                                    No report data found.
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            <div className='flex items-center justify-between px-2'>
                <div className='flex-1 text-sm text-muted-foreground'>Showing {paginatedData.length} of {data.length} records.</div>
                <div className='flex items-center space-x-2'>
                    <Button variant='outline' size='sm' onClick={() => onPaginationChange({ ...pagination, pageIndex: pagination.pageIndex - 1 })} disabled={pagination.pageIndex === 0}>Previous</Button>
                    <Button variant='outline' size='sm' onClick={() => onPaginationChange({ ...pagination, pageIndex: pagination.pageIndex + 1 })} disabled={(pagination.pageIndex + 1) * pagination.pageSize >= data.length}>Next</Button>
                </div>
            </div>
        </div>
    )
}
