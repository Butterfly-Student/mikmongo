import { useMutation } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"
import { useStore } from "@/store"
import { agentClient } from "@/lib/axios/agent-client"
import { AgentLoginResponseSchema } from "@/lib/schemas/auth"
import { toast } from "sonner"
import type { LoginFormValues } from "@/lib/schemas/auth"

export function useAgentLogin() {
  const navigate = useNavigate()
  const agentSetTokens = useStore((s) => s.agentSetTokens)
  const agentSetUser = useStore((s) => s.agentSetUser)

  return useMutation({
    mutationFn: async (credentials: LoginFormValues) => {
      const { data } = await agentClient.post("/auth/login", credentials)
      return AgentLoginResponseSchema.parse(data)
    },
    onSuccess: (response) => {
      agentSetTokens(response.data.access_token, response.data.refresh_token)
      agentSetUser({
        id: response.data.user.id,
        email: response.data.user.email,
        full_name: response.data.user.full_name,
        role: "agent",
      })
      toast.success("Login berhasil")
      navigate({ to: "/agent/dashboard" })
    },
    onError: () => {
      toast.error("Email atau password salah")
    },
  })
}

export function useAgentLogout() {
  const navigate = useNavigate()
  const agentClearAuth = useStore((s) => s.agentClearAuth)

  return useMutation({
    mutationFn: async () => {
      try {
        await agentClient.post("/auth/logout")
      } catch {
        // Ignore network errors on logout
      }
    },
    onSettled: () => {
      agentClearAuth()
      navigate({ to: "/agent/login" })
      toast.success("Logout berhasil")
    },
  })
}

export function useAgentUser() {
  return useStore((s) => s.agentUser)
}
