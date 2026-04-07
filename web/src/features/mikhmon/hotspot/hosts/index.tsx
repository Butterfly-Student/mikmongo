import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { HotspotHostsTable } from './components/hotspot-hosts-table'
import { createHotspotHostColumns } from './data/columns'
import { useHotspotHosts } from '@/hooks/use-hotspot'
import { useRouterStore } from '@/stores/router-store'

export function HotspotHosts() {
    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const selectedRouterName = useRouterStore((s) => s.selectedRouterName)
    const [search, setSearch] = useState('')

    const { data: hosts = [], isLoading } = useHotspotHosts(selectedRouterId || null)

    const filteredHosts = hosts.filter((h) => {
        if (!search) return true
        const q = search.toLowerCase()
        return (
            (h.macAddress?.toLowerCase() ?? '').includes(q) ||
            (h.address?.toLowerCase() ?? '').includes(q) ||
            (h.server?.toLowerCase() ?? '').includes(q)
        )
    })

    const columns = createHotspotHostColumns()

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
                        <h2 className='text-lg font-semibold'>Hotspot Hosts</h2>
                        <p className='text-sm text-muted-foreground'>
                            View DHCP hosts connected to the hotspot network
                            {selectedRouterName ? ` — ${selectedRouterName}` : ''}
                        </p>
                    </div>

                    {!selectedRouterId ? (
                        <div className='flex h-48 items-center justify-center rounded-md border border-dashed'>
                            <p className='text-sm text-muted-foreground'>
                                Select a router from the Routers page to view hosts
                            </p>
                        </div>
                    ) : (
                        <HotspotHostsTable
                            columns={columns}
                            data={filteredHosts}
                            isLoading={isLoading}
                            search={search}
                            onSearchChange={setSearch}
                        />
                    )}
                </div>
            </Main>
        </>
    )
}
