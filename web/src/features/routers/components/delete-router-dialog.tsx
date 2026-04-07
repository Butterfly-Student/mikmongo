import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
    DialogDescription,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Loader2, Trash2 } from 'lucide-react'
import { useDeleteRouter } from '@/hooks/use-routers'
import type { RouterResponse } from '@/lib/schemas/router'

interface DeleteRouterDialogProps {
    router: RouterResponse | null
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function DeleteRouterDialog({
    router,
    open,
    onOpenChange,
}: DeleteRouterDialogProps) {
    const { mutateAsync: deleteRouter, isPending } = useDeleteRouter()

    const handleDelete = async () => {
        if (!router) return
        try {
            await deleteRouter(router.id)
            onOpenChange(false)
        } catch {
            // Error handled by hook
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className='sm:max-w-[425px]'>
                <DialogHeader>
                    <DialogTitle className='flex items-center gap-2'>
                        <Trash2 className='size-5 text-destructive' />
                        Delete Router
                    </DialogTitle>
                    <DialogDescription>
                        Are you sure you want to delete router{' '}
                        <span className='font-semibold'>{router?.name}</span>?
                        This action cannot be undone. All associated bandwidth profiles
                        and subscriptions will be affected.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <Button
                        type='button'
                        variant='outline'
                        onClick={() => onOpenChange(false)}
                        disabled={isPending}
                    >
                        Cancel
                    </Button>
                    <Button
                        variant='destructive'
                        onClick={handleDelete}
                        disabled={isPending}
                    >
                        {isPending && (
                            <Loader2 className='mr-2 size-4 animate-spin' />
                        )}
                        Delete Router
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
