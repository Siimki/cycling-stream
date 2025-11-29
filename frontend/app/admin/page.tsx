'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Race } from '@/lib/api';
import { API_URL } from '@/lib/config';
import Link from 'next/link';

export default function AdminPage() {
  const router = useRouter();
  const [races, setRaces] = useState<Race[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('admin_token');
    if (!token) {
      router.push('/admin/login');
      return;
    }
    fetchRaces(token);
  }, [router]);

  const fetchRaces = async (token: string) => {
    try {
      const response = await fetch(`${API_URL}/races`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      if (!response.ok) throw new Error('Failed to fetch races');
      const data = await response.json();
      setRaces(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load races');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('admin_token');
    router.push('/admin/login');
  };

  if (loading) {
    return <div className="p-8">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="bg-card shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-foreground">Admin Panel</h1>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6 flex space-x-4">
          <Link
            href="/admin/races/new"
            className="inline-block px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90"
          >
            Create New Race
          </Link>
          <Link
            href="/admin/analytics"
            className="inline-block px-4 py-2 bg-success text-success-foreground rounded hover:bg-success/90"
          >
            View Analytics
          </Link>
          <Link
            href="/admin/analytics/streams"
            className="inline-block px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/80"
          >
            Stream Analytics
          </Link>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-destructive/10 border border-destructive/20 rounded text-destructive">
            {error}
          </div>
        )}

        <div className="bg-card rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-border">
            <thead className="bg-muted">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Location</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Start Date</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Price</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-card divide-y divide-border">
              {races.map((race) => (
                <tr key={race.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-foreground">
                    {race.name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-muted-foreground">
                    {race.location || '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-muted-foreground">
                    {race.start_date ? new Date(race.start_date).toLocaleDateString() : '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-muted-foreground">
                    {race.is_free ? 'Free' : `$${(race.price_cents / 100).toFixed(2)}`}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <Link
                      href={`/admin/races/${race.id}/edit`}
                      className="text-primary hover:text-primary/80 mr-4"
                    >
                      Edit
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>
    </div>
  );
}
