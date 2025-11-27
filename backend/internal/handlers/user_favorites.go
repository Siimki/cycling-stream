package handlers

import (
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type UserFavoritesHandler struct {
	favRepo *repository.UserFavoriteRepository
}

func NewUserFavoritesHandler(favRepo *repository.UserFavoriteRepository) *UserFavoritesHandler {
	return &UserFavoritesHandler{
		favRepo: favRepo,
	}
}

func (h *UserFavoritesHandler) GetFavorites(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	// Optional type filter
	favoriteType := c.Query("type")
	var typePtr *string
	if favoriteType != "" {
		// Validate favorite type
		validTypes := map[string]bool{
			"rider": true,
			"team":  true,
			"race":  true,
			"series": true,
		}
		if !validTypes[favoriteType] {
			return c.Status(fiber.StatusBadRequest).JSON(APIError{
				Error: "Invalid favorite_type. Must be 'rider', 'team', 'race', or 'series'",
			})
		}
		typePtr = &favoriteType
	}

	favorites, err := h.favRepo.GetByUserID(userID, typePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to fetch favorites",
		})
	}

	return c.Status(fiber.StatusOK).JSON(favorites)
}

func (h *UserFavoritesHandler) AddFavorite(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	var req models.AddFavoriteRequest
	if !parseBody(c, &req) {
		return nil
	}

	// Validate favorite type
	validTypes := map[string]bool{
		"rider":  true,
		"team":   true,
		"race":   true,
		"series": true,
	}
	if !validTypes[req.FavoriteType] {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "Invalid favorite_type. Must be 'rider', 'team', 'race', or 'series'",
		})
	}

	if req.FavoriteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "favorite_id is required",
		})
	}

	// Check if already exists
	exists, err := h.favRepo.Exists(userID, req.FavoriteType, req.FavoriteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to check favorite",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(APIError{
			Error: "Favorite already exists",
		})
	}

	favorite := &models.UserFavorite{
		UserID:       userID,
		FavoriteType: req.FavoriteType,
		FavoriteID:   req.FavoriteID,
	}

	if err := h.favRepo.Create(favorite); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to add favorite",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(favorite)
}

func (h *UserFavoritesHandler) RemoveFavorite(c *fiber.Ctx) error {
	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	favoriteType, ok := requireParam(c, "type", "Favorite type is required")
	if !ok {
		return nil
	}

	favoriteID, ok := requireParam(c, "id", "Favorite ID is required")
	if !ok {
		return nil
	}

	// Validate favorite type
	validTypes := map[string]bool{
		"rider":  true,
		"team":   true,
		"race":   true,
		"series": true,
	}
	if !validTypes[favoriteType] {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{
			Error: "Invalid favorite_type. Must be 'rider', 'team', 'race', or 'series'",
		})
	}

	if err := h.favRepo.Delete(userID, favoriteType, favoriteID); err != nil {
		if err.Error() == "favorite not found" {
			return c.Status(fiber.StatusNotFound).JSON(APIError{
				Error: "Favorite not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIError{
			Error: "Failed to remove favorite",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

