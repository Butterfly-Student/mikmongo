import { usePortalSubscriptions } from '@/hooks/use-customer-portal'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Wifi, WifiOff } from 'lucide-react'
import type { SubscriptionResponse } from '@/lib/schemas/subscription'

function formatDate(dateString: string | null): string {
    if (!dateString) return '-'
    return new Intl.DateTimeFormat('id-ID', {
        dateStyle: 'long',
    }).format(new Date(dateString))
}

function getStatusBadge(status: SubscriptionResponse['status']) {
    const variants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
        active: 'default',
        pending: 'outline',
        suspended: 'secondary',
        isolated: 'outline',
        expired: 'secondary',
        terminated: 'destructive',
    }
    return <Badge variant={variants[status] ?? 'secondary'}>{status.charAt(0).toUpperCase() + status.slice(1)}</Badge>
}

export function CustomerPortalSubscriptions() {
    const { data: subscriptions, isLoading } = usePortalSubscriptions()

    if (isLoading) {
        return (
            <div className="space-y-6">
                <div>
                    <h1 className="text-2xl font-semibold tracking-tight">My Subscriptions</h1>
                    <p className="text-sm text-muted-foreground">View your active and past subscriptions</p>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[1, 2, 3].map((i) => (
                        <Card key={i}>
                            <CardHeader><Skeleton className="h-6 w-32" /></CardHeader>
                            <CardContent className="space-y-2">
                                <Skeleton className="h-4 w-full" />
                                <Skeleton className="h-4 w-3/4" />
                                <Skeleton className="h-4 w-1/2" />
                            </CardContent>
                        </Card>
                    ))}
                </div>
            </div>
        )
    }

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-semibold tracking-tight">My Subscriptions</h1>
                <p className="text-sm text-muted-foreground">View your active and past subscriptions</p>
            </div>
            {(!subscriptions || subscriptions.length === 0) ? (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                    <WifiOff className="size-12 text-muted-foreground mb-4" />
                    <h3 className="text-lg font-medium">No Subscriptions Found</h3>
                    <p className="text-sm text-muted-foreground mt-1">
                        Contact your provider to set up a subscription.
                    </p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {subscriptions.map((sub) => (
                        <Card key={sub.id}>
                            <CardHeader className="flex flex-row items-center justify-between pb-2">
                                <CardTitle className="text-base font-medium flex items-center gap-2">
                                    <Wifi className="size-4" />
                                    {sub.username}
                                </CardTitle>
                                {getStatusBadge(sub.status)}
                            </CardHeader>
                            <CardContent className="space-y-2 text-sm">
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">IP Address</span>
                                    <span>{sub.static_ip ?? 'Dynamic'}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Profile</span>
                                    <span>{sub.mikrotik?.profile ?? '-'}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Expiry</span>
                                    <span>{formatDate(sub.expiry_date)}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Activated</span>
                                    <span>{formatDate(sub.activated_at)}</span>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    )
}
