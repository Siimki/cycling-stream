'use client';
/* eslint-disable react-hooks/set-state-in-effect */

import { useEffect, useMemo, useState } from 'react';

interface LevelUpOverlayProps {
  level: number;
  xpTotal: number;
  xpToNext: number;
  visible: boolean;
}

export function LevelUpOverlay({ level, xpTotal, xpToNext, visible }: LevelUpOverlayProps) {
  const [counter, setCounter] = useState(0);

  useEffect(() => {
    if (!visible) {
      setCounter(0);
      return;
    }
    setCounter(0);
    const target = xpTotal;
    const duration = 1000;
    const steps = 20;
    const increment = target / steps;
    let current = 0;
    const interval = window.setInterval(() => {
      current += increment;
      if (current >= target) {
        current = target;
        window.clearInterval(interval);
      }
      setCounter(Math.round(current));
    }, duration / steps);
    return () => window.clearInterval(interval);
  }, [xpTotal, visible]);

  const nextLevelLabel = useMemo(() => {
    if (xpToNext <= 0) {
      return 'Maxed out';
    }
    return `${xpToNext} XP to next level`;
  }, [xpToNext]);

  if (!visible) {
    return null;
  }

  return (
    <div className="levelup-overlay">
      <div className="levelup-card motion-slide-in-up">
        <div className="levelup-burst" />
        <div className="levelup-ring" />
        <div className="levelup-content">
          <p className="levelup-label">Level Up!</p>
          <div className="levelup-level">
            <span className="levelup-level-text">Level</span>
            <span className="levelup-level-number">{level}</span>
          </div>
          <div className="levelup-xp">
            <span className="levelup-xp-counter">{counter.toLocaleString()} XP</span>
            <span className="levelup-xp-next">{nextLevelLabel}</span>
          </div>
        </div>
      </div>
    </div>
  );
}

