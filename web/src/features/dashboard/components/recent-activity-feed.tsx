import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useUsers } from '@/hooks/use-users'

function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString)
  const diffMs = Date.now() - date.getTime()
  const diffSeconds = Math.floor(diffMs / 1000)
  if (diffSeconds < 60) return 'just now'
  const diffMinutes = Math.floor(diffSeconds / 60)
  if (diffMinutes < 60) return `${diffMinutes} min ago`
  const diffHours = Math.floor(diffMinutes / 60)
  if (diffHours < 24) return `${diffHours} hours ago`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays} days ago`
}

export function RecentActivityFeed() {
  const { data, isLoading } = useUsers(5, 0)
  const users = data?.users ?? []

  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='space-y-3'>
            {[1, 2, 3, 4, 5].map((i) => (
              <Skeleton key={i} className='h-10 w-full' />
            ))}
          </div>
        ) : users.length === 0 ? (
          <p className='text-sm text-muted-foreground py-4 text-center'>
            Activity will appear here as you manage your ISP.
          </p>
        ) : (
          <div>
            {users.map((user) => (
              <div
                key={user.id}
                className='flex items-center justify-between py-2 border-b last:border-0'
              >
                <div className='text-sm'>
                  New user: <span className='font-medium'>{user.full_name}</span> registered
                </div>
                <div className='text-xs text-muted-foreground ml-2 shrink-0'>
                  {formatRelativeTime(user.created_at)}
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
