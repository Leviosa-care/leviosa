package domain

import (
	"time"
)

type OTP struct {
	Email         string    `json:"-" encx:"hash_basic"`
	EmailHash     string    `json:"-"`
	Code          string    `json:"-" validate:"len=6" encx:"encrypt"`
	CodeEncrypted []byte    `json:"code_encrypted"`
	Attempts      int       `json:"attempts"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	DEK           []byte    `json:"-"`
	DEKEncrypted  []byte    `json:"dek_encrypted"`
	KeyVersion    int       `json:"key_version"`
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTP) IncrementAttempts() {
	o.Attempts++
}
