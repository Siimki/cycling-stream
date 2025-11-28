import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';
import { MissionsPanel } from '@/components/missions/MissionsPanel';
import { WeeklyOverview } from '@/components/missions/WeeklyOverview';
import { StreakCard } from '@/components/missions/StreakCard';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Missions | CyclingStream',
  description: 'Complete missions to earn points and unlock rewards',
};

export default function MissionsPage() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 pt-12 sm:pt-16 pb-6 sm:pb-8 w-full">
        <div className="mb-10 sm:mb-12">
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground/95 mb-2">
            Missions
          </h1>
          <p className="text-muted-foreground/70 text-base sm:text-lg">
            Complete challenges to earn points and unlock rewards.
          </p>
        </div>

        {/* Weekly Overview */}
        <div className="mb-8">
          <WeeklyOverview />
        </div>

        {/* Streak */}
        <div className="mb-10">
          <StreakCard />
        </div>

        <MissionsPanel />
      </main>
      <Footer />
    </div>
  );
}

