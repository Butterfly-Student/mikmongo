import { useEffect } from 'react'
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
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Pencil } from 'lucide-react'
import { createCustomerSchema, type CreateCustomer } from '../data/schema'
import { useUpdateCustomer } from '@/hooks/use-customers'
import type { CustomerResponse } from '@/lib/schemas/customer'

interface EditCustomerDialogProps {
    customer: CustomerResponse | null
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function EditCustomerDialog({
    customer,
    open,
    onOpenChange,
}: EditCustomerDialogProps) {
    const { mutateAsync: updateCustomer, isPending } = useUpdateCustomer()

    const form = useForm<CreateCustomer>({
        resolver: zodResolver(createCustomerSchema) as never,
        defaultValues: {
            full_name: customer?.full_name ?? '',
            phone: customer?.phone ?? '',
            email: customer?.email ?? '',
            address: customer?.address ?? '',
            username: customer?.username ?? '',
            password: '',
            static_ip: '',
        },
    })

    // Reset form when customer changes
    useEffect(() => {
        if (customer) {
            form.reset({
                full_name: customer.full_name ?? '',
                phone: customer.phone ?? '',
                email: customer.email ?? '',
                address: customer.address ?? '',
                username: customer.username ?? '',
                password: '',
                static_ip: '',
            })
        }
    }, [customer, form])

    const onSubmit = async (data: CreateCustomer) => {
        if (!customer) return
        try {
            const payload: Record<string, unknown> = {
                full_name: data.full_name,
                phone: data.phone,
            }
            if (data.email) payload.email = data.email
            if (data.address) payload.address = data.address
            if (data.username) payload.username = data.username
            if (data.password) payload.password = data.password
            if (data.static_ip) payload.static_ip = data.static_ip

            await updateCustomer({ id: customer.id, data: payload })
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
            <DialogContent className='sm:max-w-[500px] max-h-[90vh] overflow-y-auto'>
                <DialogHeader>
                    <DialogTitle className='flex items-center gap-2'>
                        <Pencil className='size-4' />
                        Edit Customer
                    </DialogTitle>
                    <DialogDescription>
                        Update customer details. Fields marked with * are required.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
                        <div className='grid grid-cols-2 gap-4'>
                            <FormField
                                control={form.control}
                                name='full_name'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>Full Name *</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='John Doe'
                                                {...field}
                                                disabled={isPending}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='phone'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Phone *</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='+62 812 3456 7890'
                                                {...field}
                                                disabled={isPending}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='email'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Email</FormLabel>
                                        <FormControl>
                                            <Input
                                                type='email'
                                                placeholder='john@example.com'
                                                {...field}
                                                disabled={isPending}
                                                value={field.value ?? ''}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='address'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>Address</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='123 Main St, Jakarta'
                                                {...field}
                                                disabled={isPending}
                                                value={field.value ?? ''}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='username'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Username (PPP/Hotspot)</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='john.doe'
                                                {...field}
                                                disabled={isPending}
                                                value={field.value ?? ''}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='password'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Password</FormLabel>
                                        <FormControl>
                                            <Input
                                                type='password'
                                                placeholder='Leave blank to keep current'
                                                {...field}
                                                disabled={isPending}
                                                value={field.value ?? ''}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='static_ip'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>Static IP</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='192.168.1.100'
                                                {...field}
                                                disabled={isPending}
                                                value={field.value ?? ''}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        <DialogFooter>
                            <Button
                                type='button'
                                variant='outline'
                                onClick={() => handleOpenChange(false)}
                                disabled={isPending}
                            >
                                Cancel
                            </Button>
                            <Button type='submit' disabled={isPending}>
                                {isPending ? 'Saving...' : 'Save Changes'}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
