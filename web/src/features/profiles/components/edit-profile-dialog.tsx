import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
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
import { createProfileSchema, type CreateProfile, type Profile } from '../data/schema'
import { useUpdateProfile } from '@/hooks/use-profiles'

interface EditProfileDialogProps {
  profile: Profile | null
  routerId: string | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function EditProfileDialog({ profile, routerId, open, onOpenChange }: EditProfileDialogProps) {
  const { mutateAsync: updateProfile, isPending } = useUpdateProfile()

  const form = useForm<CreateProfile>({
    resolver: zodResolver(createProfileSchema) as never,
    defaultValues: {
      profile_code: profile?.profile_code ?? '',
      name: profile?.name ?? '',
      description: profile?.description ?? '',
      download_speed: profile?.download_speed ?? 1,
      upload_speed: profile?.upload_speed ?? 1,
      price_monthly: profile?.price_monthly ?? 0,
      tax_rate: profile?.tax_rate ?? undefined,
      billing_cycle: profile?.billing_cycle ?? 'monthly',
      billing_day: profile?.billing_day ?? undefined,
      grace_period_days: profile?.grace_period_days ?? undefined,
      isolate_profile_name: profile?.isolate_profile_name ?? '',
      sort_order: profile?.sort_order ?? undefined,
      is_visible: profile?.is_visible ?? true,
      mt_local_address: profile?.mikrotik?.local_address ?? '',
      mt_remote_address: profile?.mikrotik?.remote_address ?? '',
      mt_parent_queue: profile?.mikrotik?.parent_queue ?? '',
      mt_queue_type: profile?.mikrotik?.queue_type ?? '',
      mt_dns_server: profile?.mikrotik?.dns_server ?? '',
      mt_session_timeout: profile?.mikrotik?.session_timeout ?? '',
      mt_idle_timeout: profile?.mikrotik?.idle_timeout ?? '',
    },
  })

  const onSubmit = async (data: CreateProfile) => {
    if (!routerId || !profile) return
    try {
      await updateProfile({ routerId, id: profile.id, data })
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
      <DialogContent className='sm:max-w-[425px] max-h-[90vh] overflow-y-auto'>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Pencil className="size-4" />
            Edit Bandwidth Profile
          </DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name='profile_code'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Profile Code</FormLabel>
                    <FormControl>
                      <Input placeholder='10MBPS_UNLIM' {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='name'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder='10 Mbps Unlimited' {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='download_speed'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Download (Mbps)</FormLabel>
                    <FormControl>
                      <Input type="number" {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='upload_speed'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Upload (Mbps)</FormLabel>
                    <FormControl>
                      <Input type="number" {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='price_monthly'
                render={({ field }) => (
                  <FormItem className="col-span-2">
                    <FormLabel>Monthly Price (IDR)</FormLabel>
                    <FormControl>
                      <Input type="number" {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter>
              <Button type='button' variant='outline' onClick={() => handleOpenChange(false)} disabled={isPending}>
                Cancel
              </Button>
              <Button type='submit' disabled={isPending || !routerId || !profile}>
                {isPending ? 'Saving...' : 'Save Changes'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
