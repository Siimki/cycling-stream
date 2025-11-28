import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';
import { MissionsPanel } from '@/components/missions/MissionsPanel';
import UserStats from '@/components/user/UserStats';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Missions | CyclingStream',
  description: 'Complete missions to earn points and unlock rewards',
};

export default function MissionsPage() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
        <div className="mb-6 sm:mb-8">
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground/95 mb-2">
            Missions
          </h1>
          <p className="text-muted-foreground text-base sm:text-lg">
            Complete challenges to earn points and unlock rewards
          </p>
        </div>

        {/* Weekly Goal Section */}
        <div className="mb-6 sm:mb-8">
          <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mb-4">
            This Week
          </h2>
          <UserStats compact={false} />
        </div>

        <MissionsPanel />
      </main>
      <Footer />
    </div>
  );
}

