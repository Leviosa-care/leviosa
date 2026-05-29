package otp

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockNotifier records calls to SendOTPEmail for verification.
type mockNotifier struct {
	calls []notifyCall
	err   error
}

type notifyCall struct {
	Email string
	OTP   string
}

func (m *mockNotifier) SendOTPEmail(_ context.Context, email, otp string) error {
	m.calls = append(m.calls, notifyCall{Email: email, OTP: otp})
	return m.err
}

var _ ports.NotificationService = (*mockNotifier)(nil)

// testCrypto is a lightweight encx.CryptoService implementation for unit tests.
// It uses AES-GCM for encryption and SHA-256 for hashing — no Vault dependency.
type testCrypto struct {
	key []byte // 32-byte AES key
}

func newTestCrypto(t *testing.T) encx.CryptoService {
	t.Helper()
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	return &testCrypto{key: key}
}

func (c *testCrypto) GetPepper() []byte                                          { return nil }
func (c *testCrypto) GetArgon2Params() *encx.Argon2Params                         { return nil }
func (c *testCrypto) GetAlias() string                                            { return "test" }
func (c *testCrypto) RotateKEK(_ context.Context) error                           { return nil }
func (c *testCrypto) CompareSecureHashAndValue(_ context.Context, _ any, _ string) (bool, error) {
	return false, nil
}
func (c *testCrypto) CompareBasicHashAndValue(_ context.Context, _ any, _ string) (bool, error) {
	return false, nil
}
func (c *testCrypto) EncryptStream(_ context.Context, _ io.Reader, _ io.Writer, _ []byte) error {
	return nil
}
func (c *testCrypto) DecryptStream(_ context.Context, _ io.Reader, _ io.Writer, _ []byte) error {
	return nil
}
func (c *testCrypto) HashSecure(_ context.Context, value []byte) (string, error) {
	h := sha256.Sum256(append(value, c.key...))
	return hex.EncodeToString(h[:]), nil
}

func (c *testCrypto) HashBasic(_ context.Context, data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func (c *testCrypto) GenerateDEK() ([]byte, error) {
	dek := make([]byte, 32)
	if _, err := rand.Read(dek); err != nil {
		return nil, err
	}
	return dek, nil
}

func (c *testCrypto) EncryptData(_ context.Context, plaintext []byte, dek []byte) ([]byte, error) {
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *testCrypto) DecryptData(_ context.Context, ciphertext []byte, dek []byte) ([]byte, error) {
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	return gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
}

func (c *testCrypto) EncryptDEK(_ context.Context, plaintextDEK []byte) ([]byte, error) {
	return c.EncryptData(context.Background(), plaintextDEK, c.key)
}

func (c *testCrypto) DecryptDEKWithVersion(_ context.Context, ciphertextDEK []byte, _ int) ([]byte, error) {
	return c.DecryptData(context.Background(), ciphertextDEK, c.key)
}

func (c *testCrypto) GetCurrentKEKVersion(_ context.Context, _ string) (int, error) {
	return 1, nil
}

func (c *testCrypto) GetKMSKeyIDForVersion(_ context.Context, _ string, _ int) (string, error) {
	return "test-key-id", nil
}

// Compile-time check
var _ encx.CryptoService = (*testCrypto)(nil)

// --- Unit tests ---

func TestRequestOTP_CallsNotificationService(t *testing.T) {
	ctx := context.Background()

	t.Run("sends OTP email via notification service on RequestOTP", func(t *testing.T) {
		repo := newInMemRepo()
		notif := &mockNotifier{}
		svc := &OTPService{
			repo:            repo,
			crypto:          newTestCrypto(t),
			cache:           NewMockOTPCache(defaultOTPLength, defaultOTPDuration, defaultOTPMaxAttempts),
			notificationSvc: notif,
		}

		err := svc.RequestOTP(ctx, "user@example.com")
		require.NoError(t, err)

		require.Len(t, notif.calls, 1)
		assert.Equal(t, "user@example.com", notif.calls[0].Email)
		assert.Len(t, notif.calls[0].OTP, defaultOTPLength)
	})

	t.Run("returns error when notification service fails", func(t *testing.T) {
		repo := newInMemRepo()
		notif := &mockNotifier{
			err: errs.NewExternalServiceErr(assert.AnError, "send OTP email"),
		}
		svc := &OTPService{
			repo:            repo,
			crypto:          newTestCrypto(t),
			cache:           NewMockOTPCache(defaultOTPLength, defaultOTPDuration, defaultOTPMaxAttempts),
			notificationSvc: notif,
		}

		err := svc.RequestOTP(ctx, "fail@example.com")
		assert.Error(t, err)
	})
}

func TestResendOTP_CallsNotificationService(t *testing.T) {
	ctx := context.Background()

	t.Run("sends OTP email via notification service on ResendOTP", func(t *testing.T) {
		repo := newInMemRepo()
		notif := &mockNotifier{}
		svc := &OTPService{
			repo:            repo,
			crypto:          newTestCrypto(t),
			cache:           NewMockOTPCache(defaultOTPLength, defaultOTPDuration, defaultOTPMaxAttempts),
			notificationSvc: notif,
		}

		// Create initial OTP
		err := svc.RequestOTP(ctx, "resend@example.com")
		require.NoError(t, err)
		require.Len(t, notif.calls, 1)

		notif.calls = nil

		// Resend
		err = svc.ResendOTP(ctx, "resend@example.com")
		require.NoError(t, err)

		require.Len(t, notif.calls, 1)
		assert.Equal(t, "resend@example.com", notif.calls[0].Email)
		assert.NotEmpty(t, notif.calls[0].OTP)
	})
}

func TestCreateOTP_CallsNotificationService(t *testing.T) {
	ctx := context.Background()

	t.Run("sends OTP email via notification service on CreateOTP", func(t *testing.T) {
		repo := newInMemRepo()
		notif := &mockNotifier{}
		svc := &OTPService{
			repo:            repo,
			crypto:          newTestCrypto(t),
			cache:           NewMockOTPCache(defaultOTPLength, defaultOTPDuration, defaultOTPMaxAttempts),
			notificationSvc: notif,
		}

		err := svc.CreateOTP(ctx, "create@example.com")
		require.NoError(t, err)

		require.Len(t, notif.calls, 1)
		assert.Equal(t, "create@example.com", notif.calls[0].Email)
		assert.Len(t, notif.calls[0].OTP, defaultOTPLength)
	})
}

// --- in-memory repo for unit tests ---

type inMemRepo struct {
	data map[string][]byte
}

func newInMemRepo() *inMemRepo {
	return &inMemRepo{data: make(map[string][]byte)}
}

func (r *inMemRepo) SaveOTP(_ context.Context, key string, value []byte, _ time.Duration) error {
	r.data[key] = value
	return nil
}

func (r *inMemRepo) GetOTP(_ context.Context, key string) ([]byte, error) {
	v, ok := r.data[key]
	if !ok {
		return nil, errs.ErrRepositoryNotFound
	}
	return v, nil
}

func (r *inMemRepo) InvalidateOTP(_ context.Context, key string) error {
	delete(r.data, key)
	return nil
}

func (r *inMemRepo) TouchOTP(_ context.Context, _ string, _ time.Duration) error {
	return nil
}
