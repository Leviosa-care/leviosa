package domain

import (
	"time"

	"github.com/google/uuid"
)

// RoomMetrics represents daily utilization statistics for a room
type RoomMetrics struct {
	RoomID             uuid.UUID
	Date               time.Time
	TotalMinutesOpen   int     // Operating hours (e.g., 480 minutes for 8-hour day)
	TotalMinutesBooked int     // Sum of booked availabilities
	UtilizationPercent float64 // (Booked / Open) * 100
	FragmentationCount int     // Number of gaps too small to be useful (< 30 min)
	IdleMinutes        int     // Total minutes in unusable gaps
	AverageGapMinutes  int     // Average size of gaps between bookings
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// CalculateUtilization computes the utilization percentage
// based on booked minutes vs operating hours
func (rm *RoomMetrics) CalculateUtilization() {
	if rm.TotalMinutesOpen > 0 {
		rm.UtilizationPercent = (float64(rm.TotalMinutesBooked) / float64(rm.TotalMinutesOpen)) * 100
	} else {
		rm.UtilizationPercent = 0
	}
}

// IsHighlyUtilized returns true if utilization is above 75%
func (rm *RoomMetrics) IsHighlyUtilized() bool {
	return rm.UtilizationPercent >= 75.0
}

// IsFragmented returns true if there are many small unusable gaps
func (rm *RoomMetrics) IsFragmented() bool {
	return rm.FragmentationCount > 3
}

// EfficiencyScore calculates an overall efficiency metric (0-100)
// considering both utilization and fragmentation
func (rm *RoomMetrics) EfficiencyScore() float64 {
	// Start with utilization
	score := rm.UtilizationPercent

	// Penalize for fragmentation (each fragment reduces score by 2%)
	fragmentationPenalty := float64(rm.FragmentationCount) * 2.0
	score -= fragmentationPenalty

	// Penalize for idle time (1% per 30 idle minutes)
	idlePenalty := float64(rm.IdleMinutes) / 30.0
	score -= idlePenalty

	// Clamp to 0-100 range
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}
