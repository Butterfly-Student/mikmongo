import { createFileRoute } from '@tanstack/react-router'
import InvoicesPage from '@/features/billing/invoices'

export const Route = createFileRoute('/_authenticated/invoices/')({
  component: InvoicesPage,
})
