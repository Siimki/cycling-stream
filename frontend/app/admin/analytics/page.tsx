'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import dynamic from 'next/dynamic';
import {
  getRaceAnalytics,
  getWatchTimeAnalytics,
  getRevenueAnalytics,
  exportToCSV,
  exportToJSON,
  type RaceAnalytics,
  type WatchTimeAnalytics,
  type RevenueAnalytics,
} from '@/lib/analytics';

// Dynamically import Recharts to reduce initial bundle size
// Only loads when admin analytics page is accessed
const RechartsCharts = dynamic(
  () => import('@/components/admin/RechartsCharts'),
  {
    ssr: false,
    loading: () => (
      <div className="flex items-center justify-center h-96">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    ),
  }
);

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d'];

export default function AnalyticsPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [raceAnalytics, setRaceAnalytics] = useState<RaceAnalytics[]>([]);
  const [watchTimeAnalytics, setWatchTimeAnalytics] = useState<WatchTimeAnalytics[]>([]);
  const [revenueAnalytics, setRevenueAnalytics] = useState<RevenueAnalytics[]>([]);
  const [selectedYear, setSelectedYear] = useState<number | undefined>(undefined);
  const [selectedMonth, setSelectedMonth] = useState<number | undefined>(undefined);

  useEffect(() => {
    const token = localStorage.getItem('admin_token');
    if (!token) {
      router.push('/admin/login');
      return;
    }
    fetchAllAnalytics();
  }, [router, selectedYear, selectedMonth]);

  const fetchAllAnalytics = async () => {
    try {
      setLoading(true);
      setError('');

      const [races, watchTime, revenue] = await Promise.all([
        getRaceAnalytics(),
        getWatchTimeAnalytics(selectedYear, selectedMonth),
        getRevenueAnalytics(selectedYear, selectedMonth),
      ]);

      setRaceAnalytics(races);
      setWatchTimeAnalytics(watchTime);
      setRevenueAnalytics(revenue);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load analytics');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('admin_token');
    router.push('/admin/login');
  };

  const handleExportRaceAnalytics = (format: 'csv' | 'json') => {
    if (format === 'csv') {
      exportToCSV(raceAnalytics, `race-analytics-${new Date().toISOString().split('T')[0]}.csv`);
    } else {
      exportToJSON(raceAnalytics, `race-analytics-${new Date().toISOString().split('T')[0]}.json`);
    }
  };

  const handleExportWatchTime = (format: 'csv' | 'json') => {
    if (format === 'csv') {
      exportToCSV(watchTimeAnalytics, `watch-time-analytics-${new Date().toISOString().split('T')[0]}.csv`);
    } else {
      exportToJSON(watchTimeAnalytics, `watch-time-analytics-${new Date().toISOString().split('T')[0]}.json`);
    }
  };

  const handleExportRevenue = (format: 'csv' | 'json') => {
    if (format === 'csv') {
      exportToCSV(revenueAnalytics, `revenue-analytics-${new Date().toISOString().split('T')[0]}.csv`);
    } else {
      exportToJSON(revenueAnalytics, `revenue-analytics-${new Date().toISOString().split('T')[0]}.json`);
    }
  };

  // Prepare chart data
  const viewerChartData = raceAnalytics.map((race) => ({
    name: race.race_name.length > 20 ? race.race_name.substring(0, 20) + '...' : race.race_name,
    fullName: race.race_name,
    concurrent: race.concurrent_viewers,
    unique: race.unique_viewers,
    authenticated: race.authenticated_viewers,
    anonymous: race.anonymous_viewers,
  }));

  const watchTimeChartData = watchTimeAnalytics
    .filter((wt) => wt.total_minutes > 0)
    .map((wt) => ({
      name: wt.race_name.length > 20 ? wt.race_name.substring(0, 20) + '...' : wt.race_name,
      fullName: wt.race_name,
      minutes: Math.round(wt.total_minutes),
      hours: Math.round(wt.total_minutes / 60 * 10) / 10,
      sessions: wt.session_count,
      users: wt.user_count,
    }));

  const revenueChartData = revenueAnalytics.map((rev) => ({
    name: rev.race_name.length > 20 ? rev.race_name.substring(0, 20) + '...' : rev.race_name,
    fullName: rev.race_name,
    month: `${rev.year}-${String(rev.month).padStart(2, '0')}`,
    total: rev.total_revenue_dollars,
    platform: rev.platform_share_dollars,
    organizer: rev.organizer_share_dollars,
  }));

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading analytics...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <div className="flex items-center space-x-4">
            <Link href="/admin" className="text-blue-600 hover:text-blue-800">
              ← Back to Admin
            </Link>
            <h1 className="text-2xl font-bold text-gray-900">Analytics Dashboard</h1>
          </div>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded text-red-600">
            {error}
          </div>
        )}

        {/* Filters */}
        <div className="mb-6 bg-white p-4 rounded-lg shadow">
          <div className="flex items-center space-x-4">
            <label className="text-sm font-medium text-gray-700">
              Year:
              <input
                type="number"
                value={selectedYear || ''}
                onChange={(e) => setSelectedYear(e.target.value ? parseInt(e.target.value) : undefined)}
                className="ml-2 px-3 py-1 border border-gray-300 rounded"
                placeholder="All years"
                min="2020"
                max="2100"
              />
            </label>
            <label className="text-sm font-medium text-gray-700">
              Month:
              <input
                type="number"
                value={selectedMonth || ''}
                onChange={(e) => setSelectedMonth(e.target.value ? parseInt(e.target.value) : undefined)}
                className="ml-2 px-3 py-1 border border-gray-300 rounded"
                placeholder="All months"
                min="1"
                max="12"
              />
            </label>
            <button
              onClick={() => {
                setSelectedYear(undefined);
                setSelectedMonth(undefined);
              }}
              className="px-4 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
            >
              Clear Filters
            </button>
          </div>
        </div>

        {/* Viewer Analytics Section */}
        <section className="mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-gray-900">Viewer Analytics</h2>
              <div className="space-x-2">
                <button
                  onClick={() => handleExportRaceAnalytics('csv')}
                  className="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
                >
                  Export CSV
                </button>
                <button
                  onClick={() => handleExportRaceAnalytics('json')}
                  className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
                >
                  Export JSON
                </button>
              </div>
            </div>

            {viewerChartData.length > 0 ? (
              <>
                <RechartsCharts
                  type="viewer"
                  data={viewerChartData}
                  raceAnalytics={raceAnalytics}
                />
              </>
            ) : (
              <p className="text-gray-500">No viewer data available</p>
            )}
          </div>
        </section>

        {/* Watch Time Analytics Section */}
        <section className="mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-gray-900">Watch Time Analytics</h2>
              <div className="space-x-2">
                <button
                  onClick={() => handleExportWatchTime('csv')}
                  className="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
                >
                  Export CSV
                </button>
                <button
                  onClick={() => handleExportWatchTime('json')}
                  className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
                >
                  Export JSON
                </button>
              </div>
            </div>

            {watchTimeChartData.length > 0 ? (
              <>
                <RechartsCharts
                  type="watchTime"
                  data={watchTimeChartData}
                  watchTimeAnalytics={watchTimeAnalytics}
                />
              </>
            ) : (
              <p className="text-gray-500">No watch time data available</p>
            )}
          </div>
        </section>

        {/* Revenue Analytics Section */}
        <section className="mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-gray-900">Revenue Analytics</h2>
              <div className="space-x-2">
                <button
                  onClick={() => handleExportRevenue('csv')}
                  className="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700"
                >
                  Export CSV
                </button>
                <button
                  onClick={() => handleExportRevenue('json')}
                  className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
                >
                  Export JSON
                </button>
              </div>
            </div>

            {revenueChartData.length > 0 ? (
              <>
                <RechartsCharts
                  type="revenue"
                  data={revenueChartData}
                  revenueAnalytics={revenueAnalytics}
                />
              </>
            ) : (
              <p className="text-gray-500">No revenue data available</p>
            )}
          </div>
        </section>
      </main>
    </div>
  );
}
