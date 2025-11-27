package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type UserFavoriteRepository struct {
	db *sql.DB
}

func NewUserFavoriteRepository(db *sql.DB) *UserFavoriteRepository {
	return &UserFavoriteRepository{db: db}
}

func (r *UserFavoriteRepository) GetByUserID(userID string, favoriteType *string) ([]*models.UserFavorite, error) {
	var query string
	var args []interface{}

	if favoriteType != nil {
		query = `
			SELECT id, user_id, favorite_type, favorite_id, created_at
			FROM user_favorites
			WHERE user_id = $1 AND favorite_type = $2
			ORDER BY created_at DESC
		`
		args = []interface{}{userID, *favoriteType}
	} else {
		query = `
			SELECT id, user_id, favorite_type, favorite_id, created_at
			FROM user_favorites
			WHERE user_id = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{userID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user favorites: %w", err)
	}
	defer rows.Close()

	var favorites []*models.UserFavorite
	for rows.Next() {
		var fav models.UserFavorite
		err := rows.Scan(
			&fav.ID,
			&fav.UserID,
			&fav.FavoriteType,
			&fav.FavoriteID,
			&fav.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user favorite: %w", err)
		}
		favorites = append(favorites, &fav)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user favorites: %w", err)
	}

	return favorites, nil
}

func (r *UserFavoriteRepository) Create(favorite *models.UserFavorite) error {
	favorite.ID = uuid.New().String()
	query := `
		INSERT INTO user_favorites (id, user_id, favorite_type, favorite_id)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		query,
		favorite.ID,
		favorite.UserID,
		favorite.FavoriteType,
		favorite.FavoriteID,
	).Scan(&favorite.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user favorite: %w", err)
	}

	return nil
}

func (r *UserFavoriteRepository) Delete(userID string, favoriteType string, favoriteID string) error {
	query := `
		DELETE FROM user_favorites
		WHERE user_id = $1 AND favorite_type = $2 AND favorite_id = $3
	`

	result, err := r.db.Exec(query, userID, favoriteType, favoriteID)
	if err != nil {
		return fmt.Errorf("failed to delete user favorite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("favorite not found")
	}

	return nil
}

func (r *UserFavoriteRepository) Exists(userID string, favoriteType string, favoriteID string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM user_favorites
		WHERE user_id = $1 AND favorite_type = $2 AND favorite_id = $3
	`

	var exists bool
	err := r.db.QueryRow(query, userID, favoriteType, favoriteID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if favorite exists: %w", err)
	}

	return exists, nil
}

