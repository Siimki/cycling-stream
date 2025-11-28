/**
 * API Error Handler Utility
 * Standardizes error handling for API calls
 */

import { logger } from './logger';

export interface APIError {
  message: string;
  status?: number;
  code?: string;
  details?: unknown;
}

export class APIErrorHandler {
  /**
   * Handles fetch API errors and converts them to standardized APIError
   */
  static async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      let errorData: { error?: string; message?: string; code?: string; details?: unknown };
      
      try {
        errorData = await response.json();
      } catch {
        // If response is not JSON, create a generic error
        errorData = { error: `HTTP ${response.status}: ${response.statusText}` };
      }

      const error: APIError = {
        message: errorData.error || errorData.message || `HTTP error! status: ${response.status}`,
        status: response.status,
        code: errorData.code,
        details: errorData.details,
      };

      // Don't log expected 401, 404, or 500 errors for authentication-required resources
      // These are handled gracefully by the UI
      if (response.status !== 401 && response.status !== 404 && response.status !== 500) {
        logger.error('API Error:', error);
      }
      throw error;
    }

    try {
      return await response.json();
    } catch (err) {
      logger.error('Failed to parse JSON response:', err);
      throw new Error('Invalid response format');
    }
  }

  /**
   * Converts any error to a user-friendly message
   */
  static getErrorMessage(error: unknown): string {
    if (error instanceof Error) {
      // Check if it's an APIError
      if ('status' in error) {
        const apiError = error as APIError;
        return apiError.message;
      }
      return error.message;
    }

    if (typeof error === 'string') {
      return error;
    }

    return 'An unexpected error occurred. Please try again.';
  }

  /**
   * Checks if error is a network error
   */
  static isNetworkError(error: unknown): boolean {
    if (error instanceof TypeError && error.message.includes('fetch')) {
      return true;
    }
    if (error instanceof Error && error.message.includes('network')) {
      return true;
    }
    return false;
  }

  /**
   * Checks if error is an authentication error
   */
  static isAuthError(error: unknown): boolean {
    if (error && typeof error === 'object' && 'status' in error) {
      const apiError = error as APIError;
      return apiError.status === 401 || apiError.status === 403;
    }
    return false;
  }

  /**
   * Checks if error is a not found error
   */
  static isNotFoundError(error: unknown): boolean {
    if (error && typeof error === 'object' && 'status' in error) {
      const apiError = error as APIError;
      return apiError.status === 404;
    }
    return false;
  }
}

