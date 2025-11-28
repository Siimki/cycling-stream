package handlers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cyclingstream/backend/internal/chat"
	"github.com/cyclingstream/backend/internal/logger"
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
	"github.com/cyclingstream/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ChatHandler struct {
	chatRepo        *repository.ChatRepository
	raceRepo        *repository.RaceRepository
	streamRepo      *repository.StreamRepository
	userRepo        *repository.UserRepository
	entitlementRepo *repository.EntitlementRepository
	hub             *chat.Hub
	rateLimiter     *chat.RateLimiter
	missionTriggers *services.MissionTriggers
	pollManager     *chat.PollManager
}

func NewChatHandler(
	chatRepo *repository.ChatRepository,
	raceRepo *repository.RaceRepository,
	streamRepo *repository.StreamRepository,
	userRepo *repository.UserRepository,
	entitlementRepo *repository.EntitlementRepository,
	hub *chat.Hub,
	rateLimiter *chat.RateLimiter,
	missionTriggers *services.MissionTriggers,
	pollManager *chat.PollManager,
) *ChatHandler {
	return &ChatHandler{
		chatRepo:        chatRepo,
		raceRepo:        raceRepo,
		streamRepo:      streamRepo,
		userRepo:        userRepo,
		entitlementRepo: entitlementRepo,
		hub:             hub,
		rateLimiter:     rateLimiter,
		missionTriggers: missionTriggers,
		pollManager:     pollManager,
	}
}

// HandleWebSocket handles WebSocket connections for chat
func (h *ChatHandler) HandleWebSocket(c *fiber.Ctx) error {
	// Upgrade to WebSocket immediately
	if !websocket.IsWebSocketUpgrade(c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "WebSocket upgrade required",
		})
	}

	raceID := c.Params("id")
	if raceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Race ID is required",
		})
	}

	// Basic validation: race IDs must be valid UUIDs.
	if _, err := uuid.Parse(raceID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID",
		})
	}

	// Get user info from locals (set by middleware)
	userID, _ := c.Locals("user_id").(string)
	isAdmin, _ := c.Locals("is_admin").(bool)

	// Verify race exists
	race, err := h.raceRepo.GetByID(raceID)
	if err != nil {
		logger.WithError(err).Error("Failed to fetch race for chat")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify race",
		})
	}

	if race == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Race not found",
		})
	}

	// Verify stream is live
	stream, err := h.streamRepo.GetByRaceID(raceID)
	if err != nil {
		logger.WithError(err).Error("Failed to fetch stream for chat")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify stream status",
		})
	}

	if stream == nil || stream.Status != "live" {
		return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{
			"error": "Chat is only available for live races",
		})
	}

	// Upgrade connection
	return websocket.New(func(conn *websocket.Conn) {
		var username string
		var userIDPtr *string

		var currentUser *models.User
		if userID != "" {
			userIDPtr = &userID
			// Get user to get username
			user, err := h.userRepo.GetByID(userID)
			if err == nil && user != nil {
				currentUser = user
				if user.Name != nil && *user.Name != "" {
					username = *user.Name
				} else {
					// Use email as fallback
					username = user.Email
				}
			} else {
				username = "User"
			}
		} else {
			username = "Anonymous"
		}

		// Create message handler
		messageHandler := func(client *chat.Client, msg *chat.WSMessage) {
			// Handle send_message type
			if msg.Type == string(chat.MessageTypeSendMessage) {
				h.handleSendMessage(client, raceID, msg, userIDPtr, username, currentUser)
			}
		}

		// Create onClose handler
		onClose := func(client *chat.Client) {
			// Send left message to room
			leftMsg := chat.NewLeftWSMessage(username)
			if leftBytes, err := json.Marshal(leftMsg); err == nil {
				h.hub.BroadcastToRoom(raceID, leftBytes)
			}

			// Leave room - use method directly since channels are package-private
			h.hub.LeaveRoom(client, raceID)
		}

		// Create client with message handler and onClose callback
		client := chat.NewClient(h.hub, conn, userIDPtr, username, isAdmin, messageHandler, onClose)

		// Start client (registers with hub and starts pumps)
		// Note: Start() registers the client and then blocks on readPump()
		// We need to join the room before readPump blocks, so we do it here
		// JoinRoom is thread-safe and doesn't require the client to be registered first
		h.hub.JoinRoom(client, raceID)

		// Send joined message to room
		joinedMsg := chat.NewJoinedWSMessage(username)
		if joinedBytes, err := json.Marshal(joinedMsg); err == nil {
			h.hub.BroadcastToRoom(raceID, joinedBytes)
		}

		// Start client (registers with hub and starts pumps)
		// readPump will block until connection closes
		client.Start()
	})(c)
}

