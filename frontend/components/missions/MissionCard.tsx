'use client';

import { UserMissionWithDetails } from '@/lib/api';
import { MissionProgress } from './MissionProgress';
import { Button } from '@/components/ui/button';
import { useState } from 'react';
import { claimMissionReward } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';

interface MissionCardProps {
  userMission: UserMissionWithDetails;
  onClaimed?: () => void | Promise<void>;
}

export function MissionCard({ userMission, onClaimed }: MissionCardProps) {
  const { refreshUser } = useAuth();
  const [isClaiming, setIsClaiming] = useState(false);
  const { progress, completed_at, claimed_at, mission_id, title, description, points_reward, xp_reward, target_value, mission_type } = userMission;

  const isCompleted = !!completed_at;
  const isClaimed = !!claimed_at;
  const canClaim = isCompleted && !isClaimed;

  // Determine action link based on mission type
  const getActionLink = () => {
    switch (mission_type) {
      case 'watch_time':
      case 'watch_race':
        return { href: '/races', label: 'Watch races' };
      case 'chat_message':
        return { href: '/races', label: 'Join live chat' };
      case 'predict_winner':
        return { href: '/races', label: 'View predictions' };
      default:
        return null;
    }
  };

  const actionLink = getActionLink();

  const handleClaim = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isClaiming || !canClaim) return;

    // Preserve scroll position before any async operations
    const scrollY = window.scrollY;

    setIsClaiming(true);
    try {
      await claimMissionReward(mission_id);
      
      // Reload missions first (silent refresh, no loading spinner)
      if (onClaimed) {
        await onClaimed();
      }
      
      // Refresh user to update points after missions are updated
      if (refreshUser) {
        await refreshUser();
      }
      
      // Restore scroll position after all state updates complete
      // Use double requestAnimationFrame to ensure DOM has fully updated
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          window.scrollTo(0, scrollY);
        });
      });
    } catch (error) {
      console.error('Failed to claim mission reward:', error);
      // TODO: Show error toast
    } finally {
      setIsClaiming(false);
    }
  };

  // Safety check
  if (!title) {
    return null;
  }

  return (
    <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
      <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto] gap-3 sm:gap-4">
        {/* Left: Title + Description */}
        <div className="flex-1 min-w-0">
          <h3 className="text-base sm:text-lg font-medium text-foreground/95 mb-2 truncate">
            {title.replace(/\s*\(Tier \d+\)/gi, '')}
          </h3>
          {description && (
            <p className="text-sm text-muted-foreground truncate">{description}</p>
          )}
        </div>

        {/* Right: Reward */}
        <div className="text-right sm:flex sm:items-start sm:justify-end">
          <div className="min-w-[80px]">
            <div className="text-lg sm:text-xl font-bold text-primary">
              +{points_reward}
            </div>
            <div className="text-xs text-muted-foreground">
              {xp_reward && xp_reward > 0 ? `points · +${xp_reward} XP` : 'points'}
            </div>
          </div>
        </div>
      </div>

      {/* Middle: Progress Bar */}
      <div className="mt-4">
        <MissionProgress progress={progress} target={target_value} />
      </div>

      {/* Status/Claim Button */}
      <div className="flex items-center justify-between mt-4">
        <div className="flex items-center gap-3">
          {isClaimed ? (
            <span className="text-xs text-muted-foreground font-medium">Reward claimed</span>
          ) : canClaim ? (
            <Button
              type="button"
              onClick={handleClaim}
              disabled={isClaiming}
              className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold"
              size="sm"
            >
              {isClaiming ? 'Claiming...' : `Claim +${points_reward}`}
            </Button>
          ) : isCompleted ? (
            <span className="text-sm text-primary font-medium">Completed!</span>
          ) : (
            <>
              <span className="text-sm text-muted-foreground font-medium">In progress</span>
              {actionLink && !isCompleted && (
                <a
                  href={actionLink.href}
                  className="text-xs text-primary hover:underline font-medium"
                >
                  {actionLink.label} →
                </a>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}

