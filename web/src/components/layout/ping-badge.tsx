import { Wifi, WifiOff } from 'lucide-react'
import { useRouterStore } from '@/stores/router-store'
import { usePing } from '@/hooks/use-ping'
import { cn } from '@/lib/utils'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

function getLatencyColor(latency: string, status: string): string {
  if (status === 'connecting') return 'text-muted-foreground'
  if (status === 'error' || status === 'unavailable') return 'text-muted-foreground'
  const ms = parseInt(latency)
  if (isNaN(ms)) return 'text-muted-foreground'
  if (ms < 50) return 'text-green-500'
  if (ms < 150) return 'text-yellow-500'
  return 'text-red-500'
}

function getStatusDot(status: string): string {
  if (status === 'connected') return 'bg-green-500'
  if (status === 'connecting') return 'bg-yellow-500 animate-pulse'
  return 'bg-muted-foreground'
}

export function PingBadge() {
  const { selectedRouterId, selectedRouterName } = useRouterStore()
  const { latency, status } = usePing(selectedRouterId)

  const isActive = status !== 'unavailable'
  const label = isActive ? latency : '-- ms'
  const color = getLatencyColor(latency, status)
  const tooltipText = selectedRouterName
    ? `Ping ke ${selectedRouterName} (8.8.8.8)`
    : 'Tidak ada router dipilih'

  return (
    <TooltipProvider delayDuration={300}>
      <Tooltip>
        <TooltipTrigger asChild>
          <div
            className={cn(
              'flex items-center gap-1.5 rounded-md px-2 py-1',
              'text-xs font-mono font-medium',
              'border bg-background/60',
              'select-none cursor-default',
              color
            )}
            aria-label={`Ping: ${label}`}
          >
            <span
              className={cn(
                'inline-block size-1.5 rounded-full flex-shrink-0',
                getStatusDot(status)
              )}
            />
            {isActive ? (
              <Wifi className='size-3 flex-shrink-0' />
            ) : (
              <WifiOff className='size-3 flex-shrink-0 text-muted-foreground' />
            )}
            <span className='min-w-[3ch] text-center'>{label}</span>
          </div>
        </TooltipTrigger>
        <TooltipContent side='bottom'>
          <p>{tooltipText}</p>
          {isActive && status === 'connected' && (
            <p className='text-muted-foreground text-xs'>Realtime via WebSocket</p>
          )}
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
