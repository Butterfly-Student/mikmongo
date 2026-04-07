import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Loader2 } from 'lucide-react'
import { useDeleteHotspotUser } from '@/hooks/use-hotspot'
import type { HotspotActive } from '@/lib/schemas/mikrotik'

interface DisconnectDialogProps {
    session: HotspotActive | null
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
}

export function DisconnectDialog({
    session,
    open,
    onOpenChange,
    routerId,
}: DisconnectDialogProps) {
    const { mutateAsync: removeUser, isPending } = useDeleteHotspotUser()

    const handleDisconnect = async () => {
        const id = (session as Record<string, unknown>)?.['.id'] as string | undefined
        if (!id) return
        try {
            await removeUser({ routerId, id })
            onOpenChange(false)
        } catch {
            // Error handled by hook
        }
    }

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Disconnect Session</AlertDialogTitle>
                    <AlertDialogDescription>
                        Are you sure you want to disconnect user{' '}
                        <span className='font-semibold'>{session?.user}</span> from{' '}
                        <span className='font-semibold'>{session?.address}</span>? The
                        session will be terminated immediately.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                        onClick={handleDisconnect}
                        disabled={isPending}
                        className='bg-destructive text-destructive-foreground hover:bg-destructive/90'
                    >
                        {isPending && <Loader2 className='mr-2 size-4 animate-spin' />}
                        Disconnect
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )
}
