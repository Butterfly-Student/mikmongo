import { useState } from 'react'
import { Receipt, Settings } from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { usePayments, useInitiateGatewayPayment } from '@/hooks/use-payments'
import type { PaymentResponse } from './data/schema'
import { PaymentTable } from './components/payment-table'
import { ConfirmPaymentDialog } from './components/confirm-payment-dialog'
import { RejectPaymentDialog } from './components/reject-payment-dialog'
import { RefundPaymentDialog } from './components/refund-payment-dialog'

export default function PaymentsPage() {
  const [selectedPayment, setSelectedPayment] = useState<PaymentResponse | null>(null)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [rejectOpen, setRejectOpen] = useState(false)
  const [refundOpen, setRefundOpen] = useState(false)

  const { data, isLoading } = usePayments()
  const { mutate: initiateGateway } = useInitiateGatewayPayment()

  const payments = data?.data ?? []

  function handleConfirm(payment: PaymentResponse) {
    setSelectedPayment(payment)
    setConfirmOpen(true)
  }

  function handleReject(payment: PaymentResponse) {
    setSelectedPayment(payment)
    setRejectOpen(true)
  }

  function handleRefund(payment: PaymentResponse) {
    setSelectedPayment(payment)
    setRefundOpen(true)
  }

  function handleGateway(payment: PaymentResponse) {
    initiateGateway({ id: payment.id, gateway: 'xendit' })
  }

  return (
    <>
      <Header>
        <Search />
        <div className='ms-auto flex items-center gap-4'>
          <ThemeSwitch />
          <Button
            size='icon'
            variant='ghost'
            asChild
            aria-label='Settings'
            className='rounded-full'
          >
            <Link to='/settings'>
              <Settings />
            </Link>
          </Button>
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='mb-4'>
          <h1 className='text-xl font-semibold tracking-tight'>Riwayat Pembayaran</h1>
          <p className='text-sm text-muted-foreground'>
            Konfirmasi, tolak, atau kembalikan dana pembayaran pelanggan
          </p>
        </div>
        {isLoading ? (
          <div className='space-y-2'>
            {Array.from({ length: 6 }).map((_, i) => (
              <div
                key={i}
                className='h-12 w-full animate-pulse rounded-md bg-muted'
              />
            ))}
          </div>
        ) : payments.length === 0 ? (
          <div className='flex flex-col items-center justify-center rounded-md border p-16'>
            <Receipt className='size-12 text-muted-foreground/40' />
            <div className='mt-4 text-sm text-muted-foreground'>
              Tidak ada pembayaran ditemukan.
            </div>
          </div>
        ) : (
          <PaymentTable
            data={payments}
            onConfirm={handleConfirm}
            onReject={handleReject}
            onRefund={handleRefund}
            onGateway={handleGateway}
          />
        )}
      </Main>

      <ConfirmPaymentDialog
        payment={selectedPayment}
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
      />
      <RejectPaymentDialog
        payment={selectedPayment}
        open={rejectOpen}
        onOpenChange={setRejectOpen}
      />
      <RefundPaymentDialog
        payment={selectedPayment}
        open={refundOpen}
        onOpenChange={setRefundOpen}
      />
    </>
  )
}
