import { useMutation } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"
import { useStore } from "@/store"
import { adminClient } from "@/lib/axios/admin-client"
import { AdminLoginResponseSchema } from "@/lib/schemas/auth"
import { toast } from "sonner"
import type { LoginFormValues } from "@/lib/schemas/auth"

export function useAdminLogin() {
  const navigate = useNavigate()
  const adminSetTokens = useStore((s) => s.adminSetTokens)
  const adminSetUser = useStore((s) => s.adminSetUser)

  return useMutation({
    mutationFn: async (credentials: LoginFormValues) => {
      const { data } = await adminClient.post("/auth/login", credentials)
      return AdminLoginResponseSchema.parse(data)
    },
    onSuccess: (response) => {
      adminSetTokens(response.data.access_token, response.data.refresh_token)
      adminSetUser({
        id: response.data.user.id,
        email: response.data.user.email,
        role: response.data.user.role,
        full_name: response.data.user.full_name,
      })
      toast.success("Login berhasil")
      navigate({ to: "/dashboard" })
    },
    onError: () => {
      toast.error("Email atau password salah")
    },
  })
}

export function useAdminLogout() {
  const navigate = useNavigate()
  const adminClearAuth = useStore((s) => s.adminClearAuth)

  return useMutation({
    mutationFn: async () => {
      // Best-effort logout — clear local state regardless of API response
      try {
        await adminClient.post("/auth/logout")
      } catch {
        // Ignore network errors on logout
      }
    },
    onSettled: () => {
      adminClearAuth()
      navigate({ to: "/login" })
      toast.success("Logout berhasil")
    },
  })
}

export function useAdminUser() {
  return useStore((s) => s.adminUser)
}

export function useAdminRole() {
  return useStore((s) => s.adminUser?.role ?? null)
}
