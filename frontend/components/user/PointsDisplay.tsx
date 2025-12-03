"use client"

import { useEffect, useMemo, useRef, useState } from "react"
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
  const [recentGain, setRecentGain] = useState<number | null>(null)
  const lastPointsRef = useRef<number | null>(null)

  // Tier calculation - memoized to avoid recalculation on every render
  // Must be called before any conditional returns to follow React Hooks rules
  const { currentTier, nextTier, progressToNext } = useMemo(() => {
    const tiers = POINTS_TIERS
    const current = tiers.reduce((acc, tier) => (points >= tier.min ? tier : acc), tiers[0])
    const next = tiers[tiers.indexOf(current) + 1]
    const progress = next ? ((points - current.min) / (next.min - current.min)) * 100 : 100
    return { currentTier: current, nextTier: next, progressToNext: progress }
  }, [points])

  const progressTicks = useMemo(() => {
    if (!nextTier) return []
    const span = nextTier.min - currentTier.min
    if (span <= 0) return []
    return [0.25, 0.5, 0.75].map((ratio) => ({
      ratio: ratio * 100,
      points: Math.max(1, Math.round(span * ratio)),
    }))
  }, [currentTier.min, nextTier])

  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | null = null
    if (lastPointsRef.current === null) {
      lastPointsRef.current = points
    } else if (points > (lastPointsRef.current ?? 0)) {
      setRecentGain(points - (lastPointsRef.current ?? 0))
      timer = window.setTimeout(() => setRecentGain(null), 2400)
      lastPointsRef.current = points
    } else {
      lastPointsRef.current = points
    }

    return () => {
      if (timer) {
        clearTimeout(timer)
      }
    }
  }, [points])

  if (!isAuthenticated) {
    return null
  }

  return (
    <div className="relative flex flex-col w-full bg-card/80 backdrop-blur-sm border-t border-border/30 overflow-hidden min-h-[4.5rem]">
      {recentGain && isVisible && (
        <div className="absolute top-3 left-4 z-10 px-3 py-1.5 rounded-full bg-primary/20 border border-primary/30 text-sm font-semibold text-primary shadow-lg shadow-primary/15">
          +{recentGain} pts earned
        </div>
      )}

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
        className={`px-4 md:px-5 lg:px-6 py-4 md:py-5 transition-all duration-300 ease-in-out ${
          isVisible ? "max-h-[140px] opacity-100" : "max-h-0 opacity-0 overflow-hidden py-0"
        }`}
      >
        <div className="flex flex-col gap-4">
          {/* Main badge: Tier · Points */}
          <div className="flex items-center gap-3">
            <div className="inline-flex items-center gap-2.5 px-4 py-2 bg-muted/30 rounded-lg border border-border/50">
              <Zap className="w-5 h-5 text-primary" />
              <span className={`text-base font-bold ${currentTier.color}`}>
                {currentTier.name}
              </span>
              <span className="text-muted-foreground font-medium">·</span>
              <span className="text-base font-bold text-foreground tabular-nums">
                {points.toLocaleString()} pts
              </span>
            </div>
          </div>

          <p className="text-xs text-muted-foreground/80">
            Earn points while you watch, tap reactions, and claim watch bonuses.
          </p>

          {/* Progress to next tier */}
          {nextTier && (
            <div className="flex items-center gap-4">
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-muted-foreground tracking-wide">
                    <span className="text-foreground font-semibold">{nextTier.min - points}</span> points to {nextTier.name}
                  </span>
                  <span className="text-sm font-medium text-muted-foreground tabular-nums">
                    {Math.round(progressToNext)}%
                  </span>
                </div>
                <div className="h-2.5 bg-muted/80 rounded-full overflow-hidden relative">
                  <div
                    className="h-full bg-gradient-to-r from-primary to-primary/80 rounded-full transition-all duration-500"
                    style={{ width: `${Math.min(progressToNext, 100)}%` }}
                  />
                  {progressTicks.map((tick) => (
                    <div
                      key={`tick-${tick.ratio}`}
                      className="absolute inset-y-[-4px] flex items-center justify-center"
                      style={{ left: `${tick.ratio}%` }}
                    >
                      <div className="w-[2px] h-4 bg-primary/55 rounded-full" />
                    </div>
                  ))}
                </div>
                <div className="flex items-center justify-between text-[11px] text-muted-foreground mt-2 flex-wrap gap-3">
                  <span className="font-semibold text-foreground/70">{currentTier.name}</span>
                  <div className="flex items-center gap-3 flex-wrap">
                    {progressTicks.map((tick) => (
                      <span key={`tick-label-${tick.ratio}`} className="text-xs text-foreground/70 font-medium">
                        +{tick.points.toLocaleString()} pts
                      </span>
                    ))}
                  </div>
                  <span className="font-semibold text-foreground/70">{nextTier.name}</span>
                </div>
              </div>
              <div className="flex items-center gap-2 text-sm text-muted-foreground font-medium">
                <Clock className="w-4 h-4" />
                <span className="tabular-nums">{formatTime(watchTime)}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
