import { memo } from 'react';
import { Button } from '@/components/ui/button';
import { APIErrorHandler } from '@/lib/error-handler';
import { AlertCircle } from 'lucide-react';

interface ErrorMessageProps {
  message?: string;
  error?: unknown;
  onRetry?: () => void;
  variant?: 'default' | 'inline' | 'full';
  className?: string;
}

/**
 * Standardized error message component
 * Can display errors from APIError, Error, or string
 */
function ErrorMessage({ 
  message, 
  error, 
  onRetry, 
  variant = 'default',
  className = '' 
}: ErrorMessageProps) {
  // Get error message from error object or use provided message
  const errorMessage = message || (error ? APIErrorHandler.getErrorMessage(error) : 'An error occurred');

  // Determine if it's a network error for better UX
  const isNetworkError = error ? APIErrorHandler.isNetworkError(error) : false;

  if (variant === 'inline') {
    return (
      <div className={`text-sm text-destructive ${className}`}>
        {errorMessage}
      </div>
    );
  }

  if (variant === 'full') {
    return (
      <div className={`min-h-[400px] flex items-center justify-center ${className}`}>
        <div className="max-w-md w-full text-center">
          <AlertCircle className="w-16 h-16 text-destructive mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-foreground mb-2">Error</h2>
          <p className="text-muted-foreground mb-6">{errorMessage}</p>
          {onRetry && (
            <Button onClick={onRetry} variant="default">
              Try Again
            </Button>
          )}
        </div>
      </div>
    );
  }

  // Default variant
  return (
    <div className={`bg-destructive/10 border border-destructive/20 rounded-lg p-4 ${className}`}>
      <div className="flex items-start gap-3">
        <AlertCircle className="w-5 h-5 text-destructive shrink-0 mt-0.5" />
        <div className="flex-1 min-w-0">
          <h3 className="text-destructive font-semibold mb-1">
            {isNetworkError ? 'Connection Error' : 'Error'}
          </h3>
          <p className="text-destructive/90 text-sm">{errorMessage}</p>
          {isNetworkError && (
            <p className="text-destructive/70 text-xs mt-2">
              Please check your internet connection and try again.
            </p>
          )}
        </div>
        {onRetry && (
          <Button
            onClick={onRetry}
            variant="destructive"
            size="sm"
            className="shrink-0"
          >
            Retry
          </Button>
        )}
      </div>
    </div>
  );
}

export default memo(ErrorMessage);
