'use client';

import { useEffect, useState, useMemo } from 'react';
import Link from 'next/link';
import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';
import { getLeaderboard, LeaderboardEntry } from '@/lib/api';
import ErrorMessage from '@/components/ErrorMessage';
import { Trophy, Clock, Award, User, Star, TrendingUp } from 'lucide-react';
import { Button } from '@/components/ui/button';

type SortBy = 'points' | 'time';

export default function LeaderboardPage() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<unknown>(null);
  const [sortBy, setSortBy] = useState<SortBy>('points');

  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        const data = await getLeaderboard();
        setEntries(data);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchLeaderboard();
  }, []);

  const formatWatchTime = (totalMinutes: number) => {
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;
    if (hours === 0) {
      return `${minutes} min`;
    }
    if (minutes === 0) {
      return `${hours}h`;
    }
    return `${hours}h ${minutes}m`;
  };

  const sortedEntries = useMemo(() => {
    const sorted = [...entries];
    if (sortBy === 'points') {
      sorted.sort((a, b) => {
        if (b.points !== a.points) {
          return b.points - a.points;
        }
        // If points are equal, sort by time watched
        return b.total_watch_minutes - a.total_watch_minutes;
      });
    } else {
      sorted.sort((a, b) => {
        if (b.total_watch_minutes !== a.total_watch_minutes) {
          return b.total_watch_minutes - a.total_watch_minutes;
        }
        // If time is equal, sort by points
        return b.points - a.points;
      });
    }
    return sorted;
  }, [entries, sortBy]);

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex flex-col">
        <Navigation variant="full" />
        <div className="flex-1 flex items-center justify-center">
          <div className="text-muted-foreground">Loading...</div>
        </div>
        <Footer />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex flex-col">
        <Navigation variant="full" />
        <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
          <ErrorMessage error={error} />
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8">
          {/* Header */}
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center">
                <Trophy className="w-6 h-6 text-primary" />
              </div>
              <div>
                <h1 className="text-2xl sm:text-3xl font-bold text-foreground/95">Leaderboard</h1>
                <p className="text-sm text-muted-foreground mt-1">
                  {entries.length} {entries.length === 1 ? 'user' : 'users'}
                </p>
              </div>
            </div>

            {/* Sort Toggle */}
            <div className="flex items-center gap-2">
              <Button
                variant={sortBy === 'points' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setSortBy('points')}
                className="flex items-center gap-2"
              >
                <Award className="w-4 h-4" />
                Points
              </Button>
              <Button
                variant={sortBy === 'time' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setSortBy('time')}
                className="flex items-center gap-2"
              >
                <Clock className="w-4 h-4" />
                Time Watched
              </Button>
            </div>
          </div>

          {/* Leaderboard Table */}
          {sortedEntries.length === 0 ? (
            <div className="text-center py-12 px-4">
              <p className="text-muted-foreground text-base sm:text-lg">No users found.</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-border/50">
                    <th className="text-left py-3 px-4 text-sm font-semibold text-muted-foreground">Rank</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-muted-foreground">User</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-muted-foreground">Points</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-muted-foreground">XP</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-muted-foreground">Level</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-muted-foreground">Time Watched</th>
                  </tr>
                </thead>
                <tbody>
                  {sortedEntries.map((entry, index) => {
                    const rank = index + 1;
                    const isTopThree = rank <= 3;
                    return (
                      <tr
                        key={entry.id}
                        className="border-b border-border/30 hover:bg-muted/30 transition-colors"
                      >
                        <td className="py-4 px-4">
                          <div className="flex items-center gap-2">
                            {isTopThree ? (
                              <Trophy
                                className={`w-5 h-5 ${
                                  rank === 1
                                    ? 'text-warning'
                                    : rank === 2
                                    ? 'text-muted-foreground'
                                    : 'text-warning'
                                }`}
                              />
                            ) : (
                              <span className="text-muted-foreground font-medium">#{rank}</span>
                            )}
                          </div>
                        </td>
                        <td className="py-4 px-4">
                          <Link
                            href={`/users/${entry.id}`}
                            className="flex items-center gap-3 hover:opacity-80 transition-opacity"
                          >
                            <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                              <User className="w-4 h-4 text-primary" />
                            </div>
                            <span className="font-medium text-foreground/95 hover:text-primary transition-colors">
                              {entry.name || 'Anonymous'}
                            </span>
                          </Link>
                        </td>
                        <td className="py-4 px-4 text-left">
                          <div className="flex items-center gap-1">
                            <Award className="w-4 h-4 text-muted-foreground" />
                            <span className="font-semibold text-foreground/95">
                              {entry.points.toLocaleString()}
                            </span>
                          </div>
                        </td>
                        <td className="py-4 px-4 text-left">
                          <div className="flex items-center gap-1">
                            <Star className="w-4 h-4 text-muted-foreground" />
                            <span className="font-semibold text-foreground/95">
                              {(entry.xp_total ?? 0).toLocaleString()}
                            </span>
                          </div>
                        </td>
                        <td className="py-4 px-4 text-left">
                          <div className="flex items-center gap-1">
                            <TrendingUp className="w-4 h-4 text-muted-foreground" />
                            <span className="font-semibold text-foreground/95">
                              Level {entry.level ?? 1}
                            </span>
                          </div>
                        </td>
                        <td className="py-4 px-4 text-left">
                          <div className="flex items-center gap-1">
                            <Clock className="w-4 h-4 text-muted-foreground" />
                            <span className="font-semibold text-foreground/95">
                              {formatWatchTime(entry.total_watch_minutes)}
                            </span>
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </main>
      <Footer />
    </div>
  );
}

