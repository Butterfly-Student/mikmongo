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
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { approveRegistrationSchema, type ApproveRegistration } from '../data/schema'
import { useApproveRegistration } from '@/hooks/use-customers'
import { useRouters } from '@/hooks/use-routers'
import { useProfiles } from '@/hooks/use-profiles'
import { CheckCircle2, Loader2 } from 'lucide-react'
import type { RegistrationResponse } from '@/lib/schemas/customer'

interface ApproveRegistrationDialogProps {
    registration: RegistrationResponse | null
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function ApproveRegistrationDialog({
    registration,
    open,
    onOpenChange,
}: ApproveRegistrationDialogProps) {
    const { mutateAsync: approveReg, isPending } = useApproveRegistration()
    const { data: routersData, isLoading: routersLoading } = useRouters()
    const selectedRouterId = useForm<ApproveRegistration>({
        resolver: zodResolver(approveRegistrationSchema) as never,
        defaultValues: {
            router_id: '',
            profile_id: '',
        },
    })

    const watchedRouterId = selectedRouterId.watch('router_id')
    const { data: profilesData, isLoading: profilesLoading } = useProfiles(
        watchedRouterId || null
    )

    const routers = routersData?.routers ?? []
    const profiles = profilesData?.profiles ?? []

    const form = selectedRouterId

    const onSubmit = async (data: ApproveRegistration) => {
        if (!registration) return
        try {
            await approveReg({
                id: registration.id,
                data: {
                    router_id: data.router_id,
                    ...(data.profile_id ? { profile_id: data.profile_id } : {}),
                },
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
                        <CheckCircle2 className='size-5 text-green-500' />
                        Approve Registration
                    </DialogTitle>
                    <DialogDescription>
                        {registration
                            ? `Approve registration for ${registration.full_name}`
                            : 'Approve this registration'}
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
                        <FormField
                            control={form.control}
                            name='router_id'
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Router *</FormLabel>
                                    {routersLoading ? (
                                        <Skeleton className='h-9 w-full' />
                                    ) : (
                                        <Select
                                            value={field.value}
                                            onValueChange={field.onChange}
                                            disabled={isPending}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder='Select a router' />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {routers.map((router) => (
                                                    <SelectItem key={router.id} value={router.id}>
                                                        {router.name}
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
                            name='profile_id'
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Bandwidth Profile (optional)</FormLabel>
                                    {profilesLoading && watchedRouterId ? (
                                        <Skeleton className='h-9 w-full' />
                                    ) : (
                                        <Select
                                            value={field.value}
                                            onValueChange={field.onChange}
                                            disabled={isPending}
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
                                disabled={isPending || !form.watch('router_id')}
                            >
                                {isPending && (
                                    <Loader2 className='mr-2 size-4 animate-spin' />
                                )}
                                Approve
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
