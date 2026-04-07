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
import { Loader2 } from 'lucide-react'
import { createHotspotUserSchema, type CreateHotspotUserForm } from '../data/schema'
import { useCreateHotspotUser } from '@/hooks/use-hotspot'
import type { HotspotProfile, HotspotServer } from '@/lib/schemas/mikrotik'

interface CreateUserDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
    profiles: HotspotProfile[]
    servers: HotspotServer[]
}

export function CreateUserDialog({
    open,
    onOpenChange,
    routerId,
    profiles,
    servers,
}: CreateUserDialogProps) {
    const { mutateAsync: createUser, isPending } = useCreateHotspotUser()

    const form = useForm<CreateHotspotUserForm>({
        resolver: zodResolver(createHotspotUserSchema) as never,
        defaultValues: {
            name: '',
            password: '',
            profile: '',
            server: '',
            mac_address: '',
            limit_uptime: '',
            limit_bytes_total: '',
            comment: '',
            disabled: false,
        },
    })

    const onSubmit = async (data: CreateHotspotUserForm) => {
        try {
            const payload: Record<string, unknown> = { name: data.name }
            if (data.password) payload.password = data.password
            if (data.profile) payload.profile = data.profile
            if (data.server) payload.server = data.server
            if (data.mac_address) payload.mac_address = data.mac_address
            if (data.limit_uptime) payload.limit_uptime = data.limit_uptime
            if (data.limit_bytes_total) payload.limit_bytes_total = data.limit_bytes_total
            if (data.comment) payload.comment = data.comment
            if (data.disabled) payload.disabled = data.disabled

            await createUser({ routerId, data: payload as never })
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
            <DialogContent className='sm:max-w-[550px] max-h-[90vh] overflow-y-auto'>
                <DialogHeader>
                    <DialogTitle>Add Hotspot User</DialogTitle>
                    <DialogDescription>
                        Create a new hotspot user on the MikroTik router.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
                        <div className='grid grid-cols-2 gap-4'>
                            <FormField
                                control={form.control}
                                name='name'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Username *</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='user001'
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
                                                placeholder='Password'
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
                                name='profile'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Profile</FormLabel>
                                        <Select
                                            onValueChange={field.onChange}
                                            value={field.value ?? ''}
                                            disabled={isPending}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder='Select profile' />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {profiles.map((p) => (
                                                    <SelectItem key={p['.id'] ?? p.name} value={p.name}>
                                                        {p.name}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='server'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Server</FormLabel>
                                        <Select
                                            onValueChange={field.onChange}
                                            value={field.value ?? ''}
                                            disabled={isPending}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder='Select server' />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                <SelectItem value='all'>All</SelectItem>
                                                {servers.map((s) => (
                                                    <SelectItem key={s['.id'] ?? s.name} value={s.name}>
                                                        {s.name}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='mac_address'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>MAC Address (optional)</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='AA:BB:CC:DD:EE:FF'
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
                                name='limit_uptime'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Time Limit</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. 3h, 1d'
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
                                name='limit_bytes_total'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Data Limit</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. 500M, 1G'
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
                                name='comment'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>Comment</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='Optional comment'
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
                                Add User
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
