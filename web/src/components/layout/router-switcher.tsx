import * as React from 'react'
import { ChevronsUpDown, Server, WifiOff, Wifi, HelpCircle } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import { Skeleton } from '@/components/ui/skeleton'
import { useRouters, useSelectRouter } from '@/hooks/use-routers'
import { useRouterStore } from '@/stores/router-store'
import type { RouterResponse } from '@/lib/schemas/router'

const statusConfig = {
  online: {
    label: 'Online',
    icon: Wifi,
    dot: 'bg-green-500',
    text: 'text-green-600',
  },
  offline: {
    label: 'Offline',
    icon: WifiOff,
    dot: 'bg-red-500',
    text: 'text-red-600',
  },
  unknown: {
    label: 'Unknown',
    icon: HelpCircle,
    dot: 'bg-gray-400',
    text: 'text-gray-500',
  },
} satisfies Record<RouterResponse['status'], { label: string; icon: React.ElementType; dot: string; text: string }>

function StatusDot({ status }: { status: RouterResponse['status'] }) {
  return (
    <span
      className={`inline-block size-2 shrink-0 rounded-full ${statusConfig[status].dot}`}
    />
  )
}

export function RouterSwitcher() {
  const { isMobile } = useSidebar()
  const { data, isLoading } = useRouters()
  const { selectedRouterId } = useRouterStore()
  const { mutate: selectRouter, isPending } = useSelectRouter()

  const routers = data?.routers ?? []
  const selectedRouter = routers.find((r) => r.id === selectedRouterId) ?? null

  if (isLoading) {
    return (
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size='lg' disabled>
            <div className='flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground'>
              <Server className='size-4' />
            </div>
            <div className='grid flex-1 gap-1'>
              <Skeleton className='h-3 w-24' />
              <Skeleton className='h-2.5 w-16' />
            </div>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    )
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size='lg'
              className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
              disabled={isPending}
            >
              <div className='flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground'>
                <Server className='size-4' />
              </div>
              <div className='grid flex-1 text-start text-sm leading-tight'>
                <span className='truncate font-semibold'>
                  {selectedRouter ? selectedRouter.name : 'Pilih Router'}
                </span>
                {selectedRouter ? (
                  <span className={`flex items-center gap-1 truncate text-xs ${statusConfig[selectedRouter.status].text}`}>
                    <StatusDot status={selectedRouter.status} />
                    {statusConfig[selectedRouter.status].label}
                    {selectedRouter.area ? ` · ${selectedRouter.area}` : ''}
                  </span>
                ) : (
                  <span className='truncate text-xs text-muted-foreground'>
                    {routers.length} router tersedia
                  </span>
                )}
              </div>
              <ChevronsUpDown className='ms-auto' />
            </SidebarMenuButton>
          </DropdownMenuTrigger>

          <DropdownMenuContent
            className='w-(--radix-dropdown-menu-trigger-width) min-w-64 rounded-lg'
            align='start'
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className='text-xs text-muted-foreground'>
              Router
            </DropdownMenuLabel>

            {routers.length === 0 ? (
              <>
                <DropdownMenuSeparator />
                <div className='px-2 py-4 text-center text-sm text-muted-foreground'>
                  Belum ada router terdaftar.
                </div>
              </>
            ) : (
              routers.map((router) => {
                const isActive = router.id === selectedRouterId
                return (
                  <DropdownMenuItem
                    key={router.id}
                    onClick={() => {
                      if (!isActive) selectRouter(router.id)
                    }}
                    className='gap-2 p-2'
                  >
                    <div className='flex size-6 shrink-0 items-center justify-center rounded-sm border'>
                      <StatusDot status={router.status} />
                    </div>
                    <div className='grid flex-1 text-sm leading-tight'>
                      <span className={`truncate font-medium ${isActive ? 'text-primary' : ''}`}>
                        {router.name}
                        {router.is_master && (
                          <span className='ml-1.5 text-xs font-normal text-muted-foreground'>
                            (master)
                          </span>
                        )}
                      </span>
                      <span className='truncate text-xs text-muted-foreground'>
                        {router.address}
                        {router.area ? ` · ${router.area}` : ''}
                      </span>
                    </div>
                    <span className={`text-xs ${statusConfig[router.status].text}`}>
                      {statusConfig[router.status].label}
                    </span>
                  </DropdownMenuItem>
                )
              })
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
