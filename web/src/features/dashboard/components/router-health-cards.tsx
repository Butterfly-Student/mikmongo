import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useRouters } from '@/hooks/use-routers'
import { useRouterStore } from '@/stores/router-store'
import { selectRouter } from '@/api/router'
import type { RouterResponse } from '@/lib/schemas/router'

function getStatusVariant(status: RouterResponse['status']) {
  switch (status) {
    case 'online':
      return 'default' as const
    case 'offline':
      return 'destructive' as const
    case 'unknown':
    default:
      return 'secondary' as const
  }
}

function formatRelative(dateString: string | null): string {
  if (!dateString) return 'Never'
  const date = new Date(dateString)
  const diffMs = Date.now() - date.getTime()
  const diffSeconds = Math.floor(diffMs / 1000)
  if (diffSeconds < 60) return 'just now'
  const diffMinutes = Math.floor(diffSeconds / 60)
  if (diffMinutes < 60) return `${diffMinutes} min ago`
  const diffHours = Math.floor(diffMinutes / 60)
  if (diffHours < 24) return `${diffHours}h ago`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays}d ago`
}

export function RouterHealthCards() {
  const { data, isLoading } = useRouters()
  const { selectedRouterId, setSelectedRouter } = useRouterStore()

  if (isLoading) {
    return (
      <div className='space-y-3'>
        <h2 className='text-xl font-semibold'>Router Health</h2>
        <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className='h-28 w-full rounded-lg' />
          ))}
        </div>
      </div>
    )
  }

  const routers = data?.routers ?? []

  if (routers.length === 0) {
    return (
      <div className='space-y-3'>
        <h2 className='text-xl font-semibold'>Router Health</h2>
        <Card>
          <CardContent className='flex flex-col items-center justify-center py-12 text-center'>
            <p className='text-base font-medium'>No Routers Configured</p>
            <p className='mt-1 text-sm text-muted-foreground'>
              Add your first MikroTik router to start monitoring.
            </p>
            <Button variant='outline' className='mt-4' disabled>
              Add Router
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  const displayed = routers.slice(0, 6)

  const handleClick = async (router: RouterResponse) => {
    try {
      await selectRouter(router.id)
      setSelectedRouter(router.id, router.name)
    } catch {
      // fail silently
    }
  }

  return (
    <div className='space-y-3'>
      <h2 className='text-xl font-semibold'>Router Health</h2>
      <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
        {displayed.map((router) => (
          <Card
            key={router.id}
            className={`cursor-pointer transition-all hover:shadow-md ${
              router.id === selectedRouterId ? 'ring-2 ring-primary' : ''
            }`}
            onClick={() => void handleClick(router)}
          >
            <CardHeader className='pb-2'>
              <CardTitle className='text-sm font-medium truncate'>{router.name}</CardTitle>
            </CardHeader>
            <CardContent className='space-y-1'>
              <Badge variant={getStatusVariant(router.status)} className='capitalize'>
                {router.status}
              </Badge>
              <p className='text-xs text-muted-foreground'>
                Last seen: {formatRelative(router.last_seen_at)}
              </p>
              <p className='text-xs text-muted-foreground truncate'>{router.address}</p>
            </CardContent>
          </Card>
        ))}
      </div>
      {routers.length > 6 && (
        <p className='text-sm text-muted-foreground'>
          {routers.length - 6} more routers — view all coming soon
        </p>
      )}
    </div>
  )
}
