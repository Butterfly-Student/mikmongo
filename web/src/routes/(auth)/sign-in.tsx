import { createFileRoute, redirect } from '@tanstack/react-router'
import { SignIn } from '@/features/auth/sign-in'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/(auth)/sign-in')({
  beforeLoad: () => {
    if (useAuthStore.getState().adminIsAuthenticated) {
      throw redirect({ to: '/' })
    }
  },
  component: SignIn,
})
