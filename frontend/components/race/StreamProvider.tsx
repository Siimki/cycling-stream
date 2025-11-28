'use client';

import { useState, useCallback } from 'react';
import { StreamResponse } from '@/lib/api';
import DynamicVideoPlayer from '@/components/video/DynamicVideoPlayer';
import { StreamFetcher } from './StreamFetcher';

interface StreamProviderProps {
  raceId: string;
  requiresLogin: boolean;
  initialStream: StreamResponse | null;
}

/**
 * Client component that manages stream state and fetches stream when authenticated
 */
export function StreamProvider({ raceId, requiresLogin, initialStream }: StreamProviderProps) {
  const [stream, setStream] = useState<StreamResponse | null>(initialStream);

  const handleStreamUpdate = useCallback((newStream: StreamResponse | null) => {
    setStream(newStream);
  }, []);

  const streamUrl = stream?.cdn_url || stream?.origin_url;

  return (
    <>
      <StreamFetcher
        raceId={raceId}
        requiresLogin={requiresLogin}
        initialStream={stream}
        onStreamUpdate={handleStreamUpdate}
      />
      <DynamicVideoPlayer
        streamUrl={streamUrl}
        status={stream?.status || 'offline'}
        streamType={stream?.stream_type}
        sourceId={stream?.source_id}
        requiresLogin={requiresLogin}
        raceId={raceId}
      />
    </>
  );
}

