'use client';

import { createContext, useContext, useEffect, useState, useCallback, useMemo, type ReactNode } from 'react';
import { getToken, removeToken, getProfile, User } from '@/lib/auth';
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('Auth');

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  token: string | null;
  isLoading: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    const currentToken = getToken();
    if (!currentToken) {
      setUser(null);
      setToken(null);
      setIsLoading(false);
      return;
    }

    try {
      const userData = await getProfile(currentToken);
      setUser(userData);
      setToken(currentToken);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      if (errorMessage.includes('not found') || errorMessage.includes('Unauthorized') || errorMessage.includes('401')) {
        // Token is invalid, clear it
        removeToken();
        setUser(null);
        setToken(null);
      } else {
        logger.error('Failed to fetch user profile:', err);
      }
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Initial load
  useEffect(() => {
    refreshUser();
  }, [refreshUser]);

  // Listen for storage changes (login/logout from other tabs)
  useEffect(() => {
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'auth_token') {
        refreshUser();
      }
    };

    window.addEventListener('storage', handleStorageChange);
    return () => {
      window.removeEventListener('storage', handleStorageChange);
    };
  }, [refreshUser]);

  const login = useCallback((newToken: string, userData: User) => {
    setToken(newToken);
    setUser(userData);
    setIsLoading(false);
  }, []);

  const logout = useCallback(() => {
    removeToken();
    setUser(null);
    setToken(null);
    setIsLoading(false);
  }, []);

  const value: AuthContextType = useMemo(
    () => ({
      user,
      isAuthenticated: !!user && !!token,
      token,
      isLoading,
      login,
      logout,
      refreshUser,
    }),
    [user, token, isLoading, login, logout, refreshUser]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

