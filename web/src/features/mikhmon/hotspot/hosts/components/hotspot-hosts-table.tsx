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
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import type { HotspotHost } from '@/lib/schemas/mikrotik'

interface HotspotHostsTableProps {
    columns: ColumnDef<HotspotHost>[]
    data: HotspotHost[]
    isLoading: boolean
    search: string
    onSearchChange: (search: string) => void
}

export function HotspotHostsTable({
    columns,
    data,
    isLoading,
    search,
    onSearchChange,
}: HotspotHostsTableProps) {
    const table = useReactTable({
        data,
        columns,
        getCoreRowModel: getCoreRowModel(),
    })

    return (
        <div className='space-y-4'>
            <div className='flex items-center'>
                <Input
                    placeholder='Search by MAC, IP, server...'
                    value={search}
                    onChange={(e) => onSearchChange(e.target.value)}
                    className='h-8 w-[150px] lg:w-[300px]'
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
                                        No hosts found.
                                    </div>
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            <div className='text-sm text-muted-foreground px-2'>
                {data.length} host{data.length !== 1 ? 's' : ''}
            </div>
        </div>
    )
}
