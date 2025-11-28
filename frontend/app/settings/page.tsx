'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { getToken, changePassword } from '@/lib/auth';
import { Navigation } from '@/components/layout/Navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import Footer from '@/components/layout/Footer';
import { ToggleSwitch } from '@/components/ui/toggle-switch';
import { Slider } from '@/components/ui/slider';
import { useExperience } from '@/contexts/ExperienceContext';

export default function SettingsPage() {
  const router = useRouter();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const { uiPreferences, audioPreferences, updateUIPreferences, updateAudioPreferences } = useExperience();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    // Validate passwords match
    if (newPassword !== confirmPassword) {
      setError('New passwords do not match');
      setLoading(false);
      return;
    }

    const token = getToken();
    if (!token) {
      router.push('/auth/login');
      return;
    }

    try {
      await changePassword(token, currentPassword, newPassword);
      setSuccess('Password changed successfully!');
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to change password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8">
          <h1 className="text-2xl font-bold text-foreground/95 mb-6">Settings</h1>

          {/* Change Password Section */}
          <div className="space-y-6">
            <div>
              <h2 className="text-lg font-semibold text-foreground/95 mb-4">Change Password</h2>
              
              {error && (
                <div className="mb-4 p-3 bg-destructive/10 border border-destructive/20 rounded-md text-destructive text-sm">
                  {error}
                </div>
              )}

              {success && (
                <div className="mb-4 p-3 bg-success/10 border border-success/20 rounded-md text-success text-sm">
                  {success}
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-foreground mb-2">
                    Current Password
                  </label>
                  <Input
                    type="password"
                    value={currentPassword}
                    onChange={(e) => setCurrentPassword(e.target.value)}
                    required
                    className="bg-background"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-2">
                    New Password
                  </label>
                  <Input
                    type="password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    required
                    minLength={8}
                    className="bg-background"
                  />
                  <p className="mt-1 text-xs text-muted-foreground">
                    Must be at least 8 characters with uppercase, lowercase, number, and special character
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-2">
                    Confirm New Password
                  </label>
                  <Input
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                    minLength={8}
                    className="bg-background"
                  />
                </div>

                <Button
                  type="submit"
                  disabled={loading}
                  className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground"
                >
                  {loading ? 'Changing password...' : 'Change Password'}
                </Button>
              </form>
            </div>
          </div>
        </div>

        <div className="mt-8 grid gap-6">
          <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8 space-y-4">
            <h2 className="text-lg font-semibold text-foreground/95">Motion & Animations</h2>
            <ToggleSwitch
              label="Chat animations"
              description="Controls message entry and emote motion."
              checked={uiPreferences?.chat_animations ?? true}
              onCheckedChange={(value) => updateUIPreferences({ chat_animations: value })}
            />
            <ToggleSwitch
              label="Button pulse"
              description="Hover expansion for interactive elements."
              checked={uiPreferences?.button_pulse ?? true}
              onCheckedChange={(value) => updateUIPreferences({ button_pulse: value })}
            />
            <ToggleSwitch
              label="Poll animations"
              description="Slide/fade for poll announcements."
              checked={uiPreferences?.poll_animations ?? true}
              onCheckedChange={(value) => updateUIPreferences({ poll_animations: value })}
            />
            <ToggleSwitch
              label="Reduced motion"
              description="Minimize non-essential motion."
              checked={uiPreferences?.reduced_motion ?? false}
              onCheckedChange={(value) => updateUIPreferences({ reduced_motion: value })}
            />
          </div>

          <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8 space-y-4">
            <h2 className="text-lg font-semibold text-foreground/95">Sound Feedback</h2>
            <ToggleSwitch
              label="Button clicks"
              description="Play a soft click on press."
              checked={audioPreferences?.button_clicks ?? true}
              onCheckedChange={(value) => updateAudioPreferences({ button_clicks: value })}
            />
            <ToggleSwitch
              label="Notification sounds"
              description="Polls, events, and rewards."
              checked={audioPreferences?.notification_sounds ?? true}
              onCheckedChange={(value) => updateAudioPreferences({ notification_sounds: value })}
            />
            <ToggleSwitch
              label="Mention ping"
              description="Quiet alert when someone tags you."
              checked={audioPreferences?.mention_pings ?? true}
              onCheckedChange={(value) => updateAudioPreferences({ mention_pings: value })}
            />
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="font-medium">Master volume</span>
                <span className="text-muted-foreground">
                  {Math.round((audioPreferences?.master_volume ?? 0.15) * 100)}%
                </span>
              </div>
              <Slider
                value={[Math.round((audioPreferences?.master_volume ?? 0.15) * 100)]}
                onValueChange={([value]) => updateAudioPreferences({ master_volume: value / 100 })}
                max={100}
                step={5}
              />
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}

