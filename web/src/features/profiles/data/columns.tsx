import type { ColumnDef } from '@tanstack/react-table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { MoreHorizontal, Pencil, Trash } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { Profile } from './schema'

interface ColumnActions {
  onEdit: (profile: Profile) => void
  onDelete: (profile: Profile) => void
}

export function createColumns(actions: ColumnActions): ColumnDef<Profile>[] {
  return [
    {
      accessorKey: 'name',
      header: 'Profile Name',
      cell: ({ row }) => (
        <div className="flex flex-col">
          <span className="font-medium">{row.original.name}</span>
          <span className="text-xs text-muted-foreground">
            {row.original.profile_code}
          </span>
        </div>
      ),
    },
    {
      accessorKey: 'download_speed',
      header: 'Speed (Mbps)',
      cell: ({ row }) => {
        const dl = row.original.download_speed
        const ul = row.original.upload_speed
        return <span>{dl} / {ul}</span>
      }
    },
    {
      accessorKey: 'price_monthly',
      header: 'Price (Monthly)',
      cell: ({ row }) => {
        const price = new Intl.NumberFormat('id-ID', {
          style: 'currency',
          currency: 'IDR'
        }).format(row.original.price_monthly)
        return <span>{price}</span>
      }
    },
    {
      accessorKey: 'billing_cycle',
      header: 'Cycle',
      cell: ({ row }) => (
        <span className="capitalize">{row.original.billing_cycle}</span>
      )
    },
    {
      accessorKey: 'is_active',
      header: 'Status',
      cell: ({ row }) => (
        row.original.is_active ? <Badge variant="default">Active</Badge> : <Badge variant="secondary">Inactive</Badge>
      )
    },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }) => {
        const profile = row.original

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
              <DropdownMenuItem onClick={() => actions.onEdit(profile)}>
                <Pencil className="mr-2 size-4" /> Edit
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => actions.onDelete(profile)}
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
