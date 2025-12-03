'use client';

import { UserMissionWithDetails, getUserMissions, MissionType } from '@/lib/api';
import { MissionCard } from './MissionCard';
import { useEffect, useState, useRef, useMemo } from 'react';
import LoadingSpinner from '@/components/LoadingSpinner';
import ErrorMessage from '@/components/ErrorMessage';
import Link from 'next/link';

interface MissionsPanelProps {
  limit?: number;
  showViewAll?: boolean;
}

export function MissionsPanel({ limit, showViewAll = true }: MissionsPanelProps) {
  const [missions, setMissions] = useState<UserMissionWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const isInitialLoad = useRef(true);

  const loadMissions = async (silent = false) => {
    try {
      if (!silent) {
        setLoading(true);
      }
      setError(null);
      const data = await getUserMissions();
      setMissions(data);
    } catch (err) {
      // Check if it's an auth error (401)
      if (err && typeof err === 'object' && 'status' in err && (err as { status?: number }).status === 401) {
        setError('authentication_required');
      } else {
        setError(err instanceof Error ? err.message : 'Failed to load missions');
      }
    } finally {
      if (!silent) {
        setLoading(false);
      }
      isInitialLoad.current = false;
    }
  };

  useEffect(() => {
    loadMissions();
  }, []);

  // Group missions by type
  const groupedMissions = useMemo(() => {
    const groups: Record<string, { title: string; missions: UserMissionWithDetails[] }> = {
      weekly: { title: 'Weekly Missions', missions: [] },
      chat: { title: 'Chat Missions', missions: [] },
      prediction: { title: 'Prediction Missions', missions: [] },
      special: { title: 'Special / Limited-time', missions: [] },
    };

    missions.forEach((mission) => {
      if (mission.mission_type === 'watch_time') {
        groups.weekly.missions.push(mission);
      } else if (mission.mission_type === 'chat_message') {
        groups.chat.missions.push(mission);
      } else if (mission.mission_type === 'predict_winner') {
        groups.prediction.missions.push(mission);
      } else {
        // watch_race, follow_series, streak
        groups.special.missions.push(mission);
      }
    });

    // Filter out empty groups
    return Object.entries(groups)
      .filter(([_, group]) => group.missions.length > 0)
      .map(([_, group]) => group);
  }, [missions]);

  const displayMissions = limit && missions ? missions.slice(0, limit) : (missions || []);
  const hasMore = limit && missions && missions.length > limit;

  if (loading) {
    return (
      <div className="flex justify-center py-8">
        <LoadingSpinner />
      </div>
    );
  }

  if (error === 'authentication_required') {
    return (
      <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-8 text-center">
        <h3 className="text-xl font-semibold text-foreground/95 mb-3">Log in to start missions</h3>
        <p className="text-muted-foreground mb-6">
          Earn XP, build your streak, and unlock rewards by watching races.
        </p>
        <Link
          href="/auth/login"
          className="inline-block bg-primary hover:bg-primary/90 text-primary-foreground font-semibold py-2 px-6 rounded-lg transition"
        >
          Log in
        </Link>
      </div>
    );
  }

  if (error) {
    return <ErrorMessage message={error} variant="inline" />;
  }

  if (missions.length === 0) {
    return (
      <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-8 text-center">
        <h3 className="text-lg font-semibold text-foreground/95 mb-2">No missions yet</h3>
        <p className="text-muted-foreground mb-4">
          We're not running any active missions right now. Watch live races and check back before big events.
        </p>
        <Link
          href="/races"
          className="inline-block bg-card hover:bg-card/90 border border-border text-foreground font-medium py-2 px-6 rounded-lg transition"
        >
          Browse races
        </Link>
      </div>
    );
  }

  // If limit is set, show ungrouped list (for previews)
  if (limit) {
    return (
      <div>
        <div className="space-y-4">
          {displayMissions.map((userMission) => (
            <MissionCard
              key={userMission.mission_id}
              userMission={userMission}
              onClaimed={() => loadMissions(true)}
            />
          ))}
        </div>
        {hasMore && showViewAll && (
          <div className="mt-4 text-center">
            <Link
              href="/missions"
              className="text-primary hover:underline text-sm font-medium"
            >
              View all missions →
            </Link>
          </div>
        )}
      </div>
    );
  }

  // Full page: show grouped missions
  return (
    <div className="space-y-10">
      {groupedMissions.map((group) => (
        <div key={group.title} className="space-y-4">
          <h2 className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-4">
            {group.title}
          </h2>
          <div className="space-y-4">
            {group.missions.map((userMission) => (
              <MissionCard
                key={userMission.mission_id}
                userMission={userMission}
                onClaimed={() => loadMissions(true)}
              />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

