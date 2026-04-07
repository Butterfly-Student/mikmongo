import { createFileRoute } from '@tanstack/react-router'
import CashPage from '@/features/billing/cash'

export const Route = createFileRoute('/_authenticated/cash/')({
  component: CashPage,
})
