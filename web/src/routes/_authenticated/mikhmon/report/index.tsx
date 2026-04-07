import { createFileRoute } from '@tanstack/react-router'
import { MikhmonReport } from '@/features/mikhmon/report'

export const Route = createFileRoute('/_authenticated/mikhmon/report/')({
  component: MikhmonReport,
})
