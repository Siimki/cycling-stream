"use client"

import { useState, useMemo } from "react"
import { Zap, Clock } from "lucide-react"
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

      {/* Collapsible Content - Compact badge/card format */}
      <div
        className={`px-4 md:px-5 lg:px-6 py-3 md:py-4 transition-all duration-300 ease-in-out ${
          isVisible ? "max-h-[120px] opacity-100" : "max-h-0 opacity-0 overflow-hidden py-0"
        }`}
      >
        <div className="flex flex-col gap-3">
          {/* Main badge: Tier · Points */}
          <div className="flex items-center gap-3">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 bg-muted/30 rounded-lg border border-border/50">
              <Zap className="w-4 h-4 text-primary" />
              <span className={`text-sm font-semibold ${currentTier.color}`}>
                {currentTier.name}
              </span>
              <span className="text-muted-foreground">·</span>
              <span className="text-sm font-semibold text-foreground tabular-nums">
                {points.toLocaleString()} pts
              </span>
            </div>
          </div>

          {/* Progress to next tier */}
          {nextTier && (
            <div className="flex items-center gap-3">
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-1.5">
                  <span className="text-xs-label text-muted-foreground">
                    {nextTier.min - points} to {nextTier.name}
                  </span>
                  <span className="text-xs text-muted-foreground tabular-nums">
                    {Math.round(progressToNext)}%
                  </span>
                </div>
                <div className="h-1.5 bg-muted/80 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-gradient-to-r from-primary to-primary/80 rounded-full transition-all duration-500"
                    style={{ width: `${progressToNext}%` }}
                  />
                </div>
              </div>
              <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                <Clock className="w-3.5 h-3.5" />
                <span className="tabular-nums">{formatTime(watchTime)}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
