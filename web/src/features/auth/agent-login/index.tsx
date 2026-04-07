import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { Loader2, LogIn } from 'lucide-react'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { agentLogin } from '@/api/auth'
import { AgentLoginFormSchema } from '@/lib/schemas/auth'
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

export function AgentLogin() {
  const navigate = useNavigate()
  const { agentSetToken, agentSetUser } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)

  const form = useForm<z.infer<typeof AgentLoginFormSchema>>({
    resolver: zodResolver(AgentLoginFormSchema),
    defaultValues: { username: '', password: '' },
  })

  async function onSubmit(data: z.infer<typeof AgentLoginFormSchema>) {
    setIsLoading(true)
    try {
      const result = await agentLogin(data.username, data.password)
      agentSetToken(result.token)
      agentSetUser(result.agent)
      toast.success('Login berhasil')
      navigate({ to: '/agent', replace: true })
    } catch (error: unknown) {
      const message =
        error instanceof Error
          ? error.message
          : 'Username atau password salah'
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
            Agent Portal
          </CardDescription>
          <CardTitle className='text-lg tracking-tight'>
            Masuk untuk mengelola penjualan Anda
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className='grid gap-3'>
              <FormField
                control={form.control}
                name='username'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>
                    <FormControl>
                      <Input placeholder='agent_username' {...field} />
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
                  <Checkbox id='agent-remember' />
                  <label
                    htmlFor='agent-remember'
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
