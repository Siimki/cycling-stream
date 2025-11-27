package repository

import (
	"database/sql"
	"fmt"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/google/uuid"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	payment.ID = uuid.New().String()
	query := `
		INSERT INTO payments (id, user_id, race_id, stripe_payment_intent_id, stripe_checkout_session_id, 
		                     amount_cents, currency, status, payment_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		payment.ID,
		payment.UserID,
		payment.RaceID,
		payment.StripePaymentIntentID,
		payment.StripeCheckoutSessionID,
		payment.AmountCents,
		payment.Currency,
		payment.Status,
		payment.PaymentType,
	).Scan(&payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (r *PaymentRepository) UpdateStatus(paymentIntentID string, status string) error {
	query := `
		UPDATE payments
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE stripe_payment_intent_id = $1
	`

	result, err := r.db.Exec(query, paymentIntentID, status)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

func (r *PaymentRepository) GetByCheckoutSessionID(sessionID string) (*models.Payment, error) {
	query := `
		SELECT id, user_id, race_id, stripe_payment_intent_id, stripe_checkout_session_id,
		       amount_cents, currency, status, payment_type, created_at, updated_at
		FROM payments
		WHERE stripe_checkout_session_id = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(query, sessionID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.RaceID,
		&payment.StripePaymentIntentID,
		&payment.StripeCheckoutSessionID,
		&payment.AmountCents,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentType,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}
