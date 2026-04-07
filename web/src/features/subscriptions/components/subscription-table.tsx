import {
    type ColumnDef,
    flexRender,
    getCoreRowModel,
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
import { Button } from '@/components/ui/button'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Plus } from 'lucide-react'
import type { SubscriptionResponse } from '@/lib/schemas/subscription'

interface SubscriptionTableProps {
    columns: ColumnDef<SubscriptionResponse>[]
    data: SubscriptionResponse[]
    meta?: { total: number }
    isLoading: boolean
    pagination: { pageIndex: number; pageSize: number }
    onPaginationChange: (pagination: { pageIndex: number; pageSize: number }) => void
    onAddSubscription: () => void
    statusFilter: string
    onStatusFilterChange: (status: string) => void
}

export function SubscriptionTable({
    columns,
    data,
    meta,
    isLoading,
    pagination,
    onPaginationChange,
    onAddSubscription,
    statusFilter,
    onStatusFilterChange,
}: SubscriptionTableProps) {
    const table = useReactTable({
        data,
        columns,
        getCoreRowModel: getCoreRowModel(),
        manualPagination: true,
        pageCount: Math.ceil((meta?.total ?? 0) / pagination.pageSize),
    })

    return (
        <div className='space-y-4'>
            <div className='flex items-center justify-between'>
                <div className='flex flex-1 items-center space-x-2'>
                    <Select value={statusFilter} onValueChange={onStatusFilterChange}>
                        <SelectTrigger className='h-8 w-[150px]'>
                            <SelectValue placeholder='All Status' />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value='all'>All Status</SelectItem>
                            <SelectItem value='active'>Active</SelectItem>
                            <SelectItem value='pending'>Pending</SelectItem>
                            <SelectItem value='suspended'>Suspended</SelectItem>
                            <SelectItem value='isolated'>Isolated</SelectItem>
                            <SelectItem value='expired'>Expired</SelectItem>
                            <SelectItem value='terminated'>Terminated</SelectItem>
                        </SelectContent>
                    </Select>
                </div>
                <Button onClick={onAddSubscription} size='sm' className='h-8'>
                    <Plus className='mr-2 size-4' />
                    Create Subscription
                </Button>
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
                        {isLoading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    {columns.map((_, colIndex) => (
                                        <TableCell key={colIndex}>
                                            <Skeleton className='h-6 w-full' />
                                        </TableCell>
                                    ))}
                                </TableRow>
                            ))
                        ) : table.getRowModel().rows?.length ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow
                                    key={row.id}
                                    data-state={row.getIsSelected() && 'selected'}
                                >
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
                                    <div className='flex flex-col items-center justify-center p-8'>
                                        <div className='text-sm text-muted-foreground'>
                                            No Data Available
                                        </div>
                                        <div className='text-xs text-muted-foreground mt-1'>
                                            You haven&apos;t added any subscriptions here yet.
                                            Click the &quot;Create Subscription&quot; button to
                                            get started.
                                        </div>
                                    </div>
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            <div className='flex items-center justify-between px-2'>
                <div className='flex-1 text-sm text-muted-foreground'>
                    Showing {data.length} of {meta?.total ?? 0} subscriptions.
                </div>
                <div className='flex items-center space-x-2'>
                    <Button
                        variant='outline'
                        size='sm'
                        onClick={() =>
                            onPaginationChange({
                                ...pagination,
                                pageIndex: pagination.pageIndex - 1,
                            })
                        }
                        disabled={pagination.pageIndex === 0}
                    >
                        Previous
                    </Button>
                    <Button
                        variant='outline'
                        size='sm'
                        onClick={() =>
                            onPaginationChange({
                                ...pagination,
                                pageIndex: pagination.pageIndex + 1,
                            })
                        }
                        disabled={pagination.pageIndex >= table.getPageCount() - 1}
                    >
                        Next
                    </Button>
                </div>
            </div>
        </div>
    )
}
