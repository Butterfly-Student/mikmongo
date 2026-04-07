import { MoreHorizontal, Pencil, Trash2 } from 'lucide-react'
import type { ColumnDef } from '@tanstack/react-table'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { HotspotProfile } from '@/lib/schemas/mikrotik'

interface ColumnActions {
    onEdit: (profile: HotspotProfile) => void
    onDelete: (profile: HotspotProfile) => void
}

export function createHotspotProfileColumns(
    actions: ColumnActions
): ColumnDef<HotspotProfile>[] {
    return [
        {
            accessorKey: 'name',
            header: 'Name',
            cell: ({ row }) => (
                <span className='font-medium'>{row.original.name}</span>
            ),
        },
        {
            accessorKey: 'addressPool',
            header: 'Address Pool',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.addressPool ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'sharedUsers',
            header: 'Shared Users',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.sharedUsers ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'rateLimit',
            header: 'Rate Limit',
            cell: ({ row }) => (
                <span className='font-mono text-xs text-muted-foreground'>
                    {row.original.rateLimit ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'parentQueue',
            header: 'Parent Queue',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.parentQueue ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'sessionTimeout',
            header: 'Session Timeout',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.sessionTimeout ?? '-'}
                </span>
            ),
        },
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const profile = row.original
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
                            <DropdownMenuItem onClick={() => actions.onEdit(profile)}>
                                <Pencil className='mr-2 size-4' />
                                Edit
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                                onClick={() => actions.onDelete(profile)}
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
