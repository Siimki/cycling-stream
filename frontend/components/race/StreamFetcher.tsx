'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { StreamResponse } from '@/lib/api';
import { getToken } from '@/lib/auth';
import { API_URL } from '@/lib/config';

interface StreamFetcherProps {
  raceId: string;
  requiresLogin: boolean;
  initialStream: StreamResponse | null;
  onStreamUpdate: (stream: StreamResponse | null) => void;
}

/**
 * Client component that fetches the stream when user is authenticated
 * This is needed because server-side fetch doesn't have access to auth tokens
 */
export function StreamFetcher({ raceId, requiresLogin, initialStream, onStreamUpdate }: StreamFetcherProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const [isFetching, setIsFetching] = useState(false);

  useEffect(() => {
    // Only fetch if:
    // 1. Auth check is complete
    // 2. User is authenticated
    // 3. Race requires login
    // 4. We don't already have a live stream
    if (isLoading || !isAuthenticated || !requiresLogin) {
      return;
    }

    // If we already have a live stream, don't refetch
    if (initialStream?.status === 'live') {
      return;
    }

    // Fetch stream with auth token
    const fetchStream = async () => {
      setIsFetching(true);
      try {
        const token = getToken();
        if (!token) {
          return;
        }

        const response = await fetch(`${API_URL}/races/${raceId}/stream`, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        if (response.ok) {
          const stream = await response.json();
          onStreamUpdate(stream);
        }
      } catch (error) {
        // Silently fail - we'll show offline state
        console.error('Failed to fetch stream:', error);
      } finally {
        setIsFetching(false);
      }
    };

    fetchStream();
  }, [raceId, requiresLogin, isAuthenticated, isLoading, initialStream, onStreamUpdate]);

  // This component doesn't render anything
  return null;
}

