import { useState } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { CustomerTable } from './components/customer-table'
import { RegistrationTable } from './components/registration-table'
import { CreateCustomerDialog } from './components/create-customer-dialog'
import { DeleteCustomerDialog } from './components/delete-customer-dialog'
import { EditCustomerDialog } from './components/edit-customer-dialog'
import { ApproveRegistrationDialog } from './components/approve-registration-dialog'
import { RejectRegistrationDialog } from './components/reject-registration-dialog'
import { createCustomerColumns } from './data/columns'
import { createRegistrationColumns } from './data/registration-columns'
import {
    useCustomers,
    useActivateCustomer,
    useDeactivateCustomer,
    useRegistrations,
} from '@/hooks/use-customers'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Users, UserPlus, Settings } from 'lucide-react'
import type { CustomerResponse, RegistrationResponse } from '@/lib/schemas/customer'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { ProfileDropdown } from '@/components/profile-dropdown'

export function Customers() {
    const [customerPagination, setCustomerPagination] = useState({
        pageIndex: 0,
        pageSize: 10,
    })
    const [customerSearch, setCustomerSearch] = useState('')
    const [customerStatusFilter, setCustomerStatusFilter] = useState('all')

    const [registrationPagination, setRegistrationPagination] = useState({
        pageIndex: 0,
        pageSize: 10,
    })
    const [registrationSearch, setRegistrationSearch] = useState('')
    const [registrationStatusFilter, setRegistrationStatusFilter] = useState('pending')

    const [createDialogOpen, setCreateDialogOpen] = useState(false)
    const [deleteTarget, setDeleteTarget] = useState<CustomerResponse | null>(null)
    const [editTarget, setEditTarget] = useState<CustomerResponse | null>(null)
    const [approveTarget, setApproveTarget] =
        useState<RegistrationResponse | null>(null)
    const [rejectTarget, setRejectTarget] =
        useState<RegistrationResponse | null>(null)

    const [activeTab, setActiveTab] = useState('customers')

    const { data: customersData, isLoading: customersLoading } = useCustomers(
        customerPagination.pageSize,
        customerPagination.pageIndex * customerPagination.pageSize
    )
    const { mutate: activateCustomer } = useActivateCustomer()
    const { mutate: deactivateCustomer } = useDeactivateCustomer()

    const { data: registrationsData, isLoading: registrationsLoading } =
        useRegistrations(
            registrationPagination.pageSize,
            registrationPagination.pageIndex * registrationPagination.pageSize
        )

    // Client-side filtering for customers
    const filteredCustomers = (customersData?.customers ?? []).filter(
        (customer) => {
            const matchesSearch =
                customerSearch === '' ||
                customer.full_name
                    .toLowerCase()
                    .includes(customerSearch.toLowerCase()) ||
                customer.phone
                    .toLowerCase()
                    .includes(customerSearch.toLowerCase()) ||
                (customer.email?.toLowerCase() ?? '').includes(
                    customerSearch.toLowerCase()
                )

            const matchesStatus =
                customerStatusFilter === 'all' ||
                (customerStatusFilter === 'active' && customer.is_active) ||
                (customerStatusFilter === 'inactive' && !customer.is_active)

            return matchesSearch && matchesStatus
        }
    )

    // Client-side filtering for registrations
    const filteredRegistrations = (registrationsData?.registrations ?? []).filter(
        (reg) => {
            const matchesSearch =
                registrationSearch === '' ||
                reg.full_name
                    .toLowerCase()
                    .includes(registrationSearch.toLowerCase()) ||
                reg.phone
                    .toLowerCase()
                    .includes(registrationSearch.toLowerCase())

            const matchesStatus =
                registrationStatusFilter === 'all' ||
                reg.status === registrationStatusFilter

            return matchesSearch && matchesStatus
        }
    )

    const customerColumns = createCustomerColumns({
        onActivate: (customer) => activateCustomer(customer.id),
        onDeactivate: (customer) => deactivateCustomer(customer.id),
        onDelete: (customer) => setDeleteTarget(customer),
        onEdit: (customer) => setEditTarget(customer),
    })

    const registrationColumns = createRegistrationColumns({
        onApprove: (reg) => setApproveTarget(reg),
        onReject: (reg) => setRejectTarget(reg),
    })

    return (
        <>
            {/* ===== Top Heading ===== */}
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
                        Manage customers and handle registration requests
                    </p>

                    <Tabs value={activeTab} onValueChange={setActiveTab}>
                        <TabsList>
                            <TabsTrigger value='customers' className='gap-2'>
                                <Users className='size-4' />
                                Customers
                            </TabsTrigger>
                            <TabsTrigger value='registrations' className='gap-2'>
                                <UserPlus className='size-4' />
                                Registrations
                            </TabsTrigger>
                        </TabsList>

                        <TabsContent value='customers' className='space-y-4'>
                            <CustomerTable
                                columns={customerColumns}
                                data={filteredCustomers}
                                meta={{ total: filteredCustomers.length }}
                                isLoading={customersLoading}
                                pagination={customerPagination}
                                onPaginationChange={setCustomerPagination}
                                onAddCustomer={() => setCreateDialogOpen(true)}
                                search={customerSearch}
                                onSearchChange={setCustomerSearch}
                                statusFilter={customerStatusFilter}
                                onStatusFilterChange={setCustomerStatusFilter}
                            />
                        </TabsContent>

                        <TabsContent value='registrations' className='space-y-4'>
                            <RegistrationTable
                                columns={registrationColumns}
                                data={filteredRegistrations}
                                meta={{ total: filteredRegistrations.length }}
                                isLoading={registrationsLoading}
                                pagination={registrationPagination}
                                onPaginationChange={setRegistrationPagination}
                                search={registrationSearch}
                                onSearchChange={setRegistrationSearch}
                                statusFilter={registrationStatusFilter}
                                onStatusFilterChange={setRegistrationStatusFilter}
                            />
                        </TabsContent>
                    </Tabs>
                </div>

                <CreateCustomerDialog
                    open={createDialogOpen}
                    onOpenChange={setCreateDialogOpen}
                />
                <EditCustomerDialog
                    customer={editTarget}
                    open={!!editTarget}
                    onOpenChange={(open) => {
                        if (!open) setEditTarget(null)
                    }}
                />
                <DeleteCustomerDialog
                    customer={deleteTarget}
                    open={!!deleteTarget}
                    onOpenChange={(open) => {
                        if (!open) setDeleteTarget(null)
                    }}
                />
                <ApproveRegistrationDialog
                    registration={approveTarget}
                    open={!!approveTarget}
                    onOpenChange={(open) => {
                        if (!open) setApproveTarget(null)
                    }}
                />
                <RejectRegistrationDialog
                    registration={rejectTarget}
                    open={!!rejectTarget}
                    onOpenChange={(open) => {
                        if (!open) setRejectTarget(null)
                    }}
                />
            </Main>
        </>
    )
}
