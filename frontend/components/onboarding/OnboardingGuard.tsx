'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { useExperience } from '@/contexts/ExperienceContext';

export function OnboardingGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, isLoading } = useAuth();
  const { loading: preferencesLoading, preferences, refreshPreferences } = useExperience();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    const shouldSkipRoute = pathname === '/onboarding' || pathname.startsWith('/auth/');
    if (isLoading || preferencesLoading) {
      return;
    }

    if (!isAuthenticated || shouldSkipRoute) {
      setChecking(false);
      return;
    }

    if (!preferences) {
      refreshPreferences();
      return;
    }

    if (!preferences.onboarding_completed) {
      router.push('/onboarding');
      return;
    }

    setChecking(false);
  }, [isAuthenticated, isLoading, pathname, preferences, preferencesLoading, refreshPreferences, router]);

  // Show nothing while checking
  if (checking || isLoading) {
    return null;
  }

  return <>{children}</>;
}

