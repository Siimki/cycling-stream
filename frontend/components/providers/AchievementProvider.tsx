'use client';

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react';
import { createPortal } from 'react-dom';
import { useAuth } from '@/contexts/AuthContext';
import { getUserAchievements, getUserXP, type UserAchievement } from '@/lib/api';
import { useSound } from '@/components/providers/SoundProvider';
import { LevelUpOverlay } from '@/components/achievements/LevelUpOverlay';
import { AchievementToast } from '@/components/achievements/AchievementToast';

interface AchievementContextValue {
  achievements: UserAchievement[];
}

const AchievementContext = createContext<AchievementContextValue | undefined>(undefined);

interface LevelEvent {
  level: number;
  xpTotal: number;
  xpToNext: number;
}

export function AchievementProvider({ children }: { children: ReactNode }) {
  const { user, isAuthenticated } = useAuth();
  const { play } = useSound();
  const [achievements, setAchievements] = useState<UserAchievement[]>([]);
  const [levelEvent, setLevelEvent] = useState<LevelEvent | null>(null);
  const [toastQueue, setToastQueue] = useState<UserAchievement[]>([]);
  const [activeToast, setActiveToast] = useState<UserAchievement | null>(null);
  const seenAchievementsRef = useRef<Set<string>>(new Set());
  const initializedRef = useRef(false);
  const prevLevelRef = useRef<number | null>(null);
  const overlayTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const toastTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const fetchAchievements = useCallback(async () => {
    if (!isAuthenticated) {
      setAchievements([]);
      seenAchievementsRef.current.clear();
      initializedRef.current = false;
      return;
    }
    try {
      const unlocked = await getUserAchievements();
      setAchievements(unlocked);
      unlocked.forEach((achievement) => {
      if (!seenAchievementsRef.current.has(achievement.id)) {
        seenAchievementsRef.current.add(achievement.id);
        if (initializedRef.current) {
          setToastQueue((prev) => [...prev, achievement]);
        }
      }
      });
      if (!initializedRef.current) {
        initializedRef.current = true;
      }
    } catch (error) {
      if (process.env.NODE_ENV !== 'production') {
        console.info('Failed to load achievements', error);
      }
    }
  }, [isAuthenticated]);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchAchievements();
    if (!isAuthenticated) {
      return;
    }
    const interval = window.setInterval(fetchAchievements, 40000);
    return () => window.clearInterval(interval);
  }, [fetchAchievements, isAuthenticated]);

  useEffect(() => {
    if (!user || !isAuthenticated) {
      prevLevelRef.current = null;
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setLevelEvent(null);
      return;
    }

    const currentLevel = user.level ?? 1;
    if (prevLevelRef.current !== null && currentLevel > prevLevelRef.current) {
      (async () => {
        try {
          const xp = await getUserXP();
          setLevelEvent({
            level: xp.level,
            xpTotal: xp.xp_total,
            xpToNext: xp.xp_to_next_level,
          });
          play('level-up');
          if (overlayTimeoutRef.current) {
            clearTimeout(overlayTimeoutRef.current);
          }
          overlayTimeoutRef.current = setTimeout(() => {
            setLevelEvent(null);
          }, 4000);
        } catch (error) {
          console.error('Failed to load XP after level up', error);
        }
      })();
    }

    prevLevelRef.current = currentLevel;
  }, [isAuthenticated, play, user]);

  useEffect(() => {
    if (activeToast || toastQueue.length === 0) {
      return;
    }
    const next = toastQueue[0];
    Promise.resolve().then(() => {
      setActiveToast(next);
      setToastQueue((prev) => prev.slice(1));
    });
    if (toastTimeoutRef.current) {
      clearTimeout(toastTimeoutRef.current);
    }
    toastTimeoutRef.current = setTimeout(() => {
      setActiveToast(null);
    }, 4500);
  }, [activeToast, toastQueue]);

  useEffect(() => {
    return () => {
      if (overlayTimeoutRef.current) {
        clearTimeout(overlayTimeoutRef.current);
      }
      if (toastTimeoutRef.current) {
        clearTimeout(toastTimeoutRef.current);
      }
    };
  }, []);

  const contextValue = useMemo(
    () => ({
      achievements,
    }),
    [achievements]
  );

  const portalTarget =
    typeof document !== 'undefined' ? document.body : null;

  return (
    <AchievementContext.Provider value={contextValue}>
      {children}
      {portalTarget &&
        createPortal(
          <>
            <LevelUpOverlay
              level={levelEvent?.level ?? 0}
              xpTotal={levelEvent?.xpTotal ?? 0}
              xpToNext={levelEvent?.xpToNext ?? 0}
              visible={!!levelEvent}
            />
            {activeToast && (
              <AchievementToast
                achievement={activeToast}
                onDismiss={() => setActiveToast(null)}
              />
            )}
          </>,
          portalTarget
        )}
    </AchievementContext.Provider>
  );
}

export function useAchievements() {
  const ctx = useContext(AchievementContext);
  if (!ctx) {
    throw new Error('useAchievements must be used within AchievementProvider');
  }
  return ctx;
}

