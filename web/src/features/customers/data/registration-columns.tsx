import {
    MoreHorizontal,
    CheckCircle2,
    XCircle,
} from 'lucide-react'
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
import type { RegistrationResponse } from '@/lib/schemas/customer'

interface RegistrationColumnActions {
    onApprove: (registration: RegistrationResponse) => void
    onReject: (registration: RegistrationResponse) => void
}

export function createRegistrationColumns(
    actions: RegistrationColumnActions
): ColumnDef<RegistrationResponse>[] {
    return [
        {
            accessorKey: 'full_name',
            header: 'Name',
            cell: ({ row }) => {
                const reg = row.original
                const initials = reg.full_name
                    .split(' ')
                    .map((n) => n[0])
                    .join('')
                    .slice(0, 2)
                    .toUpperCase()
                return (
                    <div className='flex items-center gap-3'>
                        <div className='flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary'>
                            {initials}
                        </div>
                        <div className='flex flex-col'>
                            <span className='font-medium'>{reg.full_name}</span>
                            <span className='text-xs text-muted-foreground'>
                                {reg.phone}
                            </span>
                        </div>
                    </div>
                )
            },
        },
        {
            accessorKey: 'email',
            header: 'Email',
            cell: ({ row }) => {
                const email = row.original.email
                return email ? (
                    <span className='text-sm text-muted-foreground'>{email}</span>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            accessorKey: 'address',
            header: 'Address',
            cell: ({ row }) => {
                const address = row.original.address
                return address ? (
                    <span className='text-sm text-muted-foreground max-w-[200px] truncate block'>
                        {address}
                    </span>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            accessorKey: 'bandwidth_profile_id',
            header: 'Requested Plan',
            cell: ({ row }) => {
                const planId = row.original.bandwidth_profile_id
                return planId ? (
                    <Badge variant='outline'>
                        {planId.slice(0, 8)}...
                    </Badge>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            accessorKey: 'status',
            header: 'Status',
            cell: ({ row }) => {
                const status = row.original.status
                switch (status) {
                    case 'pending':
                        return <Badge variant='outline'>Pending</Badge>
                    case 'approved':
                        return <Badge variant='default'>Approved</Badge>
                    case 'rejected':
                        return <Badge variant='secondary'>Rejected</Badge>
                    default:
                        return <Badge variant='secondary'>{status}</Badge>
                }
            },
        },
        {
            accessorKey: 'created_at',
            header: 'Submitted',
            cell: ({ row }) => {
                const date = new Date(row.original.created_at)
                return (
                    <span className='text-sm text-muted-foreground'>
                        {date.toLocaleDateString()}
                    </span>
                )
            },
        },
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const registration = row.original
                const isPending = registration.status === 'pending'

                if (!isPending) {
                    return (
                        <div className='flex items-center gap-2'>
                            {registration.rejection_reason && (
                                <span
                                    className='text-xs text-muted-foreground'
                                    title={registration.rejection_reason}
                                >
                                    {registration.rejection_reason}
                                </span>
                            )}
                            <span className='text-xs text-muted-foreground'>
                                {registration.status === 'approved'
                                    ? 'Approved'
                                    : 'Rejected'}
                            </span>
                        </div>
                    )
                }

                return (
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant='ghost' size='icon'>
                                <span className='sr-only'>Open menu</span>
                                <MoreHorizontal className='size-4' />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align='end'>
                            <DropdownMenuLabel>Registration Actions</DropdownMenuLabel>
                            <DropdownMenuItem
                                onClick={() => actions.onApprove(registration)}
                            >
                                <CheckCircle2 className='mr-2 size-4' />
                                Approve
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                                onClick={() => actions.onReject(registration)}
                                className='text-destructive focus:text-destructive'
                            >
                                <XCircle className='mr-2 size-4' />
                                Reject
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                )
            },
        },
    ]
}
