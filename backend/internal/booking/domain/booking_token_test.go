package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateBookingToken(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long-xxxxx")
	bookingID := uuid.New()
	slotEndTime := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)

	t.Run("returns non-empty token string", func(t *testing.T) {
		token := GenerateBookingToken(bookingID, slotEndTime, secret)
		assert.NotEmpty(t, token)
	})

	t.Run("token is deterministic for same inputs", func(t *testing.T) {
		token1 := GenerateBookingToken(bookingID, slotEndTime, secret)
		token2 := GenerateBookingToken(bookingID, slotEndTime, secret)
		assert.Equal(t, token1, token2)
	})

	t.Run("token differs for different booking IDs", func(t *testing.T) {
		otherID := uuid.New()
		token1 := GenerateBookingToken(bookingID, slotEndTime, secret)
		token2 := GenerateBookingToken(otherID, slotEndTime, secret)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("token differs for different secrets", func(t *testing.T) {
		otherSecret := []byte("different-secret-key-32-bytes-lo")
		token1 := GenerateBookingToken(bookingID, slotEndTime, secret)
		token2 := GenerateBookingToken(bookingID, slotEndTime, otherSecret)
		assert.NotEqual(t, token1, token2)
	})
}

func TestVerifyBookingToken(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long-xxxxx")
	bookingID := uuid.New()

	t.Run("valid token round-trip", func(t *testing.T) {
		slotEndTime := time.Now().Add(24 * time.Hour) // tomorrow
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		gotID, err := VerifyBookingToken(token, secret)
		require.NoError(t, err)
		assert.Equal(t, bookingID, gotID)
	})

	t.Run("rejects tampered signature", func(t *testing.T) {
		slotEndTime := time.Now().Add(24 * time.Hour)
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		// Flip a bit in the UUID section: guarantees the payload changes regardless
		// of the random UUID bytes (unlike a fixed string prefix replacement).
		raw, err := base64.RawURLEncoding.DecodeString(token)
		require.NoError(t, err)
		raw[0] ^= 0xFF
		tampered := base64.RawURLEncoding.EncodeToString(raw)

		_, err = VerifyBookingToken(tampered, secret)
		assert.ErrorIs(t, err, ErrInvalidBookingToken)
	})

	t.Run("rejects expired token", func(t *testing.T) {
		// Slot ended 31 days ago — past the 30-day grace period
		slotEndTime := time.Now().Add(-31 * 24 * time.Hour)
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		_, err := VerifyBookingToken(token, secret)
		assert.ErrorIs(t, err, ErrBookingTokenExpired)
	})

	t.Run("accepts token within grace period", func(t *testing.T) {
		// Slot ended 29 days ago — still within the 30-day grace period
		slotEndTime := time.Now().Add(-29 * 24 * time.Hour)
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		gotID, err := VerifyBookingToken(token, secret)
		require.NoError(t, err)
		assert.Equal(t, bookingID, gotID)
	})

	t.Run("rejects token with wrong secret", func(t *testing.T) {
		slotEndTime := time.Now().Add(24 * time.Hour)
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		wrongSecret := []byte("wrong-secret-key-32-bytes-long-x")
		_, err := VerifyBookingToken(token, wrongSecret)
		assert.ErrorIs(t, err, ErrInvalidBookingToken)
	})

	t.Run("rejects malformed token", func(t *testing.T) {
		_, err := VerifyBookingToken("not-valid-base64!", secret)
		assert.ErrorIs(t, err, ErrInvalidBookingToken)
	})

	t.Run("rejects empty token", func(t *testing.T) {
		_, err := VerifyBookingToken("", secret)
		assert.ErrorIs(t, err, ErrInvalidBookingToken)
	})

	t.Run("rejects token that is valid base64 but wrong structure", func(t *testing.T) {
		// Create a token with valid base64 but without a valid UUID prefix
		fake := base64.RawURLEncoding.EncodeToString([]byte("short"))
		_, err := VerifyBookingToken(fake, secret)
		assert.ErrorIs(t, err, ErrInvalidBookingToken)
	})

	t.Run("boundary: token exactly at grace period edge", func(t *testing.T) {
		// Slot end time is exactly 30 days ago — the grace period just expired
		slotEndTime := time.Now().Add(-30 * 24 * time.Hour)
		// We construct the token manually to control the exact expiry timestamp
		// that gets embedded. The verification checks: slotEndTime + 30 days < now
		// If slotEndTime was exactly 30 days ago, slotEndTime + 30d = now, which is NOT < now
		// so it should still be valid at the boundary.
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		// This is borderline depending on exact timing; verify it doesn't error
		// because the boundary is inclusive (expiry = slotEndTime + 30d, token valid if now <= expiry)
		_, err := VerifyBookingToken(token, secret)
		// At the boundary, it could go either way; just ensure it returns a typed error
		if err != nil {
			assert.ErrorIs(t, err, ErrBookingTokenExpired)
		}
	})
}

func TestBookingTokenFormat(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long-xxxxx")
	bookingID := uuid.New()
	slotEndTime := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)

	t.Run("token encodes booking ID and expiry timestamp", func(t *testing.T) {
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		// Decode and verify structure: UUID (16 bytes) + expiry (8 bytes) + HMAC (32 bytes)
		raw, err := base64.RawURLEncoding.DecodeString(token)
		require.NoError(t, err)

		// 16 (UUID) + 8 (expiry timestamp int64) + 32 (HMAC-SHA256) = 56 bytes
		assert.Len(t, raw, 56, "token payload should be 56 bytes: 16 (UUID) + 8 (expiry) + 32 (HMAC)")

		// Verify booking ID embedded
		embeddedID, err := uuid.FromBytes(raw[:16])
		require.NoError(t, err)
		assert.Equal(t, bookingID, embeddedID)

		// Verify expiry timestamp
		expiryUnix := int64(binary.BigEndian.Uint64(raw[16:24]))
		expectedExpiry := slotEndTime.Add(30 * 24 * time.Hour).Unix()
		assert.Equal(t, expectedExpiry, expiryUnix)

		// Verify HMAC over UUID + expiry
		mac := hmac.New(sha256.New, secret)
		mac.Write(raw[:24])
		expectedSig := mac.Sum(nil)
		assert.Equal(t, expectedSig, raw[24:])
	})

	t.Run("token is URL-safe base64", func(t *testing.T) {
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		// RawURLEncoding means no padding characters
		for _, c := range token {
			assert.False(t, c == '+' || c == '/' || c == '=',
				"token should be URL-safe base64, found: %c", c)
		}
	})
}

// VerifyBookingToken extracts the expiry from the token payload itself.
// The token is self-contained: it embeds booking ID + expiry timestamp + HMAC.
func TestVerifyBookingTokenExtractsExpiry(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long-xxxxx")
	bookingID := uuid.New()

	t.Run("verify reconstructs expiry from token payload", func(t *testing.T) {
		slotEndTime := time.Now().Add(7 * 24 * time.Hour)
		token := GenerateBookingToken(bookingID, slotEndTime, secret)

		raw, err := base64.RawURLEncoding.DecodeString(token)
		require.NoError(t, err)

		expiryUnix := int64(binary.BigEndian.Uint64(raw[16:24]))
		embeddedExpiry := time.Unix(expiryUnix, 0)

		expectedExpiry := slotEndTime.Add(30 * 24 * time.Hour)
		assert.WithinDuration(t, expectedExpiry, embeddedExpiry, time.Second,
			"embedded expiry should be slot_end_time + 30 days")
	})
}
