import { MoreHorizontal } from 'lucide-react'
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
import type { UserResponse } from '@/api/types'

export const roleDisplayNames: Record<string, string> = {
  superadmin: 'Super Admin',
  admin: 'Admin',
  cs: 'Customer Service',
  billing: 'Billing',
  technician: 'Technician',
  readonly: 'Read Only',
}

function formatRelativeTime(dateString: string | null): string {
  if (!dateString) return 'Never'
  const date = new Date(dateString)
  const diffMs = Date.now() - date.getTime()
  const diffSeconds = Math.floor(diffMs / 1000)
  if (diffSeconds < 60) return 'just now'
  const diffMinutes = Math.floor(diffSeconds / 60)
  if (diffMinutes < 60) return `${diffMinutes} min ago`
  const diffHours = Math.floor(diffMinutes / 60)
  if (diffHours < 24) return `${diffHours} hours ago`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays} days ago`
}

export function createColumns(
  currentUserId: string | null | undefined,
  onDeleteUser: (user: UserResponse) => void
): ColumnDef<UserResponse>[] {
  return [
    {
      accessorKey: 'full_name',
      header: 'Name',
      cell: ({ row }) => {
        const name = row.original.full_name
        const initials = name
          .split(' ')
          .map((n) => n[0])
          .join('')
          .slice(0, 2)
          .toUpperCase()
        return (
          <div className='flex items-center gap-2'>
            <div className='flex size-8 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary'>
              {initials}
            </div>
            <span className='font-medium'>{name}</span>
          </div>
        )
      },
    },
    {
      accessorKey: 'email',
      header: 'Email',
      cell: ({ row }) => (
        <span className='text-sm text-muted-foreground'>{row.original.email}</span>
      ),
    },
    {
      accessorKey: 'role',
      header: 'Role',
      cell: ({ row }) => (
        <Badge variant='default'>
          {roleDisplayNames[row.original.role] ?? row.original.role}
        </Badge>
      ),
    },
    {
      accessorKey: 'is_active',
      header: 'Status',
      cell: ({ row }) => (
        row.original.is_active ? (
          <Badge variant='default'>Active</Badge>
        ) : (
          <Badge variant='secondary'>Inactive</Badge>
        )
      ),
    },
    {
      accessorKey: 'last_login',
      header: 'Last Login',
      cell: ({ row }) => (
        <span className='text-sm text-muted-foreground'>
          {formatRelativeTime(row.original.last_login)}
        </span>
      ),
    },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }) => {
        const user = row.original
        const isSelf = currentUserId === user.id
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant='ghost' size='icon'>
                <MoreHorizontal className='size-4' />
                <span className='sr-only'>Open actions</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align='end'>
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem disabled className='text-muted-foreground'>
                Edit — Coming soon
              </DropdownMenuItem>
              {!isSelf && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className='text-destructive focus:text-destructive'
                    onClick={() => onDeleteUser(user)}
                  >
                    Delete
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
