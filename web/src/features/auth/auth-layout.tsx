import { cn } from '@/lib/utils'

type AuthLayoutProps = {
  children: React.ReactNode
  className?: string
}

export function AuthLayout({ children, className }: AuthLayoutProps) {
  return (
    <div className='container grid h-svh max-w-none items-center justify-center'>
      <div
        className={cn(
          'mx-auto flex w-full flex-col justify-center space-y-2 py-8 sm:w-[480px] sm:p-8',
          className
        )}
      >
        <div className='mb-4 flex items-center justify-center'>
          <div className='flex h-10 w-10 items-center justify-center rounded-md bg-primary text-primary-foreground text-lg font-bold'>
            M
          </div>
          <h1 className='ms-2 text-xl font-medium tracking-tight'>
            MikMongo
          </h1>
        </div>
        {children}
      </div>
    </div>
  )
}