// handleSendMessage processes a send_message request
func (h *ChatHandler) handleSendMessage(client *chat.Client, raceID string, msg *chat.WSMessage, userID *string, username string, user *models.User) {
	// Anonymous users cannot send messages
	if userID == nil {
		errorMsg := chat.NewErrorWSMessage("Authentication required to send messages")
		if errorBytes, err := json.Marshal(errorMsg); err == nil {
			client.SendMessage(errorBytes)
		}
		return
	}

	// Parse send message data
	sendData, err := chat.ParseSendMessageData(msg)
	if err != nil {
		logger.WithError(err).Error("Failed to parse send message data")
		errorMsg := chat.NewErrorWSMessage("Invalid message data")
		if errorBytes, err := json.Marshal(errorMsg); err == nil {
			client.SendMessage(errorBytes)
		}
		return
	}

	if sendData == nil {
		errorMsg := chat.NewErrorWSMessage("Invalid message format")
		if errorBytes, err := json.Marshal(errorMsg); err == nil {
			client.SendMessage(errorBytes)
		}
		return
	}

	// Validate message
	validatedMessage, err := chat.ValidateMessage(sendData.Message)
	if err != nil {
		errorMsg := chat.NewErrorWSMessage(err.Error())
		if errorBytes, err := json.Marshal(errorMsg); err == nil {
			client.SendMessage(errorBytes)
		}
		return
	}

	// Check rate limit
	identifier := *userID
	if !h.rateLimiter.CheckRateLimit(identifier) {
		errorMsg := chat.NewErrorWSMessage("Rate limit exceeded. Please wait before sending another message.")
		if errorBytes, err := json.Marshal(errorMsg); err == nil {
			client.SendMessage(errorBytes)
		}
		return
	}

	// Create chat message
	chatMsg := &models.ChatMessage{
		RaceID:   raceID,
		UserID:   userID,
		Username: username,
		Message:  validatedMessage,
	}

	h.applyMessageMetadata(chatMsg, user, client != nil && client.IsAdmin())

	// Save to database with retry logic for transient errors
	var dbErr error
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		dbErr = h.chatRepo.Create(chatMsg)
		if dbErr == nil {
			break
		}

		// Check if error is retryable (transient database errors)
		if !isRetryableDBError(dbErr) || attempt == maxRetries {
			break
		}

		// Wait before retry with exponential backoff
		waitTime := time.Duration(attempt*50) * time.Millisecond
		logger.WithFields(map[string]interface{}{
			"error":        dbErr.Error(),
			"race_id":      raceID,
			"attempt":      attempt,
			"max_retries":  maxRetries,
			"wait_time_ms": waitTime.Milliseconds(),
		}).Warn("Retrying chat message creation after transient error")
		time.Sleep(waitTime)
	}

	if dbErr != nil {
		// Log detailed error information for debugging
		logger.WithFields(map[string]interface{}{
			"error":      dbErr.Error(),
			"race_id":    raceID,
			"user_id":    userID,
			"username":   username,
			"message":    validatedMessage,
			"message_id": chatMsg.ID,
			"retryable":  isRetryableDBError(dbErr),
		}).Error("Failed to create chat message in database after retries")

		// Send a generic error message to the client to avoid leaking internal details
		errorMsg := chat.NewErrorWSMessage("Failed to send message. Please try again.")
		if errorBytes, err := json.Marshal(errorMsg); err == nil {
			client.SendMessage(errorBytes)
		}
		return
	}

	// Broadcast to room
	wsMsg := chat.NewMessageWSMessage(chatMsg)
	if wsBytes, err := json.Marshal(wsMsg); err == nil {
		h.hub.BroadcastToRoom(raceID, wsBytes)
	}

	// Trigger mission progress updates for chat messages
	if h.missionTriggers != nil && userID != nil && *userID != "" {
		// Check if stream is live
		isLive := false
		if h.streamRepo != nil {
			stream, err := h.streamRepo.GetByRaceID(raceID)
			if err == nil && stream != nil && stream.Status == "live" {
				isLive = true
			}
		}
		if err := h.missionTriggers.OnChatMessage(*userID, raceID, isLive); err != nil {
			// Log error but don't fail the request
			// Mission progress is best-effort
		}
		// Check and complete any missions that may have been completed
		if err := h.missionTriggers.CheckAndCompleteAll(*userID); err != nil {
			// Log error but don't fail the request
		}
	}
}

