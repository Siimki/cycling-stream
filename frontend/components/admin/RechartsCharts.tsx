'use client';

/**
 * Recharts wrapper component for admin analytics
 * Dynamically loaded to reduce initial bundle size
 */

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import type { RaceAnalytics, WatchTimeAnalytics, RevenueAnalytics } from '@/lib/analytics';

interface ViewerChartData {
  name: string;
  fullName: string;
  concurrent: number;
  unique: number;
  authenticated: number;
  anonymous: number;
}

interface WatchTimeChartData {
  name: string;
  fullName: string;
  minutes: number;
  hours: number;
  sessions: number;
  users: number;
}

interface RevenueChartData {
  name: string;
  fullName: string;
  month: string;
  total: number;
  platform: number;
  organizer: number;
}

interface RechartsChartsProps {
  type: 'viewer' | 'watchTime' | 'revenue';
  data: ViewerChartData[] | WatchTimeChartData[] | RevenueChartData[];
  raceAnalytics?: RaceAnalytics[];
  watchTimeAnalytics?: WatchTimeAnalytics[];
  revenueAnalytics?: RevenueAnalytics[];
}

export default function RechartsCharts({
  type,
  data,
  raceAnalytics,
  watchTimeAnalytics,
  revenueAnalytics,
}: RechartsChartsProps) {
  if (type === 'viewer' && raceAnalytics) {
    const viewerData = data as ViewerChartData[];
    return (
      <>
        <div className="mb-6" style={{ width: '100%', height: 400 }}>
          <ResponsiveContainer>
            <BarChart data={viewerData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" angle={-45} textAnchor="end" height={100} />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="concurrent" fill="#0088FE" name="Concurrent Viewers" />
              <Bar dataKey="unique" fill="#00C49F" name="Unique Viewers" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="mt-6 overflow-x-auto">
          <table className="min-w-full divide-y divide-border">
            <thead className="bg-muted">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Race</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Concurrent</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Unique</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Authenticated</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Anonymous</th>
              </tr>
            </thead>
            <tbody className="bg-card divide-y divide-border">
              {raceAnalytics.map((race) => (
                <tr key={race.race_id}>
                  <td className="px-4 py-3 text-sm text-foreground">{race.race_name}</td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">{race.concurrent_viewers}</td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">{race.unique_viewers}</td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">{race.authenticated_viewers}</td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">{race.anonymous_viewers}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </>
    );
  }

  if (type === 'watchTime' && watchTimeAnalytics) {
    const watchTimeData = data as WatchTimeChartData[];
    return (
      <>
        <div className="mb-6" style={{ width: '100%', height: 400 }}>
          <ResponsiveContainer>
            <BarChart data={watchTimeData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" angle={-45} textAnchor="end" height={100} />
              <YAxis />
              <Tooltip formatter={(value) => `${value} hours`} />
              <Legend />
              <Bar dataKey="hours" fill="#FF8042" name="Watch Time (hours)" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="mt-6 overflow-x-auto">
          <table className="min-w-full divide-y divide-border">
            <thead className="bg-muted">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Race</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Total Hours</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Total Minutes</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Sessions</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Users</th>
              </tr>
            </thead>
            <tbody className="bg-card divide-y divide-border">
              {watchTimeAnalytics
                .filter((wt) => wt.total_minutes > 0)
                .map((wt) => (
                  <tr key={wt.race_id}>
                    <td className="px-4 py-3 text-sm text-foreground">{wt.race_name}</td>
                    <td className="px-4 py-3 text-sm text-muted-foreground">
                      {Math.round(wt.total_minutes / 60 * 10) / 10}
                    </td>
                    <td className="px-4 py-3 text-sm text-muted-foreground">{Math.round(wt.total_minutes)}</td>
                    <td className="px-4 py-3 text-sm text-muted-foreground">{wt.session_count}</td>
                    <td className="px-4 py-3 text-sm text-muted-foreground">{wt.user_count}</td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      </>
    );
  }

  if (type === 'revenue' && revenueAnalytics) {
    const revenueData = data as RevenueChartData[];
    return (
      <>
        <div className="mb-6" style={{ width: '100%', height: 400 }}>
          <ResponsiveContainer>
            <BarChart data={revenueData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip formatter={(value) => `$${typeof value === 'number' ? value.toFixed(2) : value}`} />
              <Legend />
              <Bar dataKey="total" fill="#8884d8" name="Total Revenue" />
              <Bar dataKey="platform" fill="#82ca9d" name="Platform Share" />
              <Bar dataKey="organizer" fill="#ffc658" name="Organizer Share" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="mt-6 overflow-x-auto">
          <table className="min-w-full divide-y divide-border">
            <thead className="bg-muted">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Race</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Month</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Total Revenue</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Platform Share</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Organizer Share</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Watch Minutes</th>
              </tr>
            </thead>
            <tbody className="bg-card divide-y divide-border">
              {revenueAnalytics.map((rev) => (
                <tr key={rev.id}>
                  <td className="px-4 py-3 text-sm text-foreground">{rev.race_name}</td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">
                    {rev.year}-{String(rev.month).padStart(2, '0')}
                  </td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">
                    ${rev.total_revenue_dollars.toFixed(2)}
                  </td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">
                    ${rev.platform_share_dollars.toFixed(2)}
                  </td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">
                    ${rev.organizer_share_dollars.toFixed(2)}
                  </td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">
                    {Math.round(rev.total_watch_minutes)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </>
    );
  }

  return null;
}

