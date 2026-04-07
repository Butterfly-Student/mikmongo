import {
    MoreHorizontal,
    User,
    Play,
    Pause,
    ShieldOff,
    RotateCcw,
    XCircle,
    Trash2,
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
import type { SubscriptionResponse } from '@/lib/schemas/subscription'

interface SubscriptionActions {
    onActivate: (subscription: SubscriptionResponse) => void
    onSuspend: (subscription: SubscriptionResponse) => void
    onIsolate: (subscription: SubscriptionResponse) => void
    onRestore: (subscription: SubscriptionResponse) => void
    onTerminate: (subscription: SubscriptionResponse) => void
    onDelete: (subscription: SubscriptionResponse) => void
}

const statusVariant: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    active: 'default',
    pending: 'secondary',
    suspended: 'outline',
    isolated: 'outline',
    expired: 'secondary',
    terminated: 'destructive',
}

export function createSubscriptionColumns(
    actions: SubscriptionActions
): ColumnDef<SubscriptionResponse>[] {
    return [
        {
            accessorKey: 'username',
            header: 'Username',
            cell: ({ row }) => (
                <div className='flex items-center gap-3'>
                    <div className='flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary'>
                        {row.original.username.slice(0, 2).toUpperCase()}
                    </div>
                    <span className='font-medium'>{row.original.username}</span>
                </div>
            ),
        },
        {
            accessorKey: 'customer_id',
            header: 'Customer ID',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground font-mono'>
                    {row.original.customer_id.slice(0, 8)}...
                </span>
            ),
        },
        {
            accessorKey: 'plan_id',
            header: 'Plan',
            cell: ({ row }) => {
                const planId = row.original.plan_id
                return planId ? (
                    <span className='text-sm text-muted-foreground font-mono'>
                        {planId.slice(0, 8)}...
                    </span>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            accessorKey: 'status',
            header: 'Status',
            cell: ({ row }) => (
                <Badge variant={statusVariant[row.original.status] ?? 'secondary'}>
                    {row.original.status.charAt(0).toUpperCase() + row.original.status.slice(1)}
                </Badge>
            ),
        },
        {
            accessorKey: 'static_ip',
            header: 'IP Address',
            cell: ({ row }) => {
                const ip = row.original.static_ip
                return ip ? (
                    <span className='text-sm text-muted-foreground'>{ip}</span>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            accessorKey: 'expiry_date',
            header: 'Expiry',
            cell: ({ row }) => {
                const expiry = row.original.expiry_date
                return expiry ? (
                    <span className='text-sm text-muted-foreground'>{expiry}</span>
                ) : (
                    <span className='text-sm text-muted-foreground'>-</span>
                )
            },
        },
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const subscription = row.original
                const { status } = subscription

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

                            {status !== 'active' && status !== 'terminated' && (
                                <DropdownMenuItem
                                    onClick={() => actions.onActivate(subscription)}
                                >
                                    <Play className='mr-2 size-4' />
                                    Activate
                                </DropdownMenuItem>
                            )}

                            {status === 'active' && (
                                <DropdownMenuItem
                                    onClick={() => actions.onSuspend(subscription)}
                                >
                                    <Pause className='mr-2 size-4' />
                                    Suspend
                                </DropdownMenuItem>
                            )}

                            {status === 'active' && (
                                <DropdownMenuItem
                                    onClick={() => actions.onIsolate(subscription)}
                                >
                                    <ShieldOff className='mr-2 size-4' />
                                    Isolate
                                </DropdownMenuItem>
                            )}

                            {(status === 'suspended' || status === 'isolated') && (
                                <DropdownMenuItem
                                    onClick={() => actions.onRestore(subscription)}
                                >
                                    <RotateCcw className='mr-2 size-4' />
                                    Restore
                                </DropdownMenuItem>
                            )}

                            {status !== 'terminated' && (
                                <>
                                    <DropdownMenuSeparator />
                                    <DropdownMenuItem
                                        onClick={() => actions.onTerminate(subscription)}
                                        className='text-destructive focus:text-destructive'
                                    >
                                        <XCircle className='mr-2 size-4' />
                                        Terminate
                                    </DropdownMenuItem>
                                </>
                            )}

                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                                onClick={() => actions.onDelete(subscription)}
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
