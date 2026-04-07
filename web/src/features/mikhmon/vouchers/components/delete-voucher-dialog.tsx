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
import { useRemoveMikhmonVoucherBatch } from '@/hooks/use-mikhmon'
import type { VoucherResponse } from '@/lib/schemas/mikhmon'

interface DeleteVoucherDialogProps {
    voucher: VoucherResponse | null
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
}

export function DeleteVoucherDialog({
    voucher,
    open,
    onOpenChange,
    routerId,
}: DeleteVoucherDialogProps) {
    const { mutateAsync: removeBatch, isPending } = useRemoveMikhmonVoucherBatch()

    const handleDelete = async () => {
        if (!voucher?.comment) return
        try {
            await removeBatch({ routerId, comment: voucher.comment })
            onOpenChange(false)
        } catch {
            // Error handled by hook
        }
    }

    return (
        <AlertDialog open={open} onOpenChange={onOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Delete Voucher Batch</AlertDialogTitle>
                    <AlertDialogDescription>
                        Are you sure you want to delete all vouchers in batch{' '}
                        <span className='font-semibold font-mono'>{voucher?.comment}</span>?
                        This will remove all vouchers in this batch from the MikroTik router.
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
