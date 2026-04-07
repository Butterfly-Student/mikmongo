import { useState } from 'react'
import { FileText, Settings } from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { useInvoices } from '@/hooks/use-invoices'
import type { InvoiceResponse } from '@/lib/schemas/billing'
import { InvoiceTable } from './components/invoice-table'
import { InvoiceDetailSheet } from './components/invoice-detail-sheet'
import { InvoiceGenerationTrigger } from './components/invoice-generation-trigger'

export default function InvoicesPage() {
  const [selectedInvoice, setSelectedInvoice] = useState<InvoiceResponse | null>(null)
  const [sheetOpen, setSheetOpen] = useState(false)

  const { data, isLoading } = useInvoices()

  const invoices = data?.data ?? []

  function handleRowClick(invoice: InvoiceResponse) {
    setSelectedInvoice(invoice)
    setSheetOpen(true)
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
        <div className='mb-4 flex items-start justify-between'>
          <div>
            <h1 className='text-xl font-semibold tracking-tight'>Manajemen Tagihan</h1>
            <p className='text-sm text-muted-foreground'>
              Kelola tagihan pelanggan dan buat tagihan bulanan
            </p>
          </div>
          <InvoiceGenerationTrigger />
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
        ) : invoices.length === 0 ? (
          <div className='flex flex-col items-center justify-center rounded-md border p-16'>
            <FileText className='size-12 text-muted-foreground/40' />
            <div className='mt-4 text-sm text-muted-foreground'>
              Tidak ada tagihan ditemukan.
            </div>
          </div>
        ) : (
          <InvoiceTable data={invoices} onRowClick={handleRowClick} />
        )}
      </Main>

      <InvoiceDetailSheet
        invoice={selectedInvoice}
        open={sheetOpen}
        onOpenChange={setSheetOpen}
      />
    </>
  )
}
