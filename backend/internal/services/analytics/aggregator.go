package analytics

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/cyclingstream/backend/internal/models"
	"github.com/cyclingstream/backend/internal/repository"
)

const (
	defaultHeartbeatSeconds   = 15
	defaultSessionIdleMinutes = 30
)

// Aggregator builds sessions and per-stream stats from playback events.
type Aggregator struct {
	playbackRepo   *repository.PlaybackEventRepository
	statsRepo      *repository.StreamStatsRepository
	streamRepo     *repository.StreamRepository
	heartbeatSecs  int
	idleTimeoutMin int
}

func NewAggregator(
	playbackRepo *repository.PlaybackEventRepository,
	statsRepo *repository.StreamStatsRepository,
	streamRepo *repository.StreamRepository,
) *Aggregator {
	return &Aggregator{
		playbackRepo:   playbackRepo,
		statsRepo:      statsRepo,
		streamRepo:     streamRepo,
		heartbeatSecs:  defaultHeartbeatSeconds,
		idleTimeoutMin: defaultSessionIdleMinutes,
	}
}

type session struct {
	clientID          string
	country           string
	device            string
	startedAt         time.Time
	endedAt           time.Time
	lastEventAt       time.Time
	totalWatchSeconds int64
	bufferSeconds     int64
	errorCount        int
	bufferStart       *time.Time
}

// AggregateStream computes stats for a single stream and persists them.
func (a *Aggregator) AggregateStream(ctx context.Context, streamID string, since *time.Time) (*models.StreamStats, error) {
	events, err := a.playbackRepo.ListByStreamSince(ctx, streamID, since)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("no events for stream %s", streamID)
	}

	sessions := a.buildSessions(events)
	stats := a.computeStats(streamID, sessions)
	stats.LastCalculatedAt = time.Now()

	if err := a.statsRepo.Upsert(ctx, stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func (a *Aggregator) buildSessions(events []models.PlaybackEvent) []session {
	var sessions []session
	currentByClient := make(map[string]*session)
	idle := time.Duration(a.idleTimeoutMin) * time.Minute

	for _, evt := range events {
		sess := currentByClient[evt.ClientID]
		if sess == nil || evt.CreatedAt.Sub(sess.lastEventAt) > idle {
			if sess != nil {
				if sess.bufferStart != nil {
					sess.bufferSeconds += int64(sess.lastEventAt.Sub(*sess.bufferStart).Seconds())
					sess.bufferStart = nil
				}
				sess.endedAt = sess.lastEventAt
				sessions = append(sessions, *sess)
			}
			sess = &session{
				clientID:    evt.ClientID,
				country:     evt.Country,
				device:      evt.DeviceType,
				startedAt:   evt.CreatedAt,
				endedAt:     evt.CreatedAt,
				lastEventAt: evt.CreatedAt,
			}
			currentByClient[evt.ClientID] = sess
		}

		if evt.EventType == "heartbeat" || evt.EventType == "ended" || evt.EventType == "play" {
			sess.totalWatchSeconds += int64(a.heartbeatSecs)
		}

		if evt.EventType == "buffer_start" {
			t := evt.CreatedAt
			sess.bufferStart = &t
		}
		if evt.EventType == "buffer_end" && sess.bufferStart != nil {
			sess.bufferSeconds += int64(evt.CreatedAt.Sub(*sess.bufferStart).Seconds())
			sess.bufferStart = nil
		}
		if evt.EventType == "error" {
			sess.errorCount++
		}

		sess.lastEventAt = evt.CreatedAt
	}

	for _, sess := range currentByClient {
		if sess.bufferStart != nil {
			sess.bufferSeconds += int64(sess.lastEventAt.Sub(*sess.bufferStart).Seconds())
			sess.bufferStart = nil
		}
		sess.endedAt = sess.lastEventAt
		sessions = append(sessions, *sess)
	}

	return sessions
}

func (a *Aggregator) computeStats(streamID string, sessions []session) *models.StreamStats {
	uniqueViewers := len(sessions)
	var totalWatch int64
	var totalBuffer int64
	countryCounts := make(map[string]int)
	deviceCounts := make(map[string]int)

	type point struct {
		t time.Time
		d int
	}
	var points []point

	for _, s := range sessions {
		totalWatch += s.totalWatchSeconds
		totalBuffer += s.bufferSeconds
		countryCounts[s.country]++
		deviceCounts[s.device]++

		points = append(points, point{t: s.startedAt, d: 1})
		points = append(points, point{t: s.endedAt.Add(time.Second), d: -1}) // add a small delta to close interval
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].t.Before(points[j].t)
	})

	current := 0
	peak := 0
	for _, p := range points {
		current += p.d
		if current > peak {
			peak = current
		}
	}

	avg := 0
	if uniqueViewers > 0 {
		avg = int(totalWatch / int64(uniqueViewers))
	}

	errorSessions := 0
	for _, s := range sessions {
		if s.errorCount > 0 {
			errorSessions++
		}
	}

	bufferRatio := 0.0
	if totalWatch > 0 {
		bufferRatio = float64(totalBuffer) / float64(totalWatch)
	}

	errorRate := 0.0
	if uniqueViewers > 0 {
		errorRate = float64(errorSessions) / float64(uniqueViewers)
	}

	return &models.StreamStats{
		StreamID:              streamID,
		UniqueViewers:         uniqueViewers,
		TotalWatchSeconds:     totalWatch,
		AvgWatchSeconds:       avg,
		PeakConcurrentViewers: peak,
		TopCountries:          countryCounts,
		DeviceBreakdown:       deviceCounts,
		BufferSeconds:         totalBuffer,
		BufferRatio:           bufferRatio,
		ErrorRate:             errorRate,
	}
}
