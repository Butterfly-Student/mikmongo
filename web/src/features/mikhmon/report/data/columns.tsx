import type { ColumnDef } from '@tanstack/react-table'
import type { MikhmonReportResponse } from '@/lib/schemas/mikhmon'

function formatDate(dateStr: string | null | undefined): string {
    if (!dateStr) return '-'
    try {
        return new Date(dateStr).toLocaleString()
    } catch {
        return dateStr
    }
}

export function createReportColumns(): ColumnDef<MikhmonReportResponse>[] {
    return [
        {
            accessorKey: 'created_at',
            header: 'Date',
            cell: ({ row }) => (
                <span className='text-xs text-muted-foreground whitespace-nowrap'>
                    {formatDate(row.original.created_at)}
                </span>
            ),
        },
        {
            accessorKey: 'user',
            header: 'Username',
            cell: ({ row }) => (
                <span className='font-mono text-sm'>{row.original.user ?? '-'}</span>
            ),
        },
        {
            accessorKey: 'profile',
            header: 'Profile',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>{row.original.profile ?? '-'}</span>
            ),
        },
        {
            accessorKey: 'validity',
            header: 'Validity',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>{row.original.validity ?? '-'}</span>
            ),
        },
        {
            accessorKey: 'price',
            header: 'Price',
            cell: ({ row }) => {
                const price = row.original.price
                return price != null ? (
                    <span className='text-sm font-medium'>
                        {price.toLocaleString('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 })}
                    </span>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            accessorKey: 'ip',
            header: 'IP',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>{row.original.ip ?? '-'}</span>
            ),
        },
        {
            accessorKey: 'mac',
            header: 'MAC',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>{row.original.mac ?? '-'}</span>
            ),
        },
        {
            accessorKey: 'comment',
            header: 'Comment',
            cell: ({ row }) => (
                <span className='text-xs text-muted-foreground'>{row.original.comment ?? '-'}</span>
            ),
        },
    ]
}
