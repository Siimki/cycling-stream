import { OnboardingWizard } from '@/components/onboarding/OnboardingWizard';
import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';

export default function OnboardingPage() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="minimal" />
      <OnboardingWizard />
      <Footer />
    </div>
  );
}

