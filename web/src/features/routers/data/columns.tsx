import { MoreHorizontal, Router as RouterIcon, Trash, Power, RefreshCw, Activity, Pencil } from 'lucide-react'
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
import type { RouterResponse } from '@/lib/schemas/router'

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

interface ColumnActions {
  onSync: (router: RouterResponse) => void
  onTestConnection: (router: RouterResponse) => void
  onSelectActive: (router: RouterResponse) => void
  onEdit: (router: RouterResponse) => void
  onDeleteRouter: (router: RouterResponse) => void
}

export function createColumns(actions: ColumnActions): ColumnDef<RouterResponse>[] {
  return [
    {
      accessorKey: 'name',
      header: 'Router Name',
      cell: ({ row }) => {
        const router = row.original
        return (
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg border bg-background text-muted-foreground">
              <RouterIcon className="size-5" />
            </div>
            <div className="flex flex-col">
              <span className="font-medium">{router.name}</span>
              <span className="text-xs text-muted-foreground max-w-[200px] truncate">
                {router.address}
              </span>
            </div>
          </div>
        )
      },
    },
    {
      accessorKey: 'area',
      header: 'Area',
      cell: ({ row }) => {
        const area = row.original.area
        return area ? <Badge variant="outline">{area}</Badge> : <span className="text-muted-foreground">-</span>
      }
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }) => {
        const status = row.original.status
        if (status === 'online') return <Badge variant="success">Online</Badge> // using success variant assuming standard shadcn has it or using default if not. In the users table they used 'default' for active. Let's stick to default/secondary.
        return <Badge variant={status === 'online' ? 'default' : 'secondary'}>{status.charAt(0).toUpperCase() + status.slice(1)}</Badge>
      },
    },
    {
      accessorKey: 'use_ssl',
      header: 'API SSL',
      cell: ({ row }) => (
        row.original.use_ssl ? <Badge variant="default">Enabled</Badge> : <Badge variant="secondary">Disabled</Badge>
      )
    },
    {
      accessorKey: 'last_seen_at',
      header: 'Last Seen',
      cell: ({ row }) => (
        <span className='text-sm text-muted-foreground'>
          {formatRelativeTime(row.original.last_seen_at)}
        </span>
      ),
    },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }) => {
        const router = row.original

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="size-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem onClick={() => actions.onSelectActive(router)}>
                <Power className="mr-2 size-4" /> Set as Active
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => actions.onSync(router)}>
                <RefreshCw className="mr-2 size-4" /> Sync Router
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => actions.onTestConnection(router)}>
                <Activity className="mr-2 size-4" /> Test Connection
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => actions.onEdit(router)}>
                <Pencil className="mr-2 size-4" /> Edit Router
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => actions.onDeleteRouter(router)}
                className="text-destructive focus:text-destructive"
              >
                <Trash className="mr-2 size-4" /> Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
