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
import type { HotspotUser } from '@/lib/schemas/mikrotik'

interface DeleteUserDialogProps {
    user: HotspotUser | null
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
}

export function DeleteUserDialog({
    user,
    open,
    onOpenChange,
    routerId,
}: DeleteUserDialogProps) {
    const { mutateAsync: deleteUser, isPending } = useDeleteHotspotUser()

    const handleDelete = async () => {
        if (!user?.['.id']) return
        try {
            await deleteUser({ routerId, id: user['.id']! })
            onOpenChange(false)
        } catch {
            // Error handled by hook
        }
    }

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Delete Hotspot User</AlertDialogTitle>
                    <AlertDialogDescription>
                        Are you sure you want to delete the user{' '}
                        <span className='font-semibold'>{user?.name}</span>? This action
                        cannot be undone and will remove the user from the MikroTik router.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                        onClick={handleDelete}
                        disabled={isPending}
                        className='bg-destructive text-destructive-foreground hover:bg-destructive/90'
                    >
                        {isPending && <Loader2 className='mr-2 size-4 animate-spin' />}
                        Delete
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )
}
