"use client"

import { useState, useEffect } from "react"
import { Zap, Clock, Gift, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"

interface PointsDisplayProps {
  points: number
  watchTime: number
  addPoints: (amount: number) => void
}

export function PointsDisplay({ points, watchTime, addPoints }: PointsDisplayProps) {
  const [bonusReady, setBonusReady] = useState(false)
  const [cooldown, setCooldown] = useState(0)
  const [showClaimed, setShowClaimed] = useState(false)

  // Tier calculation
  const tiers = [
    { name: "Bronze", min: 0, color: "text-amber-600" },
    { name: "Silver", min: 500, color: "text-slate-400" },
    { name: "Gold", min: 1500, color: "text-yellow-500" },
    { name: "Platinum", min: 3500, color: "text-cyan-400" },
    { name: "Diamond", min: 7500, color: "text-violet-400" },
  ]

  const currentTier = tiers.reduce((acc, tier) => (points >= tier.min ? tier : acc), tiers[0])
  const nextTier = tiers[tiers.indexOf(currentTier) + 1]
  const progressToNext = nextTier ? ((points - currentTier.min) / (nextTier.min - currentTier.min)) * 100 : 100

  const formatTime = (s: number) => {
    const h = Math.floor(s / 3600)
    const m = Math.floor((s % 3600) / 60)
    return h > 0 ? `${h}h ${m}m` : `${m}m`
  }

  // Bonus ready every 5 minutes
  useEffect(() => {
    if (watchTime > 0 && watchTime % 300 === 0 && cooldown === 0) {
      setBonusReady(true)
    }
  }, [watchTime, cooldown])

  useEffect(() => {
    if (cooldown > 0) {
      const t = setTimeout(() => setCooldown(cooldown - 1), 1000)
      return () => clearTimeout(t)
    }
  }, [cooldown])

  const claimBonus = () => {
    addPoints(50)
    setBonusReady(false)
    setCooldown(300)
    setShowClaimed(true)
    setTimeout(() => setShowClaimed(false), 2000)
  }

  return (
    <div className="px-4 py-3 flex items-center justify-between gap-4 flex-wrap">
      {/* Points & Tier */}
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
            <Zap className="w-4 h-4 text-primary" />
          </div>
          <div>
            <div className="flex items-baseline gap-1.5">
              <span className="text-lg font-bold tabular-nums">{points.toLocaleString()}</span>
              <span className="text-[10px] text-muted-foreground uppercase tracking-wide">points</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className={`text-[10px] font-medium ${currentTier.color}`}>{currentTier.name}</span>
              {nextTier && (
                <>
                  <ChevronRight className="w-2.5 h-2.5 text-muted-foreground" />
                  <span className="text-[10px] text-muted-foreground">{nextTier.name}</span>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Progress bar */}
        {nextTier && (
          <div className="hidden sm:block w-32">
            <div className="h-1 bg-muted rounded-full overflow-hidden">
              <div
                className="h-full bg-primary rounded-full transition-all duration-500"
                style={{ width: `${progressToNext}%` }}
              />
            </div>
            <p className="text-[9px] text-muted-foreground mt-0.5 text-center">
              {nextTier.min - points} to {nextTier.name}
            </p>
          </div>
        )}

        {/* Watch time */}
        <div className="hidden md:flex items-center gap-1.5 text-muted-foreground">
          <Clock className="w-3.5 h-3.5" />
          <span className="text-xs tabular-nums">{formatTime(watchTime)}</span>
          <span className="text-[10px]">watched</span>
        </div>
      </div>

      {/* Earning info + Bonus */}
      <div className="flex items-center gap-3">
        <div className="hidden lg:flex items-center gap-4 text-[10px] text-muted-foreground mr-2">
          <span className="flex items-center gap-1">
            <span className="w-1 h-1 rounded-full bg-primary" />
            +10/min watching
          </span>
          <span className="flex items-center gap-1">
            <span className="w-1 h-1 rounded-full bg-amber-500" />
            +50 bonus
          </span>
        </div>

        {bonusReady ? (
          <Button
            size="sm"
            className="h-7 text-xs bg-amber-500 hover:bg-amber-600 text-black font-medium gap-1.5 glow-primary"
            onClick={claimBonus}
          >
            <Gift className="w-3 h-3" />
            Claim +50
          </Button>
        ) : showClaimed ? (
          <Button size="sm" className="h-7 text-xs bg-primary/20 text-primary font-medium" disabled>
            Claimed!
          </Button>
        ) : (
          <Button size="sm" variant="secondary" className="h-7 text-xs text-muted-foreground font-normal" disabled>
            <Gift className="w-3 h-3 mr-1.5" />
            {cooldown > 0
              ? `${Math.floor(cooldown / 60)}:${(cooldown % 60).toString().padStart(2, "0")}`
              : "Keep watching"}
          </Button>
        )}
      </div>
    </div>
  )
}
