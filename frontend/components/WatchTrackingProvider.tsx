'use client';

import { useWatchTracking } from '@/hooks/useWatchTracking';
import { ReactNode } from 'react';

interface WatchTrackingProviderProps {
  raceId: string;
  children: ReactNode;
}

export default function WatchTrackingProvider({ raceId, children }: WatchTrackingProviderProps) {
  // This component just triggers the tracking hook
  useWatchTracking(raceId);

  return <>{children}</>;
}

