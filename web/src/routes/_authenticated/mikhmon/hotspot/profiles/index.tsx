import { createFileRoute } from '@tanstack/react-router'
import { HotspotProfiles } from '@/features/mikhmon/hotspot/profiles'

export const Route = createFileRoute('/_authenticated/mikhmon/hotspot/profiles/')({
  component: HotspotProfiles,
})
