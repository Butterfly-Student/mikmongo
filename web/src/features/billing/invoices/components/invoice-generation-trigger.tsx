import { useState } from 'react'
import { FileText } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useTriggerMonthlyBilling } from '@/hooks/use-invoices'

export function InvoiceGenerationTrigger() {
  const [open, setOpen] = useState(false)
  const mutation = useTriggerMonthlyBilling()

  function handleConfirm() {
    mutation.mutate(undefined, {
      onSuccess: () => setOpen(false),
    })
  }

  return (
    <>
      <Button onClick={() => setOpen(true)} size='sm'>
        <FileText className='mr-2 size-4' />
        Buat Tagihan Bulanan
      </Button>

      <ConfirmDialog
        open={open}
        onOpenChange={setOpen}
        title='Buat Tagihan Bulanan?'
        desc='Tindakan ini akan membuat tagihan untuk semua pelanggan aktif bulan ini. Lanjutkan?'
        confirmText='Buat Tagihan'
        cancelBtnText='Batal'
        handleConfirm={handleConfirm}
        isLoading={mutation.isPending}
      />
    </>
  )
}
