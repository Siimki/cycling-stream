'use client';

/**
 * React Error Boundary Component
 * Catches JavaScript errors anywhere in the child component tree
 */

import React, { Component, type ReactNode } from 'react';
import { logger } from '@/lib/logger';
import { Button } from '@/components/ui/button';
import Link from 'next/link';

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
    };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    logger.error('ErrorBoundary caught an error:', error, errorInfo);
    
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }
  }

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
    });
  };

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <div className="min-h-screen bg-background flex items-center justify-center px-4">
          <div className="max-w-md w-full text-center">
            <div className="mb-8">
              <div className="text-6xl mb-4">⚠️</div>
              <h1 className="text-4xl font-bold text-foreground/95 mb-2">Something went wrong</h1>
              <p className="text-muted-foreground mb-4">
                We encountered an unexpected error. Please try again.
              </p>
              {process.env.NODE_ENV === 'development' && this.state.error?.message && (
                <div className="mt-4 p-4 bg-destructive/10 border border-destructive/20 rounded-lg text-left">
                  <p className="text-sm font-semibold text-destructive mb-1">Error Details:</p>
                  <p className="text-sm text-destructive/80 font-mono break-all">
                    {this.state.error.message}
                  </p>
                </div>
              )}
            </div>

            <div className="space-y-4">
              <Button
                onClick={this.handleReset}
                className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold"
              >
                Try Again
              </Button>
              <div>
                <Link
                  href="/"
                  className="text-primary hover:underline"
                >
                  Go to Homepage
                </Link>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

