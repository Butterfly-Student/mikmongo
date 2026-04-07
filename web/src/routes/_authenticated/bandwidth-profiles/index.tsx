import { createFileRoute } from '@tanstack/react-router'
import { BandwidthProfilesPage } from '@/features/bandwidth-profiles'

export const Route = createFileRoute('/_authenticated/bandwidth-profiles/')({
  component: BandwidthProfilesPage,
})
