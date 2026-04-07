import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
    DialogDescription,
} from '@/components/ui/dialog'
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { rejectRegistrationSchema, type RejectRegistration } from '../data/schema'
import { useRejectRegistration } from '@/hooks/use-customers'
import { XCircle, Loader2 } from 'lucide-react'
import type { RegistrationResponse } from '@/lib/schemas/customer'

interface RejectRegistrationDialogProps {
    registration: RegistrationResponse | null
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function RejectRegistrationDialog({
    registration,
    open,
    onOpenChange,
}: RejectRegistrationDialogProps) {
    const { mutateAsync: rejectReg, isPending } = useRejectRegistration()

    const form = useForm<RejectRegistration>({
        resolver: zodResolver(rejectRegistrationSchema) as never,
        defaultValues: {
            reason: '',
        },
    })

    const onSubmit = async (data: RejectRegistration) => {
        if (!registration) return
        try {
            await rejectReg({
                id: registration.id,
                data: { reason: data.reason },
            })
            form.reset()
            onOpenChange(false)
        } catch {
            // Error handled by hook
        }
    }

    const handleOpenChange = (newOpen: boolean) => {
        if (!newOpen) form.reset()
        onOpenChange(newOpen)
    }

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogContent className='sm:max-w-[450px]'>
                <DialogHeader>
                    <DialogTitle className='flex items-center gap-2'>
                        <XCircle className='size-5 text-destructive' />
                        Reject Registration
                    </DialogTitle>
                    <DialogDescription>
                        {registration
                            ? `Reject registration for ${registration.full_name}`
                            : 'Reject this registration'}
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
                        <FormField
                            control={form.control}
                            name='reason'
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Rejection Reason *</FormLabel>
                                    <FormControl>
                                        <Textarea
                                            placeholder='Provide a reason for rejecting this registration...'
                                            {...field}
                                            disabled={isPending}
                                            rows={3}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <DialogFooter>
                            <Button
                                type='button'
                                variant='outline'
                                onClick={() => handleOpenChange(false)}
                                disabled={isPending}
                            >
                                Cancel
                            </Button>
                            <Button
                                type='submit'
                                variant='destructive'
                                disabled={isPending || !form.watch('reason')}
                            >
                                {isPending && (
                                    <Loader2 className='mr-2 size-4 animate-spin' />
                                )}
                                Reject
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
