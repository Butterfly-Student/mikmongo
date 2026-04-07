import { Users, Wifi, Banknote, Router, Settings } from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { useReportSummary } from '@/hooks/use-report-summary'
import { useRouters } from '@/hooks/use-routers'
import { KPICard } from './components/kpi-card'
import { RevenueChart } from './components/revenue-chart'
import { RouterHealthCards } from './components/router-health-cards'
import { RecentActivityFeed } from './components/recent-activity-feed'
import { Link } from '@tanstack/react-router'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Button } from '@/components/ui/button'
import { ThemeSwitch } from '@/components/theme-switch'
import { Search } from '@/components/search'
import { TopNav } from '@/components/layout/top-nav'

function formatRupiah(n: number): string {
  return 'Rp ' + n.toLocaleString('id-ID')
}

function TotalCustomersCard() {
  const { data, isLoading } = useReportSummary()
  return (
    <KPICard
      title='Total Customers'
      value={data ? String(data.total_customers) : '--'}
      trend={data ? `+${data.new_customers} pelanggan baru` : ''}
      icon={Users}
      isLoading={isLoading}
    />
  )
}

function ActiveSubscriptionsCard() {
  const { data, isLoading } = useReportSummary()
  return (
    <KPICard
      title='Active Subscriptions'
      value={data ? String(data.active_subscriptions) : '--'}
      trend={data ? `${data.suspended_subscriptions} suspended` : ''}
      icon={Wifi}
      isLoading={isLoading}
    />
  )
}

function MonthlyRevenueCard() {
  const { data, isLoading } = useReportSummary()
  return (
    <KPICard
      title='Monthly Revenue'
      value={data ? formatRupiah(data.total_revenue) : '--'}
      trend={data ? `${data.paid_invoices} of ${data.total_invoices} invoices paid` : ''}
      icon={Banknote}
      isLoading={isLoading}
    />
  )
}

function ActiveRoutersCard() {
  const { data, isLoading } = useRouters()
  const routers = data?.routers ?? []
  const onlineCount = routers.filter((r) => r.status === 'online').length
  const offlineCount = routers.filter((r) => r.status === 'offline').length
  const totalCount = routers.length

  return (
    <KPICard
      title='Active Routers'
      value={isLoading ? '--' : `${onlineCount}/${totalCount}`}
      trend={isLoading ? '' : `${offlineCount} offline`}
      icon={Router}
      isLoading={isLoading}
    />
  )
}

export function Dashboard() {
  return (
    <>
      {/* ===== Top Heading ===== */}
      <Header>
        <TopNav links={topNav} />
        <div className='ms-auto flex items-center space-x-1 sm:space-x-4'>
          <Search />
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
          {/* KPI Widgets Grid */}
          <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
            <TotalCustomersCard />
            <ActiveSubscriptionsCard />
            <MonthlyRevenueCard />
            <ActiveRoutersCard />
          </div>

          {/* Revenue Chart */}
          <RevenueChart />

          {/* Router Health + Activity Feed */}
          <div className='grid gap-4 lg:grid-cols-7'>
            <div className='lg:col-span-4'>
              <RouterHealthCards />
            </div>
            <div className='lg:col-span-3'>
              <RecentActivityFeed />
            </div>
          </div>
        </div>
      </Main>
    </>
  )
}

const topNav = [
  {
    title: 'Overview',
    href: 'dashboard/overview',
    isActive: true,
    disabled: false,
  },
  {
    title: 'Customers',
    href: 'dashboard/customers',
    isActive: false,
    disabled: true,
  },
  {
    title: 'Products',
    href: 'dashboard/products',
    isActive: false,
    disabled: true,
  },
  {
    title: 'Settings',
    href: 'dashboard/settings',
    isActive: false,
    disabled: true,
  },
]
