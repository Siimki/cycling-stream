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
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Admin Panel</h1>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6 flex space-x-4">
          <Link
            href="/admin/races/new"
            className="inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Create New Race
          </Link>
          <Link
            href="/admin/analytics"
            className="inline-block px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
          >
            View Analytics
          </Link>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded text-red-600">
            {error}
          </div>
        )}

        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Location</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Start Date</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Price</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {races.map((race) => (
                <tr key={race.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {race.name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {race.location || '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {race.start_date ? new Date(race.start_date).toLocaleDateString() : '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {race.is_free ? 'Free' : `$${(race.price_cents / 100).toFixed(2)}`}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <Link
                      href={`/admin/races/${race.id}/edit`}
                      className="text-blue-600 hover:text-blue-900 mr-4"
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

