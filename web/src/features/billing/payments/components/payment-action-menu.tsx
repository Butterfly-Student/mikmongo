import { MoreHorizontal, CheckCircle, XCircle, RefreshCw, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { PaymentResponse } from '../data/schema'

interface PaymentActionMenuProps {
  payment: PaymentResponse
  onConfirm: (payment: PaymentResponse) => void
  onReject: (payment: PaymentResponse) => void
  onRefund: (payment: PaymentResponse) => void
  onGateway: (payment: PaymentResponse) => void
}

export function PaymentActionMenu({
  payment,
  onConfirm,
  onReject,
  onRefund,
  onGateway,
}: PaymentActionMenuProps) {
  const isPending = payment.status === 'pending'
  const isConfirmed = payment.status === 'confirmed'
  const isGateway = payment.payment_method === 'gateway'

  const hasActions = isPending || isConfirmed || isGateway

  if (!hasActions) {
    return null
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant='ghost' size='icon' className='h-8 w-8'>
          <MoreHorizontal className='h-4 w-4' />
          <span className='sr-only'>Buka menu</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end'>
        {isPending && (
          <>
            <DropdownMenuItem onClick={() => onConfirm(payment)}>
              <CheckCircle className='mr-2 h-4 w-4 text-green-600' />
              Konfirmasi Pembayaran
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => onReject(payment)}
              className='text-destructive focus:text-destructive'
            >
              <XCircle className='mr-2 h-4 w-4' />
              Tolak Pembayaran
            </DropdownMenuItem>
          </>
        )}
        {isConfirmed && (
          <DropdownMenuItem
            onClick={() => onRefund(payment)}
            className='text-destructive focus:text-destructive'
          >
            <RefreshCw className='mr-2 h-4 w-4' />
            Kembalikan Dana
          </DropdownMenuItem>
        )}
        {(isPending || isConfirmed) && isGateway && (
          <DropdownMenuSeparator />
        )}
        {isGateway && (
          <DropdownMenuItem onClick={() => onGateway(payment)}>
            <ExternalLink className='mr-2 h-4 w-4' />
            Buka Halaman Pembayaran
          </DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
