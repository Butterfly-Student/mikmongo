import type { ColumnDef } from '@tanstack/react-table'
import { Check, X } from 'lucide-react'
import { format } from 'date-fns'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import type { CashEntryResponse } from '@/lib/schemas/billing'
import { cashEntryStatuses, cashEntryTypes } from './schema'

export const cashColumns: ColumnDef<CashEntryResponse>[] = [
  {
    accessorKey: 'entry_date',
    header: 'Tanggal',
    enableSorting: true,
    cell: ({ row }) => (
      <span className='text-sm text-muted-foreground'>
        {format(new Date(row.original.entry_date), 'dd/MM/yyyy')}
      </span>
    ),
  },
  {
    accessorKey: 'type',
    header: 'Tipe',
    enableSorting: false,
    enableColumnFilter: true,
    filterFn: (row, id, value) => {
      return (value as string[]).includes(row.getValue(id))
    },
    cell: ({ row }) => {
      const typeInfo = cashEntryTypes.find((t) => t.value === row.original.type)
      return (
        <Badge className={typeInfo?.className ?? ''} variant='outline'>
          {typeInfo?.label ?? row.original.type}
        </Badge>
      )
    },
  },
  {
    accessorKey: 'source',
    header: 'Sumber',
    enableSorting: false,
    cell: ({ row }) => (
      <span className='text-sm text-muted-foreground'>{row.original.source}</span>
    ),
  },
  {
    accessorKey: 'description',
    header: 'Deskripsi',
    enableSorting: false,
    cell: ({ row }) => {
      const desc = row.original.description
      const truncated = desc.length > 40 ? desc.slice(0, 40) + '...' : desc
      return (
        <span className='text-sm' title={desc}>
          {truncated}
        </span>
      )
    },
  },
  {
    accessorKey: 'amount',
    header: 'Jumlah',
    enableSorting: true,
    cell: ({ row }) => (
      <span className='text-sm font-medium text-right block'>
        Rp {row.original.amount.toLocaleString('id-ID')}
      </span>
    ),
  },
  {
    accessorKey: 'status',
    header: 'Status',
    enableColumnFilter: true,
    filterFn: (row, id, value) => {
      return (value as string[]).includes(row.getValue(id))
    },
    cell: ({ row }) => {
      const statusInfo = cashEntryStatuses.find((s) => s.value === row.original.status)
      return (
        <Badge className={statusInfo?.className ?? ''} variant='outline'>
          {statusInfo?.label ?? row.original.status}
        </Badge>
      )
    },
  },
  {
    id: 'actions',
    header: 'Aksi',
    cell: ({ row, table }) => {
      if (row.original.status !== 'pending') return null

      const meta = table.options.meta as {
        onApprove: (id: string) => void
        onReject: (entry: CashEntryResponse) => void
      }

      return (
        <div className='flex items-center gap-1'>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant='ghost'
                size='icon'
                className='h-8 w-8 text-green-600 hover:text-green-700 hover:bg-green-50 dark:hover:bg-green-950'
                onClick={() => meta.onApprove(row.original.id)}
              >
                <Check className='size-4' />
                <span className='sr-only'>Setujui</span>
              </Button>
            </TooltipTrigger>
            <TooltipContent>Setujui</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant='ghost'
                size='icon'
                className='h-8 w-8 text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950'
                onClick={() => meta.onReject(row.original)}
              >
                <X className='size-4' />
                <span className='sr-only'>Tolak</span>
              </Button>
            </TooltipTrigger>
            <TooltipContent>Tolak</TooltipContent>
          </Tooltip>
        </div>
      )
    },
  },
]
