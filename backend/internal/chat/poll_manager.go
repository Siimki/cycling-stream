package chat

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// PollOption represents an individual poll option.
type PollOption struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Votes int    `json:"votes"`
}

// Poll represents an active or closed poll in chat.
type Poll struct {
	ID         string       `json:"id"`
	RaceID     string       `json:"race_id"`
	Question   string       `json:"question"`
	Options    []PollOption `json:"options"`
	TotalVotes int          `json:"total_votes"`
	CreatedAt  time.Time    `json:"created_at"`
	ClosesAt   *time.Time   `json:"closes_at,omitempty"`
	Closed     bool         `json:"closed"`
}

// PollManager keeps polls in memory per race.
type PollManager struct {
	mu        sync.RWMutex
	polls     map[string]*Poll
	raceIndex map[string]string            // raceID -> pollID
	votes     map[string]map[string]string // pollID -> userID -> optionID
}

func NewPollManager() *PollManager {
	return &PollManager{
		polls:     make(map[string]*Poll),
		raceIndex: make(map[string]string),
		votes:     make(map[string]map[string]string),
	}
}

// CreatePoll starts a new poll for a race, replacing any active poll.
func (pm *PollManager) CreatePoll(raceID, question string, optionLabels []string, duration time.Duration) (*Poll, error) {
	if len(optionLabels) < 2 {
		return nil, errors.New("poll requires at least two options")
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	poll := &Poll{
		ID:        uuid.NewString(),
		RaceID:    raceID,
		Question:  question,
		Options:   make([]PollOption, 0, len(optionLabels)),
		CreatedAt: time.Now().UTC(),
	}

	for _, label := range optionLabels {
		label = strings.TrimSpace(label)
		if label == "" {
			continue
		}
		poll.Options = append(poll.Options, PollOption{
			ID:    uuid.NewString(),
			Label: label,
		})
	}

	if len(poll.Options) < 2 {
		return nil, errors.New("poll requires at least two valid options")
	}

	if duration > 0 {
		closesAt := poll.CreatedAt.Add(duration)
		poll.ClosesAt = &closesAt
	}

	// Close existing poll in this race
	if existingID, ok := pm.raceIndex[raceID]; ok {
		if existing, exists := pm.polls[existingID]; exists {
			existing.Closed = true
		}
	}

	pm.polls[poll.ID] = poll
	pm.raceIndex[raceID] = poll.ID
	pm.votes[poll.ID] = make(map[string]string)

	return poll, nil
}

// Vote records a vote for the given user/option.
func (pm *PollManager) Vote(pollID, userID, optionID string) (*Poll, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	poll, ok := pm.polls[pollID]
	if !ok || poll.Closed {
		return nil, errors.New("poll not active")
	}

	if userID == "" {
		return nil, errors.New("user_id required")
	}

	optionIndex := -1
	for idx, opt := range poll.Options {
		if opt.ID == optionID {
			optionIndex = idx
			break
		}
	}

	if optionIndex == -1 {
		return nil, errors.New("invalid option")
	}

	if pm.votes[pollID] == nil {
		pm.votes[pollID] = make(map[string]string)
	}

	if prevOption, voted := pm.votes[pollID][userID]; voted {
		if prevOption == optionID {
			return poll, nil
		}
		for i := range poll.Options {
			if poll.Options[i].ID == prevOption && poll.Options[i].Votes > 0 {
				poll.Options[i].Votes--
				poll.TotalVotes--
				break
			}
		}
	}

	pm.votes[pollID][userID] = optionID
	poll.Options[optionIndex].Votes++
	poll.TotalVotes++

	return poll, nil
}

// ClosePoll terminates an active poll and removes it from the race index.
func (pm *PollManager) ClosePoll(pollID string) (*Poll, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	poll, ok := pm.polls[pollID]
	if !ok {
		return nil, errors.New("poll not found")
	}

	poll.Closed = true
	if currentID, ok := pm.raceIndex[poll.RaceID]; ok && currentID == poll.ID {
		delete(pm.raceIndex, poll.RaceID)
	}
	delete(pm.votes, poll.ID)

	return poll, nil
}

// GetActivePoll returns the active poll for a race, if any.
func (pm *PollManager) GetActivePoll(raceID string) *Poll {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if pollID, ok := pm.raceIndex[raceID]; ok {
		if poll, exists := pm.polls[pollID]; exists && !poll.Closed {
			return poll
		}
	}
	return nil
}

// GetPoll returns the poll by ID if present.
func (pm *PollManager) GetPoll(pollID string) (*Poll, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	poll, ok := pm.polls[pollID]
	return poll, ok
}
