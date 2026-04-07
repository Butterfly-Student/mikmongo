import { useState, useMemo } from 'react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { Settings } from 'lucide-react'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { SubscriptionTable } from './components/subscription-table'
import { CreateSubscriptionDialog } from './components/create-subscription-dialog'
import { ConfirmActionDialog } from '@/components/ui/confirm-action-dialog'
import { createSubscriptionColumns } from './data/columns'
import { useRouterStore } from '@/stores/router-store'
import {
    useSubscriptions,
    useActivateSubscription,
    useSuspendSubscription,
    useIsolateSubscription,
    useRestoreSubscription,
    useTerminateSubscription,
    useDeleteSubscription,
} from '@/hooks/use-subscriptions'
import { CreditCard } from 'lucide-react'
import type { SubscriptionResponse } from '@/lib/schemas/subscription'

type ConfirmAction = {
    type: 'activate' | 'suspend' | 'isolate' | 'restore' | 'terminate' | 'delete'
    subscription: SubscriptionResponse
}

const confirmActionConfig: Record<
    ConfirmAction['type'],
    { title: string; description: string; confirmLabel: string; destructive: boolean }
> = {
    activate: {
        title: 'Activate Subscription',
        description:
            'Are you sure you want to activate this subscription? The customer will regain full service access.',
        confirmLabel: 'Activate',
        destructive: false,
    },
    suspend: {
        title: 'Suspend Subscription',
        description:
            'Are you absolutely sure you want to suspend this subscription? This may impact service and billing.',
        confirmLabel: 'Suspend',
        destructive: true,
    },
    isolate: {
        title: 'Isolate Subscription',
        description:
            'Are you sure you want to isolate this subscription? Bandwidth will be limited and the customer will be quarantined.',
        confirmLabel: 'Isolate',
        destructive: true,
    },
    restore: {
        title: 'Restore Subscription',
        description:
            'Are you sure you want to restore this subscription? Service will return to normal operation.',
        confirmLabel: 'Restore',
        destructive: false,
    },
    terminate: {
        title: 'Terminate Subscription',
        description:
            'Are you absolutely sure you want to terminate this subscription? This action cannot be undone. This may impact service and billing.',
        confirmLabel: 'Terminate',
        destructive: true,
    },
    delete: {
        title: 'Delete Subscription',
        description:
            'Are you absolutely sure you want to permanently delete this subscription? This action cannot be undone.',
        confirmLabel: 'Delete',
        destructive: true,
    },
}

export function Subscriptions() {
    const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })
    const [statusFilter, setStatusFilter] = useState('all')
    const [createDialogOpen, setCreateDialogOpen] = useState(false)
    const [confirmAction, setConfirmAction] = useState<ConfirmAction | null>(null)

    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const selectedRouterName = useRouterStore((s) => s.selectedRouterName)

    const { data, isLoading } = useSubscriptions(
        selectedRouterId,
        pagination.pageSize,
        pagination.pageIndex * pagination.pageSize
    )

    const mutateActivate = useActivateSubscription(selectedRouterId ?? '')
    const mutateSuspend = useSuspendSubscription(selectedRouterId ?? '')
    const mutateIsolate = useIsolateSubscription(selectedRouterId ?? '')
    const mutateRestore = useRestoreSubscription(selectedRouterId ?? '')
    const mutateTerminate = useTerminateSubscription(selectedRouterId ?? '')
    const mutateDelete = useDeleteSubscription(selectedRouterId ?? '')

    // Client-side status filtering
    const filteredSubscriptions = useMemo(
        () =>
            (data?.subscriptions ?? []).filter((sub) => {
                if (statusFilter === 'all') return true
                return sub.status === statusFilter
            }),
        [data?.subscriptions, statusFilter]
    )

    const isPending =
        mutateActivate.isPending ||
        mutateSuspend.isPending ||
        mutateIsolate.isPending ||
        mutateRestore.isPending ||
        mutateTerminate.isPending ||
        mutateDelete.isPending

    const handleConfirm = async () => {
        if (!confirmAction) return
        const { type, subscription } = confirmAction

        try {
            switch (type) {
                case 'activate':
                    await mutateActivate.mutateAsync(subscription.id)
                    break
                case 'suspend':
                    await mutateSuspend.mutateAsync({ id: subscription.id })
                    break
                case 'isolate':
                    await mutateIsolate.mutateAsync({ id: subscription.id })
                    break
                case 'restore':
                    await mutateRestore.mutateAsync(subscription.id)
                    break
                case 'terminate':
                    await mutateTerminate.mutateAsync(subscription.id)
                    break
                case 'delete':
                    await mutateDelete.mutateAsync(subscription.id)
                    break
            }
            setConfirmAction(null)
        } catch {
            // Error handled by hooks
        }
    }

    const columns = createSubscriptionColumns({
        onActivate: (sub) => setConfirmAction({ type: 'activate', subscription: sub }),
        onSuspend: (sub) => setConfirmAction({ type: 'suspend', subscription: sub }),
        onIsolate: (sub) => setConfirmAction({ type: 'isolate', subscription: sub }),
        onRestore: (sub) => setConfirmAction({ type: 'restore', subscription: sub }),
        onTerminate: (sub) => setConfirmAction({ type: 'terminate', subscription: sub }),
        onDelete: (sub) => setConfirmAction({ type: 'delete', subscription: sub }),
    })

    const activeConfig = confirmAction
        ? confirmActionConfig[confirmAction.type]
        : null

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
                        Manage subscriptions
                        {selectedRouterName ? ` on ${selectedRouterName}` : ''}.
                        {!selectedRouterId &&
                            ' Select a router from the sidebar to view subscriptions.'}
                    </p>

                    {selectedRouterId ? (
                        <SubscriptionTable
                            columns={columns}
                            data={filteredSubscriptions}
                            meta={{ total: filteredSubscriptions.length }}
                            isLoading={isLoading}
                            pagination={pagination}
                            onPaginationChange={setPagination}
                            onAddSubscription={() => setCreateDialogOpen(true)}
                            statusFilter={statusFilter}
                            onStatusFilterChange={setStatusFilter}
                        />
                    ) : (
                        <div className='flex flex-col items-center justify-center rounded-md border p-16'>
                            <CreditCard className='size-12 text-muted-foreground/40' />
                            <div className='mt-4 text-sm text-muted-foreground'>
                                No router selected
                            </div>
                            <div className='mt-1 text-xs text-muted-foreground'>
                                Select a router from the sidebar to view its subscriptions.
                            </div>
                        </div>
                    )}
                </div>

                {activeConfig && confirmAction && (
                    <ConfirmActionDialog
                        open={!!confirmAction}
                        onOpenChange={(open) => {
                            if (!open) setConfirmAction(null)
                        }}
                        title={activeConfig.title}
                        description={activeConfig.description}
                        confirmLabel={activeConfig.confirmLabel}
                        destructive={activeConfig.destructive}
                        isPending={isPending}
                        onConfirm={handleConfirm}
                    />
                )}

                <CreateSubscriptionDialog
                    open={createDialogOpen}
                    onOpenChange={setCreateDialogOpen}
                />
            </Main>
        </>
    )
}