// GetChatHistory returns paginated chat history for a race
func (h *ChatHandler) GetChatHistory(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	// Validate raceID is a valid UUID
	if _, err := uuid.Parse(raceID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid race ID",
		})
	}

	// Get limit and offset from query params
	limit := 50
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	messages, err := h.chatRepo.GetByRaceID(raceID, limit, offset)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"race_id": raceID,
			"limit":   limit,
			"offset":  offset,
		}).Error("Failed to get chat history")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch chat history",
		})
	}

	// Ensure messages is always an array, not null
	if messages == nil {
		messages = []*models.ChatMessage{}
	}

	h.hydrateMessageMetadata(messages)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetChatStats returns chat statistics for a race
func (h *ChatHandler) GetChatStats(c *fiber.Ctx) error {
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	count, err := h.chatRepo.CountByRaceID(raceID)
	if err != nil {
		logger.WithError(err).Error("Failed to get chat stats")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch chat stats",
		})
	}

	// Get concurrent connections in room
	concurrentConnections := h.hub.GetRoomClientCount(raceID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_messages":         count,
		"concurrent_connections": concurrentConnections,
	})
}

func (h *ChatHandler) applyMessageMetadata(msg *models.ChatMessage, user *models.User, isAdmin bool) {
	if msg == nil {
		return
	}

	if msg.Badges == nil {
		msg.Badges = []string{}
	}

	role, badges := h.resolveRoleAndBadges(user, isAdmin)
	msg.Role = role
	if len(msg.Badges) == 0 && len(badges) > 0 {
		msg.Badges = badges
	} else if len(badges) > 0 {
		msg.Badges = dedupeStrings(append(msg.Badges, badges...))
	}
	msg.SpecialEmote = chat.IsSpecialEmoteMessage(msg.Message)
}

func (h *ChatHandler) hydrateMessageMetadata(messages []*models.ChatMessage) {
	if len(messages) == 0 {
		return
	}

	missing := make(map[string]struct{})
	for _, msg := range messages {
		if msg == nil || msg.UserID == nil || msg.Role != "" {
			continue
		}
		// Skip empty user IDs
		if *msg.UserID == "" {
			continue
		}
		missing[*msg.UserID] = struct{}{}
	}

	if len(missing) == 0 {
		return
	}

	userIDs := make([]string, 0, len(missing))
	for id := range missing {
		userIDs = append(userIDs, id)
	}

	users, err := h.userRepo.GetByIDs(userIDs)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":     err.Error(),
			"user_ids":  userIDs,
			"count":     len(userIDs),
		}).Warn("Failed to backfill chat metadata")
		return
	}

	// Apply metadata for each message, handling missing users gracefully
	for _, msg := range messages {
		if msg == nil || msg.UserID == nil || msg.Role != "" {
			continue
		}
		// Check if user exists in the map - if not, user was deleted or doesn't exist
		user, userExists := users[*msg.UserID]
		if !userExists {
			// Set default role/badges for messages from deleted users
			if msg.Badges == nil {
				msg.Badges = []string{}
			}
			if msg.Role == "" {
				msg.Role = "viewer"
			}
			continue
		}
		h.applyMessageMetadata(msg, user, false)
	}
}

func (h *ChatHandler) resolveRoleAndBadges(user *models.User, isAdmin bool) (string, []string) {
	if isAdmin {
		return "mod", []string{"mod"}
	}

	if user == nil {
		return "viewer", []string{}
	}

	badges := make([]string, 0, 3)
	if user.Level >= 15 {
		badges = append(badges, "vip")
	}

	hasSubscription := false
	if h.entitlementRepo != nil && user.ID != "" {
		if ok, err := h.entitlementRepo.HasActiveSubscription(user.ID); err == nil && ok {
			hasSubscription = true
			badges = append(badges, "sub")
		} else if err != nil {
			logger.WithError(err).Warn("Failed to check subscription for chat role")
		}
	}

	if !hasSubscription && user.Points >= 500 {
		badges = append(badges, "sub")
	}

	if user.BestStreakWeeks >= 4 {
		badges = append(badges, "og")
	}

	role := "viewer"
	if containsString(badges, "vip") {
		role = "vip"
	} else if containsString(badges, "sub") {
		role = "subscriber"
	}

	return role, dedupeStrings(badges)
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func dedupeStrings(items []string) []string {
	if len(items) == 0 {
		return items
	}

	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

type createPollRequest struct {
	Question        string   `json:"question"`
	Options         []string `json:"options"`
	DurationSeconds *int     `json:"duration_seconds,omitempty"`
}

type votePollRequest struct {
	OptionID string `json:"option_id"`
}

func (h *ChatHandler) CreatePoll(c *fiber.Ctx) error {
	if h.pollManager == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIError{Error: "Poll manager unavailable"})
	}
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	var req createPollRequest
	if !parseBody(c, &req) {
		return nil
	}

	if strings.TrimSpace(req.Question) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "Question is required"})
	}

	duration := time.Duration(0)
	if req.DurationSeconds != nil && *req.DurationSeconds > 0 {
		duration = time.Duration(*req.DurationSeconds) * time.Second
	}

	poll, err := h.pollManager.CreatePoll(raceID, req.Question, req.Options, duration)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: err.Error()})
	}

	if pollBytes, err := json.Marshal(chat.NewPollAnnouncementMessage(poll)); err == nil {
		h.hub.BroadcastToRoom(raceID, pollBytes)
	}

	return c.Status(fiber.StatusCreated).JSON(poll)
}

