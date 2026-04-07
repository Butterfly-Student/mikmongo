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
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'
import { RefreshCw } from 'lucide-react'
import type { HotspotActive } from '@/lib/schemas/mikrotik'

interface HotspotActiveTableProps {
    columns: ColumnDef<HotspotActive>[]
    data: HotspotActive[]
    isLoading: boolean
    search: string
    onSearchChange: (search: string) => void
    serverFilter: string
    onServerFilterChange: (server: string) => void
    servers: string[]
}

export function HotspotActiveTable({
    columns,
    data,
    isLoading,
    search,
    onSearchChange,
    serverFilter,
    onServerFilterChange,
    servers,
}: HotspotActiveTableProps) {
    const table = useReactTable({
        data,
        columns,
        getCoreRowModel: getCoreRowModel(),
    })

    return (
        <div className='space-y-4'>
            <div className='flex items-center justify-between'>
                <div className='flex flex-1 items-center space-x-2'>
                    <Input
                        placeholder='Search by user, IP, MAC...'
                        value={search}
                        onChange={(e) => onSearchChange(e.target.value)}
                        className='h-8 w-[150px] lg:w-[250px]'
                    />
                    <Select value={serverFilter} onValueChange={onServerFilterChange}>
                        <SelectTrigger className='h-8 w-[150px]'>
                            <SelectValue placeholder='All Servers' />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value='all'>All Servers</SelectItem>
                            {servers.map((s) => (
                                <SelectItem key={s} value={s}>
                                    {s}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>
                <Badge variant='outline' className='gap-1.5'>
                    <RefreshCw className='size-3 animate-spin' />
                    Auto-refreshing
                </Badge>
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
                                    <div className='flex flex-col items-center justify-center p-8'>
                                        <div className='text-sm text-muted-foreground'>
                                            No active sessions found.
                                        </div>
                                    </div>
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            <div className='flex-1 text-sm text-muted-foreground px-2'>
                {data.length} active session{data.length !== 1 ? 's' : ''}
            </div>
        </div>
    )
}
