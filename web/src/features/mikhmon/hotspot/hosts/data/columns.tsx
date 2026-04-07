import type { ColumnDef } from '@tanstack/react-table'
import type { HotspotHost } from '@/lib/schemas/mikrotik'

function formatBytes(bytes: number | null | undefined): string {
    if (!bytes) return '-'
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`
    return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`
}

export function createHotspotHostColumns(): ColumnDef<HotspotHost>[] {
    return [
        {
            accessorKey: 'macAddress',
            header: 'MAC Address',
            cell: ({ row }) => (
                <span className='font-mono text-sm'>
                    {row.original.macAddress ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'address',
            header: 'IP Address',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>
                    {row.original.address ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'toAddress',
            header: 'To Address',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>
                    {row.original.toAddress ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'server',
            header: 'Server',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.server ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'uptime',
            header: 'Uptime',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.uptime ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'bytesIn',
            header: 'Download',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {formatBytes(row.original.bytesIn)}
                </span>
            ),
        },
        {
            accessorKey: 'bytesOut',
            header: 'Upload',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {formatBytes(row.original.bytesOut)}
                </span>
            ),
        },
    ]
}
