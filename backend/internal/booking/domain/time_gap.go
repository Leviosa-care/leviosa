package domain

import (
	"time"

	"github.com/google/uuid"
)

// TimeGap represents a time interval between bookings where new availabilities could be created
type TimeGap struct {
	StartTime         time.Time
	EndTime           time.Time
	DurationMinutes   int
	IsBookable        bool         // Whether any products fit in this gap
	SuggestedProducts []uuid.UUID  // Product IDs that fit in this gap
}

// Duration calculates the gap duration in minutes
func (tg *TimeGap) Duration() int {
	return int(tg.EndTime.Sub(tg.StartTime).Minutes())
}

// CanFitProduct checks if a product (with its buffer time) can fit in this gap
func (tg *TimeGap) CanFitProduct(productDuration, bufferTime int) bool {
	requiredMinutes := productDuration + bufferTime
	return tg.DurationMinutes >= requiredMinutes
}

// CanFitMultipleSessions checks if N sessions of a product can fit in this gap
func (tg *TimeGap) CanFitMultipleSessions(productDuration, bufferTime, sessionCount int) bool {
	requiredMinutes := (productDuration + bufferTime) * sessionCount
	return tg.DurationMinutes >= requiredMinutes
}
