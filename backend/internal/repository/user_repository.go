package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/config"
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
	query := `SELECT id, email, password_hash, name, bio, points, xp_total, level, best_streak_weeks, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Bio,
		&user.Points,
		&user.XPTotal,
		&user.Level,
		&user.BestStreakWeeks,
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
	query := `SELECT id, email, password_hash, name, bio, points, xp_total, level, best_streak_weeks, created_at, updated_at FROM users WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Bio,
		&user.Points,
		&user.XPTotal,
		&user.Level,
		&user.BestStreakWeeks,
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
	query := `SELECT id, name, bio, points, xp_total, level, best_streak_weeks, created_at FROM users WHERE id = $1`

	var user models.PublicUser
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Bio,
		&user.Points,
		&user.XPTotal,
		&user.Level,
		&user.BestStreakWeeks,
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
	// Initialize XP/Level fields if not set
	if user.XPTotal == 0 {
		user.XPTotal = 0
	}
	if user.Level == 0 {
		user.Level = 1
	}
	if user.BestStreakWeeks == 0 {
		user.BestStreakWeeks = 0
	}
	query := `
		INSERT INTO users (id, email, password_hash, name, bio, points, xp_total, level, best_streak_weeks)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
		user.XPTotal,
		user.Level,
		user.BestStreakWeeks,
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

// AddXP increments a user's XP by the given amount and updates the level if needed.
// This method only updates XP; level calculation should be done by the service layer.
func (r *UserRepository) AddXP(userID string, xp int) error {
	if xp == 0 {
		return nil
	}

	query := `
		UPDATE users
		SET xp_total = xp_total + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, xp, userID)
	if err != nil {
		return fmt.Errorf("failed to update user XP: %w", err)
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

// UpdateLevel updates a user's level.
func (r *UserRepository) UpdateLevel(userID string, level int) error {
	query := `
		UPDATE users
		SET level = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, level, userID)
	if err != nil {
		return fmt.Errorf("failed to update user level: %w", err)
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

// GetLevelFromXP calculates the level from total XP using the provided config.
// Level 1: 0 to (BaseXP - 1) XP
// Level N (N > 1): Requires BaseXP + IncrementPerLevel * sum(1 to N-2)
// Progressive formula: BaseXP + IncrementPerLevel * (N-2)(N-1)/2
func GetLevelFromXP(xp int, cfg *config.LevelingConfig) int {
	if xp < cfg.BaseXP {
		return 1
	}
	// Solve: xp >= BaseXP + IncrementPerLevel * (level-2)(level-1)/2
	// This is a quadratic equation. Use binary search to find the correct level
	if cfg.IncrementPerLevel <= 0 {
		return 1
	}
	
	// Binary search for the correct level
	low := 2
	high := 1000 // Reasonable max level
	for low <= high {
		mid := (low + high) / 2
		xpForMid := GetXPForLevel(mid, cfg)
		if xpForMid <= xp {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return high
}

// GetXPForLevel returns the XP needed to REACH level N (the minimum XP for that level).
// Level 1: returns 0
// Level N (N > 1): returns BaseXP + IncrementPerLevel * sum(1 to N-2)
// Progressive formula using triangular numbers: BaseXP + IncrementPerLevel * (N-2)(N-1)/2
// This means each level requires progressively more XP:
//   Level 2: BaseXP (e.g., 100)
//   Level 3: BaseXP + IncrementPerLevel (e.g., 120)
//   Level 4: BaseXP + IncrementPerLevel * 3 (e.g., 160)
//   Level 5: BaseXP + IncrementPerLevel * 6 (e.g., 220)
func GetXPForLevel(level int, cfg *config.LevelingConfig) int {
	if level <= 1 {
		return 0
	}
	if level == 2 {
		return cfg.BaseXP
	}
	// XP needed to reach level N = BaseXP + IncrementPerLevel * (N-2)(N-1)/2
	// This uses triangular numbers: sum(1 to n) = n(n+1)/2
	// For level N, we need sum(1 to N-2) = (N-2)(N-1)/2
	n := level - 2
	return cfg.BaseXP + (cfg.IncrementPerLevel * n * (n + 1)) / 2
}

// GetXPForNextLevel returns the XP needed to reach the NEXT level (level N+1) from the current level N.
// This is the same as GetXPForLevel(level+1, cfg)
func GetXPForNextLevel(level int, cfg *config.LevelingConfig) int {
	return GetXPForLevel(level+1, cfg)
}

// UpdateBestStreak updates the best streak weeks for a user
func (r *UserRepository) UpdateBestStreak(userID string, bestStreak int) error {
	query := `
		UPDATE users
		SET best_streak_weeks = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, bestStreak, userID)
	if err != nil {
		return fmt.Errorf("failed to update best streak: %w", err)
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
