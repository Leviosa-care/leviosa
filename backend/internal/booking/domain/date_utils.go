package domain

import "time"

// DateOnly normalizes a time.Time value to midnight UTC (00:00:00.000000000).
// This is used for date-only operations where time precision is not needed,
// matching the PostgreSQL DATE column behavior which stores only year-month-day.
//
// Example:
//
//	input:  2025-12-20 14:35:22.123456789 +0100 CET
//	output: 2025-12-20 00:00:00.000000000 +0000 UTC
func DateOnly(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// DateOnlyPtr is a convenience function that returns a pointer to a DateOnly result.
// Returns nil if the input pointer is nil.
func DateOnlyPtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	normalized := DateOnly(*t)
	return &normalized
}
