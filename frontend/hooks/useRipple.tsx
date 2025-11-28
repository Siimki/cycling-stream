'use client';

import { useCallback, useEffect, useRef, useState, type MouseEvent } from 'react';

type Ripple = {
  key: number;
  size: number;
  x: number;
  y: number;
};

export function useRipple(disabled: boolean) {
  const [ripples, setRipples] = useState<Ripple[]>([]);
  const nextKey = useRef(0);

  const createRipple = useCallback(
    (event: MouseEvent<HTMLElement>) => {
      if (disabled) {
        return;
      }
      const element = event.currentTarget as HTMLElement;
      const rect = element.getBoundingClientRect();
      const size = Math.max(rect.width, rect.height);
      const x = event.clientX - rect.left - size / 2;
      const y = event.clientY - rect.top - size / 2;

      setRipples((prev) => [
        ...prev,
        {
          key: nextKey.current++,
          size,
          x,
          y,
        },
      ]);
    },
    [disabled]
  );

  useEffect(() => {
    if (ripples.length === 0) {
      return undefined;
    }

    const timer = setTimeout(() => {
      setRipples((prev) => prev.slice(1));
    }, 500);

    return () => clearTimeout(timer);
  }, [ripples]);

  const RippleContainer = (
    <span className="ripple-container" aria-hidden="true">
      {ripples.map((ripple) => (
        <span
          key={ripple.key}
          className="ripple-circle"
          style={{
            width: ripple.size,
            height: ripple.size,
            top: ripple.y,
            left: ripple.x,
          }}
        />
      ))}
    </span>
  );

  return { createRipple, RippleContainer };
}

