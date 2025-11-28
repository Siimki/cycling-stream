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
      setError(err instanceof Error ? err.message : 'Failed to load missions');
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

  if (error) {
    return <ErrorMessage message={error} variant="inline" />;
  }

  if (missions.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-muted-foreground">No active missions at the moment.</p>
        <p className="text-muted-foreground/70 text-sm mt-1">Check back soon for new challenges!</p>
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

