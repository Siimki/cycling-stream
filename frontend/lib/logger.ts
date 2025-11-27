/**
 * Centralized logging utility
 * Provides consistent logging interface for production and development
 */

type LogLevel = 'debug' | 'info' | 'warn' | 'error';

interface Logger {
  debug: (...args: unknown[]) => void;
  info: (...args: unknown[]) => void;
  warn: (...args: unknown[]) => void;
  error: (...args: unknown[]) => void;
}

const isDevelopment = process.env.NODE_ENV === 'development';

/**
 * Creates a logger instance with appropriate log levels
 */
function createLogger(context?: string): Logger {
  const prefix = context ? `[${context}]` : '';

  return {
    debug: (...args: unknown[]) => {
      if (isDevelopment) {
        console.debug(prefix, ...args);
      }
    },
    info: (...args: unknown[]) => {
      if (isDevelopment) {
        console.info(prefix, ...args);
      }
    },
    warn: (...args: unknown[]) => {
      console.warn(prefix, ...args);
    },
    error: (...args: unknown[]) => {
      console.error(prefix, ...args);
      // In production, you might want to send errors to a logging service
      // e.g., Sentry, LogRocket, etc.
    },
  };
}

// Default logger instance
export const logger = createLogger();

// Factory function to create context-specific loggers
export function createContextLogger(context: string): Logger {
  return createLogger(context);
}

