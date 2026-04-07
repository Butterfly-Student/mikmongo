import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { Loader2, LogIn } from 'lucide-react'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { adminLogin } from '@/api/auth'
import { LoginFormSchema } from '@/lib/schemas/auth'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/password-input'

export function UserAuthForm() {
  const navigate = useNavigate()
  const { adminSetTokens, adminSetUser } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)

  const form = useForm<z.infer<typeof LoginFormSchema>>({
    resolver: zodResolver(LoginFormSchema),
    defaultValues: { email: '', password: '' },
  })

  async function onSubmit(data: z.infer<typeof LoginFormSchema>) {
    setIsLoading(true)
    try {
      const result = await adminLogin(data.email, data.password)
      adminSetTokens(result.access_token, result.refresh_token)
      adminSetUser(result.user)
      toast.success('Login berhasil')
      navigate({ to: '/', replace: true })
    } catch (error: unknown) {
      const message =
        error instanceof Error
          ? error.message
          : 'Email atau password salah'
      console.log(message)
      toast.error(message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className='grid gap-3'>
        <FormField
          control={form.control}
          name='email'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder='name@example.com' {...field} />
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
                <PasswordInput placeholder='********' {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <div className='flex items-center justify-between'>
          <div className='flex items-center gap-2'>
            <Checkbox id='remember' />
            <label
              htmlFor='remember'
              className='text-sm text-muted-foreground cursor-pointer'
            >
              Ingat saya
            </label>
          </div>
          <button
            type='button'
            onClick={() => toast.info('Hubungi admin untuk reset password')}
            className='text-sm text-muted-foreground hover:text-primary'
          >
            Lupa password?
          </button>
        </div>
        <Button className='mt-2 w-full' disabled={isLoading}>
          {isLoading ? (
            <>
              <Loader2 className='animate-spin' />
              Sedang masuk...
            </>
          ) : (
            <>
              <LogIn />
              Masuk
            </>
          )}
        </Button>
      </form>
    </Form>
  )
}
