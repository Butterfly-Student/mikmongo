import {
    MoreHorizontal,
    Trash2,
    Power,
    PowerOff,
    Pencil,
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
import type { CustomerResponse } from '@/lib/schemas/customer'

interface ColumnActions {
    onActivate: (customer: CustomerResponse) => void
    onDeactivate: (customer: CustomerResponse) => void
    onDelete: (customer: CustomerResponse) => void
    onEdit: (customer: CustomerResponse) => void
}

export function createCustomerColumns(
    actions: ColumnActions
): ColumnDef<CustomerResponse>[] {
    return [
        {
            accessorKey: 'full_name',
            header: 'Name',
            cell: ({ row }) => {
                const customer = row.original
                const initials = customer.full_name
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
                            <span className='font-medium'>{customer.full_name}</span>
                            {customer.customer_code && (
                                <span className='text-xs text-muted-foreground'>
                                    {customer.customer_code}
                                </span>
                            )}
                        </div>
                    </div>
                )
            },
        },
        {
            accessorKey: 'phone',
            header: 'Phone',
            cell: ({ row }) => (
                <span className='text-sm text-muted-foreground'>
                    {row.original.phone}
                </span>
            ),
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
            accessorKey: 'is_active',
            header: 'Status',
            cell: ({ row }) =>
                row.original.is_active ? (
                    <Badge variant='default'>Active</Badge>
                ) : (
                    <Badge variant='secondary'>Inactive</Badge>
                ),
        },
        {
            accessorKey: 'tags',
            header: 'Tags',
            cell: ({ row }) => {
                const tags = row.original.tags
                if (!tags || tags.length === 0) {
                    return (
                        <span className='text-sm text-muted-foreground'>-</span>
                    )
                }
                return (
                    <div className='flex flex-wrap gap-1'>
                        {tags.slice(0, 3).map((tag) => (
                            <Badge key={tag} variant='outline'>
                                {tag}
                            </Badge>
                        ))}
                        {tags.length > 3 && (
                            <span className='text-xs text-muted-foreground'>
                                +{tags.length - 3}
                            </span>
                        )}
                    </div>
                )
            },
        },
        {
            id: 'actions',
            header: 'Actions',
            cell: ({ row }) => {
                const customer = row.original

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
                            {customer.is_active ? (
                                <DropdownMenuItem
                                    onClick={() => actions.onDeactivate(customer)}
                                >
                                    <PowerOff className='mr-2 size-4' />
                                    Deactivate
                                </DropdownMenuItem>
                            ) : (
                                <DropdownMenuItem
                                    onClick={() => actions.onActivate(customer)}
                                >
                                    <Power className='mr-2 size-4' />
                                    Activate
                                </DropdownMenuItem>
                            )}
                            <DropdownMenuItem onClick={() => actions.onEdit(customer)}>
                                <Pencil className='mr-2 size-4' />
                                Edit
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                                onClick={() => actions.onDelete(customer)}
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
