import { createFileRoute } from '@tanstack/react-router'
import { Vouchers } from '@/features/mikhmon/vouchers'

export const Route = createFileRoute('/_authenticated/mikhmon/vouchers/')({
  component: Vouchers,
})
