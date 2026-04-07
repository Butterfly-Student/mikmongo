import { ConfirmDialog } from '@/components/confirm-dialog'
import { useRefundPayment } from '@/hooks/use-payments'
import type { PaymentResponse } from '../data/schema'

interface RefundPaymentDialogProps {
  payment: PaymentResponse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function RefundPaymentDialog({
  payment,
  open,
  onOpenChange,
}: RefundPaymentDialogProps) {
  const { mutate: refundPayment, isPending } = useRefundPayment()

  function handleConfirm() {
    if (!payment) return
    refundPayment(
      { id: payment.id, amount: payment.amount, reason: 'Admin refund' },
      {
        onSuccess: () => {
          onOpenChange(false)
        },
      }
    )
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      title='Kembalikan Dana?'
      desc={
        payment
          ? `Dana sebesar Rp ${payment.amount.toLocaleString('id-ID')} akan dikembalikan kepada pelanggan. Ini tidak dapat dibatalkan.`
          : 'Dana akan dikembalikan kepada pelanggan. Ini tidak dapat dibatalkan.'
      }
      confirmText='Kembalikan Dana'
      cancelBtnText='Batal'
      destructive={true}
      handleConfirm={handleConfirm}
      isLoading={isPending}
    />
  )
}
