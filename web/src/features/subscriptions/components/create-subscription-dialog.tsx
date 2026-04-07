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
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Loader2, CreditCard } from 'lucide-react'
import { createSubscriptionSchema, type CreateSubscription } from '../data/schema'
import { useCreateSubscription } from '@/hooks/use-subscriptions'
import { useCustomers } from '@/hooks/use-customers'
import { useProfiles } from '@/hooks/use-profiles'
import { useRouterStore } from '@/stores/router-store'

interface CreateSubscriptionDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function CreateSubscriptionDialog({
    open,
    onOpenChange,
}: CreateSubscriptionDialogProps) {
    const selectedRouterId = useRouterStore((s) => s.selectedRouterId)
    const { mutateAsync: createSub, isPending } = useCreateSubscription(
        selectedRouterId ?? ''
    )
    const { data: customersData, isLoading: customersLoading } = useCustomers()
    const { data: profilesData, isLoading: profilesLoading } = useProfiles(
        selectedRouterId ?? null
    )

    const customers = customersData?.customers ?? []
    const profiles = profilesData?.profiles ?? []

    const form = useForm<CreateSubscription>({
        resolver: zodResolver(createSubscriptionSchema) as never,
        defaultValues: {
            customer_id: '',
            plan_id: '',
            username: '',
            password: '',
            static_ip: '',
            gateway: '',
            billing_day: undefined,
            auto_isolate: true,
            grace_period_days: undefined,
            notes: '',
        },
    })

    const onSubmit = async (data: CreateSubscription) => {
        if (!selectedRouterId) return
        try {
            const payload: Record<string, unknown> = {
                customer_id: data.customer_id,
                plan_id: data.plan_id,
                username: data.username,
            }
            if (data.password) payload.password = data.password
            if (data.static_ip) payload.static_ip = data.static_ip
            if (data.gateway) payload.gateway = data.gateway
            if (data.billing_day != null) payload.billing_day = data.billing_day
            if (data.auto_isolate != null) payload.auto_isolate = data.auto_isolate
            if (data.grace_period_days != null)
                payload.grace_period_days = data.grace_period_days
            if (data.notes) payload.notes = data.notes

            await createSub(payload)
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

    const isFormValid =
        form.watch('customer_id') !== '' &&
        form.watch('plan_id') !== '' &&
        form.watch('username') !== ''

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogContent className='sm:max-w-[500px] max-h-[90vh] overflow-y-auto'>
                <DialogHeader>
                    <DialogTitle className='flex items-center gap-2'>
                        <CreditCard className='size-5 text-primary' />
                        Create Subscription
                    </DialogTitle>
                    <DialogDescription>
                        Assign a bandwidth profile to a customer. Both customer and
                        profile are required.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
                        <FormField
                            control={form.control}
                            name='customer_id'
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Customer *</FormLabel>
                                    {customersLoading ? (
                                        <Skeleton className='h-9 w-full' />
                                    ) : (
                                        <Select
                                            value={field.value}
                                            onValueChange={field.onChange}
                                            disabled={isPending}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder='Select a customer' />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {customers.map((customer) => (
                                                    <SelectItem
                                                        key={customer.id}
                                                        value={customer.id}
                                                    >
                                                        {customer.full_name}{' '}
                                                        <span className='text-muted-foreground'>
                                                            ({customer.phone})
                                                        </span>
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                    )}
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name='plan_id'
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Bandwidth Profile *</FormLabel>
                                    {!selectedRouterId || (profilesLoading && selectedRouterId) ? (
                                        <Skeleton className='h-9 w-full' />
                                    ) : (
                                        <Select
                                            value={field.value}
                                            onValueChange={field.onChange}
                                            disabled={isPending || !selectedRouterId}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder='Select a profile' />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {profiles.map((profile) => (
                                                    <SelectItem
                                                        key={profile.id}
                                                        value={profile.id}
                                                    >
                                                        {profile.name}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                    )}
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name='username'
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Username *</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder='ppp-username'
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

                        <div className='grid grid-cols-2 gap-4'>
                            <FormField
                                control={form.control}
                                name='static_ip'
                                render={({ field }) => (
                                    <FormItem>
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

                            <FormField
                                control={form.control}
                                name='gateway'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Gateway</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='192.168.1.1'
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

                        <FormField
                            control={form.control}
                            name='notes'
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Notes</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder='Optional notes'
                                            {...field}
                                            disabled={isPending}
                                            value={field.value ?? ''}
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
                                disabled={isPending || !isFormValid || !selectedRouterId}
                            >
                                {isPending && (
                                    <Loader2 className='mr-2 size-4 animate-spin' />
                                )}
                                Create Subscription
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
