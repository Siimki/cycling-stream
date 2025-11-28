'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { StepCyclingLevel } from './StepCyclingLevel';
import { StepViewPreference } from './StepViewPreference';
import { StepFavorites } from './StepFavorites';
import { StepDevice } from './StepDevice';
import { StepNotifications } from './StepNotifications';
import {
  updateUserPreferences,
  completeOnboarding,
  addFavorite,
  type UpdatePreferencesRequest,
  type AddFavoriteRequest,
  type UpdateUIPreferencesRequest,
  type UpdateAudioPreferencesRequest,
} from '@/lib/api';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { useExperience } from '@/contexts/ExperienceContext';

export interface OnboardingData {
  cyclingLevel: 'new' | 'casual' | 'superfan' | null;
  viewPreference: 'clean' | 'data-rich' | null;
  favorites: AddFavoriteRequest[];
  device: ('tv' | 'laptop' | 'phone' | 'tablet')[];
  notifications: {
    challenges: boolean;
    points: boolean;
    races: boolean;
  };
}

const STORAGE_KEY = 'onboarding_data';

export function OnboardingWizard() {
  const router = useRouter();
  const { refreshPreferences } = useExperience();
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<OnboardingData>(() => {
    // Load from localStorage if available
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (saved) {
        try {
          const parsed = JSON.parse(saved);
          // Ensure device is always an array (migrate from old null/single value format)
          if (!Array.isArray(parsed.device)) {
            parsed.device = parsed.device ? [parsed.device] : [];
          }
          return parsed;
        } catch {
          // Invalid data, use defaults
        }
      }
    }
    return {
      cyclingLevel: null,
      viewPreference: null,
      favorites: [],
      device: [],
      notifications: {
        challenges: true,
        points: true,
        races: false,
      },
    };
  });

  const totalSteps = 5;

  // Save to localStorage whenever data changes
  useEffect(() => {
    if (typeof window !== 'undefined') {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
    }
  }, [data]);

  const updateData = (updates: Partial<OnboardingData>) => {
    setData((prev) => ({ ...prev, ...updates }));
  };

  const handleNext = () => {
    if (currentStep < totalSteps - 1) {
      setCurrentStep(currentStep + 1);
    }
  };

  const handleBack = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleSkip = () => {
    // Skip to next step or complete
    if (currentStep < totalSteps - 1) {
      setCurrentStep(currentStep + 1);
    } else {
      handleComplete();
    }
  };

  const handleComplete = async () => {
    setLoading(true);
    
    // Check if token exists
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
    if (!token) {
      alert('You must be logged in to complete onboarding. Please log in and try again.');
      router.push('/auth/login');
      setLoading(false);
      return;
    }
    
    try {
      // Map onboarding data to preferences
      // Map device array to primary device (prefer desktop/laptop, then mobile, then tv)
      let primaryDevice: 'tv' | 'desktop' | 'mobile' | 'tablet' | undefined;
      if (data.device.length > 0) {
        if (data.device.includes('laptop')) {
          primaryDevice = 'desktop';
        } else if (data.device.includes('phone')) {
          primaryDevice = 'mobile';
        } else if (data.device.includes('tablet')) {
          primaryDevice = 'tablet';
        } else if (data.device.includes('tv')) {
          primaryDevice = 'tv';
        }
      }

      const uiPreferences: UpdateUIPreferencesRequest = {
        chat_animations: data.viewPreference !== 'clean',
        reduced_motion: data.viewPreference === 'clean',
        button_pulse: true,
        poll_animations: data.viewPreference !== 'clean',
      };

      const audioPreferences: UpdateAudioPreferencesRequest = {
        button_clicks: true,
        notification_sounds: data.notifications.races || data.notifications.points,
        mention_pings: true,
        master_volume: 0.2,
      };

      const prefs: UpdatePreferencesRequest = {
        data_mode: data.viewPreference === 'clean' ? 'casual' : data.viewPreference === 'data-rich' ? 'pro' : 'standard',
        device_type: primaryDevice,
        notification_preferences: {
          challenges: data.notifications.challenges,
          points: data.notifications.points,
          races: data.notifications.races,
        },
        onboarding_completed: true,
        ui_preferences: uiPreferences,
        audio_preferences: audioPreferences,
      };

      // Save preferences
      await updateUserPreferences(prefs);

      // Add favorites
      for (const favorite of data.favorites) {
        try {
          await addFavorite(favorite);
        } catch (err) {
          // Continue even if some favorites fail
          console.error('Failed to add favorite:', err);
        }
      }

      // Mark onboarding as complete
      await completeOnboarding();
      await refreshPreferences();

      // Clear localStorage
      if (typeof window !== 'undefined') {
        localStorage.removeItem(STORAGE_KEY);
      }

      // Redirect to home
      router.push('/');
    } catch (error) {
      console.error('Failed to complete onboarding:', error);
      const status = typeof error === 'object' && error !== null && 'status' in error
        ? (error as { status?: number }).status
        : undefined;
      const message = error instanceof Error ? error.message : 'Unknown error';
      
      // Provide more specific error messages
      if (status === 401 || status === 403) {
        alert('Your session has expired. Please log in again and complete onboarding.');
        router.push('/auth/login');
      } else if (status === 404) {
        alert('The preferences endpoint is not available. Please refresh the page and try again.');
      } else {
        alert(`Failed to save preferences: ${message}. Please try again.`);
      }
    } finally {
      setLoading(false);
    }
  };

  const canProceed = () => {
    switch (currentStep) {
      case 0:
        return data.cyclingLevel !== null;
      case 1:
        return data.viewPreference !== null;
      case 2:
        return true; // Favorites are optional
      case 3:
        return data.device.length > 0;
      case 4:
        return true; // Notifications have defaults
      default:
        return false;
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <div className="flex-1 flex items-center justify-center px-4 py-12">
        <div className="max-w-2xl w-full bg-card/95 backdrop-blur-xl border border-border/50 rounded-lg p-6 sm:p-8">
          {/* Progress bar */}
          <div className="mb-8">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-muted-foreground">
                Step {currentStep + 1} of {totalSteps}
              </span>
              <span className="text-sm text-muted-foreground">
                {Math.round(((currentStep + 1) / totalSteps) * 100)}%
              </span>
            </div>
            <div className="w-full h-2 bg-muted rounded-full overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-primary to-primary/80 transition-all duration-300"
                style={{ width: `${((currentStep + 1) / totalSteps) * 100}%` }}
              />
            </div>
          </div>

          {/* Step content */}
          <div className="mb-8 min-h-[400px]">
            {currentStep === 0 && (
              <StepCyclingLevel
                value={data.cyclingLevel}
                onChange={(value) => updateData({ cyclingLevel: value })}
              />
            )}
            {currentStep === 1 && (
              <StepViewPreference
                value={data.viewPreference}
                onChange={(value) => updateData({ viewPreference: value })}
              />
            )}
            {currentStep === 2 && (
              <StepFavorites
                favorites={data.favorites}
                onChange={(favorites) => updateData({ favorites })}
              />
            )}
            {currentStep === 3 && (
              <StepDevice
                value={data.device}
                onChange={(value) => updateData({ device: value })}
              />
            )}
            {currentStep === 4 && (
              <StepNotifications
                value={data.notifications}
                onChange={(value) => updateData({ notifications: value })}
              />
            )}
          </div>

          {/* Navigation buttons */}
          <div className="flex justify-between items-center gap-4">
            <div>
              {currentStep > 0 && (
                <Button
                  variant="outline"
                  onClick={handleBack}
                  disabled={loading}
                  className="flex items-center gap-2"
                >
                  <ChevronLeft className="w-4 h-4" />
                  Back
                </Button>
              )}
            </div>
            <div className="flex gap-2">
              {currentStep < totalSteps - 1 ? (
                <>
                  <Button
                    variant="ghost"
                    onClick={handleSkip}
                    disabled={loading}
                  >
                    Skip
                  </Button>
                  <Button
                    onClick={handleNext}
                    disabled={!canProceed() || loading}
                    className="flex items-center gap-2 bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70"
                  >
                    Next
                    <ChevronRight className="w-4 h-4" />
                  </Button>
                </>
              ) : (
                <Button
                  onClick={handleComplete}
                  disabled={loading}
                  className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70"
                >
                  {loading ? 'Saving...' : 'Complete Setup'}
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

