package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, bio, points, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Bio,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, bio, points, created_at, updated_at FROM users WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Bio,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetPublicByID(id string) (*models.PublicUser, error) {
	query := `SELECT id, name, bio, points, created_at FROM users WHERE id = $1`

	var user models.PublicUser
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Bio,
		&user.Points,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	query := `
		INSERT INTO users (id, email, password_hash, name, bio, points)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Bio,
		user.Points,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdatePassword(userID string, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// AddPoints increments a user's points by the given delta.
// Delta can be positive (earn points) or negative (spend points).
func (r *UserRepository) AddPoints(userID string, delta int) error {
	if delta == 0 {
		return nil
	}

	query := `
		UPDATE users
		SET points = points + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, delta, userID)
	if err != nil {
		return fmt.Errorf("failed to update user points: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetLeaderboard returns all users with their points and total watch time, ordered by points DESC
func (r *UserRepository) GetLeaderboard() ([]models.LeaderboardEntry, error) {
	query := `
		SELECT 
			u.id,
			u.name,
			u.points,
			COALESCE(SUM(ws.duration_seconds) / 60, 0)::int as total_watch_minutes
		FROM users u
		LEFT JOIN watch_sessions ws ON u.id = ws.user_id AND ws.duration_seconds IS NOT NULL
		GROUP BY u.id, u.name, u.points
		ORDER BY u.points DESC, total_watch_minutes DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		var name sql.NullString

		err := rows.Scan(
			&entry.ID,
			&name,
			&entry.Points,
			&entry.TotalWatchMinutes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}

		if name.Valid {
			entry.Name = &name.String
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating leaderboard entries: %w", err)
	}

	return entries, nil
}
