'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { login, setToken } from '@/lib/auth';
import { useAuth } from '@/contexts/AuthContext';
import { Navigation } from '@/components/layout/Navigation';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import Footer from '@/components/layout/Footer';

export default function LoginPage() {
  const router = useRouter();
  const { login: authLogin } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const data = await login(email, password);
      setToken(data.token);
      authLogin(data.token, data.user);
      router.push('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="minimal" />
      <div className="flex-1 flex items-center justify-center px-4 py-12">
        <div className="max-w-md w-full bg-[oklch(0.20_0.005_260)] border-2 border-[oklch(0.30_0.005_260)] shadow-[0_8px_32px_rgba(0,0,0,0.6)] rounded-lg p-6 sm:p-8">
          <h1 className="text-2xl font-bold text-foreground/95 mb-6">Login</h1>
          <form onSubmit={handleSubmit}>
            {error && (
              <div className="mb-4 p-3 bg-destructive/10 border border-destructive/20 rounded-md text-destructive text-sm">
                {error}
              </div>
            )}
            <div className="mb-4">
              <label className="block text-sm font-medium text-foreground mb-2">
                Email
              </label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="bg-background"
              />
            </div>
            <div className="mb-6">
              <label className="block text-sm font-medium text-foreground mb-2">
                Password
              </label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="bg-background"
              />
            </div>
            <Button
              type="submit"
              disabled={loading}
              className="w-full bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground"
            >
              {loading ? 'Logging in...' : 'Login'}
            </Button>
          </form>
          <p className="mt-4 text-center text-sm text-muted-foreground">
            Don&apos;t have an account?{' '}
            <Link href="/auth/register" className="text-primary hover:underline font-medium">
              Register
            </Link>
          </p>
        </div>
      </div>
      <Footer />
    </div>
  );
}

