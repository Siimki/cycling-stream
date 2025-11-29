'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { getStreamAnalytics, getStreamAnalyticsSummary, exportToCSV } from '@/lib/analytics';

export default function StreamAnalyticsPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [streams, setStreams] = useState<any[]>([]);
  const [summary, setSummary] = useState<any | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('admin_token');
    if (!token) {
      router.push('/admin/login');
      return;
    }
    const load = async () => {
      try {
        setError(null);
        setLoading(true);
        const [streamData, summaryData] = await Promise.all([
          getStreamAnalytics(),
          getStreamAnalyticsSummary(),
        ]);
        setStreams(streamData);
        setSummary(summaryData);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load stream analytics');
      } finally {
        setLoading(false);
      }
    };
    void load();
  }, [router]);

  const handleExport = () => {
    if (!streams.length) return;
    exportToCSV(
      streams.map((s) => ({
        stream_id: s.stream_id,
        unique_viewers: s.unique_viewers,
        total_watch_seconds: s.total_watch_seconds,
        avg_watch_seconds: s.avg_watch_seconds,
        peak_concurrent_viewers: s.peak_concurrent_viewers,
        buffer_seconds: s.buffer_seconds,
        buffer_ratio: s.buffer_ratio,
        error_rate: s.error_rate,
      })),
      `stream-analytics-${new Date().toISOString().split('T')[0]}.csv`
    );
  };

  const formatTop = (obj?: Record<string, number>) => {
    if (!obj) return '-';
    const entries = Object.entries(obj).sort((a, b) => b[1] - a[1]).slice(0, 3);
    if (!entries.length) return '-';
    return entries.map(([k, v]) => `${k}:${v}`).join(', ');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">Loading stream analytics...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="bg-card shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <div className="flex items-center space-x-4">
            <Link href="/admin/analytics" className="text-primary hover:text-primary/80">
              ← Back
            </Link>
            <h1 className="text-2xl font-bold text-foreground">Stream Analytics</h1>
          </div>
          <div className="flex items-center space-x-3">
            <button
              onClick={handleExport}
              className="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80"
            >
              Export CSV
            </button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
        {error && (
          <div className="p-4 bg-destructive/10 border border-destructive/20 rounded text-destructive">
            {error}
          </div>
        )}

        {summary && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <SummaryCard label="Streams" value={summary.stream_count} />
            <SummaryCard label="Unique viewers" value={summary.total_unique_viewers} />
            <SummaryCard
              label="Total watch (hrs)"
              value={(summary.total_watch_seconds / 3600).toFixed(1)}
            />
            <SummaryCard
              label="Avg peak concurrent"
              value={Number(summary.avg_peak_concurrent).toFixed(1)}
            />
          </div>
        )}

        <div className="bg-card rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-border">
            <thead className="bg-muted">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Stream ID
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Unique
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Watch (hrs)
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Avg (min)
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Peak CCV
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Buffer %
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Error %
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Top Countries
                </th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">
                  Devices
                </th>
              </tr>
            </thead>
            <tbody className="bg-card divide-y divide-border">
              {streams.map((s) => (
                <tr key={s.stream_id}>
                  <td className="px-4 py-3 text-sm font-mono text-foreground truncate max-w-[160px]">
                    {s.stream_id}
                  </td>
                  <td className="px-4 py-3 text-sm text-foreground">{s.unique_viewers}</td>
                  <td className="px-4 py-3 text-sm text-foreground">
                    {(s.total_watch_seconds / 3600).toFixed(1)}
                  </td>
                  <td className="px-4 py-3 text-sm text-foreground">
                    {(s.avg_watch_seconds / 60).toFixed(1)}
                  </td>
                  <td className="px-4 py-3 text-sm text-foreground">{s.peak_concurrent_viewers}</td>
                  <td className="px-4 py-3 text-sm text-foreground">
                    {s.buffer_ratio ? (s.buffer_ratio * 100).toFixed(1) + '%' : '-'}
                  </td>
                  <td className="px-4 py-3 text-sm text-foreground">
                    {s.error_rate ? (s.error_rate * 100).toFixed(1) + '%' : '-'}
                  </td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">{formatTop(s.top_countries)}</td>
                  <td className="px-4 py-3 text-sm text-muted-foreground">{formatTop(s.device_breakdown)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>
    </div>
  );
}

function SummaryCard({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="bg-card border border-border rounded-lg p-4 shadow-sm">
      <p className="text-sm text-muted-foreground">{label}</p>
      <p className="text-2xl font-semibold text-foreground">{value}</p>
    </div>
  );
}
