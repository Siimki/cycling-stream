package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/webhook"
)

type PaymentHandler struct {
	paymentRepo     *repository.PaymentRepository
	entitlementRepo *repository.EntitlementRepository
	raceRepo        *repository.RaceRepository
	stripeKey       string
	webhookSecret   string
}

func NewPaymentHandler(
	paymentRepo *repository.PaymentRepository,
	entitlementRepo *repository.EntitlementRepository,
	raceRepo *repository.RaceRepository,
	stripeKey string,
	webhookSecret string,
) *PaymentHandler {
	return &PaymentHandler{
		paymentRepo:     paymentRepo,
		entitlementRepo: entitlementRepo,
		raceRepo:        raceRepo,
		stripeKey:       stripeKey,
		webhookSecret:   webhookSecret,
	}
}

type CreateCheckoutRequest struct {
	RaceID string `json:"race_id"`
}

func (h *PaymentHandler) CreateCheckout(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	var req CreateCheckoutRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Get race
	race, err := h.raceRepo.GetByID(req.RaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch race",
		})
	}

	if race == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Race not found",
		})
	}

	if race.IsFree {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Race is free, no payment required",
		})
	}

	// Check if user already has access
	hasAccess, err := h.entitlementRepo.HasAccess(userID, req.RaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check access",
		})
	}

	if hasAccess {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You already have access to this race",
		})
	}

	// Create Stripe checkout session
	stripe.Key = h.stripeKey

	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(race.Name),
					},
					UnitAmount: stripe.Int64(int64(race.PriceCents)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(fmt.Sprintf("%s/races/%s/watch?payment=success", baseURL, req.RaceID)),
		CancelURL:  stripe.String(fmt.Sprintf("%s/races/%s?payment=cancelled", baseURL, req.RaceID)),
		Metadata: map[string]string{
			"user_id": userID,
			"race_id": req.RaceID,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create checkout session",
		})
	}

	// Create payment record
	payment := &models.Payment{
		UserID:                  userID,
		RaceID:                  &req.RaceID,
		StripeCheckoutSessionID: &sess.ID,
		AmountCents:             race.PriceCents,
		Currency:                "usd",
		Status:                  "pending",
		PaymentType:             "ticket",
	}

	if err := h.paymentRepo.Create(payment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment record",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"checkout_url": sess.URL,
		"session_id":   sess.ID,
	})
}

func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, h.webhookSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid webhook signature",
		})
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse session",
			})
		}

		// Get payment record
		payment, err := h.paymentRepo.GetByCheckoutSessionID(session.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get payment",
			})
		}

		if payment == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Payment not found",
			})
		}

		// Update payment status
		if session.PaymentStatus == "paid" {
			if err := h.paymentRepo.UpdateStatus(*payment.StripePaymentIntentID, "succeeded"); err != nil {
				// Log error but don't fail webhook
				fmt.Printf("Failed to update payment status: %v\n", err)
			}

			// Create entitlement
			if payment.RaceID != nil {
				entitlement := &models.Entitlement{
					UserID:    payment.UserID,
					RaceID:    *payment.RaceID,
					Type:      "ticket",
					ExpiresAt: nil, // No expiration for one-time tickets
				}

				if err := h.entitlementRepo.Create(entitlement); err != nil {
					fmt.Printf("Failed to create entitlement: %v\n", err)
				}
			}
		}

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse payment intent",
			})
		}

		// Update payment status
		if err := h.paymentRepo.UpdateStatus(paymentIntent.ID, "succeeded"); err != nil {
			fmt.Printf("Failed to update payment status: %v\n", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"received": true,
	})
}
