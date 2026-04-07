import { createFileRoute } from '@tanstack/react-router'
import { Routers } from '@/features/routers'

export const Route = createFileRoute('/_authenticated/routers/')({
  component: Routers,
})
