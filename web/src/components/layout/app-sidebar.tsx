import { useLayout } from '@/context/layout-provider'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/components/ui/sidebar'
import { sidebarData } from './data/sidebar-data'
import { NavGroup } from './nav-group'
import { NavUser } from './nav-user'
import { RouterSwitcher } from './router-switcher'
import { useAuthStore } from '@/stores/auth-store'

export function AppSidebar() {
  const { collapsible, variant } = useLayout()
  const adminUser = useAuthStore((s) => s.adminUser)

  const navUser = {
    name: adminUser?.full_name ?? sidebarData.user.name,
    email: adminUser?.email ?? sidebarData.user.email,
    avatar: sidebarData.user.avatar,
  }

  return (
    <Sidebar collapsible={collapsible} variant={variant}>
      <SidebarHeader>
        <RouterSwitcher />
      </SidebarHeader>
      <SidebarContent>
        {sidebarData.navGroups.map((props) => (
          <NavGroup key={props.title} {...props} />
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={navUser} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
