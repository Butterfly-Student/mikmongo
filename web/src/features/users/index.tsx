import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { Settings } from 'lucide-react'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { UserTable } from './components/user-table'
import { CreateUserDialog } from './components/create-user-dialog'
import { DeleteUserDialog } from './components/delete-user-dialog'
import { createColumns } from './data/columns'
import { useUsers } from '@/hooks/use-users'
import { useAuthStore } from '@/stores/auth-store'
import type { UserResponse } from '@/api/types'

export function Users() {
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState('all')
  
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<UserResponse | null>(null)

  const { data, isLoading } = useUsers(pagination.pageSize, pagination.pageIndex * pagination.pageSize)
  const { adminUser } = useAuthStore()

  // Client-side filtering logic as API is limited to pagination
  const filteredUsers = (data?.users ?? []).filter((user) => {
    const matchesSearch =
      search === '' ||
      user.full_name.toLowerCase().includes(search.toLowerCase()) ||
      user.email.toLowerCase().includes(search.toLowerCase())
    
    const matchesRole = roleFilter === 'all' || user.role === roleFilter

    return matchesSearch && matchesRole
  })

  // We construct columns dynamically passing the current user ID and delete handler
  const columns = createColumns(adminUser?.id, (user) => setDeleteTarget(user))

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
          <p className='text-sm text-muted-foreground'>Manage admin users and their roles</p>
          <UserTable
            columns={columns}
            data={filteredUsers}
            meta={{ total: filteredUsers.length }} // since client side filtering, pass filtered length if filtering is active, though API meta is technically from server. 
            isLoading={isLoading}
            pagination={pagination}
            onPaginationChange={setPagination}
            onAddUser={() => setCreateDialogOpen(true)}
            search={search}
            onSearchChange={setSearch}
            roleFilter={roleFilter}
            onRoleFilterChange={setRoleFilter}
          />
        </div>
        <CreateUserDialog open={createDialogOpen} onOpenChange={setCreateDialogOpen} />
        <DeleteUserDialog user={deleteTarget} onClose={() => setDeleteTarget(null)} />
      </Main>
    </>
  )
}
