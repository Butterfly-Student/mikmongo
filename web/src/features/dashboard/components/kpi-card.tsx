import { Skeleton } from '@/components/ui/skeleton'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

interface KPICardProps {
  title: string
  value: string
  trend: string
  icon: React.ElementType
  isLoading?: boolean
}

export function KPICard({ title, value, trend, icon: Icon, isLoading }: KPICardProps) {
  return (
    <Card>
      <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
        <CardTitle className='text-sm font-semibold text-muted-foreground'>
          {title}
        </CardTitle>
        <Icon className='h-4 w-4 text-muted-foreground' />
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <>
            <Skeleton className='h-8 w-24' />
            <Skeleton className='mt-1 h-4 w-32' />
          </>
        ) : (
          <>
            <div className='text-2xl font-semibold'>{value}</div>
            <p className='text-sm text-muted-foreground'>{trend}</p>
          </>
        )}
      </CardContent>
    </Card>
  )
}
