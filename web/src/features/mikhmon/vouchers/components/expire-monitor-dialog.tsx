import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '@/components/ui/dialog'
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Loader2, Clock, CheckCircle2, XCircle } from 'lucide-react'
import { useExpirationStatus, useSetupExpiration, useDisableExpiration } from '@/hooks/use-mikhmon'

interface ExpireMonitorDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
}

export function ExpireMonitorDialog({
    open,
    onOpenChange,
    routerId,
}: ExpireMonitorDialogProps) {
    const { data: status, isLoading: statusLoading } = useExpirationStatus(routerId || null)
    const { mutateAsync: setupExpiration, isPending: setupPending } = useSetupExpiration()
    const { mutateAsync: disableExpiration, isPending: disablePending } = useDisableExpiration()

    const isEnabled = status?.enabled ?? false
    const isPending = setupPending || disablePending

    const handleSetup = async () => {
        try {
            await setupExpiration({ routerId, data: {} })
        } catch {
            // Error handled by hook
        }
    }

    const handleDisable = async () => {
        try {
            await disableExpiration({ routerId })
        } catch {
            // Error handled by hook
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className='sm:max-w-[460px]'>
                <DialogHeader>
                    <DialogTitle className='flex items-center gap-2'>
                        <Clock className='size-5' />
                        Expire Monitor
                    </DialogTitle>
                    <DialogDescription>
                        The expire monitor runs every minute to remove or notify expired
                        hotspot users based on their profile settings.
                    </DialogDescription>
                </DialogHeader>

                <div className='space-y-4 py-2'>
                    {statusLoading ? (
                        <div className='flex items-center justify-center py-6'>
                            <Loader2 className='size-6 animate-spin text-muted-foreground' />
                        </div>
                    ) : (
                        <>
                            <div className='flex items-center justify-between rounded-lg border p-4'>
                                <div className='space-y-1'>
                                    <p className='text-sm font-medium'>Status</p>
                                    <p className='text-xs text-muted-foreground'>
                                        MikroTik scheduler task
                                    </p>
                                </div>
                                {isEnabled ? (
                                    <Badge variant='default' className='gap-1.5'>
                                        <CheckCircle2 className='size-3.5' />
                                        Active
                                    </Badge>
                                ) : (
                                    <Badge variant='secondary' className='gap-1.5'>
                                        <XCircle className='size-3.5' />
                                        Inactive
                                    </Badge>
                                )}
                            </div>

                            {status && (
                                <div className='rounded-lg border p-4 space-y-2 text-sm'>
                                    {status.last_run && (
                                        <div className='flex justify-between'>
                                            <span className='text-muted-foreground'>Last run</span>
                                            <span>{status.last_run}</span>
                                        </div>
                                    )}
                                    {status.next_run && (
                                        <div className='flex justify-between'>
                                            <span className='text-muted-foreground'>Next run</span>
                                            <span>{status.next_run}</span>
                                        </div>
                                    )}
                                    {status.user_count != null && (
                                        <div className='flex justify-between'>
                                            <span className='text-muted-foreground'>Users monitored</span>
                                            <span>{status.user_count}</span>
                                        </div>
                                    )}
                                </div>
                            )}
                        </>
                    )}
                </div>

                <DialogFooter className='gap-2'>
                    <Button
                        variant='outline'
                        onClick={() => onOpenChange(false)}
                        disabled={isPending}
                    >
                        Close
                    </Button>
                    {isEnabled ? (
                        <AlertDialog>
                            <AlertDialogTrigger asChild>
                                <Button variant='destructive' disabled={isPending}>
                                    {disablePending && (
                                        <Loader2 className='mr-2 size-4 animate-spin' />
                                    )}
                                    Disable Monitor
                                </Button>
                            </AlertDialogTrigger>
                            <AlertDialogContent>
                                <AlertDialogHeader>
                                    <AlertDialogTitle>Disable Expire Monitor?</AlertDialogTitle>
                                    <AlertDialogDescription>
                                        This will remove the scheduler task from MikroTik.
                                        Expired users will no longer be automatically removed
                                        or notified. You can re-enable it at any time.
                                    </AlertDialogDescription>
                                </AlertDialogHeader>
                                <AlertDialogFooter>
                                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                                    <AlertDialogAction
                                        onClick={handleDisable}
                                        className='bg-destructive text-destructive-foreground hover:bg-destructive/90'
                                    >
                                        Disable
                                    </AlertDialogAction>
                                </AlertDialogFooter>
                            </AlertDialogContent>
                        </AlertDialog>
                    ) : (
                        <Button onClick={handleSetup} disabled={isPending}>
                            {setupPending && (
                                <Loader2 className='mr-2 size-4 animate-spin' />
                            )}
                            Enable Monitor
                        </Button>
                    )}
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
