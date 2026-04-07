import { WifiOff } from 'lucide-react'
import type { ColumnDef } from '@tanstack/react-table'
import { Button } from '@/components/ui/button'
import type { HotspotActive } from '@/lib/schemas/mikrotik'

function formatBytes(bytes: number | null | undefined): string {
    if (!bytes) return '-'
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`
    return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`
}

interface ColumnActions {
    onDisconnect: (session: HotspotActive) => void
}

export function createHotspotActiveColumns(
    actions: ColumnActions
): ColumnDef<HotspotActive>[] {
    return [
        {
            accessorKey: 'user',
            header: 'User',
            cell: ({ row }) => (
                <span className='font-medium'>{row.original.user ?? '-'}</span>
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
            accessorKey: 'address',
            header: 'IP Address',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>
                    {row.original.address ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'macAddress',
            header: 'MAC Address',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>
                    {row.original.macAddress ?? '-'}
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
            accessorKey: 'sessionTimeLeft',
            header: 'Time Left',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.sessionTimeLeft ?? '-'}
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
        {
            accessorKey: 'loginBy',
            header: 'Login By',
            cell: ({ row }) => (
                <span className='text-xs text-muted-foreground'>
                    {row.original.loginBy ?? '-'}
                </span>
            ),
        },
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const session = row.original
                return (
                    <Button
                        variant='ghost'
                        size='sm'
                        className='text-destructive hover:text-destructive'
                        onClick={() => actions.onDisconnect(session)}
                    >
                        <WifiOff className='mr-1 size-4' />
                        Disconnect
                    </Button>
                )
            },
        },
    ]
}
