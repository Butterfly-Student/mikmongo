import { useMutation } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"
import { useStore } from "@/store"
import { customerClient } from "@/lib/axios/customer-client"
import { CustomerLoginResponseSchema } from "@/lib/schemas/auth"
import { toast } from "sonner"
import type { LoginFormValues } from "@/lib/schemas/auth"

export function useCustomerLogin() {
  const navigate = useNavigate()
  const customerSetTokens = useStore((s) => s.customerSetTokens)
  const customerSetUser = useStore((s) => s.customerSetUser)

  return useMutation({
    mutationFn: async (credentials: LoginFormValues) => {
      const { data } = await customerClient.post("/auth/login", credentials)
      return CustomerLoginResponseSchema.parse(data)
    },
    onSuccess: (response) => {
      customerSetTokens(response.data.access_token, response.data.refresh_token)
      customerSetUser({
        id: response.data.user.id,
        email: response.data.user.email,
        full_name: response.data.user.full_name,
        role: "customer",
      })
      toast.success("Login berhasil")
      navigate({ to: "/customer/dashboard" })
    },
    onError: () => {
      toast.error("Email atau password salah")
    },
  })
}

export function useCustomerLogout() {
  const navigate = useNavigate()
  const customerClearAuth = useStore((s) => s.customerClearAuth)

  return useMutation({
    mutationFn: async () => {
      try {
        await customerClient.post("/auth/logout")
      } catch {
        // Ignore network errors on logout
      }
    },
    onSettled: () => {
      customerClearAuth()
      navigate({ to: "/customer/login" })
      toast.success("Logout berhasil")
    },
  })
}

export function useCustomerUser() {
  return useStore((s) => s.customerUser)
}
