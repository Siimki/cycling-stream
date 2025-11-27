'use client';

import { useEffect, useState, useCallback } from 'react';
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
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    ),
  }
);

export default function AnalyticsPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [raceAnalytics, setRaceAnalytics] = useState<RaceAnalytics[]>([]);
  const [watchTimeAnalytics, setWatchTimeAnalytics] = useState<WatchTimeAnalytics[]>([]);
  const [revenueAnalytics, setRevenueAnalytics] = useState<RevenueAnalytics[]>([]);
  const [selectedYear, setSelectedYear] = useState<number | undefined>(undefined);
  const [selectedMonth, setSelectedMonth] = useState<number | undefined>(undefined);

  // Memoize fetchAllAnalytics to avoid infinite loops and satisfy useEffect dependency
  const fetchAllAnalytics = useCallback(async () => {
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
  }, [selectedYear, selectedMonth]);

  useEffect(() => {
    const token = localStorage.getItem('admin_token');
    if (!token) {
      router.push('/admin/login');
      return;
    }
    fetchAllAnalytics();
  }, [router, fetchAllAnalytics]);

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
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">Loading analytics...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="bg-card shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <div className="flex items-center space-x-4">
            <Link href="/admin" className="text-primary hover:text-primary/80">
              ← Back to Admin
            </Link>
            <h1 className="text-2xl font-bold text-foreground">Analytics Dashboard</h1>
          </div>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {error && (
          <div className="mb-4 p-4 bg-destructive/10 border border-destructive/20 rounded text-destructive">
            {error}
          </div>
        )}

        {/* Filters */}
        <div className="mb-6 bg-card p-4 rounded-lg shadow">
          <div className="flex items-center space-x-4">
            <label className="text-sm font-medium text-foreground">
              Year:
              <input
                type="number"
                value={selectedYear || ''}
                onChange={(e) => setSelectedYear(e.target.value ? parseInt(e.target.value) : undefined)}
                className="ml-2 px-3 py-1 border border-border rounded bg-background text-foreground"
                placeholder="All years"
                min="2020"
                max="2100"
              />
            </label>
            <label className="text-sm font-medium text-foreground">
              Month:
              <input
                type="number"
                value={selectedMonth || ''}
                onChange={(e) => setSelectedMonth(e.target.value ? parseInt(e.target.value) : undefined)}
                className="ml-2 px-3 py-1 border border-border rounded bg-background text-foreground"
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
              className="px-4 py-1 bg-muted text-foreground rounded hover:bg-muted/80"
            >
              Clear Filters
            </button>
          </div>
        </div>

        {/* Viewer Analytics Section */}
        <section className="mb-8">
          <div className="bg-card rounded-lg shadow p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-foreground">Viewer Analytics</h2>
              <div className="space-x-2">
                <button
                  onClick={() => handleExportRaceAnalytics('csv')}
                  className="px-3 py-1 text-sm bg-success text-success-foreground rounded hover:bg-success/90"
                >
                  Export CSV
                </button>
                <button
                  onClick={() => handleExportRaceAnalytics('json')}
                  className="px-3 py-1 text-sm bg-info text-info-foreground rounded hover:bg-info/90"
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
              <p className="text-muted-foreground">No viewer data available</p>
            )}
          </div>
        </section>

        {/* Watch Time Analytics Section */}
        <section className="mb-8">
          <div className="bg-card rounded-lg shadow p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-foreground">Watch Time Analytics</h2>
              <div className="space-x-2">
                <button
                  onClick={() => handleExportWatchTime('csv')}
                  className="px-3 py-1 text-sm bg-success text-success-foreground rounded hover:bg-success/90"
                >
                  Export CSV
                </button>
                <button
                  onClick={() => handleExportWatchTime('json')}
                  className="px-3 py-1 text-sm bg-info text-info-foreground rounded hover:bg-info/90"
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
              <p className="text-muted-foreground">No watch time data available</p>
            )}
          </div>
        </section>

        {/* Revenue Analytics Section */}
        <section className="mb-8">
          <div className="bg-card rounded-lg shadow p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-foreground">Revenue Analytics</h2>
              <div className="space-x-2">
                <button
                  onClick={() => handleExportRevenue('csv')}
                  className="px-3 py-1 text-sm bg-success text-success-foreground rounded hover:bg-success/90"
                >
                  Export CSV
                </button>
                <button
                  onClick={() => handleExportRevenue('json')}
                  className="px-3 py-1 text-sm bg-info text-info-foreground rounded hover:bg-info/90"
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
              <p className="text-muted-foreground">No revenue data available</p>
            )}
          </div>
        </section>
      </main>
    </div>
  );
}
