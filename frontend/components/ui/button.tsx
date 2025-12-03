'use client';

import * as React from "react"
import { cn } from "@/lib/utils"
import { useMotionPref } from "@/motion"
import { useRipple } from "@/hooks/useRipple"
import { useSound } from "@/components/providers/SoundProvider"
import { type SoundId } from "@/lib/sound/sound-manager"

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "default" | "destructive" | "outline" | "secondary" | "ghost" | "link"
  size?: "default" | "sm" | "lg" | "icon"
  disableMotion?: boolean
  disableRipple?: boolean
  soundId?: SoundId | null
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className,
      variant = "default",
      size = "default",
      disableMotion = false,
      disableRipple = false,
      soundId = 'button-click',
      onClick,
      children,
      ...props
    },
    ref
  ) => {
    const { resolved } = useMotionPref()
    const motionEnabled = !disableMotion && !resolved.reduced_motion && resolved.button_pulse
    const rippleEnabled = motionEnabled && !disableRipple && !props.disabled
    const { createRipple, RippleContainer } = useRipple(!rippleEnabled)
    const { play } = useSound()

    const baseStyles =
      "relative inline-flex items-center justify-center whitespace-nowrap rounded-[var(--control-radius)] font-medium ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 overflow-hidden transition-colors text-[var(--font-size-sm)]"

    const motionStyles = motionEnabled ? "button-motion" : ""

    const variants = {
      default: "bg-primary text-primary-foreground hover:bg-primary/90",
      destructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
      outline: "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
      secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80",
      ghost: "hover:bg-accent hover:text-accent-foreground",
      link: "text-primary underline-offset-4 hover:underline",
    } as const

    const sizes = {
      default: "h-[var(--control-height-md)] min-h-[var(--control-height-md)] px-[var(--space-4)]",
      sm: "h-[var(--control-height-sm)] min-h-[var(--control-height-sm)] px-[var(--space-3)]",
      lg: "h-[var(--control-height-lg)] min-h-[var(--control-height-lg)] px-[var(--space-5)] text-[var(--font-size-md)]",
      icon: "h-[var(--control-height-md)] w-[var(--control-height-md)] min-h-[var(--control-height-md)] p-0",
    } as const

    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
      if (soundId) {
        play(soundId)
      }
      if (rippleEnabled) {
        createRipple(event)
      }
      onClick?.(event)
    }

    return (
      <button
        className={cn(baseStyles, motionStyles, variants[variant], sizes[size], className)}
        ref={ref}
        onClick={handleClick}
        {...props}
      >
        {RippleContainer}
        <span className="relative z-10 inline-flex items-center gap-[var(--space-2)]">
          {children}
        </span>
      </button>
    )
  }
)

Button.displayName = "Button"

export { Button }
