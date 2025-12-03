"use client"

import { useState } from "react"
import { cn } from "@/lib/utils"

interface FloatingReactionsProps {
  enabled?: boolean
}

type Reaction = {
  icon: string
  label: string
  className: string
}

type FloatingReaction = Reaction & {
  id: number
  offset: number
}

const REACTIONS: Reaction[] = [
  { icon: "❤️", label: "Cheer", className: "text-rose-200 drop-shadow-[0_8px_24px_rgba(244,63,94,0.35)]" },
  { icon: "🚴", label: "Ride", className: "text-emerald-200 drop-shadow-[0_8px_24px_rgba(16,185,129,0.35)]" },
  { icon: "🔥", label: "Push", className: "text-orange-200 drop-shadow-[0_8px_24px_rgba(251,146,60,0.35)]" },
]

export function FloatingReactions({ enabled = true }: FloatingReactionsProps) {
  const [floating, setFloating] = useState<FloatingReaction[]>([])

  if (!enabled) {
    return null
  }

  const handleReact = (reaction: Reaction) => {
    const id = Date.now() + Math.random()
    const offset = Math.random() * 60
    const nextReaction: FloatingReaction = { ...reaction, id, offset }

    setFloating((prev) => [...prev.slice(-10), nextReaction])
    window.setTimeout(() => {
      setFloating((prev) => prev.filter((item) => item.id !== id))
    }, 1700)
  }

  return (
    <div className="pointer-events-none absolute inset-0">
      <div className="absolute right-4 bottom-6 sm:bottom-8 flex flex-col gap-2 pointer-events-auto">
        {REACTIONS.map((reaction) => (
          <button
            key={reaction.label}
            onClick={() => handleReact(reaction)}
            className="flex items-center gap-2 rounded-full bg-background/70 border border-border/50 px-3 py-1.5 text-xs font-semibold text-foreground/80 hover:text-foreground hover:border-border/80 backdrop-blur pointer-events-auto transition-colors"
            aria-label={`Send ${reaction.label} reaction`}
          >
            <span className="text-lg leading-none">{reaction.icon}</span>
            <span className="uppercase tracking-wide">{reaction.label}</span>
          </button>
        ))}
      </div>

      {floating.map((reaction) => (
        <span
          key={reaction.id}
          className={cn(
            "reaction-float absolute text-2xl will-change-transform select-none",
            reaction.className,
          )}
          style={{
            right: `${16 + reaction.offset}px`,
            bottom: `${64 + reaction.offset / 2}px`,
            animationDelay: "20ms",
          }}
          aria-hidden
        >
          {reaction.icon}
        </span>
      ))}
    </div>
  )
}
