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
import { Loader2 } from 'lucide-react'
import { createCustomerSchema, type CreateCustomer } from '../data/schema'
import { useCreateCustomer } from '@/hooks/use-customers'

interface CreateCustomerDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function CreateCustomerDialog({
    open,
    onOpenChange,
}: CreateCustomerDialogProps) {
    const { mutateAsync: createCustomer, isPending } = useCreateCustomer()

    const form = useForm<CreateCustomer>({
        resolver: zodResolver(createCustomerSchema) as never,
        defaultValues: {
            full_name: '',
            phone: '',
            email: '',
            address: '',
            username: '',
            password: '',
            static_ip: '',
        },
    })

    const onSubmit = async (data: CreateCustomer) => {
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

            await createCustomer(payload)
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
                    <DialogTitle>Add New Customer</DialogTitle>
                    <DialogDescription>
                        Create a new customer in the system. Fields marked with * are
                        required.
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
                                                placeholder='Secure password'
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
                                {isPending && (
                                    <Loader2 className='mr-2 size-4 animate-spin' />
                                )}
                                Create Customer
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
