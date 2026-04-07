import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { agentLogin } from '@/api/auth'
import type { AgentLoginFormValues } from '@/api/types'

export function useAgentLogin() {
  const navigate = useNavigate()
  const { agentSetToken, agentSetUser } = useAuthStore()

  return useMutation({
    mutationFn: async (data: AgentLoginFormValues) => {
      return agentLogin(data.username, data.password)
    },
    onSuccess: (response) => {
      agentSetToken(response.token)
      agentSetUser(response.agent)
      toast.success('Login berhasil')
      navigate({ to: '/agent', replace: true })
    },
    onError: () => {
      toast.error('Username atau password salah')
    },
  })
}

export function useAgentLogout() {
  const navigate = useNavigate()
  const { agentClearAuth } = useAuthStore()

  return useMutation({
    mutationFn: async () => {
      // Agent portal has no explicit logout endpoint in OpenAPI
      // Just clear auth locally
    },
    onSettled: () => {
      agentClearAuth()
      navigate({ to: '/agent-login', replace: true })
      toast.success('Logout berhasil')
    },
  })
}

export function useAgentUser() {
  return useAuthStore((s) => s.agentUser)
}
