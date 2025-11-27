package models

import "time"

type UserPreferences struct {
	ID                     string                 `json:"id" db:"id"`
	UserID                 string                 `json:"user_id" db:"user_id"`
	DataMode               string                 `json:"data_mode" db:"data_mode"`
	PreferredUnits         string                 `json:"preferred_units" db:"preferred_units"`
	Theme                  string                 `json:"theme" db:"theme"`
	AccentColor            *string                `json:"accent_color,omitempty" db:"accent_color"`
	DeviceType             *string                `json:"device_type,omitempty" db:"device_type"`
	NotificationPreferences map[string]interface{} `json:"notification_preferences" db:"-"`
	OnboardingCompleted    bool                   `json:"onboarding_completed" db:"onboarding_completed"`
	CreatedAt              time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at" db:"updated_at"`
}

// GetDefaultPreferences returns default preferences for a new user
func GetDefaultPreferences() *UserPreferences {
	return &UserPreferences{
		DataMode:               "standard",
		PreferredUnits:         "metric",
		Theme:                  "auto",
		NotificationPreferences: make(map[string]interface{}),
		OnboardingCompleted:    false,
	}
}

type UpdatePreferencesRequest struct {
	DataMode               *string                `json:"data_mode,omitempty"`
	PreferredUnits         *string                `json:"preferred_units,omitempty"`
	Theme                  *string                `json:"theme,omitempty"`
	AccentColor            *string                `json:"accent_color,omitempty"`
	DeviceType             *string                `json:"device_type,omitempty"`
	NotificationPreferences *map[string]interface{} `json:"notification_preferences,omitempty"`
	OnboardingCompleted    *bool                  `json:"onboarding_completed,omitempty"`
}

