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
import { useDeleteHotspotProfile } from '@/hooks/use-hotspot'
import type { HotspotProfile } from '@/lib/schemas/mikrotik'

interface DeleteProfileDialogProps {
    profile: HotspotProfile | null
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
}

export function DeleteProfileDialog({
    profile,
    open,
    onOpenChange,
    routerId,
}: DeleteProfileDialogProps) {
    const { mutateAsync: deleteProfile, isPending } = useDeleteHotspotProfile()

    const handleDelete = async () => {
        if (!profile?.['.id']) return
        try {
            await deleteProfile({ routerId, id: profile['.id']! })
            onOpenChange(false)
        } catch {
            // Error handled by hook
        }
    }

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Delete Hotspot Profile</AlertDialogTitle>
                    <AlertDialogDescription>
                        Are you sure you want to delete profile{' '}
                        <span className='font-semibold'>{profile?.name}</span>? This will
                        remove the profile from MikroTik. Users assigned to this profile
                        may be affected.
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
