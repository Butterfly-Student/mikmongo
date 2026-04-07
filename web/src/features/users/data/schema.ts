import { z } from 'zod'

export const userTableFilterSchema = z.object({
  search: z.string().optional(),
  role: z.string().optional(),
})

export type UserTableFilter = z.infer<typeof userTableFilterSchema>
