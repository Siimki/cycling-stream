'use client';

import Link from 'next/link';
import { UserMissionWithDetails, WeeklyGoalProgress } from '@/lib/api';
import { Button } from '@/components/ui/button';

interface GamificationSectionProps {
  missions?: UserMissionWithDetails[];
  weeklyData?: WeeklyGoalProgress;
  level?: number;
}

export function GamificationSection({ missions, weeklyData, level }: GamificationSectionProps) {
  const activeMissions = missions?.filter(m => !m.completed_at) || [];
  const displayMissions = activeMissions.slice(0, 2);
  const currentStreak = weeklyData?.current_streak_weeks || 0;

  return (
    <section className="py-12 mb-12 bg-muted/20">
      <div className="grid md:grid-cols-2 gap-12 items-center">
        <div>
          <div className="inline-block bg-primary/20 text-primary text-sm font-bold px-3 py-1 rounded-full mb-4">
            LEVEL UP YOUR EXPERIENCE
          </div>
          <h2 className="text-4xl sm:text-5xl font-black mb-6 leading-tight">
            Don&apos;t Just Watch.<br />
            <span className="text-primary">Compete.</span>
          </h2>
          <p className="text-muted-foreground text-lg mb-8 leading-relaxed">
            Earn XP by watching live, chatting with fans, and predicting winners. Climb the leaderboard and unlock exclusive badges that showcase your dedication.
          </p>
          <Link href="/missions">
            <Button className="bg-card hover:bg-card/90 text-foreground font-bold py-3 px-8 rounded-full transition">
              View Missions
            </Button>
          </Link>
        </div>

        <div className="bg-card/80 backdrop-blur-sm border-2 border-primary/40 rounded-2xl p-8 hover:border-primary/50 hover:bg-card/90 hover:shadow-[0_4px_20px_rgba(0,0,0,0.3)] hover:shadow-primary/20 transition-all duration-200">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-bold">Your Missions</h3>
            {level !== undefined && (
              <div className="bg-primary text-primary-foreground text-xs font-bold px-3 py-1 rounded-full">
                Level {level}
              </div>
            )}
          </div>

          {/* Mission Progress */}
          {displayMissions.length > 0 ? (
            <div className="mb-6">
              {displayMissions.map((mission, index) => {
                const progressPercent = mission.target_value > 0
                  ? Math.min((mission.progress / mission.target_value) * 100, 100)
                  : 0;

                return (
                  <div key={mission.id || mission.mission_id || `mission-${index}`} className="mb-6">
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-sm text-foreground font-medium">
                        {mission.title}
                      </span>
                      <span className="text-primary text-sm font-bold">
                        +{mission.points_reward} XP
                      </span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2 mb-1 relative">
                      <div
                        className="bg-primary h-2 rounded-full relative transition-all"
                        style={{ width: `${progressPercent}%` }}
                      >
                        <div className="absolute right-0 top-1/2 -translate-y-1/2 w-2 h-2 bg-foreground rounded-full" />
                      </div>
                    </div>
                    <div className="flex justify-between text-[10px] text-muted-foreground font-mono">
                      <span>IN PROGRESS</span>
                      <span>
                        {Math.round(progressPercent)}% ({mission.progress}/{mission.target_value})
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground mb-6">No active missions at the moment.</p>
          )}

          {/* Streak Counter */}
          {weeklyData && (
            <div className="bg-muted rounded-xl p-4 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="text-3xl">🔥</span>
                <div>
                  <div className="text-sm font-medium text-foreground">Current Streak</div>
                  <div className="text-2xl font-bold text-primary">
                    {currentStreak} Week{currentStreak !== 1 ? 's' : ''}
                  </div>
                </div>
              </div>
              <div className="text-xs text-muted-foreground">
                Keep going!
              </div>
            </div>
          )}
        </div>
      </div>
    </section>
  );
}

