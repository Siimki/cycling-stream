import { useMemo } from 'react';
import { useExperience } from '@/contexts/ExperienceContext';
import { motion as motionTokensFromDesign } from '@/constants/design-tokens';

export const motionTokens = motionTokensFromDesign;

const PRESET_CLASSES = {
  'chat-message-entry': 'chat-message-enter',
  'vip-ring': 'motion-pulse-ring',
  'emote-bounce': 'motion-bounce',
  'poll-announcement': 'motion-slide-in-up motion-fade-in',
  'button-hover': 'motion-fade-in',
  'overlay-fade': 'motion-overlay-fade',
  'stat-emphasis': 'motion-pop-in',
  'sweep-highlight': 'motion-sweep-line',
  'glow-ambient': 'motion-glow-ambient',
  'slide-right': 'motion-slide-in-right',
  'slide-left': 'motion-slide-in-left',
  'card-pop': 'motion-pop-in',
  'flip-up': 'motion-flip-up',
} as const;

export type MotionPresetName = keyof typeof PRESET_CLASSES;
export type MotionOptions = { disabled?: boolean };

export function useMotionPref() {
  const { resolvedUIPreferences, uiPreferences } = useExperience();
  return { resolved: resolvedUIPreferences, raw: uiPreferences };
}

export function useMotionPreset(preset: MotionPresetName, options?: MotionOptions) {
  const { resolvedUIPreferences } = useExperience();
  const disabledByPrefs = (() => {
    switch (preset) {
      case 'chat-message-entry':
      case 'emote-bounce':
      case 'vip-ring':
      case 'poll-announcement':
      case 'stat-emphasis':
      case 'sweep-highlight':
      case 'glow-ambient':
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

export function useSlideIn(direction: 'up' | 'left' | 'right' = 'up', options?: MotionOptions) {
  const mapping: Record<typeof direction, MotionPresetName> = {
    up: 'chat-message-entry',
    left: 'slide-left',
    right: 'slide-right',
  };
  return useMotionPreset(mapping[direction], options);
}

export function usePulseRing(options?: MotionOptions) {
  return useMotionPreset('vip-ring', options);
}

export function useSweepHighlight(options?: MotionOptions) {
  return useMotionPreset('sweep-highlight', options);
}

export function useAmbientGlow(options?: MotionOptions) {
  return useMotionPreset('glow-ambient', options);
}

export function usePopAccent(options?: MotionOptions) {
  return useMotionPreset('stat-emphasis', options);
}

export function useFlipEntrance(options?: MotionOptions) {
  return useMotionPreset('flip-up', options);
}
