import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { Settings, Gauge } from 'lucide-react'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ProfileTable } from '@/features/profiles/components/profile-table'
import { CreateProfileDialog } from '@/features/profiles/components/create-profile-dialog'
import { EditProfileDialog } from '@/features/profiles/components/edit-profile-dialog'
import { createColumns } from '@/features/profiles/data/columns'
import { useProfiles, useDeleteProfile } from '@/hooks/use-profiles'
import { useRouterStore } from '@/stores/router-store'
import type { Profile } from '@/features/profiles/data/schema'

export function BandwidthProfilesPage() {
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })
  const [search, setSearch] = useState('')
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Profile | null>(null)

  const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
  const selectedRouterName = useRouterStore((s) => s.selectedRouterName)

  const { data, isLoading } = useProfiles(selectedRouterId)
  const { mutate: deleteProfile } = useDeleteProfile()

  const filteredProfiles = (data?.profiles ?? []).filter((profile) =>
    search === '' ||
    profile.name.toLowerCase().includes(search.toLowerCase()) ||
    profile.profile_code.toLowerCase().includes(search.toLowerCase())
  )

  const columns = createColumns({
    onEdit: (profile) => setEditTarget(profile),
    onDelete: (profile) => {
      if (selectedRouterId) {
        deleteProfile({ routerId: selectedRouterId, id: profile.id })
      }
    },
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
          <p className='text-sm text-muted-foreground'>
            Manage bandwidth profiles
            {selectedRouterName ? ` for router: ${selectedRouterName}` : ''}.
            {!selectedRouterId && ' Select a router from the Routers page to view its profiles.'}
          </p>

          {selectedRouterId ? (
            <ProfileTable
              columns={columns}
              data={filteredProfiles}
              meta={{ total: filteredProfiles.length }}
              isLoading={isLoading}
              pagination={pagination}
              onPaginationChange={setPagination}
              onAddProfile={() => setCreateDialogOpen(true)}
              search={search}
              onSearchChange={setSearch}
            />
          ) : (
            <div className='flex flex-col items-center justify-center rounded-md border p-16'>
              <Gauge className='size-12 text-muted-foreground/40' />
              <div className='mt-4 text-sm text-muted-foreground'>No router selected</div>
              <div className='mt-1 text-xs text-muted-foreground'>
                Go to the{' '}
                <Link to='/routers' className='text-primary underline-offset-4 hover:underline'>
                  Routers page
                </Link>{' '}
                and set an active router to manage its bandwidth profiles.
              </div>
            </div>
          )}
        </div>

        <CreateProfileDialog
          open={createDialogOpen}
          onOpenChange={setCreateDialogOpen}
          routerId={selectedRouterId}
        />
        <EditProfileDialog
          profile={editTarget}
          routerId={selectedRouterId}
          open={!!editTarget}
          onOpenChange={(open) => { if (!open) setEditTarget(null) }}
        />
      </Main>
    </>
  )
}
