package domain

import (
	"time"
)

type OTP struct {
	Email     string    `json:"-" encx:"hash_basic"`
	Code      string    `json:"-" validate:"len=6" encx:"encrypt"`
	Attempts  int       `json:"attempts"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTP) IncrementAttempts() {
	o.Attempts++
}
