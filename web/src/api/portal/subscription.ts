import { customerClient } from '@/lib/axios/customer-client'
import { SubscriptionResponseSchema } from '@/lib/schemas/subscription'
import { z } from 'zod'

const PortalSubscriptionListResponseSchema = z.object({
    success: z.boolean(),
    data: z.array(SubscriptionResponseSchema),
})

export async function listPortalSubscriptions(): Promise<z.infer<typeof PortalSubscriptionListResponseSchema>['data']> {
    const response = await customerClient.get('/subscriptions')
    const parsed = PortalSubscriptionListResponseSchema.parse(response.data)
    return parsed.data
}
