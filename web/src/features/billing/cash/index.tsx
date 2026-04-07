import { useState } from 'react'
import { Banknote, Settings } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Link } from '@tanstack/react-router'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { useCashEntries, useApproveCashEntry, usePettyCashFunds } from '@/hooks/use-cash'
import type { CashEntryResponse } from '@/lib/schemas/billing'
import { CashTable } from './components/cash-table'
import { CreateCashEntryDialog } from './components/create-cash-entry-dialog'
import { CashRejectDialog } from './components/cash-reject-dialog'
import { PettyCashCard } from './components/petty-cash-card'
import { TopUpDialog } from './components/top-up-dialog'

export default function CashPage() {
  const [createOpen, setCreateOpen] = useState(false)
  const [rejectTarget, setRejectTarget] = useState<CashEntryResponse | null>(null)
  const [topUpOpen, setTopUpOpen] = useState(false)

  const { data: cashData, isLoading: cashLoading } = useCashEntries()
  const { data: fundsData } = usePettyCashFunds()
  const { mutate: approveEntry } = useApproveCashEntry()

  const entries = cashData?.data ?? []
  const fund = fundsData?.data?.[0] ?? null

  const handleApprove = (id: string) => {
    approveEntry(id)
  }

  const handleReject = (entry: CashEntryResponse) => {
    setRejectTarget(entry)
  }

  const openTopUp = () => setTopUpOpen(true)
  const openCreate = () => setCreateOpen(true)

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
        <div className='space-y-6'>
          <div className='flex items-start justify-between'>
            <div>
              <h1 className='text-xl font-semibold tracking-tight'>Kas & Dana Kecil</h1>
              <p className='text-sm text-muted-foreground'>
                Kelola entri kas dan saldo dana kecil
              </p>
            </div>
            <Button onClick={() => setCreateOpen(true)}>
              Tambah Entri Kas
            </Button>
          </div>
          <PettyCashCard fund={fund} onTopUp={openTopUp} onCreate={openCreate} />

          {cashLoading ? (
            <div className='space-y-2'>
              {Array.from({ length: 6 }).map((_, i) => (
                <div
                  key={i}
                  className='h-12 w-full animate-pulse rounded-md bg-muted'
                />
              ))}
            </div>
          ) : entries.length === 0 ? (
            <div className='flex flex-col items-center justify-center rounded-md border p-16'>
              <Banknote className='size-12 text-muted-foreground/40' />
              <div className='mt-4 text-sm text-muted-foreground'>
                Tidak ada entri kas ditemukan.
              </div>
            </div>
          ) : (
            <CashTable
              data={entries}
              onApprove={handleApprove}
              onReject={handleReject}
            />
          )}
        </div>
      </Main>

      <CreateCashEntryDialog open={createOpen} onOpenChange={setCreateOpen} />

      <CashRejectDialog
        entry={rejectTarget}
        open={rejectTarget !== null}
        onOpenChange={(open) => {
          if (!open) setRejectTarget(null)
        }}
      />

      <TopUpDialog fund={fund} open={topUpOpen} onOpenChange={setTopUpOpen} />
    </>
  )
}
