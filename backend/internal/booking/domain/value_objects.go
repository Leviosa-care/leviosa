package domain

import (
	"time"
)

// Date represents a date value object
type Date struct {
	time.Time
}

// NewDate creates a new Date
func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// DateTime represents a datetime value object
type DateTime struct {
	time.Time
}

// NewDateTime creates a new DateTime
func NewDateTime(t time.Time) DateTime {
	return DateTime{t}
}
