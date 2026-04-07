import { useState } from 'react'
import { ProfileTable } from './components/profile-table'
import { CreateProfileDialog } from './components/create-profile-dialog'
import { EditProfileDialog } from './components/edit-profile-dialog'
import { createColumns } from './data/columns'
import { useProfiles, useDeleteProfile } from '@/hooks/use-profiles'
import { useRouterStore } from '@/stores/router-store'
import type { Profile } from './data/schema'

export function BandwidthProfiles() {
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })
  const [search, setSearch] = useState('')

  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Profile | null>(null)

  const { selectedRouterId, selectedRouterName } = useRouterStore()

  const { data, isLoading } = useProfiles(selectedRouterId)
  const { mutate: deleteProfile } = useDeleteProfile()

  // Client-side filtering logic
  const filteredProfiles = (data?.profiles ?? []).filter((profile) => {
    return (
      search === '' ||
      profile.name.toLowerCase().includes(search.toLowerCase()) ||
      profile.profile_code.toLowerCase().includes(search.toLowerCase())
    )
  })

  const columns = createColumns({
    onEdit: (profile) => setEditTarget(profile),
    onDelete: (profile) => {
      if (selectedRouterId) {
        deleteProfile({ routerId: selectedRouterId, id: profile.id })
      }
    },
  })

  if (!selectedRouterId) {
    return null
  }

  return (
    <div className='mt-8 pt-8 border-t'>
      <div className='mb-6'>
        <h2 className='text-lg font-semibold tracking-tight'>Bandwidth Profiles</h2>
        <p className='text-sm text-muted-foreground'>
          Managing plans for router: <span className="font-medium text-foreground">{selectedRouterName}</span>
        </p>
      </div>
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
    </div>
  )
}
