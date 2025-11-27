package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(message *models.ChatMessage) error {
	message.ID = uuid.New().String()
	query := `
		INSERT INTO chat_messages (id, race_id, user_id, username, message)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	// Handle nullable user_id
	var userID interface{}
	if message.UserID != nil {
		userID = *message.UserID
	} else {
		userID = nil
	}

	err := r.db.QueryRow(
		query,
		message.ID,
		message.RaceID,
		userID,
		message.Username,
		message.Message,
	).Scan(&message.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create chat message (id=%s, race_id=%s, user_id=%v, username=%s): %w",
			message.ID, message.RaceID, userID, message.Username, err)
	}

	return nil
}

func (r *ChatRepository) GetByRaceID(raceID string, limit, offset int) ([]*models.ChatMessage, error) {
	query := `
		SELECT id, race_id, user_id, username, message, created_at
		FROM chat_messages
		WHERE race_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, raceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		var userID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.RaceID,
			&userID,
			&msg.Username,
			&msg.Message,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat message: %w", err)
		}

		if userID.Valid {
			msg.UserID = &userID.String
		}

		messages = append(messages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chat messages: %w", err)
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *ChatRepository) GetRecentByRaceID(raceID string, limit int) ([]*models.ChatMessage, error) {
	query := `
		SELECT id, race_id, user_id, username, message, created_at
		FROM chat_messages
		WHERE race_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, raceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		var userID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.RaceID,
			&userID,
			&msg.Username,
			&msg.Message,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat message: %w", err)
		}

		if userID.Valid {
			msg.UserID = &userID.String
		}

		messages = append(messages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chat messages: %w", err)
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *ChatRepository) CountByRaceID(raceID string) (int, error) {
	query := `SELECT COUNT(*) FROM chat_messages WHERE race_id = $1`

	var count int
	err := r.db.QueryRow(query, raceID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count chat messages: %w", err)
	}

	return count, nil
}
