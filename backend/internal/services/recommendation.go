package services

import (
	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
)

type RecommendationService struct {
	raceRepo           *repository.RaceRepository
	watchHistoryRepo   *repository.WatchHistoryRepository
	userFavRepo        *repository.UserFavoriteRepository
	streamRepo         *repository.StreamRepository
}

func NewRecommendationService(
	raceRepo *repository.RaceRepository,
	watchHistoryRepo *repository.WatchHistoryRepository,
	userFavRepo *repository.UserFavoriteRepository,
	streamRepo *repository.StreamRepository,
) *RecommendationService {
	return &RecommendationService{
		raceRepo:         raceRepo,
		watchHistoryRepo: watchHistoryRepo,
		userFavRepo:      userFavRepo,
		streamRepo:       streamRepo,
	}
}

// GetContinueWatching returns races where user watched >10min but didn't finish
func (s *RecommendationService) GetContinueWatching(userID string, limit int) ([]models.Race, error) {
	history, err := s.watchHistoryRepo.GetByUserID(userID, limit, 0)
	if err != nil {
		return nil, err
	}

	// Filter for races that are likely not completed
	var raceIDs []string
	for _, entry := range history {
		if !entry.LikelyCompleted {
			raceIDs = append(raceIDs, entry.RaceID)
		}
	}

	if len(raceIDs) == 0 {
		return []models.Race{}, nil
	}

	// Get race details
	var races []models.Race
	for _, raceID := range raceIDs {
		race, err := s.raceRepo.GetByID(raceID)
		if err != nil || race == nil {
			continue
		}
		races = append(races, *race)
	}

	return races, nil
}

// GetUpcomingRacesForUser returns upcoming races matching user preferences
func (s *RecommendationService) GetUpcomingRacesForUser(userID string, limit int) ([]models.Race, error) {
	// Get user favorites
	favorites, err := s.userFavRepo.GetByUserID(userID, nil)
	if err != nil {
		return nil, err
	}

	// Get watch history to find preferred categories
	history, err := s.watchHistoryRepo.GetByUserID(userID, 10, 0)
	if err != nil {
		return nil, err
	}

	// Collect preferred categories from watch history
	categoryMap := make(map[string]bool)
	for _, entry := range history {
		if entry.RaceCategory != nil {
			categoryMap[*entry.RaceCategory] = true
		}
	}

	// Get upcoming races
	upcoming, err := s.raceRepo.GetUpcomingRaces(limit * 2) // Get more to filter
	if err != nil {
		return nil, err
	}

	// Prioritize races from favorite series/categories
	var prioritized []models.Race
	var others []models.Race

	for _, race := range upcoming {
		// Check if race matches favorites (for now, we'll match by name/category)
		matches := false
		for _, fav := range favorites {
			if fav.FavoriteType == "series" || fav.FavoriteType == "race" {
				// Simple name matching (can be enhanced later)
				if race.Name == fav.FavoriteID || (race.Category != nil && *race.Category == fav.FavoriteID) {
					matches = true
					break
				}
			}
		}

		// Check if race matches preferred categories
		if !matches && race.Category != nil && categoryMap[*race.Category] {
			matches = true
		}

		if matches {
			prioritized = append(prioritized, race)
		} else {
			others = append(others, race)
		}
	}

	// Combine prioritized first, then others
	result := append(prioritized, others...)
	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// GetRecommendedReplays returns recommended replays based on watch history
func (s *RecommendationService) GetRecommendedReplays(userID string, limit int) ([]models.Race, error) {
	// Get watch history
	history, err := s.watchHistoryRepo.GetByUserID(userID, 10, 0)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		// No history, return popular races (for now, just upcoming)
		return s.raceRepo.GetUpcomingRaces(limit)
	}

	// Get categories from watch history
	categoryMap := make(map[string]bool)
	for _, entry := range history {
		if entry.RaceCategory != nil {
			categoryMap[*entry.RaceCategory] = true
		}
	}

	// Get similar races for each watched race
	var recommended []models.Race
	seen := make(map[string]bool)

	for _, entry := range history {
		similar, err := s.raceRepo.GetSimilarRaces(entry.RaceID, 3)
		if err != nil {
			continue
		}

		for _, race := range similar {
			if !seen[race.ID] {
				seen[race.ID] = true
				recommended = append(recommended, race)
			}
		}
	}

	// If not enough recommendations, add races from same categories
	if len(recommended) < limit {
		for category := range categoryMap {
			races, err := s.raceRepo.GetRacesByCategory(category, limit)
			if err != nil {
				continue
			}
			for _, race := range races {
				if !seen[race.ID] {
					seen[race.ID] = true
					recommended = append(recommended, race)
					if len(recommended) >= limit {
						break
					}
				}
			}
			if len(recommended) >= limit {
				break
			}
		}
	}

	if len(recommended) > limit {
		recommended = recommended[:limit]
	}

	return recommended, nil
}

// GetAllRecommendations returns all recommendation types
type AllRecommendations struct {
	ContinueWatching []models.Race `json:"continue_watching"`
	Upcoming         []models.Race `json:"upcoming"`
	Replays          []models.Race `json:"replays"`
}

func (s *RecommendationService) GetAllRecommendations(userID string) (*AllRecommendations, error) {
	continueWatching, _ := s.GetContinueWatching(userID, 5)
	upcoming, _ := s.GetUpcomingRacesForUser(userID, 5)
	replays, _ := s.GetRecommendedReplays(userID, 5)

	return &AllRecommendations{
		ContinueWatching: continueWatching,
		Upcoming:         upcoming,
		Replays:          replays,
	}, nil
}

