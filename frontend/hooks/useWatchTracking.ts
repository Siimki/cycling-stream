 'use client';
 
 import { useEffect, useRef, useState } from 'react';
 import { getToken, awardWatchPoints } from '@/lib/auth';
 import { startWatchSession, endWatchSession } from '@/lib/watch';
 import { useAuth } from '@/contexts/AuthContext';
 import { POINTS_ACCRUAL_INTERVAL_MS } from '@/constants/intervals';
 import { createContextLogger } from '@/lib/logger';

 const logger = createContextLogger('WatchTracking');
 
 export function useWatchTracking(raceId: string | null) {
   const { refreshUser } = useAuth();
   const [sessionId, setSessionId] = useState<string | null>(null);
   const [isTracking, setIsTracking] = useState(false);
   const sessionIdRef = useRef<string | null>(null);
   const raceIdRef = useRef<string | null>(raceId);
   const tickIntervalRef = useRef<number | null>(null);
  // Update raceId ref when it changes
  useEffect(() => {
    raceIdRef.current = raceId;
  }, [raceId]);

  const startTracking = async () => {
    if (!raceId || isTracking) return;

    const token = getToken();
    if (!token) {
      logger.debug('No auth token, skipping watch tracking');
      return;
    }

    try {
      const session = await startWatchSession(raceId, token);
      setSessionId(session.id);
      sessionIdRef.current = session.id;
      setIsTracking(true);
    } catch (error) {
      logger.error('Failed to start watch tracking:', error);
    }
  };

  const stopTracking = async () => {
    if (!sessionIdRef.current || !isTracking) return;

    const token = getToken();
    if (!token) return;

    try {
      await endWatchSession(sessionIdRef.current, token);
      setSessionId(null);
      sessionIdRef.current = null;
      setIsTracking(false);
    } catch (error) {
      logger.error('Failed to stop watch tracking:', error);
    }
  };

  // Start tracking when component mounts and raceId is available
  useEffect(() => {
    if (raceId) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      startTracking();
    }

    // Stop tracking on unmount
    return () => {
      if (sessionIdRef.current) {
        stopTracking();
      }
    };
  }, [raceId]);

  // Handle page visibility changes (tab switch, minimize, etc.)
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (document.hidden) {
        // Page is hidden, pause tracking
        if (sessionIdRef.current && isTracking) {
          stopTracking();
        }
      } else {
        // Page is visible again, resume tracking
        if (raceIdRef.current && !isTracking) {
          startTracking();
        }
      }
    };
 
    document.addEventListener('visibilitychange', handleVisibilityChange);
 
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [isTracking, raceId]);

  // Live points accrual: award bonus every 10 seconds while tracking.
  useEffect(() => {
    const startTicks = () => {
      if (!raceIdRef.current || !isTracking || tickIntervalRef.current !== null) return;

      const token = getToken();
      if (!token) return;

      const intervalId = window.setInterval(async () => {
        try {
          const currentToken = getToken();
          if (!currentToken || !raceIdRef.current || !sessionIdRef.current) {
            return;
          }

          // Award 10 points via watch tick endpoint
          await awardWatchPoints(currentToken);
          // Refresh user data in auth context
          refreshUser();
        } catch (error) {
          logger.error('Failed to tick watch points:', error);
        }
      }, POINTS_ACCRUAL_INTERVAL_MS);

      tickIntervalRef.current = intervalId;
    };

    const stopTicks = () => {
      if (tickIntervalRef.current !== null) {
        window.clearInterval(tickIntervalRef.current);
        tickIntervalRef.current = null;
      }
    };

    if (isTracking) {
      startTicks();
    } else {
      stopTicks();
    }

    return () => {
      stopTicks();
    };
  }, [isTracking]);

  // Handle page unload (browser close, navigation away)
  useEffect(() => {
    const handleBeforeUnload = () => {
      if (sessionIdRef.current) {
        // Use sendBeacon for reliable delivery on page unload
        const token = getToken();
        if (token) {
          const data = JSON.stringify({ session_id: sessionIdRef.current });
          const url = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/users/watch/sessions/end`; // Note: sendBeacon doesn't support custom headers
          // Note: sendBeacon doesn't support custom headers, so we'll need to handle this differently
          // For now, we'll rely on the cleanup function
        }
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, []);

  return {
    isTracking,
    sessionId,
    startTracking,
    stopTracking,
  };
}

