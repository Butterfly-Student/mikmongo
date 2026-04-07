import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { customerLogin } from '@/api/auth'
import type { PortalLoginFormValues } from '@/api/types'

export function useCustomerLogin() {
  const navigate = useNavigate()
  const { customerSetToken, customerSetUser } = useAuthStore()

  return useMutation({
    mutationFn: async (data: PortalLoginFormValues) => {
      return customerLogin(data.identifier, data.password)
    },
    onSuccess: (response) => {
      customerSetToken(response.token)
      customerSetUser(response.customer)
      toast.success('Login berhasil')
      navigate({ to: '/customer', replace: true })
    },
    onError: () => {
      toast.error('Email, telepon, atau password salah')
    },
  })
}

export function useCustomerLogout() {
  const navigate = useNavigate()
  const { customerClearAuth } = useAuthStore()

  return useMutation({
    mutationFn: async () => {
      // Customer portal has no explicit logout endpoint in OpenAPI
      // Just clear auth locally
    },
    onSettled: () => {
      customerClearAuth()
      navigate({ to: '/customer-login', replace: true })
      toast.success('Logout berhasil')
    },
  })
}

export function useCustomerUser() {
  return useAuthStore((s) => s.customerUser)
}
