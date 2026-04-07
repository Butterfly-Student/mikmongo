import { createFileRoute } from '@tanstack/react-router'
import BusinessReportsPage from '@/features/reports'

export const Route = createFileRoute('/_authenticated/reports/')({
  component: BusinessReportsPage,
})
