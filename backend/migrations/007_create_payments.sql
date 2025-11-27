-- Create payments table to track payment transactions
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    race_id UUID REFERENCES races(id) ON DELETE SET NULL,
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    stripe_checkout_session_id VARCHAR(255) UNIQUE,
    amount_cents INTEGER NOT NULL,
    currency VARCHAR(3) DEFAULT 'usd',
    status VARCHAR(50) DEFAULT 'pending', -- pending, succeeded, failed, refunded
    payment_type VARCHAR(50) DEFAULT 'ticket', -- ticket, subscription
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_race_id ON payments(race_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_stripe_payment_intent_id ON payments(stripe_payment_intent_id);

