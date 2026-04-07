import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { InvoiceResponse } from '@/lib/schemas/billing'
import { invoiceStatuses } from '../data/schema'

interface InvoiceDetailSheetProps {
  invoice: InvoiceResponse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

function formatRp(amount: number): string {
  return `Rp ${amount.toLocaleString('id-ID')}`
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID')
}

const invoiceTypeLabels: Record<string, string> = {
  recurring: 'Bulanan',
  installation: 'Instalasi',
  additional: 'Tambahan',
  refund: 'Refund',
}

export function InvoiceDetailSheet({ invoice, open, onOpenChange }: InvoiceDetailSheetProps) {
  const statusInfo = invoice
    ? invoiceStatuses.find((s) => s.value === invoice.status)
    : null

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className='w-full sm:max-w-[540px] overflow-y-auto'>
        <SheetHeader>
          <SheetTitle>Detail Tagihan #{invoice?.invoice_number ?? ''}</SheetTitle>
          <SheetDescription>
            Informasi lengkap tagihan pelanggan
          </SheetDescription>
        </SheetHeader>

        {invoice && (
          <div className='px-4 space-y-6'>
            {/* Section 1: Info Tagihan */}
            <div>
              <h3 className='text-sm font-semibold mb-3 text-foreground'>Info Tagihan</h3>
              <div className='grid grid-cols-2 gap-4'>
                <div>
                  <p className='text-xs text-muted-foreground'>Pelanggan</p>
                  <p className='text-sm font-mono'>{invoice.customer_id.slice(0, 8)}...</p>
                </div>
                <div>
                  <p className='text-xs text-muted-foreground'>Periode</p>
                  <p className='text-sm'>
                    {formatDate(invoice.billing_period_start)} -{' '}
                    {formatDate(invoice.billing_period_end)}
                  </p>
                </div>
                <div>
                  <p className='text-xs text-muted-foreground'>Tanggal Tagihan</p>
                  <p className='text-sm'>{formatDate(invoice.issue_date)}</p>
                </div>
                <div>
                  <p className='text-xs text-muted-foreground'>Jatuh Tempo</p>
                  <p className='text-sm'>{formatDate(invoice.due_date)}</p>
                </div>
                <div>
                  <p className='text-xs text-muted-foreground'>Jenis Tagihan</p>
                  <p className='text-sm'>{invoiceTypeLabels[invoice.invoice_type] ?? invoice.invoice_type}</p>
                </div>
              </div>
            </div>

            {/* Section 2: Jumlah */}
            <div>
              <h3 className='text-sm font-semibold mb-3 text-foreground'>Jumlah</h3>
              <div className='space-y-2'>
                <div className='flex justify-between text-sm'>
                  <span className='text-muted-foreground'>Subtotal</span>
                  <span>{formatRp(invoice.subtotal)}</span>
                </div>
                {invoice.tax_amount > 0 && (
                  <div className='flex justify-between text-sm'>
                    <span className='text-muted-foreground'>Pajak</span>
                    <span>{formatRp(invoice.tax_amount)}</span>
                  </div>
                )}
                {invoice.discount_amount > 0 && (
                  <div className='flex justify-between text-sm'>
                    <span className='text-muted-foreground'>Diskon</span>
                    <span className='text-green-600'>-{formatRp(invoice.discount_amount)}</span>
                  </div>
                )}
                {invoice.late_fee > 0 && (
                  <div className='flex justify-between text-sm'>
                    <span className='text-muted-foreground'>Denda Keterlambatan</span>
                    <span className='text-red-600'>{formatRp(invoice.late_fee)}</span>
                  </div>
                )}
                <div className='border-t pt-2 flex justify-between'>
                  <span className='font-semibold'>Total</span>
                  <span className='text-xl font-bold'>{formatRp(invoice.total_amount)}</span>
                </div>
              </div>
            </div>

            {/* Section 3: Status */}
            <div>
              <h3 className='text-sm font-semibold mb-3 text-foreground'>Status</h3>
              <Badge className={statusInfo?.className ?? ''} variant='outline'>
                {statusInfo?.label ?? invoice.status}
              </Badge>
            </div>

            {/* Section 4: Riwayat Pembayaran */}
            <div>
              <h3 className='text-sm font-semibold mb-3 text-foreground'>Riwayat Pembayaran</h3>
              <p className='text-sm text-muted-foreground'>
                Lihat halaman pembayaran untuk detail
              </p>
            </div>
          </div>
        )}

        <SheetFooter>
          <Button variant='outline' onClick={() => onOpenChange(false)}>
            Tutup
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  )
}
