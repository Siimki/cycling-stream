package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type UserPreferencesRepository struct {
	db *sql.DB
}

func NewUserPreferencesRepository(db *sql.DB) *UserPreferencesRepository {
	return &UserPreferencesRepository{db: db}
}

func (r *UserPreferencesRepository) GetByUserID(userID string) (*models.UserPreferences, error) {
	query := `
		SELECT id, user_id, data_mode, preferred_units, theme, accent_color, device_type,
		       notification_preferences, onboarding_completed, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	var prefs models.UserPreferences
	var notificationPrefsJSON []byte
	var accentColor sql.NullString
	var deviceType sql.NullString

	err := r.db.QueryRow(query, userID).Scan(
		&prefs.ID,
		&prefs.UserID,
		&prefs.DataMode,
		&prefs.PreferredUnits,
		&prefs.Theme,
		&accentColor,
		&deviceType,
		&notificationPrefsJSON,
		&prefs.OnboardingCompleted,
		&prefs.CreatedAt,
		&prefs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return defaults if preferences don't exist
		defaults := models.GetDefaultPreferences()
		defaults.UserID = userID
		return defaults, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Parse JSONB notification preferences
	if len(notificationPrefsJSON) > 0 {
		if err := json.Unmarshal(notificationPrefsJSON, &prefs.NotificationPreferences); err != nil {
			prefs.NotificationPreferences = make(map[string]interface{})
		}
	} else {
		prefs.NotificationPreferences = make(map[string]interface{})
	}

	if accentColor.Valid {
		prefs.AccentColor = &accentColor.String
	}
	if deviceType.Valid {
		prefs.DeviceType = &deviceType.String
	}

	return &prefs, nil
}

func (r *UserPreferencesRepository) Create(prefs *models.UserPreferences) error {
	prefs.ID = uuid.New().String()

	// Convert notification preferences to JSONB
	notificationPrefsJSON, err := json.Marshal(prefs.NotificationPreferences)
	if err != nil {
		return fmt.Errorf("failed to marshal notification preferences: %w", err)
	}

	query := `
		INSERT INTO user_preferences (
			id, user_id, data_mode, preferred_units, theme, accent_color, device_type,
			notification_preferences, onboarding_completed
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	err = r.db.QueryRow(
		query,
		prefs.ID,
		prefs.UserID,
		prefs.DataMode,
		prefs.PreferredUnits,
		prefs.Theme,
		prefs.AccentColor,
		prefs.DeviceType,
		notificationPrefsJSON,
		prefs.OnboardingCompleted,
	).Scan(&prefs.CreatedAt, &prefs.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user preferences: %w", err)
	}

	return nil
}

func (r *UserPreferencesRepository) Update(userID string, req *models.UpdatePreferencesRequest) (*models.UserPreferences, error) {
	// Get existing preferences first
	existing, err := r.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.DataMode != nil {
		existing.DataMode = *req.DataMode
	}
	if req.PreferredUnits != nil {
		existing.PreferredUnits = *req.PreferredUnits
	}
	if req.Theme != nil {
		existing.Theme = *req.Theme
	}
	if req.AccentColor != nil {
		existing.AccentColor = req.AccentColor
	}
	if req.DeviceType != nil {
		existing.DeviceType = req.DeviceType
	}
	if req.NotificationPreferences != nil {
		existing.NotificationPreferences = *req.NotificationPreferences
	}
	if req.OnboardingCompleted != nil {
		existing.OnboardingCompleted = *req.OnboardingCompleted
	}

	// Convert notification preferences to JSONB
	notificationPrefsJSON, err := json.Marshal(existing.NotificationPreferences)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification preferences: %w", err)
	}

	// If preferences don't exist, create them; otherwise update
	if existing.ID == "" {
		existing.UserID = userID
		if err := r.Create(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	query := `
		UPDATE user_preferences
		SET data_mode = $1, preferred_units = $2, theme = $3, accent_color = $4,
		    device_type = $5, notification_preferences = $6, onboarding_completed = $7,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $8
		RETURNING updated_at
	`

	err = r.db.QueryRow(
		query,
		existing.DataMode,
		existing.PreferredUnits,
		existing.Theme,
		existing.AccentColor,
		existing.DeviceType,
		notificationPrefsJSON,
		existing.OnboardingCompleted,
		userID,
	).Scan(&existing.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update user preferences: %w", err)
	}

	return existing, nil
}

