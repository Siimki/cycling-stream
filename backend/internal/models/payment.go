package models

import "time"

type Payment struct {
	ID                      string    `json:"id" db:"id"`
	UserID                  string    `json:"user_id" db:"user_id"`
	RaceID                  *string   `json:"race_id,omitempty" db:"race_id"`
	StripePaymentIntentID   *string   `json:"stripe_payment_intent_id,omitempty" db:"stripe_payment_intent_id"`
	StripeCheckoutSessionID *string   `json:"stripe_checkout_session_id,omitempty" db:"stripe_checkout_session_id"`
	AmountCents             int       `json:"amount_cents" db:"amount_cents"`
	Currency                string    `json:"currency" db:"currency"`
	Status                  string    `json:"status" db:"status"`
	PaymentType             string    `json:"payment_type" db:"payment_type"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}
