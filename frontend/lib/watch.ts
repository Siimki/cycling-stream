import { API_URL } from './config';

export interface WatchSession {
  id: string;
  user_id: string;
  race_id: string;
  started_at: string;
  ended_at?: string;
  duration_seconds?: number;
  created_at: string;
}

export interface WatchTimeStats {
  user_id: string;
  race_id: string;
  session_count: number;
  total_seconds: number;
  total_minutes: number;
  first_watched: string;
  last_watched?: string;
}

export async function startWatchSession(raceId: string, token: string): Promise<WatchSession> {
  const response = await fetch(`${API_URL}/users/watch/sessions/start`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ race_id: raceId }),
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to start watch session' }));
    throw new Error(error.error || 'Failed to start watch session');
  }

  return response.json();
}

export async function endWatchSession(sessionId: string, token: string): Promise<void> {
  const response = await fetch(`${API_URL}/users/watch/sessions/end`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ session_id: sessionId }),
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to end watch session' }));
    throw new Error(error.error || 'Failed to end watch session');
  }
}

export async function getWatchStats(raceId: string, token: string): Promise<WatchTimeStats> {
  const response = await fetch(`${API_URL}/users/watch/sessions/stats/${raceId}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to get watch stats' }));
    throw new Error(error.error || 'Failed to get watch stats');
  }

  return response.json();
}

