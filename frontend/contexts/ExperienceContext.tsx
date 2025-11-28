'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState, type ReactNode } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import {
  getUserPreferences,
  updateUserPreferences,
  type AudioPreferences,
  type UIPreferences,
  type UpdateAudioPreferencesRequest,
  type UpdateUIPreferencesRequest,
  type UserPreferences,
} from '@/lib/api';
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('Experience');

const DEFAULT_UI_PREFS: UIPreferences = {
  chat_animations: true,
  reduced_motion: false,
  button_pulse: true,
  poll_animations: true,
};

const DEFAULT_AUDIO_PREFS: AudioPreferences = {
  button_clicks: true,
  notification_sounds: true,
  mention_pings: true,
  master_volume: 0.15,
};

function useSystemReducedMotion() {
  const [reducedMotion, setReducedMotion] = useState<boolean>(() => {
    if (typeof window === 'undefined' || !window.matchMedia) {
      return false;
    }
    return window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  });

  useEffect(() => {
    if (typeof window === 'undefined' || !window.matchMedia) {
      return;
    }
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    const handleChange = (event: MediaQueryListEvent) => {
      setReducedMotion(event.matches);
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  return reducedMotion;
}

interface ExperienceContextValue {
  loading: boolean;
  preferences: UserPreferences | null;
  uiPreferences: UIPreferences;
  resolvedUIPreferences: UIPreferences;
  audioPreferences: AudioPreferences;
  refreshPreferences: () => Promise<void>;
  updateUIPreferences: (updates: UpdateUIPreferencesRequest) => Promise<void>;
  updateAudioPreferences: (updates: UpdateAudioPreferencesRequest) => Promise<void>;
}

const ExperienceContext = createContext<ExperienceContextValue | undefined>(undefined);

export function ExperienceProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [loading, setLoading] = useState(true);
  const [preferences, setPreferences] = useState<UserPreferences | null>(null);
  const [uiPreferences, setUIPreferences] = useState<UIPreferences>(DEFAULT_UI_PREFS);
  const [audioPreferences, setAudioPreferences] = useState<AudioPreferences>(DEFAULT_AUDIO_PREFS);
  const hydratingRef = useRef(false);
  const systemReducedMotion = useSystemReducedMotion();

  const hydrate = useCallback(async () => {
    // Prevent concurrent hydration calls
    if (hydratingRef.current) {
      return;
    }

    // Wait for auth to finish loading before making API calls
    if (authLoading) {
      return;
    }

    if (!isAuthenticated) {
      setPreferences(null);
      setUIPreferences(DEFAULT_UI_PREFS);
      setAudioPreferences(DEFAULT_AUDIO_PREFS);
      setLoading(false);
      return;
    }

    hydratingRef.current = true;
    setLoading(true);
    try {
      const prefs = await getUserPreferences();
      setPreferences(prefs);
      setUIPreferences(prefs.ui_preferences ?? DEFAULT_UI_PREFS);
      setAudioPreferences(prefs.audio_preferences ?? DEFAULT_AUDIO_PREFS);
    } catch (error) {
      // Silently handle auth errors (401, 404, 500) - user is not authenticated or preferences don't exist
      const errorStatus = (error as { status?: number })?.status;
      if (errorStatus === 401 || errorStatus === 404 || errorStatus === 500) {
        // User is not authenticated, preferences don't exist, or token is invalid - use defaults
        // Don't log these errors as they're expected
        setPreferences(null);
        setUIPreferences(DEFAULT_UI_PREFS);
        setAudioPreferences(DEFAULT_AUDIO_PREFS);
      } else {
        // Only log unexpected errors (but don't spam the console)
        if (errorStatus !== undefined) {
          logger.error('Failed to load preferences', error);
        }
        setPreferences(null);
        setUIPreferences(DEFAULT_UI_PREFS);
        setAudioPreferences(DEFAULT_AUDIO_PREFS);
      }
    } finally {
      setLoading(false);
      hydratingRef.current = false;
    }
  }, [isAuthenticated, authLoading]);

  useEffect(() => {
    hydrate();
  }, [hydrate]);

  const resolvedUIPreferences = useMemo(() => {
    if (systemReducedMotion) {
      return {
        ...uiPreferences,
        reduced_motion: true,
        chat_animations: false,
        poll_animations: false,
        button_pulse: false,
      };
    }
    return uiPreferences;
  }, [uiPreferences, systemReducedMotion]);

  const updateUIPreferencesHandler = useCallback(
    async (updates: UpdateUIPreferencesRequest) => {
      setUIPreferences((prev) => ({
        ...prev,
        ...(updates.chat_animations !== undefined ? { chat_animations: updates.chat_animations } : null),
        ...(updates.reduced_motion !== undefined ? { reduced_motion: updates.reduced_motion } : null),
        ...(updates.button_pulse !== undefined ? { button_pulse: updates.button_pulse } : null),
        ...(updates.poll_animations !== undefined ? { poll_animations: updates.poll_animations } : null),
      }));

      try {
        const updated = await updateUserPreferences({
          ui_preferences: updates,
        });
        setPreferences(updated);
        setUIPreferences(updated.ui_preferences ?? DEFAULT_UI_PREFS);
        setAudioPreferences(updated.audio_preferences ?? DEFAULT_AUDIO_PREFS);
      } catch (error) {
        logger.error('Failed to update UI preferences', error);
        await hydrate();
        throw error;
      }
    },
    [hydrate]
  );

  const updateAudioPreferencesHandler = useCallback(
    async (updates: UpdateAudioPreferencesRequest) => {
      setAudioPreferences((prev) => ({
        ...prev,
        ...(updates.button_clicks !== undefined ? { button_clicks: updates.button_clicks } : null),
        ...(updates.notification_sounds !== undefined ? { notification_sounds: updates.notification_sounds } : null),
        ...(updates.mention_pings !== undefined ? { mention_pings: updates.mention_pings } : null),
        ...(updates.master_volume !== undefined ? { master_volume: updates.master_volume } : null),
      }));

      try {
        const updated = await updateUserPreferences({
          audio_preferences: updates,
        });
        setPreferences(updated);
        setUIPreferences(updated.ui_preferences ?? DEFAULT_UI_PREFS);
        setAudioPreferences(updated.audio_preferences ?? DEFAULT_AUDIO_PREFS);
      } catch (error) {
        logger.error('Failed to update audio preferences', error);
        await hydrate();
        throw error;
      }
    },
    [hydrate]
  );

  const value: ExperienceContextValue = useMemo(
    () => ({
      loading,
      preferences,
      uiPreferences,
      resolvedUIPreferences,
      audioPreferences,
      refreshPreferences: hydrate,
      updateUIPreferences: updateUIPreferencesHandler,
      updateAudioPreferences: updateAudioPreferencesHandler,
    }),
    [
      audioPreferences,
      hydrate,
      loading,
      preferences,
      resolvedUIPreferences,
      uiPreferences,
      updateAudioPreferencesHandler,
      updateUIPreferencesHandler,
    ]
  );

  return <ExperienceContext.Provider value={value}>{children}</ExperienceContext.Provider>;
}

export function useExperience() {
  const context = useContext(ExperienceContext);
  if (!context) {
    throw new Error('useExperience must be used within an ExperienceProvider');
  }
  return context;
}

