import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { AuthLayout } from '../auth-layout'
import { adminChangePassword } from '@/api/auth'
import { ChangePasswordSchema } from '@/lib/schemas/auth'
import type { ChangePasswordValues } from '@/api/types'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { PasswordInput } from '@/components/password-input'

export function ChangePassword() {
  const navigate = useNavigate()

  const form = useForm<ChangePasswordValues>({
    resolver: zodResolver(ChangePasswordSchema),
    defaultValues: {
      old_password: '',
      new_password: '',
      confirm_password: '',
    },
  })

  async function onSubmit(data: ChangePasswordValues) {
    try {
      await adminChangePassword(data.old_password, data.new_password)
      toast.success('Password berhasil diubah')
      setTimeout(() => navigate({ to: '/' }), 1500)
    } catch (error: unknown) {
      toast.error('Password lama tidak cocok')
    }
  }

  return (
    <AuthLayout>
      <Card className='gap-4'>
        <CardHeader className='text-center'>
          <CardTitle className='text-lg tracking-tight'>
            Ubah Password
          </CardTitle>
          <CardDescription>
            Perbarui password akun Anda
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className='grid gap-3'>
              <FormField
                control={form.control}
                name='old_password'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password Lama</FormLabel>
                    <FormControl>
                      <PasswordInput placeholder='********' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='new_password'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password Baru</FormLabel>
                    <FormControl>
                      <PasswordInput placeholder='********' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='confirm_password'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Konfirmasi Password Baru</FormLabel>
                    <FormControl>
                      <PasswordInput placeholder='********' {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button
                type='submit'
                className='mt-2 w-full'
                disabled={form.formState.isSubmitting}
              >
                {form.formState.isSubmitting ? (
                  <>
                    <Loader2 className='animate-spin' />
                    Memperbarui...
                  </>
                ) : (
                  'Perbarui Password'
                )}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
