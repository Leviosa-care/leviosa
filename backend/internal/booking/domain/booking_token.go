package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"

	"github.com/google/uuid"
)

// BookingTokenGracePeriod is the time after slot_end_time during which the
// booking token remains valid. Allows guests to access their booking for
// 30 days after the appointment.
const BookingTokenGracePeriod = 30 * 24 * time.Hour

// GenerateBookingToken creates a signed, time-limited token scoped to a
// single booking. The token encodes the booking ID and expiry timestamp
// (slot_end_time + 30-day grace period), then appends an HMAC-SHA256
// signature over both fields.
//
// The token is URL-safe base64 (no padding) and can be embedded in links.
func GenerateBookingToken(bookingID uuid.UUID, slotEndTime time.Time, secret []byte) string {
	// Token layout: UUID (16 bytes) | expiry_unix (8 bytes) | HMAC-SHA256 (32 bytes)
	buf := make([]byte, 0, 56)

	// 1. Booking ID (16 bytes)
	buf = append(buf, bookingID[:]...)

	// 2. Expiry timestamp = slotEndTime + grace period (8 bytes, big-endian int64)
	expiry := slotEndTime.Add(BookingTokenGracePeriod).Unix()
	buf = binary.BigEndian.AppendUint64(buf, uint64(expiry))

	// 3. HMAC-SHA256 signature over booking ID + expiry
	mac := hmac.New(sha256.New, secret)
	mac.Write(buf)
	sig := mac.Sum(nil)
	buf = append(buf, sig...)

	return base64.RawURLEncoding.EncodeToString(buf)
}

// VerifyBookingToken validates a booking token's signature and expiry.
// Returns the booking ID on success, or a typed error on failure.
func VerifyBookingToken(token string, secret []byte) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, ErrInvalidBookingToken
	}

	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return uuid.Nil, ErrInvalidBookingToken
	}

	// Must be exactly 56 bytes: 16 (UUID) + 8 (expiry) + 32 (HMAC)
	if len(raw) != 56 {
		return uuid.Nil, ErrInvalidBookingToken
	}

	payload := raw[:24] // UUID + expiry
	signature := raw[24:] // HMAC-SHA256

	// Verify HMAC signature
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	expectedSig := mac.Sum(nil)
	if !hmac.Equal(signature, expectedSig) {
		return uuid.Nil, ErrInvalidBookingToken
	}

	// Extract booking ID
	bookingID, err := uuid.FromBytes(payload[:16])
	if err != nil {
		return uuid.Nil, ErrInvalidBookingToken
	}

	// Check expiry
	expiryUnix := int64(binary.BigEndian.Uint64(payload[16:24]))
	expiry := time.Unix(expiryUnix, 0)
	if time.Now().After(expiry) {
		return uuid.Nil, ErrBookingTokenExpired
	}

	return bookingID, nil
}

// IsBookingTokenError returns true if the error is a booking token error
// (either invalid or expired).
func IsBookingTokenError(err error) bool {
	return errors.Is(err, ErrInvalidBookingToken) || errors.Is(err, ErrBookingTokenExpired)
}
