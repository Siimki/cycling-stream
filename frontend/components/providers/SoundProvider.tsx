'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';
import { useExperience } from '@/contexts/ExperienceContext';
import { soundManager, type PlaySoundOptions, type SoundId } from '@/lib/sound/sound-manager';

interface SoundContextValue {
  ready: boolean;
  play: (soundId: SoundId, options?: PlaySoundOptions) => void;
}

const SoundContext = createContext<SoundContextValue | undefined>(undefined);

const SOUND_TOGGLE_MAP: Record<SoundId, (prefs: ReturnType<typeof useExperience>['audioPreferences']) => boolean> = {
  'button-click': (prefs) => prefs.button_clicks,
  notification: (prefs) => prefs.notification_sounds,
  'chat-mention': (prefs) => prefs.mention_pings,
  'level-up': (prefs) => prefs.notification_sounds,
};

export function SoundProvider({ children }: { children: ReactNode }) {
  const { audioPreferences } = useExperience();
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const enabledSounds = Object.entries(SOUND_TOGGLE_MAP)
      .filter(([, predicate]) => predicate(audioPreferences))
      .map(([key]) => key as SoundId);

    if (enabledSounds.length > 0) {
      soundManager.preload(enabledSounds).finally(() => setReady(true));
    } else {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setReady(true);
    }
  }, [audioPreferences]);

  const play = useCallback(
    (soundId: SoundId, options?: PlaySoundOptions) => {
      const predicate = SOUND_TOGGLE_MAP[soundId];
      if (predicate && !predicate(audioPreferences)) {
        return;
      }

      soundManager.play(soundId, {
        ...options,
        masterVolume: audioPreferences.master_volume,
      });
    },
    [audioPreferences]
  );

  const value = useMemo(
    () => ({
      ready,
      play,
    }),
    [play, ready]
  );

  return <SoundContext.Provider value={value}>{children}</SoundContext.Provider>;
}

export function useSound() {
  const context = useContext(SoundContext);
  if (!context) {
    throw new Error('useSound must be used within a SoundProvider');
  }
  return context;
}

