/**
 * Custom hook for handling video player keyboard shortcuts
 */

import { useEffect } from 'react';

interface UseVideoKeyboardShortcutsProps {
  videoRef: React.RefObject<HTMLVideoElement | null>;
  togglePlay: () => void;
  toggleFullscreen: () => Promise<void>;
  toggleMute: () => void;
}

export function useVideoKeyboardShortcuts({
  videoRef,
  togglePlay,
  toggleFullscreen,
  toggleMute,
}: UseVideoKeyboardShortcutsProps) {
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      if (!videoRef.current) return;

      // Don't trigger shortcuts when typing in inputs
      if (
        e.target instanceof HTMLInputElement ||
        e.target instanceof HTMLTextAreaElement
      ) {
        return;
      }

      switch (e.key) {
        case ' ':
        case 'k':
          e.preventDefault();
          togglePlay();
          break;
        case 'f':
          e.preventDefault();
          toggleFullscreen();
          break;
        case 'm':
          e.preventDefault();
          toggleMute();
          break;
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => {
      window.removeEventListener('keydown', handleKeyPress);
    };
  }, [videoRef, togglePlay, toggleFullscreen, toggleMute]);
}

