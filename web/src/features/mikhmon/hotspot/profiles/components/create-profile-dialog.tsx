import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
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
import { useCreateHotspotProfile } from '@/hooks/use-hotspot'

const createProfileSchema = z.object({
    name: z.string().min(1, 'Profile name is required'),
    address_pool: z.string().optional(),
    shared_users: z.string().optional(),
    rate_limit: z.string().optional(),
    parent_queue: z.string().optional(),
    session_timeout: z.string().optional(),
})

type CreateProfileForm = z.infer<typeof createProfileSchema>

interface CreateProfileDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
}

export function CreateProfileDialog({
    open,
    onOpenChange,
    routerId,
}: CreateProfileDialogProps) {
    const { mutateAsync: createProfile, isPending } = useCreateHotspotProfile()

    const form = useForm<CreateProfileForm>({
        resolver: zodResolver(createProfileSchema) as never,
        defaultValues: {
            name: '',
            address_pool: '',
            shared_users: '',
            rate_limit: '',
            parent_queue: '',
            session_timeout: '',
        },
    })

    const onSubmit = async (data: CreateProfileForm) => {
        try {
            const payload: Record<string, unknown> = { name: data.name }
            if (data.address_pool) payload.address_pool = data.address_pool
            if (data.shared_users) payload.shared_users = data.shared_users
            if (data.rate_limit) payload.rate_limit = data.rate_limit
            if (data.parent_queue) payload.parent_queue = data.parent_queue
            if (data.session_timeout) payload.session_timeout = data.session_timeout

            await createProfile({ routerId, data: payload as never })
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
            <DialogContent className='sm:max-w-[500px]'>
                <DialogHeader>
                    <DialogTitle>Add Hotspot Profile</DialogTitle>
                    <DialogDescription>
                        Create a new hotspot user profile on the MikroTik router.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
                        <div className='grid grid-cols-2 gap-4'>
                            <FormField
                                control={form.control}
                                name='name'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>Profile Name *</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. 1day, 7days'
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
                                name='address_pool'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Address Pool</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. hs-pool'
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
                                name='shared_users'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Shared Users</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. 1'
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
                                name='rate_limit'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Rate Limit</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. 5M/5M'
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
                                name='session_timeout'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Session Timeout</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. 1d, 3h'
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
                                name='parent_queue'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>Parent Queue</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='Optional queue name'
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
                                {isPending && <Loader2 className='mr-2 size-4 animate-spin' />}
                                Add Profile
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
