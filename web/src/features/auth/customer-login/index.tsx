import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { Loader2, LogIn } from 'lucide-react'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { customerLogin } from '@/api/auth'
import { PortalLoginFormSchema } from '@/lib/schemas/auth'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
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
import { AuthLayout } from '../auth-layout'

export function CustomerLogin() {
  const navigate = useNavigate()
  const { customerSetToken, customerSetUser } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)

  const form = useForm<z.infer<typeof PortalLoginFormSchema>>({
    resolver: zodResolver(PortalLoginFormSchema),
    defaultValues: { identifier: '', password: '' },
  })

  async function onSubmit(data: z.infer<typeof PortalLoginFormSchema>) {
    setIsLoading(true)
    try {
      const result = await customerLogin(data.identifier, data.password)
      customerSetToken(result.token)
      customerSetUser(result.customer)
      toast.success('Login berhasil')
      navigate({ to: '/customer', replace: true })
    } catch (error: unknown) {
      const message =
        error instanceof Error
          ? error.message
          : 'Email, telepon, atau password salah'
      toast.error(message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <AuthLayout>
      <Card className='gap-4'>
        <CardHeader className='text-center'>
          <CardDescription className='text-base'>
            Customer Portal
          </CardDescription>
          <CardTitle className='text-lg tracking-tight'>
            Masuk untuk mengelola akun Anda
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className='grid gap-3'>
              <FormField
                control={form.control}
                name='identifier'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email, Telepon, atau Username</FormLabel>
                    <FormControl>
                      <Input placeholder='email@contoh.com' {...field} />
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
                  <Checkbox id='customer-remember' />
                  <label
                    htmlFor='customer-remember'
                    className='text-sm text-muted-foreground cursor-pointer'
                  >
                    Ingat saya
                  </label>
                </div>
                <button
                  type='button'
                  onClick={() =>
                    toast.info('Hubungi admin untuk reset password')
                  }
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
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
