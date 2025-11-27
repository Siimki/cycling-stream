'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';
import { getPublicUser, PublicUser } from '@/lib/api';
import { Calendar, Clock, Award, User } from 'lucide-react';

export default function PublicProfilePage() {
  const params = useParams();
  const userId = params.id as string;
  const [user, setUser] = useState<PublicUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!userId) return;

    const fetchUser = async () => {
      try {
        const userData = await getPublicUser(userId);
        setUser(userData);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load user profile');
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [userId]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

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
        <div className="flex-1 p-8">
          <div className="max-w-2xl mx-auto">
            <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4 text-destructive">
              {error}
            </div>
          </div>
        </div>
        <Footer />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-background flex flex-col">
        <Navigation variant="full" />
        <div className="flex-1 p-8">
          <div className="max-w-2xl mx-auto text-center">
            <h1 className="text-2xl font-bold text-foreground/95 mb-4">User Not Found</h1>
            <p className="text-muted-foreground">The user you&apos;re looking for doesn&apos;t exist.</p>
          </div>
        </div>
        <Footer />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8">
          {/* User Header */}
          <div className="flex items-center gap-4 mb-6">
            <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
              <User className="w-8 h-8 text-primary" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-foreground/95">
                {user.name || 'Anonymous User'}
              </h1>
              {user.bio && (
                <p className="text-muted-foreground mt-1 text-sm">{user.bio}</p>
              )}
            </div>
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            {/* Member Since */}
            <div className="bg-muted/30 rounded-lg p-4">
              <div className="flex items-center gap-2 text-muted-foreground mb-1">
                <Calendar className="w-4 h-4" />
                <span className="text-xs font-medium uppercase tracking-wide">Member Since</span>
              </div>
              <p className="text-foreground/95 font-semibold">
                {formatDate(user.created_at)}
              </p>
            </div>

            {/* Watch Time */}
            <div className="bg-muted/30 rounded-lg p-4">
              <div className="flex items-center gap-2 text-muted-foreground mb-1">
                <Clock className="w-4 h-4" />
                <span className="text-xs font-medium uppercase tracking-wide">Watch Time</span>
              </div>
              <p className="text-foreground/95 font-semibold">
                {formatWatchTime(user.total_watch_minutes)}
              </p>
            </div>

            {/* Points */}
            <div className="bg-muted/30 rounded-lg p-4">
              <div className="flex items-center gap-2 text-muted-foreground mb-1">
                <Award className="w-4 h-4" />
                <span className="text-xs font-medium uppercase tracking-wide">Points</span>
              </div>
              <p className="text-foreground/95 font-semibold">
                {user.points.toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}

