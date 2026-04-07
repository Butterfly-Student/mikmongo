import { CalendarIcon, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'

interface DateRangeFilterProps {
  dateFrom: string
  dateTo: string
  onDateFromChange: (val: string) => void
  onDateToChange: (val: string) => void
}

export function DateRangeFilter({
  dateFrom,
  dateTo,
  onDateFromChange,
  onDateToChange,
}: DateRangeFilterProps) {
  const hasFilter = dateFrom || dateTo

  const label =
    dateFrom && dateTo
      ? `${dateFrom} – ${dateTo}`
      : dateFrom
        ? `Dari ${dateFrom}`
        : dateTo
          ? `s/d ${dateTo}`
          : 'Tanggal Pembayaran'

  function handleReset() {
    onDateFromChange('')
    onDateToChange('')
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          size='sm'
          className='h-8 border-dashed'
        >
          <CalendarIcon className='mr-2 h-4 w-4' />
          {label}
        </Button>
      </PopoverTrigger>
      <PopoverContent className='w-72' align='start'>
        <div className='space-y-3'>
          <div className='space-y-1'>
            <Label htmlFor='date-from' className='text-xs'>
              Dari
            </Label>
            <Input
              id='date-from'
              type='date'
              value={dateFrom}
              onChange={(e) => onDateFromChange(e.target.value)}
              className='h-8 text-sm'
            />
          </div>
          <div className='space-y-1'>
            <Label htmlFor='date-to' className='text-xs'>
              Sampai
            </Label>
            <Input
              id='date-to'
              type='date'
              value={dateTo}
              onChange={(e) => onDateToChange(e.target.value)}
              className='h-8 text-sm'
            />
          </div>
          {hasFilter && (
            <Button
              variant='ghost'
              size='sm'
              className='h-8 w-full'
              onClick={handleReset}
            >
              <X className='mr-2 h-4 w-4' />
              Reset Tanggal
            </Button>
          )}
        </div>
      </PopoverContent>
    </Popover>
  )
}
