import { createFileRoute } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { format, startOfMonth, endOfMonth } from "date-fns"
import { Users, Wifi, BadgeDollarSign, AlertTriangle } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { fetchSummary } from "@/api/reports"

export const Route = createFileRoute("/_admin/")({
  component: OverviewPage,
})

interface StatCardProps {
  title: string
  value: number | undefined
  icon: React.ElementType
  loading: boolean
  prefix?: string
  formatter?: (n: number) => string
}

function StatCard({ title, value, icon: Icon, loading, prefix = "", formatter }: StatCardProps) {
  const display =
    value === undefined
      ? "—"
      : formatter
        ? formatter(value)
        : `${prefix}${value.toLocaleString("id-ID")}`

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        {loading ? (
          <Skeleton className="h-8 w-24" />
        ) : (
          <p className="text-2xl font-bold">{display}</p>
        )}
      </CardContent>
    </Card>
  )
}

function formatIDR(n: number): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(n)
}

function OverviewPage() {
  const now = new Date()
  const from = format(startOfMonth(now), "yyyy-MM-dd")
  const to = format(endOfMonth(now), "yyyy-MM-dd")

  const { data, isLoading, isError } = useQuery({
    queryKey: ["reports", "summary", from, to],
    queryFn: () => fetchSummary(from, to),
    staleTime: 1000 * 60 * 5, // 5 minutes
  })

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold">Overview</h2>
        <p className="text-sm text-muted-foreground">
          {format(now, "MMMM yyyy")} summary
        </p>
      </div>

      {isError && (
        <div className="rounded-md border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          Failed to load summary data. Please refresh the page.
        </div>
      )}

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          title="Total Customers"
          value={data?.total_customers}
          icon={Users}
          loading={isLoading}
        />
        <StatCard
          title="Active Subscriptions"
          value={data?.active_subscriptions}
          icon={Wifi}
          loading={isLoading}
        />
        <StatCard
          title="Revenue This Month"
          value={data?.revenue_this_month}
          icon={BadgeDollarSign}
          loading={isLoading}
          formatter={formatIDR}
        />
        <StatCard
          title="Overdue Invoices"
          value={data?.overdue_invoices}
          icon={AlertTriangle}
          loading={isLoading}
        />
      </div>
    </div>
  )
}
