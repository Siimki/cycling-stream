import { API_URL } from './config';

export interface RaceAnalytics {
  race_id: string;
  race_name: string;
  concurrent_viewers: number;
  authenticated_viewers: number;
  anonymous_viewers: number;
  unique_viewers: number;
  unique_authenticated: number;
  unique_anonymous: number;
}

export interface WatchTimeAnalytics {
  race_id: string;
  race_name: string;
  total_seconds: number;
  total_minutes: number;
  session_count: number;
  user_count: number;
  year?: number;
  month?: number;
}

export interface RevenueAnalytics {
  id: string;
  race_id: string;
  race_name: string;
  year: number;
  month: number;
  total_revenue_cents: number;
  total_revenue_dollars: number;
  total_watch_minutes: number;
  platform_share_cents: number;
  platform_share_dollars: number;
  organizer_share_cents: number;
  organizer_share_dollars: number;
  calculated_at: string;
  created_at: string;
  updated_at: string;
}

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('admin_token') : null;
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options?.headers as Record<string, string> || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  const data = await response.json();
  return data.data || data;
}

export async function getRaceAnalytics(): Promise<RaceAnalytics[]> {
  return fetchAPI<RaceAnalytics[]>('/admin/analytics/races');
}

export async function getWatchTimeAnalytics(year?: number, month?: number): Promise<WatchTimeAnalytics[]> {
  const params = new URLSearchParams();
  if (year) params.append('year', year.toString());
  if (month) params.append('month', month.toString());
  
  const query = params.toString();
  const endpoint = query ? `/admin/analytics/watch-time?${query}` : '/admin/analytics/watch-time';
  
  return fetchAPI<WatchTimeAnalytics[]>(endpoint);
}

export async function getRevenueAnalytics(year?: number, month?: number): Promise<RevenueAnalytics[]> {
  const params = new URLSearchParams();
  if (year) params.append('year', year.toString());
  if (month) params.append('month', month.toString());
  
  const query = params.toString();
  const endpoint = query ? `/admin/analytics/revenue?${query}` : '/admin/analytics/revenue';
  
  return fetchAPI<RevenueAnalytics[]>(endpoint);
}

// Export functions for CSV/JSON
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function exportToCSV<T extends Record<string, any>>(data: T[], filename: string): void {
  if (data.length === 0) return;

  const headers = Object.keys(data[0]);
  const csvRows = [
    headers.join(','),
    ...data.map(row =>
      headers.map(header => {
        const value = row[header];
        // Handle values that might contain commas or quotes
        if (value === null || value === undefined) return '';
        const stringValue = String(value);
        if (stringValue.includes(',') || stringValue.includes('"') || stringValue.includes('\n')) {
          return `"${stringValue.replace(/"/g, '""')}"`;
        }
        return stringValue;
      }).join(',')
    ),
  ];

  const csvContent = csvRows.join('\n');
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);
  
  link.setAttribute('href', url);
  link.setAttribute('download', filename);
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

export function exportToJSON<T>(data: T[], filename: string): void {
  const jsonContent = JSON.stringify(data, null, 2);
  const blob = new Blob([jsonContent], { type: 'application/json;charset=utf-8;' });
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);
  
  link.setAttribute('href', url);
  link.setAttribute('download', filename);
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

