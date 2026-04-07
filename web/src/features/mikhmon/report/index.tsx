import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ReportSummary } from './components/report-summary'
import { ReportTable } from './components/report-table'
import { createReportColumns } from './data/columns'
import { useMikhmonReports, useMikhmonReportSummary } from '@/hooks/use-mikhmon'
import { useRouterStore } from '@/stores/router-store'

export function MikhmonReport() {
    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const selectedRouterName = useRouterStore((s) => s.selectedRouterName)
    const [search, setSearch] = useState('')
    const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 15 })

    const { data: reports = [], isLoading } = useMikhmonReports(selectedRouterId || null)
    const { data: summary, isLoading: summaryLoading } = useMikhmonReportSummary(selectedRouterId || null)

    const filteredReports = reports.filter((r) => {
        if (!search) return true
        const q = search.toLowerCase()
        return (
            (r.user?.toLowerCase() ?? '').includes(q) ||
            (r.profile?.toLowerCase() ?? '').includes(q) ||
            (r.ip?.toLowerCase() ?? '').includes(q) ||
            (r.mac?.toLowerCase() ?? '').includes(q) ||
            (r.comment?.toLowerCase() ?? '').includes(q)
        )
    })

    const columns = createReportColumns()

    return (
        <>
            <Header>
                <Search />
                <div className='ms-auto flex items-center gap-4'>
                    <ThemeSwitch />
                    <ProfileDropdown />
                </div>
            </Header>
            <Main>
                <div className='space-y-6'>
                    <div>
                        <h2 className='text-lg font-semibold'>MikHmon Report</h2>
                        <p className='text-sm text-muted-foreground'>
                            Sales and income report from hotspot voucher transactions
                            {selectedRouterName ? ` — ${selectedRouterName}` : ''}
                        </p>
                    </div>

                    {!selectedRouterId ? (
                        <div className='flex h-48 items-center justify-center rounded-md border border-dashed'>
                            <p className='text-sm text-muted-foreground'>
                                Select a router from the Routers page to view reports
                            </p>
                        </div>
                    ) : (
                        <>
                            <ReportSummary summary={summary} isLoading={summaryLoading} />
                            <ReportTable
                                columns={columns}
                                data={filteredReports}
                                isLoading={isLoading}
                                pagination={pagination}
                                onPaginationChange={setPagination}
                                search={search}
                                onSearchChange={(s) => { setSearch(s); setPagination((p) => ({ ...p, pageIndex: 0 })) }}
                            />
                        </>
                    )}
                </div>
            </Main>
        </>
    )
}
