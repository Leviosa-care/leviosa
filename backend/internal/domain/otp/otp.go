package otpService

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/hengadev/leviosa/internal/domain"
)

const (
	OTPDURATION    = 15 * time.Minute
	MaxOTPAttempts = 3
)

type OTP struct {
	Email         string    `json:"email" encx:"hash_basic"`
	EmailHash     string    `json:"-"`
	Code          string    `json:"code" validate:"len=6" encx:"encrypt"`
	CodeEncrypted []byte    `json:"-"`
	Attempts      int       `json:"attempts"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	DEK           []byte    `json:"-"`
	DEKEncrypted  []byte    `json:"-"`
	KeyVersion    int       `json:"-"`
}

type Data struct {
	CodeEncrypted []byte    `json:"-"`
	Attempts      int       `json:"attempts"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	DEKEncrypted  []byte    `json:"-"`
	KeyVersion    int       `json:"-"`
}

func (o *OTP) Data() *Data {
	return &Data{
		CodeEncrypted: o.CodeEncrypted,
		Attempts:      o.Attempts,
		ExpiresAt:     o.ExpiresAt,
		CreatedAt:     o.CreatedAt,
		DEKEncrypted:  o.DEKEncrypted,
		KeyVersion:    o.KeyVersion,
	}
}

func (s *service) newOTP(email string) (*OTP, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate secure random number: %w", err)
	}
	num := int(binary.BigEndian.Uint32(bytes) % 100000000)
	dek, err := s.crypto.GenerateDEK()
	if err != nil {
		return nil, fmt.Errorf("failed to generate DEK for OTP: %w", err)
	}
	return &OTP{
		Email:     email,
		Code:      fmt.Sprintf("%06d", num),
		Attempts:  1,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(OTPDURATION),
		DEK:       dek,
	}, nil
}

func (o *OTP) IncreaseAttempt() error {
	if o.Attempts+1 >= MaxOTPAttempts {
		return domain.NewInvalidValueErr("max attempts reached for provided OTP")
	}
	o.Attempts++
	return nil
}
