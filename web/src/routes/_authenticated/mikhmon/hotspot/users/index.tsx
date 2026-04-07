import { createFileRoute } from '@tanstack/react-router'
import { HotspotUsers } from '@/features/mikhmon/hotspot/users'

export const Route = createFileRoute('/_authenticated/mikhmon/hotspot/users/')({
  component: HotspotUsers,
})
