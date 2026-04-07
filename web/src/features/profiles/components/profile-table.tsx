import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Plus } from 'lucide-react'
import type { Profile } from '../data/schema'

interface ProfileTableProps {
  columns: ColumnDef<Profile>[]
  data: Profile[]
  meta?: { total: number }
  isLoading: boolean
  pagination: { pageIndex: number; pageSize: number }
  onPaginationChange: (pagination: { pageIndex: number; pageSize: number }) => void
  onAddProfile: () => void
  search: string
  onSearchChange: (search: string) => void
}

export function ProfileTable({
  columns,
  data,
  meta,
  isLoading,
  pagination,
  onPaginationChange,
  onAddProfile,
  search,
  onSearchChange,
}: ProfileTableProps) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    pageCount: Math.ceil((meta?.total ?? 0) / pagination.pageSize),
  })

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between'>
        <div className='flex flex-1 items-center space-x-2'>
          <Input
            placeholder='Search profiles...'
            value={search}
            onChange={(event) => onSearchChange(event.target.value)}
            className='h-8 w-[150px] lg:w-[250px]'
          />
        </div>
        <Button onClick={onAddProfile} size='sm' className='h-8'>
          <Plus className='mr-2 size-4' />
          Add Profile
        </Button>
      </div>
      <div className='rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={i}>
                  {columns.map((_, colIndex) => (
                    <TableCell key={colIndex}>
                      <Skeleton className='h-6 w-full' />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center'
                >
                  <div className='flex flex-col items-center justify-center p-8'>
                    <div className='text-sm text-muted-foreground'>
                      No profiles found. Create a new bandwidth profile.
                    </div>
                  </div>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className='flex items-center justify-between px-2'>
        <div className='flex-1 text-sm text-muted-foreground'>
          Showing {data.length} of {meta?.total ?? 0} profiles.
        </div>
        <div className='flex items-center space-x-2'>
          <Button
            variant='outline'
            size='sm'
            onClick={() => onPaginationChange({ ...pagination, pageIndex: pagination.pageIndex - 1 })}
            disabled={pagination.pageIndex === 0}
          >
            Previous
          </Button>
          <Button
            variant='outline'
            size='sm'
            onClick={() => onPaginationChange({ ...pagination, pageIndex: pagination.pageIndex + 1 })}
            disabled={pagination.pageIndex >= table.getPageCount() - 1}
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  )
}
