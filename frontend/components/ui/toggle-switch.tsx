'use client';

import * as React from 'react';
import { cn } from '@/lib/utils';
import { useMotionPref } from '@/motion';

interface ToggleSwitchProps extends Omit<React.ButtonHTMLAttributes<HTMLButtonElement>, 'onChange'> {
  checked: boolean;
  onCheckedChange?: (checked: boolean) => void;
  label?: string;
  description?: string;
}

export function ToggleSwitch({
  checked,
  onCheckedChange,
  disabled,
  label,
  description,
  className,
  ...props
}: ToggleSwitchProps) {
  const { resolved } = useMotionPref();
  const motionEnabled = !resolved.reduced_motion && resolved.button_pulse;

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    if (disabled) {
      return;
    }
    onCheckedChange?.(!checked);
    props.onClick?.(event);
  };

  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={handleClick}
      className={cn(
        'toggle-root',
        checked ? 'toggle-root-on' : 'toggle-root-off',
        disabled && 'opacity-60 cursor-not-allowed',
        className
      )}
      {...props}
    >
      <span
        className={cn(
          'toggle-thumb',
          checked ? 'translate-x-5' : 'translate-x-0'
        )}
        data-motion={motionEnabled ? 'on' : 'off'}
      />
      {(label || description) && (
        <span className="ml-[var(--space-3)] text-left">
          {label && <span className="block text-[var(--font-size-sm)] font-medium text-foreground leading-[var(--line-height-tight)]">{label}</span>}
          {description && <span className="block text-[var(--font-size-xs)] text-muted-foreground leading-[var(--line-height-base)]">{description}</span>}
        </span>
      )}
    </button>
  );
}
