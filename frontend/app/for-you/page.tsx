'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getRecommendations, type RecommendationsResponse } from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';
import { ContinueWatching } from '@/components/for-you/ContinueWatching';
import { UpcomingRaces } from '@/components/for-you/UpcomingRaces';
import { RecommendedReplays } from '@/components/for-you/RecommendedReplays';
import { PersonalStats } from '@/components/for-you/PersonalStats';
import ErrorMessage from '@/components/ErrorMessage';
import { useAuth } from '@/contexts/AuthContext';
import Link from 'next/link';
import { Button } from '@/components/ui/button';

export default function ForYouPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();
  const [recommendations, setRecommendations] = useState<RecommendationsResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      setLoading(false);
      return;
    }

    if (isAuthenticated) {
      getRecommendations()
        .then((data) => {
          setRecommendations(data);
          setError(null);
        })
        .catch((err) => {
          const errorMsg = APIErrorHandler.getErrorMessage(err);
          setError(errorMsg);
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [isAuthenticated, isLoading]);

  if (isLoading || loading) {
    return (
      <div className="min-h-screen bg-background flex flex-col">
        <Navigation variant="full" />
        <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading...</p>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-background flex flex-col">
        <Navigation variant="full" />
        <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
          <div className="text-center py-12">
            <h1 className="text-3xl font-bold text-foreground mb-4">For You</h1>
            <p className="text-muted-foreground mb-6">
              Sign in to see personalized race recommendations, continue watching, and your stats.
            </p>
            <div className="flex gap-4 justify-center">
              <Link href="/auth/login">
                <Button>Log In</Button>
              </Link>
              <Link href="/auth/register">
                <Button variant="outline">Sign Up</Button>
              </Link>
            </div>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
        <h1 className="text-3xl font-bold text-foreground mb-8">For You</h1>
        
        {error ? (
          <ErrorMessage message={error} />
        ) : recommendations ? (
          <>
            <PersonalStats />
            <ContinueWatching races={recommendations.continue_watching || []} />
            <UpcomingRaces races={recommendations.upcoming || []} />
            <RecommendedReplays races={recommendations.replays || []} />
          </>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading recommendations...</p>
          </div>
        )}
      </main>
      <Footer />
    </div>
  );
}

