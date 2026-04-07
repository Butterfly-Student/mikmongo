import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { BarChart3, TrendingUp, Receipt } from 'lucide-react'
import type { MikhmonReportSummary } from '@/lib/schemas/mikhmon'

function formatCurrency(amount: number | null | undefined): string {
    if (amount == null) return '-'
    return amount.toLocaleString('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 })
}

interface ReportSummaryProps {
    summary: MikhmonReportSummary | undefined
    isLoading: boolean
}

export function ReportSummary({ summary, isLoading }: ReportSummaryProps) {
    return (
        <div className='grid gap-4 md:grid-cols-3'>
            <Card>
                <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                    <CardTitle className='text-sm font-medium'>Total Transactions</CardTitle>
                    <Receipt className='size-4 text-muted-foreground' />
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <Skeleton className='h-8 w-24' />
                    ) : (
                        <div className='text-2xl font-bold'>
                            {summary?.totalCount ?? 0}
                        </div>
                    )}
                </CardContent>
            </Card>
            <Card>
                <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                    <CardTitle className='text-sm font-medium'>Total Sales</CardTitle>
                    <BarChart3 className='size-4 text-muted-foreground' />
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <Skeleton className='h-8 w-32' />
                    ) : (
                        <div className='text-2xl font-bold'>
                            {formatCurrency(summary?.totalSales)}
                        </div>
                    )}
                </CardContent>
            </Card>
            <Card>
                <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                    <CardTitle className='text-sm font-medium'>Total Revenue</CardTitle>
                    <TrendingUp className='size-4 text-muted-foreground' />
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <Skeleton className='h-8 w-32' />
                    ) : (
                        <div className='text-2xl font-bold'>
                            {formatCurrency(summary?.totalRevenue)}
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}
