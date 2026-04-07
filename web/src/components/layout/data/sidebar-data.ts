import {
  LayoutDashboard,
  Users,
  UsersRound,
  Server,
  CreditCard,
  FileText,
  Wallet,
  Banknote,
  UserCog,
  Flame,
  Network,
  Wifi,
  Gauge,
  Cpu,
  ScrollText,
  Activity,
  BarChart3,
  Settings,
  EthernetPort,
  MonitorSmartphone,
  Ticket,
  PrinterCheck,
  CalendarClock,
  PencilRuler,
} from 'lucide-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  user: {
    name: 'Admin',
    email: '',
    avatar: '',
  },
  teams: [
    {
      name: 'MikMongo ISP',
      logo: Server,
      plan: 'ISP Management',
    },
  ],
  navGroups: [
    {
      title: 'Overview',
      items: [
        {
          title: 'Dashboard',
          url: '/',
          icon: LayoutDashboard,
        },
      ],
    },
    {
      title: 'Management',
      items: [
        {
          title: 'Users',
          url: '/users',
          icon: Users,
        },
        {
          title: 'Customers',
          url: '/customers',
          icon: UsersRound,
        },
        {
          title: 'Routers',
          url: '/routers',
          icon: Server,
        },
        {
          title: 'Subscriptions',
          url: '/subscriptions',
          icon: CreditCard,
        },
        {
          title: 'Bandwidth Profiles',
          url: '/bandwidth-profiles',
          icon: Gauge,
        },
      ],
    },
    {
      title: 'Billing',
      items: [
        {
          title: 'Invoices',
          url: '/invoices',
          icon: FileText,
        },
        {
          title: 'Payments',
          url: '/payments',
          icon: Wallet,
        },
        {
          title: 'Cash',
          url: '/cash',
          icon: Banknote,
        },
      ],
    },
    // {
    //   title: 'Sales',
    //   items: [
    //     {
    //       title: 'Agents',
    //       url: '#',
    //       icon: UserCog,
    //     },
    //     {
    //       title: 'Hotspot Sales',
    //       url: '#',
    //       icon: Flame,
    //     },
    //   ],
    // },
    // {
    //   title: 'MikroTik',
    //   items: [
    //     {
    //       title: 'PPP',
    //       url: '#',
    //       icon: Network,
    //     },
    //     {
    //       title: 'Hotspot',
    //       url: '#',
    //       icon: Wifi,
    //     },
    //     {
    //       title: 'Network',
    //       url: '#',
    //       icon: Globe,
    //     },
    //   ],
    // },
    // {
    //   title: 'Monitor',
    //   items: [
    //     {
    //       title: 'System Resources',
    //       url: '#',
    //       icon: Cpu,
    //     },
    //     {
    //       title: 'Interfaces',
    //       url: '#',
    //       icon: EthernetPort,
    //     },
    //     {
    //       title: 'Logs',
    //       url: '#',
    //       icon: ScrollText,
    //     },
    //     {
    //       title: 'Ping',
    //       url: '#',
    //       icon: Activity,
    //     },
    //   ],
    // },
    {
      title: 'MikHmon',
      items: [
        {
          title: 'Hotspot',
          icon: Wifi,
          items: [
            {
              title: 'Users',
              url: '/mikhmon/hotspot/users',
            },
            {
              title: 'User Profiles',
              url: '/mikhmon/hotspot/profiles',
            },
            {
              title: 'Active Sessions',
              url: '/mikhmon/hotspot/active',
            },
            {
              title: 'Hosts',
              url: '/mikhmon/hotspot/hosts',
            },
          ],
        },
        {
          title: 'Vouchers',
          url: '/mikhmon/vouchers',
          icon: Ticket,
        },
        {
          title: 'Report',
          url: '/mikhmon/report',
          icon: BarChart3,
        },
      ],
    },
    {
      title: 'Reports',

      items: [
        {
          title: 'Business Reports',
          url: '/reports',
          icon: BarChart3,
        },
      ],
    },
    {
      title: 'Settings',
      items: [
        {
          title: 'Settings',
          url: '/settings',
          icon: Settings,
        },
      ],
    },
  ],
}
