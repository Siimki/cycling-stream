package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type PredictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository(db *sql.DB) *PredictionRepository {
	return &PredictionRepository{db: db}
}

// CreateMarket creates a new prediction market
func (r *PredictionRepository) CreateMarket(market *models.PredictionMarket) error {
	market.ID = uuid.New().String()

	// Convert options to JSONB
	optionsJSON, err := json.Marshal(market.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	query := `
		INSERT INTO prediction_markets (id, race_id, question, options, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	err = r.db.QueryRow(
		query,
		market.ID,
		market.RaceID,
		market.Question,
		optionsJSON,
		market.Status,
	).Scan(&market.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create prediction market: %w", err)
	}

	return nil
}

// GetMarketByID gets a prediction market by ID
func (r *PredictionRepository) GetMarketByID(id string) (*models.PredictionMarket, error) {
	query := `
		SELECT id, race_id, question, options, status, settled_option_id, created_at, settled_at
		FROM prediction_markets
		WHERE id = $1
	`

	var market models.PredictionMarket
	var optionsJSON []byte
	var settledOptionID sql.NullString
	var settledAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&market.ID,
		&market.RaceID,
		&market.Question,
		&optionsJSON,
		&market.Status,
		&settledOptionID,
		&market.CreatedAt,
		&settledAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction market: %w", err)
	}

	// Unmarshal options
	if err := json.Unmarshal(optionsJSON, &market.Options); err != nil {
		return nil, fmt.Errorf("failed to unmarshal options: %w", err)
	}

	if settledOptionID.Valid {
		market.SettledOptionID = &settledOptionID.String
	}
	if settledAt.Valid {
		market.SettledAt = &settledAt.Time
	}

	return &market, nil
}

// GetMarketsByRace gets all prediction markets for a race
func (r *PredictionRepository) GetMarketsByRace(raceID string) ([]models.PredictionMarket, error) {
	query := `
		SELECT id, race_id, question, options, status, settled_option_id, created_at, settled_at
		FROM prediction_markets
		WHERE race_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, raceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction markets: %w", err)
	}
	defer rows.Close()

	var markets []models.PredictionMarket
	for rows.Next() {
		var market models.PredictionMarket
		var optionsJSON []byte
		var settledOptionID sql.NullString
		var settledAt sql.NullTime

		err := rows.Scan(
			&market.ID,
			&market.RaceID,
			&market.Question,
			&optionsJSON,
			&market.Status,
			&settledOptionID,
			&market.CreatedAt,
			&settledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prediction market: %w", err)
		}

		// Unmarshal options
		if err := json.Unmarshal(optionsJSON, &market.Options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %w", err)
		}

		if settledOptionID.Valid {
			market.SettledOptionID = &settledOptionID.String
		}
		if settledAt.Valid {
			market.SettledAt = &settledAt.Time
		}

		markets = append(markets, market)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating prediction markets: %w", err)
	}

	return markets, nil
}

