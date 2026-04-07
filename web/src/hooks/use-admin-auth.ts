import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { adminLogin, adminLogout, adminChangePassword } from '@/api/auth'
import type { LoginFormValues, ChangePasswordValues } from '@/api/types'

export function useAdminLogin() {
  const navigate = useNavigate()
  const { adminSetTokens, adminSetUser } = useAuthStore()

  return useMutation({
    mutationFn: async (data: LoginFormValues) => {
      return adminLogin(data.email, data.password)
    },
    onSuccess: (response) => {
      adminSetTokens(response.access_token, response.refresh_token)
      adminSetUser(response.user)
      toast.success('Login berhasil')
      navigate({ to: '/', replace: true })
    },
    onError: () => {
      toast.error('Email atau password salah')
    },
  })
}

export function useAdminLogout() {
  const navigate = useNavigate()
  const { adminClearAuth } = useAuthStore()

  return useMutation({
    mutationFn: async () => {
      try {
        await adminLogout()
      } catch {
        // Swallow error -- clear auth regardless
      }
    },
    onSettled: () => {
      adminClearAuth()
      navigate({ to: '/sign-in', replace: true })
      toast.success('Logout berhasil')
    },
  })
}

export function useAdminChangePassword() {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: async (data: ChangePasswordValues) => {
      return adminChangePassword(data.old_password, data.new_password)
    },
    onSuccess: () => {
      toast.success('Password berhasil diubah')
      setTimeout(() => navigate({ to: '/' }), 1500)
    },
    onError: () => {
      toast.error('Password lama tidak cocok')
    },
  })
}
