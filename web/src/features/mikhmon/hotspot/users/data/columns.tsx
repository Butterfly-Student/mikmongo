import { MoreHorizontal, Pencil, Trash2 } from 'lucide-react'
import type { ColumnDef } from '@tanstack/react-table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { HotspotUser } from '@/lib/schemas/mikrotik'

interface ColumnActions {
    onEdit: (user: HotspotUser) => void
    onDelete: (user: HotspotUser) => void
}

export function createHotspotUserColumns(actions: ColumnActions): ColumnDef<HotspotUser>[] {
    return [
        {
            accessorKey: 'name',
            header: 'Username',
            cell: ({ row }) => (
                <span className='font-medium'>{row.original.name}</span>
            ),
        },
        {
            accessorKey: 'profile',
            header: 'Profile',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.profile ?? '-'}
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
            accessorKey: 'macAddress',
            header: 'MAC Address',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>
                    {row.original.macAddress ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'limitUptime',
            header: 'Time Limit',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.limitUptime ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'comment',
            header: 'Comment',
            cell: ({ row }) => {
                const comment = row.original.comment
                return comment ? (
                    <span className='max-w-[150px] truncate block text-xs text-muted-foreground'>
                        {comment}
                    </span>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            accessorKey: 'disabled',
            header: 'Status',
            cell: ({ row }) =>
                row.original.disabled ? (
                    <Badge variant='secondary'>Disabled</Badge>
                ) : (
                    <Badge variant='default'>Active</Badge>
                ),
        },
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const user = row.original
                return (
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant='ghost' size='icon'>
                                <span className='sr-only'>Open menu</span>
                                <MoreHorizontal className='size-4' />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align='end'>
                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                            <DropdownMenuItem onClick={() => actions.onEdit(user)}>
                                <Pencil className='mr-2 size-4' />
                                Edit
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                                onClick={() => actions.onDelete(user)}
                                className='text-destructive focus:text-destructive'
                            >
                                <Trash2 className='mr-2 size-4' />
                                Delete
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                )
            },
        },
    ]
}
