import { useState } from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Loader2 } from 'lucide-react'
import type { PettyCashFundResponse } from '@/lib/schemas/billing'
import { useTopUpPettyCashFund } from '@/hooks/use-cash'

interface TopUpDialogProps {
  fund: PettyCashFundResponse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function TopUpDialog({ fund, open, onOpenChange }: TopUpDialogProps) {
  const [amount, setAmount] = useState<number>(0)
  const { mutateAsync: topUp, isPending } = useTopUpPettyCashFund()

  const handleConfirm = async () => {
    if (!fund || !amount || amount <= 0) return

    try {
      await topUp({ id: fund.id, amount })
      setAmount(0)
      onOpenChange(false)
    } catch {
      // Error handled by hook
    }
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) setAmount(0)
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className='sm:max-w-[400px]'>
        <DialogHeader>
          <DialogTitle>Tambah Saldo Dana Kecil</DialogTitle>
          <DialogDescription>
            Masukkan jumlah dan keterangan penambahan saldo.
          </DialogDescription>
        </DialogHeader>
        <div className='space-y-3 py-2'>
          <div className='space-y-2'>
            <Label htmlFor='top-up-amount'>Jumlah (Rp) *</Label>
            <Input
              id='top-up-amount'
              type='number'
              min={1}
              placeholder='0'
              value={amount || ''}
              onChange={(e) => setAmount(Number(e.target.value))}
              disabled={isPending}
            />
          </div>
        </div>
        <DialogFooter>
          <Button
            type='button'
            variant='ghost'
            onClick={() => handleOpenChange(false)}
            disabled={isPending}
          >
            Batal
          </Button>
          <Button
            onClick={handleConfirm}
            disabled={!amount || amount <= 0 || isPending}
          >
            {isPending && <Loader2 className='mr-2 size-4 animate-spin' />}
            Tambah Saldo
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
