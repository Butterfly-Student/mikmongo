import { createFileRoute } from '@tanstack/react-router'
import { HotspotActiveSessions } from '@/features/mikhmon/hotspot/active'

export const Route = createFileRoute('/_authenticated/mikhmon/hotspot/active/')({
  component: HotspotActiveSessions,
})
