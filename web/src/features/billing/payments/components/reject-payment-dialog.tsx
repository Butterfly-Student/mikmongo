import { useState } from 'react'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { useRejectPayment } from '@/hooks/use-payments'
import type { PaymentResponse } from '../data/schema'

interface RejectPaymentDialogProps {
  payment: PaymentResponse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function RejectPaymentDialog({
  payment,
  open,
  onOpenChange,
}: RejectPaymentDialogProps) {
  const [reason, setReason] = useState('')
  const { mutate: rejectPayment, isPending } = useRejectPayment()

  function handleConfirm() {
    if (!payment || !reason.trim()) return
    rejectPayment(
      { id: payment.id, reason },
      {
        onSuccess: () => {
          setReason('')
          onOpenChange(false)
        },
      }
    )
  }

  function handleOpenChange(open: boolean) {
    if (!open) {
      setReason('')
    }
    onOpenChange(open)
  }

  return (
    <AlertDialog open={open} onOpenChange={handleOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader className='text-start'>
          <AlertDialogTitle>Tolak Pembayaran?</AlertDialogTitle>
          <AlertDialogDescription>
            Berikan alasan penolakan agar pelanggan dapat mengirimkan ulang
            bukti yang benar.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <div className='space-y-2'>
          <Label htmlFor='reject-reason'>Alasan</Label>
          <Textarea
            id='reject-reason'
            placeholder='Masukkan alasan penolakan...'
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            required
            rows={3}
          />
        </div>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Batal</AlertDialogCancel>
          <Button
            variant='destructive'
            onClick={handleConfirm}
            disabled={!reason.trim() || isPending}
          >
            Tolak
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
