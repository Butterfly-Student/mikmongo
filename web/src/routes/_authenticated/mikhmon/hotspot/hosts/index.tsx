import { createFileRoute } from '@tanstack/react-router'
import { HotspotHosts } from '@/features/mikhmon/hotspot/hosts'

export const Route = createFileRoute('/_authenticated/mikhmon/hotspot/hosts/')({
  component: HotspotHosts,
})
