import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { AuthLayout } from '../auth-layout'
import { UserAuthForm } from './components/user-auth-form'

export function SignIn() {
  return (
    <AuthLayout>
      <Card className='gap-4'>
        <CardHeader className='text-center'>
          <CardDescription className='text-base'>
            Admin Portal
          </CardDescription>
          <CardTitle className='text-lg tracking-tight'>
            Masuk untuk melanjutkan
          </CardTitle>
        </CardHeader>
        <CardContent>
          <UserAuthForm />
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
