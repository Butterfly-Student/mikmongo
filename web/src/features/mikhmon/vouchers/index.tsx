import { useState, useMemo, useCallback } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { VouchersTable } from './components/vouchers-table'
import { GenerateVoucherDialog } from './components/generate-voucher-dialog'
import { ExpireMonitorDialog } from './components/expire-monitor-dialog'
import { DeleteVoucherDialog } from './components/delete-voucher-dialog'
import { createVoucherColumns } from './data/columns'
import { useMikhmonVouchers } from '@/hooks/use-mikhmon'
import { useHotspotProfiles, useHotspotServers } from '@/hooks/use-hotspot'
import { useRouterStore } from '@/stores/router-store'
import type { VoucherResponse } from '@/lib/schemas/mikhmon'

export function Vouchers() {
    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const selectedRouterName = useRouterStore((s) => s.selectedRouterName)
    const [batchComment, setBatchComment] = useState('')
    const [committedComment, setCommittedComment] = useState('')
    const [search, setSearch] = useState('')
    const [profileFilter, setProfileFilter] = useState('all')
    const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 15 })
    const [generateOpen, setGenerateOpen] = useState(false)
    const [expireMonitorOpen, setExpireMonitorOpen] = useState(false)
    const [deleteTarget, setDeleteTarget] = useState<VoucherResponse | null>(null)

    const { data: vouchers = [], isLoading } = useMikhmonVouchers(selectedRouterId || null, committedComment)
    const { data: profiles = [] } = useHotspotProfiles(selectedRouterId || null)
    const { data: servers = [] } = useHotspotServers(selectedRouterId || null)

    const uniqueProfiles = useMemo(() => 
        [...new Set(vouchers.map((v) => v.profile).filter(Boolean))] as string[]
    , [vouchers])

    const filteredVouchers = useMemo(() => {
        if (!search && profileFilter === 'all') return vouchers
        return vouchers.filter((v) => {
            const matchesProfile = profileFilter === 'all' || v.profile === profileFilter
            if (!search) return matchesProfile
            const q = search.toLowerCase()
            return matchesProfile && (
                v.name.toLowerCase().includes(q) ||
                (v.comment?.toLowerCase() ?? '').includes(q) ||
                (v.code?.toLowerCase() ?? '').includes(q)
            )
        })
    }, [vouchers, search, profileFilter])

    const columns = useMemo(() => createVoucherColumns({
        onDelete: (voucher) => setDeleteTarget(voucher),
    }), [])

    const handleSearchChange = useCallback((s: string) => { 
        setSearch(s) 
        setPagination((p) => ({ ...p, pageIndex: 0 })) 
    }, [])

    const handleProfileFilterChange = useCallback((profile: string) => {
        setProfileFilter(profile)
    }, [])

    const handlePaginationChange = useCallback((newPagination: { pageIndex: number; pageSize: number }) => {
        setPagination(newPagination)
    }, [])

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
                <div className='space-y-4'>
                    <div>
                        <h2 className='text-lg font-semibold'>Vouchers</h2>
                        <p className='text-sm text-muted-foreground'>
                            Generate and manage hotspot vouchers
                            {selectedRouterName ? ` — ${selectedRouterName}` : ''}
                        </p>
                    </div>

                    {!selectedRouterId ? (
                        <div className='flex h-48 items-center justify-center rounded-md border border-dashed'>
                            <p className='text-sm text-muted-foreground'>
                                Select a router from the Routers page to view vouchers
                            </p>
                        </div>
                    ) : (
                        <>
                        <div className='flex gap-2'>
                            <Input
                                placeholder='Enter batch code to search vouchers...'
                                value={batchComment}
                                onChange={(e) => setBatchComment(e.target.value)}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') setCommittedComment(batchComment)
                                }}
                                className='max-w-sm'
                            />
                            <Button variant='outline' onClick={() => setCommittedComment(batchComment)}>
                                Search
                            </Button>
                        </div>
                        <VouchersTable
                            columns={columns}
                            data={filteredVouchers}
                            isLoading={isLoading}
                            pagination={pagination}
                            onPaginationChange={handlePaginationChange}
                            onGenerate={() => setGenerateOpen(true)}
                            onExpireMonitor={() => setExpireMonitorOpen(true)}
                            search={search}
                            onSearchChange={handleSearchChange}
                            profileFilter={profileFilter}
                            onProfileFilterChange={handleProfileFilterChange}
                            profiles={uniqueProfiles}
                        />
                        </>
                    )}
                </div>

                {selectedRouterId && (
                    <>
                        <GenerateVoucherDialog
                            open={generateOpen}
                            onOpenChange={setGenerateOpen}
                            routerId={selectedRouterId}
                            profiles={profiles}
                            servers={servers}
                            onGenerated={(comment) => {
                                setBatchComment(comment)
                                setCommittedComment(comment)
                            }}
                        />
                        <ExpireMonitorDialog
                            open={expireMonitorOpen}
                            onOpenChange={setExpireMonitorOpen}
                            routerId={selectedRouterId}
                        />
                        <DeleteVoucherDialog
                            voucher={deleteTarget}
                            open={!!deleteTarget}
                            onOpenChange={(open) => { if (!open) setDeleteTarget(null) }}
                            routerId={selectedRouterId}
                        />
                    </>
                )}
            </Main>
        </>
    )
}
