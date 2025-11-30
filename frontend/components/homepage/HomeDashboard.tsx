'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { 
  getRecommendations, 
  getUserXP, 
  getUserWeekly, 
  getUserMissions,
  getWatchHistory,
  RecommendationsResponse,
  XPProgress,
  WeeklyGoalProgress,
  UserMissionWithDetails,
  WatchHistoryResponse,
} from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import { WelcomeHeader } from './WelcomeHeader';
import { QuickStatsGrid } from './QuickStatsGrid';
import { ContinueWatchingSection } from './ContinueWatchingSection';
import { UpcomingFromFavorites } from './UpcomingFromFavorites';
import { GamificationSection } from './GamificationSection';
import { FairModelSection } from './FairModelSection';

interface HomeDashboardProps {
  liveRacesCount?: number;
}

export function HomeDashboard({ liveRacesCount = 0 }: HomeDashboardProps) {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Data state
  const [recommendations, setRecommendations] = useState<RecommendationsResponse | null>(null);
  const [xpData, setXpData] = useState<XPProgress | null>(null);
  const [weeklyData, setWeeklyData] = useState<WeeklyGoalProgress | null>(null);
  const [missions, setMissions] = useState<UserMissionWithDetails[]>([]);
  const [watchHistory, setWatchHistory] = useState<WatchHistoryResponse | null>(null);

  useEffect(() => {
    if (authLoading || !isAuthenticated) {
      setLoading(false);
      return;
    }

    async function loadDashboardData() {
      setLoading(true);
      setError(null);

      try {
        // Load all data in parallel
        const [
          recsResult,
          xpResult,
          weeklyResult,
          missionsResult,
          historyResult,
        ] = await Promise.allSettled([
          getRecommendations().catch(() => null),
          getUserXP().catch(() => null),
          getUserWeekly().catch(() => null),
          getUserMissions().catch(() => null),
          getWatchHistory(20, 0).catch(() => null),
        ]);

        if (recsResult.status === 'fulfilled' && recsResult.value) {
          setRecommendations(recsResult.value);
        }
        if (xpResult.status === 'fulfilled' && xpResult.value) {
          setXpData(xpResult.value);
        }
        if (weeklyResult.status === 'fulfilled' && weeklyResult.value) {
          setWeeklyData(weeklyResult.value);
        }
        if (missionsResult.status === 'fulfilled' && missionsResult.value) {
          setMissions(Array.isArray(missionsResult.value) ? missionsResult.value : []);
        }
        if (historyResult.status === 'fulfilled' && historyResult.value) {
          setWatchHistory(historyResult.value);
        }
      } catch (err) {
        const errorMsg = APIErrorHandler.getErrorMessage(err);
        setError(errorMsg);
      } finally {
        setLoading(false);
      }
    }

    loadDashboardData();
  }, [isAuthenticated, authLoading]);


  if (authLoading || loading) {
    return (
      <div className="flex-1 max-w-7xl mx-auto px-6 lg:px-8 pt-8 pb-6 sm:pb-8 w-full">
        <div className="text-center py-12">
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex-1 max-w-7xl mx-auto px-6 lg:px-8 pt-8 pb-6 sm:pb-8 w-full">
        <div className="text-center py-12">
          <p className="text-muted-foreground">Error loading dashboard: {error}</p>
        </div>
      </div>
    );
  }

  // Calculate stats
  const watchTimeMinutes = (watchHistory?.entries ?? []).reduce((sum, entry) => sum + entry.total_minutes, 0) || 0;
  const watchTimeHours = Math.floor(watchTimeMinutes / 60);
  const racesWatchedCount = (watchHistory?.entries ?? []).length || 0;
  
  // Calculate streak from weekly data
  const streakDays = weeklyData?.current_streak_weeks ? weeklyData.current_streak_weeks * 7 : 0;
  const longestStreak = weeklyData?.best_streak_weeks ? weeklyData.best_streak_weeks * 7 : undefined;

  // Calculate impact (placeholder - would need revenue share API)
  const impactAmount = 0; // TODO: Calculate from revenue share data

  return (
    <div className="flex-1 max-w-7xl mx-auto px-6 lg:px-8 pt-8 pb-6 sm:pb-8 w-full">
      <WelcomeHeader 
        liveRacesCount={liveRacesCount}
        upcomingText={recommendations?.upcoming?.[0] ? `Your favorite series starts soon` : undefined}
      />
      
      <QuickStatsGrid
        watchTimeHours={watchTimeHours}
        watchTimeMinutes={Math.round(watchTimeMinutes % 60)}
        streakDays={streakDays}
        longestStreak={longestStreak}
        racesWatched={racesWatchedCount}
        racesWatchedLabel="Across multiple series"
        impactAmount={impactAmount}
      />

      {recommendations && (
        <>
          <ContinueWatchingSection races={recommendations.continue_watching || []} />
          <UpcomingFromFavorites races={recommendations.upcoming || []} />
        </>
      )}

      <GamificationSection
        missions={missions}
        weeklyData={weeklyData || undefined}
        level={xpData?.level}
      />

      <FairModelSection />
    </div>
  );
}

