import { ConfirmDialog } from '@/components/confirm-dialog'
import { useConfirmPayment } from '@/hooks/use-payments'
import type { PaymentResponse } from '../data/schema'

interface ConfirmPaymentDialogProps {
  payment: PaymentResponse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function ConfirmPaymentDialog({
  payment,
  open,
  onOpenChange,
}: ConfirmPaymentDialogProps) {
  const { mutate: confirmPayment, isPending } = useConfirmPayment()

  function handleConfirm() {
    if (!payment) return
    confirmPayment(payment.id, {
      onSuccess: () => {
        onOpenChange(false)
      },
    })
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      title='Konfirmasi Pembayaran?'
      desc='Pembayaran ini akan ditandai sebagai dikonfirmasi. Tindakan ini tidak dapat dibatalkan.'
      confirmText='Konfirmasi'
      cancelBtnText='Batal'
      destructive={false}
      handleConfirm={handleConfirm}
      isLoading={isPending}
    />
  )
}
