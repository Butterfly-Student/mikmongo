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
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Loader2 } from 'lucide-react'
import { generateVoucherSchema, type GenerateVoucherForm } from '../data/schema'
import { useGenerateMikhmonVouchers } from '@/hooks/use-mikhmon'
import type { HotspotServer } from '@/lib/schemas/mikrotik'

interface ProfileOption {
    id?: string
    name: string
}

interface GenerateVoucherDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    routerId: string
    profiles: ProfileOption[]
    servers: HotspotServer[]
    onGenerated?: (comment: string) => void
}

const CHAR_SETS = [
    { value: 'lower', label: 'Lowercase letters (abc...)' },
    { value: 'upper', label: 'Uppercase letters (ABC...)' },
    { value: 'lower1', label: 'Lowercase + numbers (abc123...)' },
    { value: 'upper1', label: 'Uppercase + numbers (ABC123...)' },
    { value: 'num', label: 'Numbers only (123...)' },
    { value: 'mix', label: 'Mixed (aAbB...)' },
    { value: 'mix1', label: 'Mixed + numbers (aAbB123...)' },
]

export function GenerateVoucherDialog({
    open,
    onOpenChange,
    routerId,
    profiles,
    servers,
    onGenerated,
}: GenerateVoucherDialogProps) {
    const { mutateAsync: generateVouchers, isPending } = useGenerateMikhmonVouchers()

    const form = useForm<GenerateVoucherForm>({
        resolver: zodResolver(generateVoucherSchema) as never,
        defaultValues: {
            quantity: 10,
            profile: '',
            mode: 'vc',
            char_set: 'lower1',
            name_length: 6,
            prefix: '',
            server: 'all',
            time_limit: '',
            data_limit: '',
            comment: '',
        },
    })

    const onSubmit = async (data: GenerateVoucherForm) => {
        try {
            const payload: Record<string, unknown> = {
                quantity: data.quantity,
                profile: data.profile,
                mode: data.mode,
                char_set: data.char_set,
            }
            if (data.name_length) payload.name_length = data.name_length
            if (data.prefix) payload.prefix = data.prefix
            if (data.server) payload.server = data.server
            if (data.time_limit) payload.time_limit = data.time_limit
            if (data.data_limit) payload.data_limit = data.data_limit
            if (data.comment) payload.comment = data.comment

            await generateVouchers({ routerId, data: payload as never })
            form.reset()
            onOpenChange(false)
            if (data.comment) onGenerated?.(data.comment)
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
            <DialogContent className='sm:max-w-[580px] max-h-[90vh] overflow-y-auto'>
                <DialogHeader>
                    <DialogTitle>Generate Vouchers</DialogTitle>
                    <DialogDescription>
                        Create a batch of hotspot vouchers on the MikroTik router.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
                        <div className='grid grid-cols-2 gap-4'>
                            <FormField
                                control={form.control}
                                name='quantity'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Quantity *</FormLabel>
                                        <FormControl>
                                            <Input
                                                type='number'
                                                min={1}
                                                max={500}
                                                placeholder='10'
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
                                name='profile'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Profile *</FormLabel>
                                        <Select
                                            onValueChange={field.onChange}
                                            value={field.value}
                                            disabled={isPending}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder='Select profile' />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {profiles.map((p) => (
                                                    <SelectItem key={p.id ?? p.name} value={p.name}>
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
                                name='mode'
                                render={({ field }) => (
                                    <FormItem className='col-span-2'>
                                        <FormLabel>Mode</FormLabel>
                                        <FormControl>
                                            <RadioGroup
                                                value={field.value}
                                                onValueChange={field.onChange}
                                                className='flex gap-6'
                                                disabled={isPending}
                                            >
                                                <div className='flex items-center gap-2'>
                                                    <RadioGroupItem value='vc' id='mode-vc' />
                                                    <Label htmlFor='mode-vc'>
                                                        Voucher (username = password)
                                                    </Label>
                                                </div>
                                                <div className='flex items-center gap-2'>
                                                    <RadioGroupItem value='up' id='mode-up' />
                                                    <Label htmlFor='mode-up'>
                                                        User/Pass (separate)
                                                    </Label>
                                                </div>
                                            </RadioGroup>
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name='char_set'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Character Set</FormLabel>
                                        <Select
                                            onValueChange={field.onChange}
                                            value={field.value}
                                            disabled={isPending}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {CHAR_SETS.map((cs) => (
                                                    <SelectItem key={cs.value} value={cs.value}>
                                                        {cs.label}
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
                                name='name_length'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Name Length</FormLabel>
                                        <FormControl>
                                            <Input
                                                type='number'
                                                min={3}
                                                max={12}
                                                placeholder='6'
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
                                name='prefix'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Prefix</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='e.g. hs-'
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
                                name='server'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Server</FormLabel>
                                        <Select
                                            onValueChange={field.onChange}
                                            value={field.value ?? 'all'}
                                            disabled={isPending}
                                        >
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue />
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
                                name='time_limit'
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
                                name='data_limit'
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
                                        <FormLabel>Comment / Tag</FormLabel>
                                        <FormControl>
                                            <Input
                                                placeholder='Optional batch tag (e.g. promo-apr)'
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
                                Generate
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
