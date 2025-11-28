import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';
import { HomeDashboard } from '@/components/homepage/HomeDashboard';
import { CTASection } from '@/components/homepage/CTASection';
import { HomePageClient } from './HomePageClient';

export const metadata = {
  title: 'PelotonLive - The Future of Grassroots Racing',
  description: 'Watch live cycling races, support grassroots racing, and compete in missions',
};

export default function Home() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <HomePageClient />
      <Footer />
    </div>
  );
}
