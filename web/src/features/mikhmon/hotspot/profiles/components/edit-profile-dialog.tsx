import { useEffect } from 'react'
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
import { useUpdateHotspotProfile } from '@/hooks/use-hotspot'
import type { HotspotProfile } from '@/lib/schemas/mikrotik'

const editProfileSchema = z.object({
    name: z.string().min(1, 'Profile name is required'),
    address_pool: z.string().optional(),
    shared_users: z.string().optional(),
    rate_limit: z.string().optional(),
    parent_queue: z.string().optional(),
    session_timeout: z.string().optional(),
})

type EditProfileForm = z.infer<typeof editProfileSchema>

interface EditProfileDialogProps {
    profile: HotspotProfile | null
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
}

export function EditProfileDialog({
    profile,
    open,
    onOpenChange,
    routerId,
}: EditProfileDialogProps) {
    const { mutateAsync: updateProfile, isPending } = useUpdateHotspotProfile()

    const form = useForm<EditProfileForm>({
        resolver: zodResolver(editProfileSchema) as never,
        defaultValues: {
            name: '',
            address_pool: '',
            shared_users: '',
            rate_limit: '',
            parent_queue: '',
            session_timeout: '',
        },
    })

    useEffect(() => {
        if (profile) {
            form.reset({
                name: profile.name,
                address_pool: profile.addressPool ?? '',
                shared_users: profile.sharedUsers != null ? String(profile.sharedUsers) : '',
                rate_limit: profile.rateLimit ?? '',
                parent_queue: profile.parentQueue ?? '',
                session_timeout: profile.sessionTimeout ?? '',
            })
        }
    }, [profile, form])

    const onSubmit = async (data: EditProfileForm) => {
        if (!profile?.['.id']) return
        try {
            const payload: Record<string, unknown> = {}
            if (data.address_pool !== undefined) payload.address_pool = data.address_pool
            if (data.shared_users) payload.shared_users = data.shared_users
            if (data.rate_limit !== undefined) payload.rate_limit = data.rate_limit
            if (data.parent_queue !== undefined) payload.parent_queue = data.parent_queue
            if (data.session_timeout !== undefined) payload.session_timeout = data.session_timeout

            await updateProfile({ routerId, id: profile['.id']!, data: payload as never })
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
                    <DialogTitle>Edit Hotspot Profile</DialogTitle>
                    <DialogDescription>
                        Update hotspot profile settings on the MikroTik router.
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
                                        <FormLabel>Profile Name</FormLabel>
                                        <FormControl>
                                            <Input {...field} disabled className='bg-muted' />
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
                                Save Changes
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
