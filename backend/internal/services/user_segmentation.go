package services

// UserSegment represents a user segment type
type UserSegment string

const (
	SegmentCasualViewer  UserSegment = "CasualViewer"
	SegmentHardcoreFan   UserSegment = "HardcoreFan"
	SegmentDataNerd      UserSegment = "DataNerd"
	SegmentGamblerFantasy UserSegment = "GamblerFantasy"
	SegmentSocialViewer  UserSegment = "SocialViewer"
)

// SegmentationData contains data used for user segmentation
type SegmentationData struct {
	CyclingLevel      string // "new", "casual", "superfan"
	ViewPreference    string // "clean", "data-rich"
	ChatParticipation int    // Number of chat messages
	WatchTime         int    // Total watch time in minutes
	Points            int    // Total points
}

// DetermineUserSegment determines the user segment based on onboarding answers and behavior
func DetermineUserSegment(data SegmentationData) UserSegment {
	// Primary segmentation based on onboarding
	if data.CyclingLevel == "new" || data.ViewPreference == "clean" {
		return SegmentCasualViewer
	}

	if data.CyclingLevel == "superfan" && data.ViewPreference == "data-rich" {
		return SegmentHardcoreFan
	}

	if data.ViewPreference == "data-rich" {
		return SegmentDataNerd
	}

	// Secondary segmentation based on behavior
	if data.ChatParticipation > 100 {
		return SegmentSocialViewer
	}

	// Default to casual viewer
	return SegmentCasualViewer
}

// GetSegmentDefaults returns segment-specific defaults
func GetSegmentDefaults(segment UserSegment) map[string]interface{} {
	switch segment {
	case SegmentCasualViewer:
		return map[string]interface{}{
			"data_mode":       "casual",
			"show_minimal_ui": true,
			"default_view_mode": "clean",
		}
	case SegmentHardcoreFan:
		return map[string]interface{}{
			"data_mode":       "pro",
			"show_minimal_ui": false,
			"default_view_mode": "data-rich",
		}
	case SegmentDataNerd:
		return map[string]interface{}{
			"data_mode":       "pro",
			"show_minimal_ui": false,
			"default_view_mode": "data-rich",
		}
	case SegmentGamblerFantasy:
		return map[string]interface{}{
			"data_mode":       "standard",
			"show_minimal_ui": false,
			"default_view_mode": "standard",
		}
	case SegmentSocialViewer:
		return map[string]interface{}{
			"data_mode":       "standard",
			"show_minimal_ui": false,
			"default_view_mode": "standard",
		}
	default:
		return map[string]interface{}{
			"data_mode":       "standard",
			"show_minimal_ui": false,
			"default_view_mode": "standard",
		}
	}
}

