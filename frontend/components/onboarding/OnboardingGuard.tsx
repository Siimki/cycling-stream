'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { getUserPreferences } from '@/lib/api';

export function OnboardingGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, isLoading } = useAuth();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    const checkOnboarding = async () => {
      // Skip check if not authenticated or already on onboarding page
      if (!isAuthenticated || pathname === '/onboarding' || pathname.startsWith('/auth/')) {
        setChecking(false);
        return;
      }

      try {
        const prefs = await getUserPreferences();
        if (!prefs.onboarding_completed) {
          router.push('/onboarding');
          return;
        }
      } catch (error: any) {
        // If we get a 404, preferences don't exist yet - need onboarding
        // If we get 401, user isn't authenticated - skip check
        if (error?.status === 404) {
          router.push('/onboarding');
          return;
        }
        // For other errors (like 401), don't redirect - let auth handle it
        if (error?.status !== 401) {
          console.error('Error checking onboarding status:', error);
        }
      } finally {
        setChecking(false);
      }
    };

    if (!isLoading) {
      checkOnboarding();
    }
  }, [isAuthenticated, isLoading, pathname, router]);

  // Show nothing while checking
  if (checking || isLoading) {
    return null;
  }

  return <>{children}</>;
}

