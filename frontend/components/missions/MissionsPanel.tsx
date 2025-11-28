'use client';

import { UserMissionWithDetails, getUserMissions } from '@/lib/api';
import { MissionCard } from './MissionCard';
import { useEffect, useState, useRef } from 'react';
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

