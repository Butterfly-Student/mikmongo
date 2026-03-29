import { Link, useRouterState } from "@tanstack/react-router"
import {
  LayoutDashboard,
  Users,
  Router,
  FileText,
  CreditCard,
  ClipboardList,
  UserCheck,
  Wallet,
  BarChart2,
  Wifi,
  Settings,
  Shield,
  ChevronRight,
} from "lucide-react"
import { cn } from "@/lib/utils"
import { useStore } from "@/store"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"

interface NavItem {
  label: string
  to: string
  icon: React.ElementType
  superadminOnly?: boolean
}

const NAV_ITEMS: NavItem[] = [
  { label: "Dashboard", to: "/", icon: LayoutDashboard },
  { label: "Customers", to: "/customers", icon: Users },
  { label: "Routers", to: "/routers", icon: Router },
  { label: "Subscriptions", to: "/subscriptions", icon: Wifi },
  { label: "Invoices", to: "/invoices", icon: FileText },
  { label: "Payments", to: "/payments", icon: CreditCard },
  { label: "Registrations", to: "/registrations", icon: ClipboardList },
  { label: "Sales Agents", to: "/sales-agents", icon: UserCheck },
  { label: "Cash Management", to: "/cash-entries", icon: Wallet },
  { label: "Reports", to: "/reports", icon: BarChart2 },
  { label: "MikroTik Live", to: "/mikrotik/monitor", icon: Wifi },
  { label: "Settings", to: "/settings", icon: Settings },
  { label: "Users", to: "/users", icon: Shield, superadminOnly: true },
]

interface SidebarProps {
  onNavClick?: () => void
}

export function Sidebar({ onNavClick }: SidebarProps) {
  const adminUser = useStore((s) => s.adminUser)
  const isSuperadmin = adminUser?.role === "superadmin"
  const pathname = useRouterState({ select: (s) => s.location.pathname })

  const visibleItems = NAV_ITEMS.filter((item) => !item.superadminOnly || isSuperadmin)

  const isActive = (to: string) => {
    if (to === "/") return pathname === "/"
    return pathname.startsWith(to)
  }

  return (
    <div className="flex h-full flex-col">
      {/* Brand */}
      <div className="flex h-14 items-center gap-2 border-b px-4">
        <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary text-primary-foreground text-sm font-bold">
          M
        </div>
        <span className="font-semibold text-sm">MikMongo</span>
      </div>

      {/* Navigation */}
      <ScrollArea className="flex-1 py-3">
        <nav className="space-y-0.5 px-2">
          {visibleItems.map((item) => {
            const Icon = item.icon
            const active = isActive(item.to)
            return (
              <Link
                key={item.to}
                to={item.to}
                onClick={onNavClick}
                className={cn(
                  "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                  active
                    ? "bg-primary text-primary-foreground"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                )}
              >
                <Icon className="h-4 w-4 shrink-0" />
                <span className="truncate">{item.label}</span>
                {active && <ChevronRight className="ml-auto h-3 w-3" />}
              </Link>
            )
          })}
        </nav>
      </ScrollArea>

      {/* User footer */}
      <Separator />
      <div className="px-4 py-3">
        <p className="text-xs font-medium truncate">{adminUser?.full_name ?? "—"}</p>
        <p className="text-xs text-muted-foreground truncate capitalize">{adminUser?.role ?? "—"}</p>
      </div>
    </div>
  )
}
