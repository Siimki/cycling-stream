package services

import (
	"fmt"

	"github.com/cyclingstream/backend/internal/config"
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
)

type PredictionService struct {
	predictionRepo   *repository.PredictionRepository
	userRepo         *repository.UserRepository
	xpService        *XPService
	missionTriggers  *MissionTriggers
	xpConfig         *config.XPConfig
}

func NewPredictionService(
	predictionRepo *repository.PredictionRepository,
	userRepo *repository.UserRepository,
	xpService *XPService,
	missionTriggers *MissionTriggers,
	xpConfig *config.XPConfig,
) *PredictionService {
	return &PredictionService{
		predictionRepo:  predictionRepo,
		userRepo:        userRepo,
		xpService:       xpService,
		missionTriggers: missionTriggers,
		xpConfig:        xpConfig,
	}
}

// PlaceBet places a bet on a prediction market
func (s *PredictionService) PlaceBet(userID, marketID, optionID string, stakePoints int) error {
	// Validate bet
	if err := s.ValidateBet(userID, marketID, optionID, stakePoints); err != nil {
		return err
	}

	// Get market to find option and calculate payout
	market, err := s.predictionRepo.GetMarketByID(marketID)
	if err != nil {
		return fmt.Errorf("failed to get market: %w", err)
	}
	if market == nil {
		return fmt.Errorf("market not found")
	}

	if market.Status != "open" {
		return fmt.Errorf("market is not open for betting")
	}

	// Find the option and its odds
	var selectedOption *models.PredictionOption
	for i := range market.Options {
		if market.Options[i].ID == optionID {
			selectedOption = &market.Options[i]
			break
		}
	}

	if selectedOption == nil {
		return fmt.Errorf("option not found in market")
	}

	// Calculate potential payout (stake * odds, rounded to nearest integer)
	potentialPayout := int(float64(stakePoints) * selectedOption.Odds)

	// Deduct stake from user's points
	if err := s.userRepo.AddPoints(userID, -stakePoints); err != nil {
		return fmt.Errorf("failed to deduct points: %w", err)
	}

	// Create bet
	bet := &models.PredictionBet{
		UserID:          userID,
		MarketID:        marketID,
		OptionID:        optionID,
		StakePoints:     stakePoints,
		PotentialPayout: potentialPayout,
		Result:          "pending",
	}

	if err := s.predictionRepo.CreateBet(bet); err != nil {
		// Refund points if bet creation fails
		s.userRepo.AddPoints(userID, stakePoints)
		return fmt.Errorf("failed to create bet: %w", err)
	}

	// Award XP for placing bet (from config)
	if s.xpService != nil && s.xpConfig != nil {
		xpReward := s.xpConfig.Awards.Prediction.XPForPlacing
		if xpReward > 0 {
			if err := s.xpService.AwardXP(userID, xpReward, "prediction_placed"); err != nil {
				// Log error but don't fail
				fmt.Printf("Failed to award XP for bet placement: %v\n", err)
			}
		}
	}

	// Trigger mission progress for prediction placed
	if s.missionTriggers != nil {
		if err := s.missionTriggers.OnPredictionPlaced(userID); err != nil {
			// Log error but don't fail
			fmt.Printf("Failed to update mission progress: %v\n", err)
		}
	}

	return nil
}

// ValidateBet validates a bet before placement
func (s *PredictionService) ValidateBet(userID, marketID, optionID string, stakePoints int) error {
	if stakePoints <= 0 {
		return fmt.Errorf("stake must be greater than 0")
	}

	// Get user to check balance
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check sufficient balance
	if user.Points < stakePoints {
		return fmt.Errorf("insufficient points balance")
	}

	// Check max stake: 10% of balance or 1000 points, whichever is lower
	maxStakePercent := user.Points / 10
	maxStake := 1000
	if maxStakePercent < maxStake {
		maxStake = maxStakePercent
	}

	if stakePoints > maxStake {
		return fmt.Errorf("stake exceeds maximum allowed (%d points)", maxStake)
	}

	// Validate market exists and is open
	market, err := s.predictionRepo.GetMarketByID(marketID)
	if err != nil {
		return fmt.Errorf("failed to get market: %w", err)
	}
	if market == nil {
		return fmt.Errorf("market not found")
	}

	if market.Status != "open" {
		return fmt.Errorf("market is not open for betting")
	}

	// Validate option exists in market
	optionFound := false
	for _, opt := range market.Options {
		if opt.ID == optionID {
			optionFound = true
			break
		}
	}

	if !optionFound {
		return fmt.Errorf("option not found in market")
	}

	return nil
}

// SettleMarket settles a prediction market and all its bets
func (s *PredictionService) SettleMarket(marketID, winningOptionID string) error {
	// Get market
	market, err := s.predictionRepo.GetMarketByID(marketID)
	if err != nil {
		return fmt.Errorf("failed to get market: %w", err)
	}
	if market == nil {
		return fmt.Errorf("market not found")
	}

	if market.Status != "open" {
		return fmt.Errorf("market is not open for settlement")
	}

	// Validate winning option exists
	optionFound := false
	for _, opt := range market.Options {
		if opt.ID == winningOptionID {
			optionFound = true
			break
		}
	}

	if !optionFound {
		return fmt.Errorf("winning option not found in market")
	}

	// Settle the market
	if err := s.predictionRepo.SettleMarket(marketID, winningOptionID); err != nil {
		return fmt.Errorf("failed to settle market: %w", err)
	}

	// Get all bets for this market
	bets, err := s.predictionRepo.GetBetsByMarket(marketID)
	if err != nil {
		return fmt.Errorf("failed to get bets: %w", err)
	}

	// Settle each bet
	for _, bet := range bets {
		won := bet.OptionID == winningOptionID
		payout := 0
		if won {
			payout = bet.PotentialPayout
			// Award payout to user
			if err := s.userRepo.AddPoints(bet.UserID, payout); err != nil {
				// Log error but continue
				fmt.Printf("Failed to award payout to user %s: %v\n", bet.UserID, err)
			}

			// Award XP for winning bet (from config)
			if s.xpService != nil && s.xpConfig != nil {
				xpReward := s.xpConfig.Awards.Prediction.XPForWinning
				if xpReward > 0 {
					if err := s.xpService.AwardXP(bet.UserID, xpReward, "prediction_won"); err != nil {
						// Log error but don't fail
						fmt.Printf("Failed to award XP for bet win: %v\n", err)
					}
				}
			}

			// Trigger mission progress for prediction won
			if s.missionTriggers != nil {
				if err := s.missionTriggers.OnPredictionWon(bet.UserID); err != nil {
					// Log error but don't fail
					fmt.Printf("Failed to update mission progress: %v\n", err)
				}
			}
		}

		// Settle the bet
		if err := s.predictionRepo.SettleBet(bet.ID, won, payout); err != nil {
			// Log error but continue
			fmt.Printf("Failed to settle bet %s: %v\n", bet.ID, err)
		}
	}

	return nil
}

// GetUserBetBalance returns the available points balance for betting
func (s *PredictionService) GetUserBetBalance(userID string) (int, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return 0, fmt.Errorf("user not found")
	}

	return user.Points, nil
}

// GetMarketsByRace gets all prediction markets for a race
func (s *PredictionService) GetMarketsByRace(raceID string) ([]models.PredictionMarket, error) {
	return s.predictionRepo.GetMarketsByRace(raceID)
}

// GetBetsByUser gets all bets for a user
func (s *PredictionService) GetBetsByUser(userID string) ([]models.PredictionBet, error) {
	return s.predictionRepo.GetBetsByUser(userID)
}

