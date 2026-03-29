// Raw API call functions for authentication — called directly by login pages.
// These use plain axios (not the portal clients) to avoid interceptor-auth loops.
import axios from "axios"
import {
  AdminLoginResponseSchema,
  AgentLoginResponseSchema,
  CustomerLoginResponseSchema,
} from "@/lib/schemas/auth"

export async function adminLogin(email: string, password: string) {
  const { data } = await axios.post("/api/v1/auth/login", { email, password })
  return AdminLoginResponseSchema.parse(data)
}

export async function agentLogin(email: string, password: string) {
  const { data } = await axios.post("/agent-portal/v1/login", { email, password })
  return AgentLoginResponseSchema.parse(data)
}

export async function customerLogin(email: string, password: string) {
  const { data } = await axios.post("/portal/v1/login", { email, password })
  return CustomerLoginResponseSchema.parse(data)
}
