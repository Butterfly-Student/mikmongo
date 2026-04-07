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
import { createProfileSchema, type CreateProfile } from '../data/schema'
import { useCreateProfile } from '@/hooks/use-profiles'

interface CreateProfileDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  routerId: string | null
}

export function CreateProfileDialog({ open, onOpenChange, routerId }: CreateProfileDialogProps) {
  const { mutateAsync: createProfile, isPending } = useCreateProfile()

  const form = useForm<CreateProfile>({
    resolver: zodResolver(createProfileSchema) as any,
    defaultValues: {
      profile_code: '',
      name: '',
      description: '',
      download_speed: 10,
      upload_speed: 10,
      price_monthly: 0,
      billing_cycle: 'monthly',
      is_visible: true,
    },
  })

  const onSubmit = async (data: CreateProfile) => {
    if (!routerId) return
    try {
      await createProfile({ routerId, data })
      form.reset()
      onOpenChange(false)
    } catch (error) {
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
          <DialogTitle>Add New Profile</DialogTitle>
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
              <Button type='submit' disabled={isPending || !routerId}>
                {isPending ? 'Creating...' : 'Add Profile'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
