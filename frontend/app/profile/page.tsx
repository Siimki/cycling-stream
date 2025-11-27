'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { Navigation } from '@/components/layout/Navigation';
import { Button } from '@/components/ui/button';
import Footer from '@/components/layout/Footer';

export default function ProfilePage() {
  const router = useRouter();
  const { user, isLoading: loading, logout } = useAuth();

  useEffect(() => {
    if (!loading && !user) {
      router.push('/auth/login');
    }
  }, [user, loading, router]);

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex flex-col">
        <Navigation variant="full" />
        <div className="flex-1 flex items-center justify-center">
          <div className="text-muted-foreground">Loading...</div>
        </div>
        <Footer />
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8">
          <h1 className="text-2xl font-bold text-foreground/95 mb-6">Profile</h1>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-muted-foreground mb-1">Email</label>
              <p className="text-foreground/95">{user.email}</p>
            </div>
            {user.name && (
              <div>
                <label className="block text-sm font-medium text-muted-foreground mb-1">Name</label>
                <p className="text-foreground/95">{user.name}</p>
              </div>
            )}
            <div>
              <label className="block text-sm font-medium text-muted-foreground mb-1">Bio</label>
              <p className="text-foreground/95">{user.bio || <span className="text-muted-foreground italic">No bio set</span>}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-muted-foreground mb-1">Points</label>
              <p className="text-foreground/95">{(user.points ?? 0).toLocaleString()}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-muted-foreground mb-1">Member since</label>
              <p className="text-foreground/95">
                {new Date(user.created_at).toLocaleDateString()}
              </p>
            </div>
          </div>

          <div className="mt-8 pt-6 border-t border-border/50">
            <Button
              onClick={handleLogout}
              variant="destructive"
              className="px-4 py-2"
            >
              Logout
            </Button>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}

