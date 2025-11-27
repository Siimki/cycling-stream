package handlers

import (
	"time"

	"github.com/cyclingstream/backend/internal/middleware"
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     *string `json:"name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate email
	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	req.Email = middleware.SanitizeString(req.Email, 255)
	if !middleware.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	// Validate password
	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}
	if valid, msg := middleware.ValidatePassword(req.Password); !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	// Sanitize name if provided
	if req.Name != nil {
		sanitized := middleware.SanitizeString(*req.Name, 255)
		req.Name = &sanitized
	}

	// Sanitize and validate bio if provided (max 120 chars)
	bio := ""
	if req.Bio != nil {
		bio = middleware.SanitizeString(*req.Bio, 120)
	}

	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check user existence",
		})
	}
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Bio:          bio,
		Points:       0, // New users start with 0 points
	}

	if err := h.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate token
	token, err := h.generateToken(user.ID, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Don't return password hash
	user.PasswordHash = ""

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate and sanitize email
	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	req.Email = middleware.SanitizeString(req.Email, 255)
	if !middleware.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	// Validate password (basic check - don't reveal too much)
	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}
	if len(req.Password) > 128 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check for admin credentials first
	if req.Email == "admin@cyclingstream.local" && req.Password == "admin123" {
		tokenString, err := h.generateToken("admin", true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(AuthResponse{
			Token: tokenString,
			User: &models.User{
				ID:    "admin",
				Email: "admin@cyclingstream.local",
			},
		})
	}

	// Check database for regular users
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate token
	token, err := h.generateToken(user.ID, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Don't return password hash
	user.PasswordHash = ""

	return c.Status(fiber.StatusOK).JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Don't return password hash
	user.PasswordHash = ""

	return c.Status(fiber.StatusOK).JSON(user)
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	var req ChangePasswordRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate current password
	if req.CurrentPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current password is required",
		})
	}

	// Validate new password
	if req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password is required",
		})
	}
	if valid, msg := middleware.ValidatePassword(req.NewPassword); !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	// Get user to verify current password
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Current password is incorrect",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Update password
	if err := h.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

// AwardBonusPoints grants a fixed bonus of points to the authenticated user.
// This is used for actions like collecting a watch bonus in the UI.
func (h *AuthHandler) AwardBonusPoints(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	const bonusPoints = 10

	if err := h.userRepo.AddPoints(userID, bonusPoints); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to award bonus points",
		})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load updated user points",
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Bonus points awarded",
		"bonus_points": bonusPoints,
		"total_points": user.Points,
	})
}

func (h *AuthHandler) generateToken(userID string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
