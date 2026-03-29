import { z } from "zod"
import { adminClient } from "@/lib/axios/admin-client"

export const SummarySchema = z.object({
  success: z.boolean(),
  data: z.object({
    total_customers: z.number(),
    active_subscriptions: z.number(),
    revenue_this_month: z.number(),
    overdue_invoices: z.number(),
  }),
})

export type Summary = z.infer<typeof SummarySchema>["data"]

export async function fetchSummary(from: string, to: string): Promise<Summary> {
  const { data } = await adminClient.get("/reports/summary", { params: { from, to } })
  const parsed = SummarySchema.parse(data)
  return parsed.data
}
