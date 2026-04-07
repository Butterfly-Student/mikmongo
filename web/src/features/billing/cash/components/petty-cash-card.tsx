import { Button } from '@/components/ui/button'
import type { PettyCashFundResponse } from '@/lib/schemas/billing'

interface PettyCashCardProps {
  fund: PettyCashFundResponse | null
  onTopUp: () => void
  onCreate: () => void
}

export function PettyCashCard({ fund, onTopUp, onCreate }: PettyCashCardProps) {
  if (!fund) {
    return (
      <div className='bg-card border rounded-lg p-6'>
        <p className='text-sm text-muted-foreground'>
          Dana kecil belum dikonfigurasi. Buat dana kecil baru.
        </p>
        <Button className='mt-4' onClick={onCreate}>
          Buat Dana Kecil
        </Button>
      </div>
    )
  }

  return (
    <div className='bg-card border rounded-lg p-6'>
      <div className='flex flex-col md:flex-row md:items-center md:justify-between gap-4'>
        <div>
          <p className='text-sm text-muted-foreground'>Saldo Dana Kecil Saat Ini</p>
          <p className='text-3xl font-bold mt-1'>
            Rp {fund.current_balance.toLocaleString('id-ID')}
          </p>
          {fund.fund_name && (
            <p className='text-xs text-muted-foreground mt-1'>{fund.fund_name}</p>
          )}
        </div>
        <Button onClick={onTopUp} variant='outline'>
          Tambah Saldo
        </Button>
      </div>
    </div>
  )
}
