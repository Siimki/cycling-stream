"use client"

import { useState, useMemo } from "react"
import { Zap, Clock, ChevronRight } from "lucide-react"
import { HudToggleButton } from "@/components/user/HudToggleButton"
import { useHudStats } from "@/components/user/HudStatsProvider"
import { useAuth } from "@/contexts/AuthContext"
import { formatTime } from "@/lib/formatters"
import { POINTS_TIERS } from "@/constants/tiers"

export function PointsDisplay() {
  const [isVisible, setIsVisible] = useState(true);
  const { isAuthenticated } = useAuth()
  const { points, watchTime } = useHudStats()

  // Tier calculation - memoized to avoid recalculation on every render
  // Must be called before any conditional returns to follow React Hooks rules
  const { currentTier, nextTier, progressToNext } = useMemo(() => {
    const tiers = POINTS_TIERS
    const current = tiers.reduce((acc, tier) => (points >= tier.min ? tier : acc), tiers[0])
    const next = tiers[tiers.indexOf(current) + 1]
    const progress = next ? ((points - current.min) / (next.min - current.min)) * 100 : 100
    return { currentTier: current, nextTier: next, progressToNext: progress }
  }, [points])

  if (!isAuthenticated) {
    return null
  }

  return (
    <div className="relative flex flex-col w-full bg-card/80 backdrop-blur-sm border-t border-border/30 overflow-hidden min-h-[4.5rem]">
      {/* Placeholder bar when collapsed - Always maintains space */}
      <div className={`transition-all duration-300 ease-in-out ${!isVisible ? "h-14 opacity-100" : "h-0 opacity-0 overflow-hidden"}`}>
        <div className="h-14 w-full flex items-center px-4">
          <span className="text-xs text-muted-foreground/50">Points collapsed</span>
        </div>
      </div>

      {/* Toggle Button - Always visible, positioned above content */}
      <HudToggleButton
        isActive={isVisible}
        label="Points info"
        onToggle={() => setIsVisible(!isVisible)}
        className="top-2 right-3"
      />

      {/* Collapsible Content */}
      <div
        className={`px-4 md:px-5 lg:px-6 py-2.5 md:py-3 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-4 overflow-x-auto transition-all duration-300 ease-in-out ${
          isVisible ? "max-h-[200px] opacity-100" : "max-h-0 opacity-0 overflow-hidden py-0"
        }`}
      >
      {/* Points & Tier - Always visible */}
      <div className="flex items-center gap-3 sm:gap-5 lg:gap-6">
        <div className="flex items-center gap-2.5 sm:gap-3">
          <div className="w-9 h-9 sm:w-10 sm:h-10 rounded-lg bg-gradient-to-br from-primary/20 to-primary/10 flex items-center justify-center shrink-0 ring-1 ring-primary/20">
            <Zap className="w-4 h-4 sm:w-5 sm:h-5 text-primary" />
          </div>
          <div>
            <div className="flex items-baseline gap-1 sm:gap-1.5">
              <span className="text-xl sm:text-2xl font-bold tabular-nums tracking-tight text-foreground/95">{points.toLocaleString()}</span>
              <span className="text-[0.65rem] sm:text-xs font-semibold text-muted-foreground uppercase tracking-wider">pts</span>
            </div>
            <div className="flex items-center gap-1 sm:gap-1.5">
              <span className={`text-[0.65rem] sm:text-xs font-bold ${currentTier.color}`}>{currentTier.name}</span>
              {nextTier && (
                <>
                  <ChevronRight className="w-2.5 h-2.5 sm:w-3 sm:h-3 text-muted-foreground/50" />
                  <span className="text-[0.65rem] sm:text-xs font-medium text-muted-foreground/70">{nextTier.name}</span>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Progress bar - Now shown on all sizes, just smaller on mobile */}
        {nextTier && (
          <div className="w-20 sm:w-28 lg:w-36">
            <div className="h-1 sm:h-1.5 bg-muted/80 rounded-full overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-primary to-primary/80 rounded-full transition-all duration-500"
                style={{ width: `${progressToNext}%` }}
              />
            </div>
            <p className="text-[0.55rem] sm:text-[0.65rem] font-semibold text-muted-foreground mt-0.5 sm:mt-1 text-center uppercase tracking-wider">
              <span className="text-foreground/90">{nextTier.min - points}</span> to {nextTier.name}
            </p>
          </div>
        )}

        {/* Watch time - Compact on mobile, full on desktop */}
        <div className="flex items-center gap-1.5 text-muted-foreground sm:border-l sm:border-border/30 sm:pl-5 lg:pl-6">
          <Clock className="w-3.5 h-3.5 text-primary/70" />
          <span className="text-xs sm:text-sm font-semibold tabular-nums text-foreground/90">{formatTime(watchTime)}</span>
          <span className="hidden sm:inline text-xs font-medium text-muted-foreground">watched</span>
        </div>
      </div>

      </div>
    </div>
  )
}