func (h *ChatHandler) CastPollVote(c *fiber.Ctx) error {
	if h.pollManager == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIError{Error: "Poll manager unavailable"})
	}
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	pollID, ok := requireParam(c, "pollId", "Poll ID is required")
	if !ok {
		return nil
	}

	userID, ok := requireUserID(c, "Authentication required")
	if !ok {
		return nil
	}

	var req votePollRequest
	if !parseBody(c, &req) {
		return nil
	}
	if req.OptionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: "option_id is required"})
	}

	poll, exists := h.pollManager.GetPoll(pollID)
	if !exists || poll.RaceID != raceID {
		return c.Status(fiber.StatusNotFound).JSON(APIError{Error: "Poll not found"})
	}

	updated, err := h.pollManager.Vote(pollID, userID, req.OptionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: err.Error()})
	}

	if pollBytes, err := json.Marshal(chat.NewPollUpdateMessage(updated)); err == nil {
		h.hub.BroadcastToRoom(raceID, pollBytes)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (h *ChatHandler) ClosePoll(c *fiber.Ctx) error {
	if h.pollManager == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIError{Error: "Poll manager unavailable"})
	}
	raceID, ok := requireParam(c, "id", "Race ID is required")
	if !ok {
		return nil
	}

	pollID, ok := requireParam(c, "pollId", "Poll ID is required")
	if !ok {
		return nil
	}

	poll, err := h.pollManager.ClosePoll(pollID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIError{Error: err.Error()})
	}

	if pollBytes, err := json.Marshal(chat.NewPollClosedMessage(poll)); err == nil {
		h.hub.BroadcastToRoom(raceID, pollBytes)
	}

	return c.Status(fiber.StatusOK).JSON(poll)
}

// isRetryableDBError checks if a database error is retryable (transient)
// Returns true for connection errors, deadlocks, and other transient issues
func isRetryableDBError(err error) bool {
	if err == nil {
		return false
	}

	// Check for PostgreSQL error codes
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		// Retryable PostgreSQL error codes:
		// - 40001: serialization_failure (deadlock)
		// - 40P01: deadlock_detected
		// - 08003: connection_does_not_exist
		// - 08006: connection_failure
		// - 08001: sqlclient_unable_to_establish_sqlconnection
		// - 08004: sqlserver_rejected_establishment_of_sqlconnection
		// - 57P01: admin_shutdown
		// - 57P02: crash_shutdown
		// - 57P03: cannot_connect_now
		retryableCodes := []string{
			"40001", // serialization_failure
			"40P01", // deadlock_detected
			"08003", // connection_does_not_exist
			"08006", // connection_failure
			"08001", // sqlclient_unable_to_establish_sqlconnection
			"08004", // sqlserver_rejected_establishment_of_sqlconnection
			"57P01", // admin_shutdown
			"57P02", // crash_shutdown
			"57P03", // cannot_connect_now
		}
		for _, code := range retryableCodes {
			if pqErr.Code == pq.ErrorCode(code) {
				return true
			}
		}
	}

	// Check for connection-related error messages (network issues, timeouts)
	errStr := strings.ToLower(err.Error())
	retryablePatterns := []string{
		"connection",
		"timeout",
		"deadlock",
		"temporary",
		"retry",
		"network",
		"broken pipe",
		"connection reset",
	}
	for _, pattern := range retryablePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}
