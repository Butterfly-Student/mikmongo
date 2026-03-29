import { useState, useCallback } from "react"
import { Sheet, SheetContent } from "@/components/ui/sheet"
import { Sidebar } from "./Sidebar"
import { Topbar } from "./Topbar"

interface AppShellProps {
  children: React.ReactNode
}

export function AppShell({ children }: AppShellProps) {
  const [mobileOpen, setMobileOpen] = useState(false)

  const handleMenuClick = useCallback(() => setMobileOpen(true), [])
  const handleNavClick = useCallback(() => setMobileOpen(false), [])

  return (
    <div className="flex h-svh overflow-hidden bg-background">
      {/* Desktop sidebar — fixed, 240px */}
      <aside className="hidden lg:flex lg:w-60 lg:shrink-0 lg:flex-col border-r bg-background">
        <Sidebar />
      </aside>

      {/* Mobile sidebar — Sheet overlay */}
      <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
        <SheetContent side="left" className="w-60 p-0">
          <Sidebar onNavClick={handleNavClick} />
        </SheetContent>
      </Sheet>

      {/* Main content */}
      <div className="flex flex-1 flex-col overflow-hidden">
        <Topbar onMenuClick={handleMenuClick} />
        <main className="flex-1 overflow-auto p-4 md:p-6">{children}</main>
      </div>
    </div>
  )
}
