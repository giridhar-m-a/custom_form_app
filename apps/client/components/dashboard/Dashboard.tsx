'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useDashboard } from '@/hooks/queryHooks/useDashboard'
import { MONTHS_SHORT } from '@/lib/constants/constants'
import { cn } from '@/lib/utils'
import { getMonth } from 'date-fns'
import { Activity, ArrowUpRight, FileText, Lock, Mail, MessageSquare, TrendingUp } from 'lucide-react'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from 'recharts'

const StatCard = ({
  title,
  value,
  icon: Icon,
  description,
  trend,
  className
}: {
  title: string
  value: number | string
  icon: any
  description?: string
  trend?: string
  className?: string
}) => (
  <Card className={cn('overflow-hidden transition-all hover:shadow-md', className)}>
    <CardHeader className="flex flex-row items-center justify-between space-y-0">
      <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
      <div className="p-2 rounded-lg bg-primary/10">
        <Icon className="w-4 h-4 text-primary" />
      </div>
    </CardHeader>
    <CardContent>
      <div className="text-4xl font-bold">{value}</div>
      {description && <p className="text-sm text-muted-foreground mt-1">{description}</p>}
      {trend && (
        <div className="flex items-center mt-2 text-xs text-green-500 font-medium">
          <ArrowUpRight className="w-3 h-3 mr-1" />
          {trend}
        </div>
      )}
    </CardContent>
  </Card>
)

const DashboardLoading = () => (
  <div className="space-y-8 animate-in fade-in duration-500">
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
      {[...Array(5)].map((_, i) => (
        <Card key={i} className="overflow-hidden">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-8 w-8 rounded-lg" />
          </CardHeader>
          <CardContent>
            <Skeleton className="h-8 w-16 mb-2" />
            <Skeleton className="h-3 w-32" />
          </CardContent>
        </Card>
      ))}
    </div>
    <Card className="col-span-4">
      <CardHeader>
        <Skeleton className="h-6 w-48 mb-2" />
        <Skeleton className="h-4 w-64" />
      </CardHeader>
      <CardContent className="h-[350px]">
        <Skeleton className="h-full w-full rounded-lg" />
      </CardContent>
    </Card>
  </div>
)

const Dashboard = () => {
  const { data: dashboardData, isLoading } = useDashboard()

  if (isLoading) return <DashboardLoading />
  if (!dashboardData?.data) return <div>No data available</div>

  const data = dashboardData.data

  return (
    <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">Welcome back! Here's what's happening with your forms today.</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
        <StatCard
          title="Total Forms"
          value={data.totalForms}
          icon={FileText}
          description="Forms created in your account"
        />
        <StatCard
          title="Total Submissions"
          value={data.totalSubmissions}
          icon={MessageSquare}
          description="Total responses received"
        />
        <StatCard
          title="Active Forms"
          value={data.totalActiveForms}
          icon={Activity}
          className="border-green-500/20"
          description="Currently accepting responses"
        />
        <StatCard title="Closed Forms" value={data.totalClosedForms} icon={Lock} description="Forms no longer active" />
        <StatCard
          title="Total Invitations"
          value={data.totalInvitations}
          icon={Mail}
          description="Sent to participants"
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-full lg:col-span-4">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Submission Trends</CardTitle>
                <CardDescription>Your submission activity over the last few months.</CardDescription>
              </div>
              <div className="p-2 rounded-full bg-primary/5">
                <TrendingUp className="w-5 h-5 text-primary" />
              </div>
            </div>
          </CardHeader>
          <CardContent className="pl-2 h-[350px]">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart
                data={data.submissionsByMonth}
                margin={{
                  top: 10,
                  right: 30,
                  left: 0,
                  bottom: 0
                }}>
                <CartesianGrid strokeDasharray="5 5" stroke="var(--border)" />
                <XAxis
                  dataKey="month"
                  stroke="var(--muted-foreground)"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                  tickFormatter={value => `${MONTHS_SHORT[getMonth(value)]} ${value.split('-')[0]}`}
                  
                />
                <YAxis
                  stroke="var(--muted-foreground)"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                  tickFormatter={value => `${value}`}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'var(--card)',
                    borderColor: 'var(--border)',
                    borderRadius: '8px',
                    fontSize: '12px'
                  }}
                  itemStyle={{ color: 'var(--primary)' }}
                />
                <Line
                  type="monotone"
                  dataKey="totalSubmissions"
                  name="Total Submissions"
                  stroke="var(--primary)"
                  strokeWidth={2}
                  dot={{
                    fill: 'var(--card)',
                    stroke: 'var(--primary)',
                    strokeWidth: 2,
                    r: 4
                  }}
                  activeDot={{
                    r: 6,
                    strokeWidth: 0,
                    fill: 'var(--primary)'
                  }}
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card className="col-span-full lg:col-span-3">
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>Quick glance at form distribution.</CardDescription>
          </CardHeader>
          <CardContent className="h-[350px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={[
                  { name: 'Active', value: data.totalActiveForms },
                  { name: 'Closed', value: data.totalClosedForms }
                ]}
                layout="vertical"
                margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" horizontal={false} stroke="var(--border)" />
                <XAxis type="number" hide />
                <YAxis
                  dataKey="name"
                  type="category"
                  stroke="var(--muted-foreground)"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                />
                <Tooltip
                  cursor={{ fill: 'transparent' }}
                  contentStyle={{
                    backgroundColor: 'var(--card)',
                    borderColor: 'var(--border)',
                    borderRadius: '8px',
                    fontSize: '12px',
                    color: 'var(--muted-foreground)'
                  }}
                  itemStyle={{ color: 'var(--primary)' }}
                />
                <Bar dataKey="value" radius={[0, 4, 4, 0]} barSize={40}>
                  {[0, 1].map((_, index) => (
                    <Cell key={`cell-${index}`} fill={index === 0 ? 'var(--primary)' : '#b91c1c'} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default Dashboard
