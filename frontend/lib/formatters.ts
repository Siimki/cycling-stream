/**
 * Formatting utility functions
 */

/**
 * Formats a date string to a readable format
 * @param dateString - ISO date string or undefined
 * @param options - Formatting options
 * @returns Formatted date string or 'TBD' if date is not provided
 */
export function formatDate(
  dateString?: string,
  options?: {
    includeTime?: boolean;
    format?: 'short' | 'long';
  }
): string {
  if (!dateString) return 'TBD';

  const date = new Date(dateString);
  const formatOptions: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: options?.format === 'short' ? 'short' : 'long',
    day: 'numeric',
  };

  if (options?.includeTime) {
    formatOptions.hour = '2-digit';
    formatOptions.minute = '2-digit';
  }

  return date.toLocaleDateString('en-US', formatOptions);
}

/**
 * Formats seconds to a readable time format (e.g., "1h 30m" or "45m")
 * @param seconds - Total seconds
 * @returns Formatted time string
 */
export function formatTime(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  return h > 0 ? `${h}h ${m}m` : `${m}m`;
}

/**
 * Formats seconds to HH:MM:SS format
 * @param seconds - Total seconds
 * @returns Formatted time string in HH:MM:SS format
 */
export function formatTimeDetailed(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
}

