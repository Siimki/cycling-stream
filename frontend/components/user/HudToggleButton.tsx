"use client"

import { Minus, List } from "lucide-react"
import { cn } from "@/lib/utils"
import type { MouseEvent } from "react"

interface HudToggleButtonProps {
  isActive: boolean
  label: string
  onToggle: () => void
  className?: string
}

export function HudToggleButton({
  isActive,
  label,
  onToggle,
  className,
}: HudToggleButtonProps) {
  const handleClick = (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault()
    event.stopPropagation()
    onToggle()
  }

  return (
    <button
      type="button"
      onClick={handleClick}
      className={cn(
        "absolute top-2.5 right-2.5 z-30 inline-flex items-center gap-2 rounded-full border border-white/15 bg-black/40 px-3 py-1.5 text-[0.55rem] font-semibold tracking-[0.25em] text-white/80 shadow-md shadow-black/40 backdrop-blur transition hover:border-white/50 hover:text-white focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white/40",
        className,
      )}
      aria-pressed={isActive}
      aria-label={isActive ? `Hide ${label}` : `Show ${label}`}
    >
      <span className="flex h-8 w-8 items-center justify-center rounded-full bg-white/10 text-white">
        {isActive ? <Minus className="h-4 w-4" /> : <List className="h-4 w-4" />}
      </span>
      <span className="text-[0.5rem] tracking-[0.3em] min-w-[4.5rem] text-center">{label}</span>
      <span
        className={cn(
          "h-2 w-2 rounded-full transition-all duration-200",
          isActive
            ? "bg-emerald-400 shadow-[0_0_12px_rgba(52,211,153,0.6)]"
            : "bg-white/40",
        )}
        aria-hidden="true"
      />
    </button>
  )
}


