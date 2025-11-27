'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { getToken, changePassword } from '@/lib/auth';
import { Navigation } from '@/components/layout/Navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import Footer from '@/components/layout/Footer';

export default function SettingsPage() {
  const router = useRouter();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);

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
                <div className="mb-4 p-3 bg-green-500/10 border border-green-500/20 rounded-md text-green-600 dark:text-green-400 text-sm">
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
      </main>
      <Footer />
    </div>
  );
}

