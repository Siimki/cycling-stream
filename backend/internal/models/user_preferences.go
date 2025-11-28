package models

import "time"

type UIPreferences struct {
	ChatAnimations bool `json:"chat_animations"`
	ReducedMotion  bool `json:"reduced_motion"`
	ButtonPulse    bool `json:"button_pulse"`
	PollAnimations bool `json:"poll_animations"`
}

func DefaultUIPreferences() UIPreferences {
	return UIPreferences{
		ChatAnimations: true,
		ReducedMotion:  false,
		ButtonPulse:    true,
		PollAnimations: true,
	}
}

type AudioPreferences struct {
	ButtonClicks       bool    `json:"button_clicks"`
	NotificationSounds bool    `json:"notification_sounds"`
	MentionPings       bool    `json:"mention_pings"`
	MasterVolume       float64 `json:"master_volume"`
}

func DefaultAudioPreferences() AudioPreferences {
	return AudioPreferences{
		ButtonClicks:       true,
		NotificationSounds: true,
		MentionPings:       true,
		MasterVolume:       0.15,
	}
}

type UserPreferences struct {
	ID                      string                 `json:"id" db:"id"`
	UserID                  string                 `json:"user_id" db:"user_id"`
	DataMode                string                 `json:"data_mode" db:"data_mode"`
	PreferredUnits          string                 `json:"preferred_units" db:"preferred_units"`
	Theme                   string                 `json:"theme" db:"theme"`
	AccentColor             *string                `json:"accent_color,omitempty" db:"accent_color"`
	DeviceType              *string                `json:"device_type,omitempty" db:"device_type"`
	NotificationPreferences map[string]interface{} `json:"notification_preferences" db:"-"`
	UIPreferences           UIPreferences          `json:"ui_preferences" db:"-"`
	AudioPreferences        AudioPreferences       `json:"audio_preferences" db:"-"`
	OnboardingCompleted     bool                   `json:"onboarding_completed" db:"onboarding_completed"`
	CreatedAt               time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at" db:"updated_at"`
}

// GetDefaultPreferences returns default preferences for a new user
func GetDefaultPreferences() *UserPreferences {
	return &UserPreferences{
		DataMode:                "standard",
		PreferredUnits:          "metric",
		Theme:                   "auto",
		NotificationPreferences: make(map[string]interface{}),
		UIPreferences:           DefaultUIPreferences(),
		AudioPreferences:        DefaultAudioPreferences(),
		OnboardingCompleted:     false,
	}
}

type UpdatePreferencesRequest struct {
	DataMode                *string                        `json:"data_mode,omitempty"`
	PreferredUnits          *string                        `json:"preferred_units,omitempty"`
	Theme                   *string                        `json:"theme,omitempty"`
	AccentColor             *string                        `json:"accent_color,omitempty"`
	DeviceType              *string                        `json:"device_type,omitempty"`
	NotificationPreferences *map[string]interface{}        `json:"notification_preferences,omitempty"`
	OnboardingCompleted     *bool                          `json:"onboarding_completed,omitempty"`
	UIPreferences           *UpdateUIPreferencesRequest    `json:"ui_preferences,omitempty"`
	AudioPreferences        *UpdateAudioPreferencesRequest `json:"audio_preferences,omitempty"`
}

type UpdateUIPreferencesRequest struct {
	ChatAnimations *bool `json:"chat_animations,omitempty"`
	ReducedMotion  *bool `json:"reduced_motion,omitempty"`
	ButtonPulse    *bool `json:"button_pulse,omitempty"`
	PollAnimations *bool `json:"poll_animations,omitempty"`
}

type UpdateAudioPreferencesRequest struct {
	ButtonClicks       *bool    `json:"button_clicks,omitempty"`
	NotificationSounds *bool    `json:"notification_sounds,omitempty"`
	MentionPings       *bool    `json:"mention_pings,omitempty"`
	MasterVolume       *float64 `json:"master_volume,omitempty"`
}
