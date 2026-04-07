import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Link } from '@tanstack/react-router'
import { Settings, RefreshCw } from 'lucide-react'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { RouterTable } from './components/router-table'
import { CreateRouterDialog } from './components/create-router-dialog'
import { EditRouterDialog } from './components/edit-router-dialog'
import { DeleteRouterDialog } from './components/delete-router-dialog'
import { createColumns } from './data/columns'
import { useRouters, useSyncRouter, useTestRouterConnection, useSelectRouter, useSyncAllRouters } from '@/hooks/use-routers'
import { Button } from '@/components/ui/button'
import type { RouterResponse } from '@/lib/schemas/router'

export function Routers() {
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')

  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<RouterResponse | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<RouterResponse | null>(null)

  const { data, isLoading } = useRouters()
  const { mutate: syncRouter } = useSyncRouter()
  const { mutate: testConnection } = useTestRouterConnection()
  const { mutate: selectRouter } = useSelectRouter()
  const { mutate: syncAllRouters, isPending: isSyncingAll } = useSyncAllRouters()

  // Client-side filtering logic
  const filteredRouters = (data?.routers ?? []).filter((router) => {
    const matchesSearch =
      search === '' ||
      router.name.toLowerCase().includes(search.toLowerCase()) ||
      router.address.toLowerCase().includes(search.toLowerCase())

    const matchesStatus = statusFilter === 'all' || router.status === statusFilter

    return matchesSearch && matchesStatus
  })

  // We construct columns dynamically passing actions
  const columns = createColumns({
    onSync: (router) => syncRouter(router.id),
    onTestConnection: (router) => testConnection(router.id),
    onSelectActive: (router) => selectRouter(router.id),
    onEdit: (router) => setEditTarget(router),
    onDeleteRouter: (router) => setDeleteTarget(router),
  })

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
        <div className='space-y-4'>
          <div className="flex items-center justify-between">
            <p className='text-sm text-muted-foreground'>Manage connected MikroTik routers</p>
            <Button
              variant="outline"
              size="sm"
              onClick={() => syncAllRouters()}
              disabled={isSyncingAll}
            >
              <RefreshCw className={`mr-2 size-4 ${isSyncingAll ? 'animate-spin' : ''}`} />
              {isSyncingAll ? 'Syncing...' : 'Sync All'}
            </Button>
          </div>
          <RouterTable
            columns={columns}
            data={filteredRouters}
            meta={{ total: filteredRouters.length }}
            isLoading={isLoading}
            pagination={pagination}
            onPaginationChange={setPagination}
            onAddRouter={() => setCreateDialogOpen(true)}
            search={search}
            onSearchChange={setSearch}
            statusFilter={statusFilter}
            onStatusFilterChange={setStatusFilter}
          />
        </div>
        <CreateRouterDialog open={createDialogOpen} onOpenChange={setCreateDialogOpen} />
        <EditRouterDialog
          router={editTarget}
          open={!!editTarget}
          onOpenChange={(open) => { if (!open) setEditTarget(null) }}
        />
        <DeleteRouterDialog
          router={deleteTarget}
          open={!!deleteTarget}
          onOpenChange={(open) => { if (!open) setDeleteTarget(null) }}
        />
      </Main>
    </>
  )
}
