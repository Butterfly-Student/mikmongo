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
import { Switch } from '@/components/ui/switch'
import { PasswordInput } from '@/components/password-input'
import { createRouterSchema, type CreateRouter } from '../data/schema'
import { useCreateRouter } from '@/hooks/use-routers'

interface CreateRouterDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateRouterDialog({ open, onOpenChange }: CreateRouterDialogProps) {
  const { mutateAsync: createRouter, isPending } = useCreateRouter()

  const form = useForm<CreateRouter>({
    resolver: zodResolver(createRouterSchema) as any,
    defaultValues: {
      name: '',
      address: '',
      username: '',
      password: '',
      area: '',
      api_port: 8728,
      rest_port: 80,
      use_ssl: false,
      is_master: false,
      notes: '',
    },
  })

  const onSubmit = async (data: CreateRouter) => {
    try {
      await createRouter(data)
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
      <DialogContent className='sm:max-w-[500px] max-h-[90vh] overflow-y-auto'>
        <DialogHeader>
          <DialogTitle>Add New Router</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name='name'
                render={({ field }) => (
                  <FormItem className="col-span-2">
                    <FormLabel>Router Name</FormLabel>
                    <FormControl>
                      <Input placeholder='Core Router 01' {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='address'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>IP Address / Hostname</FormLabel>
                    <FormControl>
                      <Input placeholder='192.168.1.1 or vpn.example.com' {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='area'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Area / Location</FormLabel>
                    <FormControl>
                      <Input placeholder='Jakarta South' {...field} disabled={isPending} value={field.value || ''} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='api_port'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>API Port</FormLabel>
                    <FormControl>
                      <Input type="number" {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='rest_port'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>REST API Port</FormLabel>
                    <FormControl>
                      <Input type="number" {...field} disabled={isPending} />
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
                    <FormLabel>API Username</FormLabel>
                    <FormControl>
                      <Input placeholder='admin' {...field} disabled={isPending} />
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
                    <FormLabel>API Password</FormLabel>
                    <FormControl>
                      <PasswordInput placeholder='Router secret' {...field} disabled={isPending} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='notes'
                render={({ field }) => (
                  <FormItem className="col-span-2">
                    <FormLabel>Notes (Optional)</FormLabel>
                    <FormControl>
                      <Input placeholder='Additional context' {...field} disabled={isPending} value={field.value || ''} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='use_ssl'
                render={({ field }) => (
                  <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
                    <div className="space-y-0.5">
                      <FormLabel>Use SSL Protocol</FormLabel>
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        disabled={isPending}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='is_master'
                render={({ field }) => (
                  <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
                    <div className="space-y-0.5">
                      <FormLabel>Is Master Router</FormLabel>
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        disabled={isPending}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter>
              <Button type='button' variant='outline' onClick={() => handleOpenChange(false)} disabled={isPending}>
                Cancel
              </Button>
              <Button type='submit' disabled={isPending}>
                {isPending ? 'Creating...' : 'Add Router'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
