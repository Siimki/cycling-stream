package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AchievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) UpsertAchievement(seed models.AchievementSeed) error {
	query := `
		INSERT INTO achievements (id, slug, title, description, icon, points)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (slug) DO UPDATE
		SET title = EXCLUDED.title,
		    description = EXCLUDED.description,
		    icon = EXCLUDED.icon,
		    points = EXCLUDED.points
	`

	id := uuid.New().String()
	_, err := r.db.Exec(query, id, seed.Slug, seed.Title, seed.Description, seed.Icon, seed.Points)
	if err != nil {
		return fmt.Errorf("failed to upsert achievement: %w", err)
	}
	return nil
}

func (r *AchievementRepository) EnsureDefaults(seeds []models.AchievementSeed) error {
	for _, seed := range seeds {
		if err := r.UpsertAchievement(seed); err != nil {
			return err
		}
	}
	return nil
}

func (r *AchievementRepository) GetBySlug(slug string) (*models.Achievement, error) {
	query := `
		SELECT id, slug, title, description, icon, points, created_at
		FROM achievements
		WHERE slug = $1
	`

	var achievement models.Achievement
	err := r.db.QueryRow(query, slug).Scan(
		&achievement.ID,
		&achievement.Slug,
		&achievement.Title,
		&achievement.Description,
		&achievement.Icon,
		&achievement.Points,
		&achievement.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get achievement: %w", err)
	}

	return &achievement, nil
}

func (r *AchievementRepository) Unlock(userID string, achievement *models.Achievement, metadata map[string]interface{}) (bool, error) {
	if achievement == nil {
		return false, fmt.Errorf("achievement is nil")
	}

	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		return false, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO user_achievements (id, user_id, achievement_id, metadata)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`

	id := uuid.New().String()
	result, err := r.db.Exec(query, id, userID, achievement.ID, metaJSON)
	if err != nil {
		return false, fmt.Errorf("failed to unlock achievement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to check inserted achievement row: %w", err)
	}
	return rowsAffected > 0, nil
}

func (r *AchievementRepository) HasUserAchievement(userID, achievementID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_achievements
			WHERE user_id = $1 AND achievement_id = $2
		)
	`

	var exists bool
	if err := r.db.QueryRow(query, userID, achievementID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check user achievement: %w", err)
	}
	return exists, nil
}

func (r *AchievementRepository) GetUserAchievements(userID string) ([]models.UserAchievement, error) {
	query := `
		SELECT 
			ua.id,
			ua.user_id,
			ua.achievement_id,
			a.slug,
			a.title,
			a.description,
			a.icon,
			a.points,
			ua.unlocked_at,
			ua.metadata
		FROM user_achievements ua
		INNER JOIN achievements a ON ua.achievement_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.unlocked_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user achievements: %w", err)
	}
	defer rows.Close()

	var results []models.UserAchievement
	for rows.Next() {
		var ua models.UserAchievement
		var metaJSON []byte
		err := rows.Scan(
			&ua.ID,
			&ua.UserID,
			&ua.AchievementID,
			&ua.Slug,
			&ua.Title,
			&ua.Description,
			&ua.Icon,
			&ua.Points,
			&ua.UnlockedAt,
			&metaJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user achievement: %w", err)
		}

		if len(metaJSON) > 0 {
			if err := json.Unmarshal(metaJSON, &ua.Metadata); err != nil {
				ua.Metadata = make(map[string]interface{})
			}
		} else {
			ua.Metadata = make(map[string]interface{})
		}

		results = append(results, ua)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user achievements: %w", err)
	}

	return results, nil
}

func (r *AchievementRepository) UnlockBySlug(userID, slug string, metadata map[string]interface{}) (bool, *models.Achievement, error) {
	achievement, err := r.GetBySlug(slug)
	if err != nil {
		return false, nil, err
	}
	if achievement == nil {
		return false, nil, fmt.Errorf("achievement slug %s not found", slug)
	}

	unlocked, err := r.Unlock(userID, achievement, metadata)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return false, achievement, nil
		}
		return false, nil, err
	}

	return unlocked, achievement, nil
}
