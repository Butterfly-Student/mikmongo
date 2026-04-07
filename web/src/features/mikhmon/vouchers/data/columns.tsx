import { Trash2 } from 'lucide-react'
import type { ColumnDef } from '@tanstack/react-table'
import { Button } from '@/components/ui/button'
import type { VoucherResponse } from '@/lib/schemas/mikhmon'

interface ColumnActions {
    onDelete: (voucher: VoucherResponse) => void
}

export function createVoucherColumns(actions: ColumnActions): ColumnDef<VoucherResponse>[] {
    return [
        {
            accessorKey: 'name',
            header: 'Username',
            cell: ({ row }) => (
                <span className='font-mono text-sm font-medium'>{row.original.name}</span>
            ),
        },
        {
            accessorKey: 'password',
            header: 'Password',
            cell: ({ row }) => {
                const pwd = row.original.password
                return pwd ? (
                    <span className='font-mono text-sm text-muted-foreground'>{pwd}</span>
                ) : (
                    <span className='text-xs text-muted-foreground italic'>same as username</span>
                )
            },
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
            accessorKey: 'code',
            header: 'Batch Code',
            cell: ({ row }) => (
                <span className='text-xs text-muted-foreground'>
                    {row.original.code ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'date',
            header: 'Date',
            cell: ({ row }) => (
                <span className='text-xs text-muted-foreground'>
                    {row.original.date ?? '-'}
                </span>
            ),
        },
        {
            accessorKey: 'comment',
            header: 'Comment',
            cell: ({ row }) => {
                const comment = row.original.comment
                return comment ? (
                    <span className='max-w-[120px] truncate block text-xs text-muted-foreground'>
                        {comment}
                    </span>
                ) : (
                    <span className='text-xs text-muted-foreground'>-</span>
                )
            },
        },
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const voucher = row.original
                return (
                    <Button
                        variant='ghost'
                        size='icon'
                        className='text-destructive hover:text-destructive'
                        onClick={() => actions.onDelete(voucher)}
                    >
                        <Trash2 className='size-4' />
                    </Button>
                )
            },
        },
    ]
}
