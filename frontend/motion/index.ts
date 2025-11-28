import { useMemo } from 'react';
import { useExperience } from '@/contexts/ExperienceContext';
import { motion as motionTokensFromDesign } from '@/constants/design-tokens';

export const motionTokens = motionTokensFromDesign;

const PRESET_CLASSES = {
  'chat-message-entry': 'motion-slide-in-up motion-fade-in',
  'vip-ring': 'motion-pulse-ring',
  'emote-bounce': 'motion-bounce',
  'poll-announcement': 'motion-slide-in-up motion-fade-in',
  'button-hover': 'motion-fade-in',
} as const;

export type MotionPresetName = keyof typeof PRESET_CLASSES;

export function useMotionPref() {
  const { resolvedUIPreferences, uiPreferences } = useExperience();
  return { resolved: resolvedUIPreferences, raw: uiPreferences };
}

export function useMotionPreset(preset: MotionPresetName, options?: { disabled?: boolean }) {
  const { resolvedUIPreferences } = useExperience();
  const disabledByPrefs = (() => {
    switch (preset) {
      case 'chat-message-entry':
      case 'emote-bounce':
      case 'vip-ring':
      case 'poll-announcement':
        return !resolvedUIPreferences.chat_animations || resolvedUIPreferences.reduced_motion;
      case 'button-hover':
        return !resolvedUIPreferences.button_pulse || resolvedUIPreferences.reduced_motion;
      default:
        return resolvedUIPreferences.reduced_motion;
    }
  })();

  const disabled = options?.disabled ?? disabledByPrefs;

  return useMemo(() => (disabled ? '' : PRESET_CLASSES[preset] ?? ''), [disabled, preset]);
}

