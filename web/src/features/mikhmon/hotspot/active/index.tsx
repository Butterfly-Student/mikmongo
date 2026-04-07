import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { HotspotActiveTable } from './components/hotspot-active-table'
import { DisconnectDialog } from './components/disconnect-dialog'
import { createHotspotActiveColumns } from './data/columns'
import { useHotspotActive } from '@/hooks/use-hotspot'
import { useRouterStore } from '@/stores/router-store'
import type { HotspotActive } from '@/lib/schemas/mikrotik'

export function HotspotActiveSessions() {
    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const selectedRouterName = useRouterStore((s) => s.selectedRouterName)
    const [search, setSearch] = useState('')
    const [serverFilter, setServerFilter] = useState('all')
    const [disconnectTarget, setDisconnectTarget] = useState<HotspotActive | null>(null)

    const { data: sessions = [], isLoading } = useHotspotActive(selectedRouterId || null)

    const uniqueServers = [...new Set(sessions.map((s) => s.server).filter(Boolean))] as string[]

    const filteredSessions = sessions.filter((s) => {
        const matchesServer = serverFilter === 'all' || s.server === serverFilter
        if (!search) return matchesServer
        const q = search.toLowerCase()
        return (
            matchesServer &&
            (
                (s.user?.toLowerCase() ?? '').includes(q) ||
                (s.address?.toLowerCase() ?? '').includes(q) ||
                (s.macAddress?.toLowerCase() ?? '').includes(q)
            )
        )
    })

    const columns = createHotspotActiveColumns({
        onDisconnect: (session) => setDisconnectTarget(session),
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
                    <div>
                        <h2 className='text-lg font-semibold'>Active Sessions</h2>
                        <p className='text-sm text-muted-foreground'>
                            View and manage currently active hotspot sessions (auto-refreshes every 30s)
                            {selectedRouterName ? ` — ${selectedRouterName}` : ''}
                        </p>
                    </div>

                    {!selectedRouterId ? (
                        <div className='flex h-48 items-center justify-center rounded-md border border-dashed'>
                            <p className='text-sm text-muted-foreground'>
                                Select a router from the Routers page to view active sessions
                            </p>
                        </div>
                    ) : (
                        <HotspotActiveTable
                            columns={columns}
                            data={filteredSessions}
                            isLoading={isLoading}
                            search={search}
                            onSearchChange={setSearch}
                            serverFilter={serverFilter}
                            onServerFilterChange={setServerFilter}
                            servers={uniqueServers}
                        />
                    )}
                </div>

                {selectedRouterId && (
                    <DisconnectDialog
                        session={disconnectTarget}
                        open={!!disconnectTarget}
                        onOpenChange={(open) => {
                            if (!open) setDisconnectTarget(null)
                        }}
                        routerId={selectedRouterId}
                    />
                )}
            </Main>
        </>
    )
}
