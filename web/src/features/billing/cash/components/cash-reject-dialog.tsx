import { useState } from 'react'
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogCancel,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Loader2 } from 'lucide-react'
import type { CashEntryResponse } from '@/lib/schemas/billing'
import { useRejectCashEntry } from '@/hooks/use-cash'

interface CashRejectDialogProps {
  entry: CashEntryResponse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CashRejectDialog({ entry, open, onOpenChange }: CashRejectDialogProps) {
  const [reason, setReason] = useState('')
  const { mutateAsync: rejectEntry, isPending } = useRejectCashEntry()

  const handleConfirm = async () => {
    if (!entry || !reason.trim()) return

    try {
      await rejectEntry({ id: entry.id, reason: reason.trim() })
      setReason('')
      onOpenChange(false)
    } catch {
      // Error handled by hook
    }
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) setReason('')
    onOpenChange(newOpen)
  }

  return (
    <AlertDialog open={open} onOpenChange={handleOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Tolak Entri Kas</AlertDialogTitle>
          <AlertDialogDescription>
            Berikan alasan penolakan untuk entri ini.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <div className='space-y-2 py-2'>
          <Label htmlFor='reject-reason'>Alasan</Label>
          <Textarea
            id='reject-reason'
            placeholder='Masukkan alasan penolakan...'
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            rows={3}
            disabled={isPending}
          />
        </div>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Batal</AlertDialogCancel>
          <Button
            variant='destructive'
            onClick={handleConfirm}
            disabled={!reason.trim() || isPending}
          >
            {isPending && <Loader2 className='mr-2 size-4 animate-spin' />}
            Tolak
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
