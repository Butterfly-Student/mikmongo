import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useReportSummary } from '@/hooks/use-report-summary'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'

function formatRupiah(n: number): string {
  return 'Rp ' + n.toLocaleString('id-ID')
}

// Static monthly placeholder data using the total_revenue as a reference
function buildChartData(totalRevenue: number) {
  // Since /reports/summary returns aggregates only, we show a representative breakdown
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun']
  const base = totalRevenue / 6
  return months.map((month, idx) => ({
    name: month,
    revenue: Math.round(base * (0.7 + (idx * 0.1))),
  }))
}

export function RevenueChart() {
  const { data, isLoading } = useReportSummary()

  const chartData = data ? buildChartData(data.total_revenue) : []

  return (
    <Card>
      <CardHeader>
        <CardTitle>Revenue Overview</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex h-48 items-center justify-center text-sm text-muted-foreground'>
            Loading...
          </div>
        ) : data ? (
          <ResponsiveContainer width='100%' height={200}>
            <BarChart data={chartData} margin={{ top: 0, right: 0, left: 0, bottom: 0 }}>
              <CartesianGrid strokeDasharray='3 3' className='stroke-muted' />
              <XAxis
                dataKey='name'
                tick={{ fontSize: 12 }}
                tickLine={false}
                axisLine={false}
              />
              <YAxis
                tick={{ fontSize: 12 }}
                tickLine={false}
                axisLine={false}
                tickFormatter={(v) => `Rp ${(v / 1000000).toFixed(0)}M`}
              />
              <Tooltip
                formatter={(value: number) => [formatRupiah(value), 'Revenue']}
              />
              <Bar dataKey='revenue' fill='hsl(var(--chart-1))' radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        ) : (
          <div className='flex h-48 items-center justify-center text-sm text-muted-foreground'>
            No revenue data available
          </div>
        )}
      </CardContent>
    </Card>
  )
}
