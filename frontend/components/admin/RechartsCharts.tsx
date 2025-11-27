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
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Race</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Concurrent</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Unique</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Authenticated</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Anonymous</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {raceAnalytics.map((race) => (
                <tr key={race.race_id}>
                  <td className="px-4 py-3 text-sm text-gray-900">{race.race_name}</td>
                  <td className="px-4 py-3 text-sm text-gray-500">{race.concurrent_viewers}</td>
                  <td className="px-4 py-3 text-sm text-gray-500">{race.unique_viewers}</td>
                  <td className="px-4 py-3 text-sm text-gray-500">{race.authenticated_viewers}</td>
                  <td className="px-4 py-3 text-sm text-gray-500">{race.anonymous_viewers}</td>
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
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Race</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total Hours</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total Minutes</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Sessions</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Users</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {watchTimeAnalytics
                .filter((wt) => wt.total_minutes > 0)
                .map((wt) => (
                  <tr key={wt.race_id}>
                    <td className="px-4 py-3 text-sm text-gray-900">{wt.race_name}</td>
                    <td className="px-4 py-3 text-sm text-gray-500">
                      {Math.round(wt.total_minutes / 60 * 10) / 10}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-500">{Math.round(wt.total_minutes)}</td>
                    <td className="px-4 py-3 text-sm text-gray-500">{wt.session_count}</td>
                    <td className="px-4 py-3 text-sm text-gray-500">{wt.user_count}</td>
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
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Race</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Month</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total Revenue</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Platform Share</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Organizer Share</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Watch Minutes</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {revenueAnalytics.map((rev) => (
                <tr key={rev.id}>
                  <td className="px-4 py-3 text-sm text-gray-900">{rev.race_name}</td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    {rev.year}-{String(rev.month).padStart(2, '0')}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    ${rev.total_revenue_dollars.toFixed(2)}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    ${rev.platform_share_dollars.toFixed(2)}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    ${rev.organizer_share_dollars.toFixed(2)}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500">
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

