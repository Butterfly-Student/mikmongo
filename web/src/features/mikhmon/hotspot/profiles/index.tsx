import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { HotspotProfilesTable } from './components/hotspot-profiles-table'
import { CreateProfileDialog } from './components/create-profile-dialog'
import { EditProfileDialog } from './components/edit-profile-dialog'
import { DeleteProfileDialog } from './components/delete-profile-dialog'
import { createHotspotProfileColumns } from './data/columns'
import { useHotspotProfiles } from '@/hooks/use-hotspot'
import { useRouterStore } from '@/stores/router-store'
import type { HotspotProfile } from '@/lib/schemas/mikrotik'

export function HotspotProfiles() {
    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const selectedRouterName = useRouterStore((s) => s.selectedRouterName)
    const [search, setSearch] = useState('')
    const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 15 })
    const [createOpen, setCreateOpen] = useState(false)
    const [editTarget, setEditTarget] = useState<HotspotProfile | null>(null)
    const [deleteTarget, setDeleteTarget] = useState<HotspotProfile | null>(null)

    const { data: profiles = [], isLoading } = useHotspotProfiles(selectedRouterId || null)

    const filteredProfiles = profiles.filter((p) => {
        if (!search) return true
        const q = search.toLowerCase()
        return (
            p.name.toLowerCase().includes(q) ||
            (p.rateLimit?.toLowerCase() ?? '').includes(q) ||
            (p.addressPool?.toLowerCase() ?? '').includes(q)
        )
    })

    const columns = createHotspotProfileColumns({
        onEdit: (profile) => setEditTarget(profile),
        onDelete: (profile) => setDeleteTarget(profile),
    })

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
                    <div className='flex items-center justify-between'>
                        <div>
                            <h2 className='text-lg font-semibold'>Hotspot User Profiles</h2>
                            <p className='text-sm text-muted-foreground'>
                                Manage hotspot user profiles on your MikroTik router
                                {selectedRouterName ? ` — ${selectedRouterName}` : ''}
                            </p>
                        </div>
                    </div>

                    {!selectedRouterId ? (
                        <div className='flex h-48 items-center justify-center rounded-md border border-dashed'>
                            <p className='text-sm text-muted-foreground'>
                                Select a router from the Routers page to view hotspot profiles
                            </p>
                        </div>
                    ) : (
                        <HotspotProfilesTable
                            columns={columns}
                            data={filteredProfiles}
                            isLoading={isLoading}
                            pagination={pagination}
                            onPaginationChange={setPagination}
                            onAddProfile={() => setCreateOpen(true)}
                            search={search}
                            onSearchChange={(s) => {
                                setSearch(s)
                                setPagination((p) => ({ ...p, pageIndex: 0 }))
                            }}
                        />
                    )}
                </div>

                {selectedRouterId && (
                    <>
                        <CreateProfileDialog
                            open={createOpen}
                            onOpenChange={setCreateOpen}
                            routerId={selectedRouterId}
                        />
                        <EditProfileDialog
                            profile={editTarget}
                            open={!!editTarget}
                            onOpenChange={(open) => {
                                if (!open) setEditTarget(null)
                            }}
                            routerId={selectedRouterId}
                        />
                        <DeleteProfileDialog
                            profile={deleteTarget}
                            open={!!deleteTarget}
                            onOpenChange={(open) => {
                                if (!open) setDeleteTarget(null)
                            }}
                            routerId={selectedRouterId}
                        />
                    </>
                )}
            </Main>
        </>
    )
}
