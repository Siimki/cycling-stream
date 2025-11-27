export type UserSegment = 'CasualViewer' | 'HardcoreFan' | 'DataNerd' | 'GamblerFantasy' | 'SocialViewer';

export interface SegmentationData {
  cyclingLevel: 'new' | 'casual' | 'superfan' | null;
  viewPreference: 'clean' | 'data-rich' | null;
  chatParticipation?: number; // Number of chat messages
  watchTime?: number; // Total watch time in minutes
  points?: number; // Total points
}

/**
 * Determines user segment based on onboarding answers and behavior
 */
export function determineUserSegment(data: SegmentationData): UserSegment {
  // Primary segmentation based on onboarding
  if (data.cyclingLevel === 'new' || data.viewPreference === 'clean') {
    return 'CasualViewer';
  }

  if (data.cyclingLevel === 'superfan' && data.viewPreference === 'data-rich') {
    return 'HardcoreFan';
  }

  if (data.viewPreference === 'data-rich') {
    return 'DataNerd';
  }

  // Secondary segmentation based on behavior (if available)
  if (data.chatParticipation && data.chatParticipation > 100) {
    return 'SocialViewer';
  }

  // Default to casual viewer
  return 'CasualViewer';
}

/**
 * Get segment-specific defaults
 */
export function getSegmentDefaults(segment: UserSegment) {
  switch (segment) {
    case 'CasualViewer':
      return {
        dataMode: 'casual' as const,
        showMinimalUI: true,
        defaultViewMode: 'clean' as const,
      };
    case 'HardcoreFan':
      return {
        dataMode: 'pro' as const,
        showMinimalUI: false,
        defaultViewMode: 'data-rich' as const,
      };
    case 'DataNerd':
      return {
        dataMode: 'pro' as const,
        showMinimalUI: false,
        defaultViewMode: 'data-rich' as const,
      };
    case 'GamblerFantasy':
      return {
        dataMode: 'standard' as const,
        showMinimalUI: false,
        defaultViewMode: 'standard' as const,
      };
    case 'SocialViewer':
      return {
        dataMode: 'standard' as const,
        showMinimalUI: false,
        defaultViewMode: 'standard' as const,
      };
    default:
      return {
        dataMode: 'standard' as const,
        showMinimalUI: false,
        defaultViewMode: 'standard' as const,
      };
  }
}

