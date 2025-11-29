'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { API_URL } from '@/lib/config';
import { createContextLogger } from '@/lib/logger';

type OutboundEvent = {
  type: string;
  videoTime?: number;
  extra?: Record<string, unknown>;
};

const logger = createContextLogger('useAnalyticsTracking');
const CLIENT_ID_KEY = 'cs_analytics_client_id';
const FLUSH_INTERVAL_MS = 5000;
const HEARTBEAT_BATCH_SIZE = 10;

export function useAnalyticsTracking(streamId?: string) {
  const [clientId, setClientId] = useState<string | null>(null);
  const queueRef = useRef<OutboundEvent[]>([]);
  const flushTimeoutRef = useRef<number | null>(null);

  // Initialize clientId once on mount
  useEffect(() => {
    try {
      const existing = window.localStorage.getItem(CLIENT_ID_KEY);
      if (existing) {
        setClientId(existing);
        return;
      }
      const generated = crypto.randomUUID();
      window.localStorage.setItem(CLIENT_ID_KEY, generated);
      setClientId(generated);
    } catch (error) {
      logger.error('Failed to initialize analytics client id', error);
      setClientId(null);
    }
  }, []);

  const flush = useCallback(
    async (options?: { keepalive?: boolean }) => {
      if (!streamId || !clientId || queueRef.current.length === 0) {
        return;
      }

      const events = queueRef.current.splice(0, queueRef.current.length);
      const payload = {
        streamId,
        clientId,
        events: events.map((evt) => ({
          type: evt.type,
          videoTime: evt.videoTime,
          extra: evt.extra,
        })),
      };

      try {
        await fetch(`${API_URL}/analytics/events`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(payload),
          keepalive: options?.keepalive === true,
        });
      } catch (error) {
        logger.error('Failed to flush analytics events', error);
      }
    },
    [clientId, streamId]
  );

  // Flush queued events as soon as we have a clientId
  useEffect(() => {
    if (clientId) {
      void flush();
    }
  }, [clientId, flush]);

  // Clear timers and flush on unmount
  useEffect(() => {
    return () => {
      if (flushTimeoutRef.current) {
        window.clearTimeout(flushTimeoutRef.current);
      }
      void flush({ keepalive: true });
    };
  }, [flush]);

  const enqueue = useCallback(
    (type: string, videoTime?: number, extra?: Record<string, unknown>) => {
      queueRef.current.push({ type, videoTime, extra });

      if (queueRef.current.length >= HEARTBEAT_BATCH_SIZE) {
        void flush();
      } else if (flushTimeoutRef.current === null) {
        flushTimeoutRef.current = window.setTimeout(() => {
          void flush();
          flushTimeoutRef.current = null;
        }, FLUSH_INTERVAL_MS);
      }
    },
    [flush]
  );

  const trackPlay = useCallback((videoTime?: number) => enqueue('play', videoTime), [enqueue]);
  const trackPause = useCallback((videoTime?: number) => enqueue('pause', videoTime), [enqueue]);
  const trackHeartbeat = useCallback((videoTime?: number) => enqueue('heartbeat', videoTime), [enqueue]);
  const trackEnded = useCallback((videoTime?: number) => enqueue('ended', videoTime), [enqueue]);
  const trackError = useCallback(
    (videoTime?: number, extra?: Record<string, unknown>) => enqueue('error', videoTime, extra),
    [enqueue]
  );
  const trackBufferStart = useCallback((videoTime?: number) => enqueue('buffer_start', videoTime), [enqueue]);
  const trackBufferEnd = useCallback((videoTime?: number) => enqueue('buffer_end', videoTime), [enqueue]);

  return {
    clientId,
    trackPlay,
    trackPause,
    trackHeartbeat,
    trackEnded,
    trackError,
    trackBufferStart,
    trackBufferEnd,
    flush,
  };
}
