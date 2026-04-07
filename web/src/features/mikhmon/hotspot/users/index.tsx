import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { HotspotUsersTable } from './components/hotspot-users-table'
import { CreateUserDialog } from './components/create-user-dialog'
import { EditUserDialog } from './components/edit-user-dialog'
import { DeleteUserDialog } from './components/delete-user-dialog'
import { createHotspotUserColumns } from './data/columns'
import {
    useHotspotUsers,
    useHotspotProfiles,
    useHotspotServers,
} from '@/hooks/use-hotspot'
import { useRouterStore } from '@/stores/router-store'
import type { HotspotUser } from '@/lib/schemas/mikrotik'

export function HotspotUsers() {
    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const selectedRouterName = useRouterStore((s) => s.selectedRouterName)
    const [search, setSearch] = useState('')
    const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 15 })
    const [createOpen, setCreateOpen] = useState(false)
    const [editTarget, setEditTarget] = useState<HotspotUser | null>(null)
    const [deleteTarget, setDeleteTarget] = useState<HotspotUser | null>(null)

    const { data: users = [], isLoading } = useHotspotUsers(selectedRouterId || null)
    const { data: profiles = [] } = useHotspotProfiles(selectedRouterId || null)
    const { data: servers = [] } = useHotspotServers(selectedRouterId || null)

    const filteredUsers = users.filter((u) => {
        if (!search) return true
        const q = search.toLowerCase()
        return (
            u.name.toLowerCase().includes(q) ||
            (u.profile?.toLowerCase() ?? '').includes(q) ||
            (u.macAddress?.toLowerCase() ?? '').includes(q) ||
            (u.comment?.toLowerCase() ?? '').includes(q)
        )
    })

    const columns = createHotspotUserColumns({
        onEdit: (user) => setEditTarget(user),
        onDelete: (user) => setDeleteTarget(user),
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
                            <h2 className='text-lg font-semibold'>Hotspot Users</h2>
                            <p className='text-sm text-muted-foreground'>
                                Manage hotspot users on your MikroTik router
                                {selectedRouterName ? ` — ${selectedRouterName}` : ''}
                            </p>
                        </div>
                    </div>

                    {!selectedRouterId ? (
                        <div className='flex h-48 items-center justify-center rounded-md border border-dashed'>
                            <p className='text-sm text-muted-foreground'>
                                Select a router from the Routers page to view hotspot users
                            </p>
                        </div>
                    ) : (
                        <HotspotUsersTable
                            columns={columns}
                            data={filteredUsers}
                            isLoading={isLoading}
                            pagination={pagination}
                            onPaginationChange={setPagination}
                            onAddUser={() => setCreateOpen(true)}
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
                        <CreateUserDialog
                            open={createOpen}
                            onOpenChange={setCreateOpen}
                            routerId={selectedRouterId}
                            profiles={profiles}
                            servers={servers}
                        />
                        <EditUserDialog
                            user={editTarget}
                            open={!!editTarget}
                            onOpenChange={(open) => {
                                if (!open) setEditTarget(null)
                            }}
                            routerId={selectedRouterId}
                            profiles={profiles}
                            servers={servers}
                        />
                        <DeleteUserDialog
                            user={deleteTarget}
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
