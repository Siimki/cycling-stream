"use client"

import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useCallback,
  useState,
  type ReactNode,
} from "react"
import { awardBonusPoints, getProfile } from "@/lib/auth"
import { useAuth } from "@/contexts/AuthContext"
import { WATCH_TIME_UPDATE_INTERVAL_MS, POINTS_REFRESH_INTERVAL_MS, BONUS_COOLDOWN_SECONDS } from "@/constants/intervals"
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('HudStats');

interface HudStatsContextValue {
  points: number
  watchTime: number
  bonusReady: boolean
  claimBonus: () => void
}

const HudStatsContext = createContext<HudStatsContextValue | null>(null)

export function HudStatsProvider({ children }: { children: ReactNode }) {
  const { token, refreshUser } = useAuth()
  const [points, setPoints] = useState(0)
  const [watchTime, setWatchTime] = useState(0)
  const [bonusReady, setBonusReady] = useState(false)
  const [cooldown, setCooldown] = useState(0)

  // Accumulate local watch time (UI only, not persisted)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => {
    const interval = setInterval(() => {
      // setState in interval callback is valid - not synchronous setState in effect
      setWatchTime((prev) => prev + 1)
    }, WATCH_TIME_UPDATE_INTERVAL_MS)

    return () => clearInterval(interval)
  }, []) // Empty deps intentional - interval should run once

  // Load real points from backend and refresh periodically
  useEffect(() => {
    if (!token) {
      setPoints(0)
      return
    }

    const fetchPoints = async () => {
      try {
        const user = await getProfile(token)
        setPoints(user.points ?? 0)
      } catch (err) {
        logger.error("Failed to fetch user points:", err)
      }
    }

    fetchPoints()

    const interval = setInterval(fetchPoints, POINTS_REFRESH_INTERVAL_MS)
    return () => clearInterval(interval)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]) // getProfile is stable, no need to include in deps

  // Surface bonus every 30 seconds when cooldown is clear (fast iteration)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => {
    if (watchTime > 0 && watchTime % BONUS_COOLDOWN_SECONDS === 0 && cooldown === 0) {
      setBonusReady(true)
    }
  }, [watchTime, cooldown]) // Intentional setState in effect for derived state

  // Cooldown countdown
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => {
    if (cooldown <= 0) return

    const timer = setTimeout(() => {
      setCooldown((prev) => (prev > 0 ? prev - 1 : 0))
    }, 1000)

    return () => clearTimeout(timer)
  }, [cooldown]) // Intentional setState in effect for timer-based state

  const claimBonus = useCallback(() => {
    if (!bonusReady || !token) {
      setBonusReady(false)
      setCooldown(0)
      return
    }

    // Optimistically hide bonus and start cooldown
    setBonusReady(false)
    setCooldown(BONUS_COOLDOWN_SECONDS)

    awardBonusPoints(token)
      .then((totalPoints) => {
        const safeTotal = totalPoints ?? 0
        setPoints(safeTotal)
        // Refresh user data in auth context
        refreshUser()
      })
      .catch((err) => {
        logger.error("Failed to award bonus points:", err)
      })
  }, [bonusReady, token, refreshUser])

  const value = useMemo(
    () => ({
      points,
      watchTime,
      bonusReady,
      claimBonus,
    }),
    [points, watchTime, bonusReady, claimBonus],
  )

  return <HudStatsContext.Provider value={value}>{children}</HudStatsContext.Provider>
}

export function useHudStats() {
  const context = useContext(HudStatsContext)
  if (!context) {
    throw new Error("useHudStats must be used within a HudStatsProvider")
  }
  return context
}


