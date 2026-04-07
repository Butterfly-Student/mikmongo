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
import { useDeleteCustomer } from '@/hooks/use-customers'
import type { CustomerResponse } from '@/lib/schemas/customer'

interface DeleteCustomerDialogProps {
    customer: CustomerResponse | null
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function DeleteCustomerDialog({
    customer,
    open,
    onOpenChange,
}: DeleteCustomerDialogProps) {
    const { mutateAsync: deleteCustomer, isPending } = useDeleteCustomer()

    const handleDelete = async () => {
        if (!customer) return
        try {
            await deleteCustomer(customer.id)
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
                        Delete Customer
                    </DialogTitle>
                    <DialogDescription>
                        Are you sure you want to delete{' '}
                        <span className='font-semibold'>{customer?.full_name}</span>?
                        This action cannot be undone. All associated data including
                        subscriptions and invoices will be permanently removed.
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
                        Delete Customer
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
