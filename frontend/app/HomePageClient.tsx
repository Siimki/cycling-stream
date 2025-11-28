'use client';

import { useAuth } from '@/contexts/AuthContext';
import { HomeDashboard } from '@/components/homepage/HomeDashboard';
import { CTASection } from '@/components/homepage/CTASection';
import { FairModelSection } from '@/components/homepage/FairModelSection';
import { LiveRacesSection } from '@/components/homepage/LiveRacesSection';
import { getRaces, getRaceStream, Race, StreamResponse } from '@/lib/api';
import { useEffect, useState } from 'react';

export function HomePageClient() {
  const { isAuthenticated, isLoading } = useAuth();
  const [liveRacesData, setLiveRacesData] = useState<Array<{
    race: Race;
    stream?: StreamResponse;
  }>>([]);

  useEffect(() => {
    async function loadLiveRaces() {
      try {
        const races = await getRaces();
        if (!races || races.length === 0) return;

        // Check stream status for each race (limit to first 10 for performance)
        const livePromises = races.slice(0, 10).map(async (race) => {
          try {
            const stream = await getRaceStream(race.id, true);
            return { race, stream };
          } catch {
            return { race, stream: undefined };
          }
        });

        const results = await Promise.allSettled(livePromises);
        const liveData = results
          .filter((r): r is PromiseFulfilledResult<{ race: Race; stream?: StreamResponse }> => 
            r.status === 'fulfilled'
          )
          .map(r => r.value)
          .filter(item => item.stream?.status === 'live' || item.stream?.status === 'planned');
        
        setLiveRacesData(liveData);
      } catch (err) {
        console.error('Failed to load live races:', err);
      }
    }

    loadLiveRaces();
  }, []);

  if (isLoading) {
    return (
      <main className="flex-1 max-w-7xl mx-auto px-6 lg:px-8 pt-8 pb-6 sm:pb-8 w-full">
        <div className="text-center py-12">
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </main>
    );
  }

  return (
    <main className="flex-1">
      {isAuthenticated ? (
        <>
          <HomeDashboard liveRacesCount={liveRacesData.filter(item => item.stream?.status === 'live').length} />
          <div className="max-w-7xl mx-auto px-6 lg:px-8">
            <LiveRacesSection races={liveRacesData.map(item => ({
              race: item.race,
              stream: item.stream,
            }))} />
          </div>
        </>
      ) : (
        <>
          {/* Non-authenticated view */}
          <section className="py-24 mb-12 bg-gradient-to-b from-background to-muted/20 relative overflow-hidden">
            <div className="absolute inset-0 opacity-10">
              <div className="absolute top-0 left-1/4 w-96 h-96 bg-primary rounded-full blur-3xl" />
              <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-primary rounded-full blur-3xl" />
            </div>

            <div className="max-w-4xl mx-auto px-6 text-center relative z-10 pt-24">
              <h1 className="text-5xl sm:text-6xl font-black mb-6">PelotonLive</h1>
              <p className="text-2xl sm:text-3xl font-bold text-foreground mb-4">
                The Future of Grassroots Racing
              </p>
              <p className="text-muted-foreground text-xl mb-10 max-w-2xl mx-auto">
                Watch live cycling races, support grassroots racing, and compete with fans around the world.
              </p>
              <CTASection />
            </div>
          </section>

          <div className="max-w-7xl mx-auto px-6 lg:px-8">
            <LiveRacesSection races={liveRacesData.map(item => ({
              race: item.race,
              stream: item.stream,
            }))} />
          </div>

          <FairModelSection />
          <CTASection />
        </>
      )}
    </main>
  );
}

