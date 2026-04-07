import { useState } from 'react'
import { BarChart3, Settings } from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Link } from '@tanstack/react-router'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { useReportSummary } from '@/hooks/use-report-summary'

function formatRupiah(n: number): string {
  return 'Rp ' + n.toLocaleString('id-ID')
}

function currentMonthRange(): { from: string; to: string } {
  const now = new Date()
  const from = new Date(now.getFullYear(), now.getMonth(), 1)
    .toISOString()
    .split('T')[0]
  const to = new Date(now.getFullYear(), now.getMonth() + 1, 0)
    .toISOString()
    .split('T')[0]
  return { from, to }
}

interface StatCardProps {
  title: string
  value: string
  description?: string
  isLoading: boolean
}

function StatCard({ title, value, description, isLoading }: StatCardProps) {
  return (
    <Card>
      <CardHeader className='pb-2'>
        <CardTitle className='text-sm font-medium text-muted-foreground'>
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <Skeleton className='h-8 w-24' />
        ) : (
          <>
            <div className='text-2xl font-bold'>{value}</div>
            {description && (
              <p className='mt-1 text-xs text-muted-foreground'>{description}</p>
            )}
          </>
        )}
      </CardContent>
    </Card>
  )
}

export default function BusinessReportsPage() {
  const defaults = currentMonthRange()
  const [from, setFrom] = useState(defaults.from)
  const [to, setTo] = useState(defaults.to)
  const [appliedFrom, setAppliedFrom] = useState(defaults.from)
  const [appliedTo, setAppliedTo] = useState(defaults.to)

  const { data, isLoading } = useReportSummary(appliedFrom, appliedTo)

  function handleApply() {
    setAppliedFrom(from)
    setAppliedTo(to)
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
        <div className='mb-6 flex items-start justify-between'>
          <div>
            <h1 className='text-xl font-semibold tracking-tight'>Laporan Bisnis</h1>
            <p className='text-sm text-muted-foreground'>
              Ringkasan pendapatan, tagihan, pelanggan, dan langganan
            </p>
          </div>
        </div>

        {/* Date range filter */}
        <div className='mb-6 flex flex-wrap items-end gap-4 rounded-md border p-4'>
          <div className='grid gap-1.5'>
            <Label htmlFor='from'>Dari</Label>
            <Input
              id='from'
              type='date'
              value={from}
              onChange={(e) => setFrom(e.target.value)}
              className='w-40'
            />
          </div>
          <div className='grid gap-1.5'>
            <Label htmlFor='to'>Sampai</Label>
            <Input
              id='to'
              type='date'
              value={to}
              onChange={(e) => setTo(e.target.value)}
              className='w-40'
            />
          </div>
          <Button onClick={handleApply} disabled={isLoading}>
            Tampilkan
          </Button>
          {data && (
            <p className='text-xs text-muted-foreground'>
              Periode: {data.period_start} — {data.period_end}
            </p>
          )}
        </div>

        {/* No data state */}
        {!isLoading && !data && (
          <div className='flex flex-col items-center justify-center rounded-md border p-16'>
            <BarChart3 className='size-12 text-muted-foreground/40' />
            <div className='mt-4 text-sm text-muted-foreground'>
              Tidak ada data laporan untuk periode yang dipilih.
            </div>
          </div>
        )}

        {(isLoading || data) && (
          <div className='space-y-6'>
            {/* Revenue */}
            <section>
              <h2 className='mb-3 text-sm font-semibold uppercase tracking-wider text-muted-foreground'>
                Pendapatan
              </h2>
              <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
                <StatCard
                  title='Total Pendapatan'
                  value={data ? formatRupiah(data.total_revenue) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Total Ditagihkan'
                  value={data ? formatRupiah(data.total_invoiced) : '--'}
                  isLoading={isLoading}
                />
              </div>
            </section>

            {/* Invoices */}
            <section>
              <h2 className='mb-3 text-sm font-semibold uppercase tracking-wider text-muted-foreground'>
                Tagihan
              </h2>
              <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
                <StatCard
                  title='Total Tagihan'
                  value={data ? String(data.total_invoices) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Tagihan Lunas'
                  value={data ? String(data.paid_invoices) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Tagihan Belum Lunas'
                  value={data ? String(data.unpaid_invoices) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Tagihan Jatuh Tempo'
                  value={data ? String(data.overdue_invoices) : '--'}
                  description={data ? `${data.total_payments} pembayaran diterima` : undefined}
                  isLoading={isLoading}
                />
              </div>
            </section>

            {/* Customers */}
            <section>
              <h2 className='mb-3 text-sm font-semibold uppercase tracking-wider text-muted-foreground'>
                Pelanggan
              </h2>
              <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
                <StatCard
                  title='Total Pelanggan'
                  value={data ? String(data.total_customers) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Pelanggan Aktif'
                  value={data ? String(data.active_customers) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Pelanggan Baru'
                  value={data ? String(data.new_customers) : '--'}
                  isLoading={isLoading}
                />
              </div>
            </section>

            {/* Subscriptions */}
            <section>
              <h2 className='mb-3 text-sm font-semibold uppercase tracking-wider text-muted-foreground'>
                Langganan
              </h2>
              <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
                <StatCard
                  title='Langganan Aktif'
                  value={data ? String(data.active_subscriptions) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Langganan Diisolir'
                  value={data ? String(data.isolated_subscriptions) : '--'}
                  isLoading={isLoading}
                />
                <StatCard
                  title='Langganan Disuspend'
                  value={data ? String(data.suspended_subscriptions) : '--'}
                  isLoading={isLoading}
                />
              </div>
            </section>
          </div>
        )}
      </Main>
    </>
  )
}