// SettleMarket settles a prediction market with a winning option
func (r *PredictionRepository) SettleMarket(marketID, winningOptionID string) error {
	now := time.Now()
	query := `
		UPDATE prediction_markets
		SET status = 'settled', settled_option_id = $1, settled_at = $2
		WHERE id = $3 AND status = 'open'
	`

	result, err := r.db.Exec(query, winningOptionID, now, marketID)
	if err != nil {
		return fmt.Errorf("failed to settle market: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("market not found or already settled")
	}

	return nil
}

// CreateBet creates a new prediction bet
func (r *PredictionRepository) CreateBet(bet *models.PredictionBet) error {
	bet.ID = uuid.New().String()

	query := `
		INSERT INTO prediction_bets (id, user_id, market_id, option_id, stake_points, potential_payout, result)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		query,
		bet.ID,
		bet.UserID,
		bet.MarketID,
		bet.OptionID,
		bet.StakePoints,
		bet.PotentialPayout,
		bet.Result,
	).Scan(&bet.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create prediction bet: %w", err)
	}

	return nil
}

// GetBetByID gets a prediction bet by ID
func (r *PredictionRepository) GetBetByID(id string) (*models.PredictionBet, error) {
	query := `
		SELECT id, user_id, market_id, option_id, stake_points, potential_payout, result, payout_points, created_at, settled_at
		FROM prediction_bets
		WHERE id = $1
	`

	var bet models.PredictionBet
	var payoutPoints sql.NullInt64
	var settledAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&bet.ID,
		&bet.UserID,
		&bet.MarketID,
		&bet.OptionID,
		&bet.StakePoints,
		&bet.PotentialPayout,
		&bet.Result,
		&payoutPoints,
		&bet.CreatedAt,
		&settledAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction bet: %w", err)
	}

	if payoutPoints.Valid {
		payout := int(payoutPoints.Int64)
		bet.PayoutPoints = &payout
	}
	if settledAt.Valid {
		bet.SettledAt = &settledAt.Time
	}

	return &bet, nil
}

// GetBetsByUser gets all bets for a user
func (r *PredictionRepository) GetBetsByUser(userID string) ([]models.PredictionBet, error) {
	query := `
		SELECT id, user_id, market_id, option_id, stake_points, potential_payout, result, payout_points, created_at, settled_at
		FROM prediction_bets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bets: %w", err)
	}
	defer rows.Close()

	var bets []models.PredictionBet
	for rows.Next() {
		var bet models.PredictionBet
		var payoutPoints sql.NullInt64
		var settledAt sql.NullTime

		err := rows.Scan(
			&bet.ID,
			&bet.UserID,
			&bet.MarketID,
			&bet.OptionID,
			&bet.StakePoints,
			&bet.PotentialPayout,
			&bet.Result,
			&payoutPoints,
			&bet.CreatedAt,
			&settledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bet: %w", err)
		}

		if payoutPoints.Valid {
			payout := int(payoutPoints.Int64)
			bet.PayoutPoints = &payout
		}
		if settledAt.Valid {
			bet.SettledAt = &settledAt.Time
		}

		bets = append(bets, bet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bets: %w", err)
	}

	return bets, nil
}

// GetBetsByMarket gets all bets for a market
func (r *PredictionRepository) GetBetsByMarket(marketID string) ([]models.PredictionBet, error) {
	query := `
		SELECT id, user_id, market_id, option_id, stake_points, potential_payout, result, payout_points, created_at, settled_at
		FROM prediction_bets
		WHERE market_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, marketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get market bets: %w", err)
	}
	defer rows.Close()

	var bets []models.PredictionBet
	for rows.Next() {
		var bet models.PredictionBet
		var payoutPoints sql.NullInt64
		var settledAt sql.NullTime

		err := rows.Scan(
			&bet.ID,
			&bet.UserID,
			&bet.MarketID,
			&bet.OptionID,
			&bet.StakePoints,
			&bet.PotentialPayout,
			&bet.Result,
			&payoutPoints,
			&bet.CreatedAt,
			&settledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bet: %w", err)
		}

		if payoutPoints.Valid {
			payout := int(payoutPoints.Int64)
			bet.PayoutPoints = &payout
		}
		if settledAt.Valid {
			bet.SettledAt = &settledAt.Time
		}

		bets = append(bets, bet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bets: %w", err)
	}

	return bets, nil
}

// SettleBet settles a bet (marks as won or lost and sets payout)
func (r *PredictionRepository) SettleBet(betID string, won bool, payout int) error {
	now := time.Now()
	result := "lost"
	if won {
		result = "won"
	}

	query := `
		UPDATE prediction_bets
		SET result = $1, payout_points = $2, settled_at = $3
		WHERE id = $4 AND result = 'pending'
	`

	dbResult, err := r.db.Exec(query, result, payout, now, betID)
	if err != nil {
		return fmt.Errorf("failed to settle bet: %w", err)
	}

	rowsAffected, err := dbResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("bet not found or already settled")
	}

	return nil
}


