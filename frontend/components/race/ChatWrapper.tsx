'use client';

import { useAuth } from '@/contexts/AuthContext';
import { ChatProvider } from '@/components/chat/ChatProvider';
import dynamic from 'next/dynamic';

const Chat = dynamic(() => import('@/components/chat/Chat'), {
  loading: () => (
    <div className="flex flex-col h-full bg-card/50 border-l border-border min-h-0">
      {/* Header skeleton */}
      <div className="px-4 py-3 border-b border-border/50 flex items-center justify-between shrink-0 h-12">
        <div className="h-4 w-24 bg-muted animate-pulse rounded" />
        <div className="h-8 w-8 bg-muted animate-pulse rounded" />
      </div>
      {/* Messages skeleton */}
      <div className="flex-1 overflow-hidden px-4 py-3 space-y-2">
        <div className="h-4 w-3/4 bg-muted animate-pulse rounded" />
        <div className="h-4 w-full bg-muted animate-pulse rounded" />
        <div className="h-4 w-5/6 bg-muted animate-pulse rounded" />
        <div className="h-4 w-2/3 bg-muted animate-pulse rounded" />
      </div>
      {/* Input skeleton */}
      <div className="p-3 border-t border-border/50 shrink-0 bg-card/30">
        <div className="h-10 w-full bg-muted animate-pulse rounded" />
      </div>
    </div>
  ),
});

interface ChatWrapperProps {
  raceId: string;
  requiresLogin: boolean;
  isLive: boolean;
}

/**
 * Wrapper component that conditionally shows chat based on authentication
 * For login-required races, chat only shows when user is authenticated
 * For non-login-required races, chat always shows
 */
export function ChatWrapper({ raceId, requiresLogin, isLive }: ChatWrapperProps) {
  const { isAuthenticated, isLoading } = useAuth();

  // If race requires login, only show chat when authenticated
  if (requiresLogin) {
    if (isLoading) {
      // Show loading state while checking auth
      return (
        <div className="lg:w-80 xl:w-96 2xl:w-[400px] border-t lg:border-t-0 lg:border-l border-border flex flex-col h-[300px] sm:h-[350px] lg:h-[calc(100vh-4rem)] shrink-0 bg-background">
          <div className="flex flex-col h-full bg-card/50 border-l border-border min-h-0">
            <div className="px-4 py-3 border-b border-border/50 flex items-center justify-between shrink-0 h-12">
              <div className="h-4 w-24 bg-muted animate-pulse rounded" />
            </div>
            <div className="flex-1 flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          </div>
        </div>
      );
    }

    if (!isAuthenticated) {
      // Don't show chat for unauthenticated users on login-required races
      return (
        <div className="lg:w-80 xl:w-96 2xl:w-[400px] border-t lg:border-t-0 lg:border-l border-border flex flex-col h-[300px] sm:h-[350px] lg:h-[calc(100vh-4rem)] shrink-0 bg-background">
          <div className="flex flex-col h-full bg-card/50 border-l border-border min-h-0 items-center justify-center p-4">
            <div className="text-center text-muted-foreground">
              <div className="text-4xl mb-2">💬</div>
              <p className="text-sm">Chat is only available for logged-in users</p>
            </div>
          </div>
        </div>
      );
    }

    // User is authenticated - show chat (enabled if stream is live OR user is authenticated)
    return (
      <div className="lg:w-80 xl:w-96 2xl:w-[400px] border-t lg:border-t-0 lg:border-l border-border flex flex-col h-[300px] sm:h-[350px] lg:h-[calc(100vh-4rem)] shrink-0 bg-background relative z-0">
        <ChatProvider raceId={raceId} enabled={isLive || isAuthenticated}>
          <Chat />
        </ChatProvider>
      </div>
    );
  }

  // Show chat for races that don't require login (enabled if stream is live)
  return (
    <div className="lg:w-80 xl:w-96 2xl:w-[400px] border-t lg:border-t-0 lg:border-l border-border flex flex-col h-[300px] sm:h-[350px] lg:h-[calc(100vh-4rem)] shrink-0 bg-background relative z-0">
      <ChatProvider raceId={raceId} enabled={isLive}>
        <Chat />
      </ChatProvider>
    </div>
  );
}
