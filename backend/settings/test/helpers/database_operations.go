package helpers

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/settings/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	EncryptedDataExists = "encrypted_data_exists"
)

func InsertSettingString(t *testing.T, ctx context.Context, setting *domain.Setting[string], pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at;
	`
	err := pool.QueryRow(ctx, query, setting.Key, setting.Value).Scan(
		&setting.ID,
		setting.CreatedAt,
		setting.UpdatedAt,
	)
	require.NoError(t, err)
}

func InsertSettingInt(t *testing.T, ctx context.Context, setting *domain.Setting[int], pool *pgxpool.Pool) {
	t.Helper()
	valueStr := fmt.Sprintf("%d", setting.Value)
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at;
	`
	err := pool.QueryRow(ctx, query, setting.Key, valueStr).Scan(
		&setting.ID,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)
	require.NoError(t, err)
}

func InsertEncryptedSettingString(t *testing.T, ctx context.Context, setting *domain.SettingEncryptedEncx, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.encrypted (key, value_encrypted, dek_encrypted, key_version, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at;
	`
	err := pool.QueryRow(
		ctx,
		query,
		setting.Key,
		setting.ValueEncrypted,
		setting.DEKEncrypted,
		setting.KeyVersion,
		setting.Metadata,
	).Scan(
		&setting.ID,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)

	require.NoError(t, err)
}

func GetSettingIntByKey(t *testing.T, ctx context.Context, key string, pool *pgxpool.Pool) *domain.Setting[int] {
	t.Helper()
	var valueStr string
	var setting domain.Setting[int]

	query := `
	SELECT
	id,
	value,
	created_at,
	updated_at
	FROM settings.plain
	WHERE key = $1;
	`
	err := pool.QueryRow(ctx, query, key).Scan(
		&setting.ID,
		&valueStr,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)
	require.NoError(t, err)

	value, err := strconv.Atoi(valueStr)
	require.NoError(t, err)

	setting.Key = key
	setting.Value = value
	return &setting
}

func GetSettingStringByKey(t *testing.T, ctx context.Context, key string, pool *pgxpool.Pool) *domain.Setting[string] {
	t.Helper()
	var setting domain.Setting[string]

	query := `
	SELECT
	id,
	value,
	created_at,
	updated_at
	FROM settings.plain
	WHERE key = $1;
	`
	err := pool.QueryRow(ctx, query, key).Scan(
		&setting.ID,
		&setting.Value,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)
	require.NoError(t, err)

	setting.Key = key
	return &setting
}

// Company Name specific helpers
func InsertCompanyName(t *testing.T, ctx context.Context, name string, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, settings.CompanyName, name)
	require.NoError(t, err)
}

func GetCompanyNameFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (string, error) {
	t.Helper()
	var name string
	query := `SELECT value FROM settings.plain WHERE key = $1`
	err := pool.QueryRow(ctx, query, settings.CompanyName).Scan(&name)
	return name, err
}

// Company Email specific helpers
func InsertCompanyEmail(t *testing.T, ctx context.Context, email string, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, settings.CompanyEmail, email)
	require.NoError(t, err)
}

func GetCompanyEmailFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (string, error) {
	t.Helper()
	var email string
	query := `SELECT value FROM settings.plain WHERE key = $1`
	err := pool.QueryRow(ctx, query, settings.CompanyEmail).Scan(&email)
	return email, err
}

// OTP Duration specific helpers
func InsertOTPDuration(t *testing.T, ctx context.Context, duration int, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, settings.OTPDuration, fmt.Sprintf("%d", duration))
	require.NoError(t, err)
}

func GetOTPDurationFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, error) {
	t.Helper()
	var durationStr string
	query := `SELECT value FROM settings.plain WHERE key = $1`
	err := pool.QueryRow(ctx, query, settings.OTPDuration).Scan(&durationStr)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(durationStr)
}

func GetEncryptedSettingByKey(t *testing.T, ctx context.Context, key string, pool *pgxpool.Pool) *domain.SettingEncryptedEncx {
	t.Helper()
	var setting domain.SettingEncryptedEncx

	query := `
	SELECT
	id,
	value_encrypted,
	created_at,
	updated_at,
	dek_encrypted,
	key_version
	FROM settings.encrypted
	WHERE key = $1;
	`
	err := pool.QueryRow(ctx, query, key).Scan(
		&setting.ID,
		&setting.ValueEncrypted,
		&setting.CreatedAt,
		&setting.UpdatedAt,
		&setting.DEKEncrypted,
		&setting.KeyVersion,
	)
	require.NoError(t, err)

	setting.Key = key
	return &setting
}

// Company Address specific helpers
func InsertCompanyAddress(t *testing.T, ctx context.Context, address string, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, settings.CompanyLegalAddress, address)
	require.NoError(t, err)
}

func GetCompanyAddressFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (string, error) {
	t.Helper()
	var address string
	query := `SELECT value FROM settings.plain WHERE key = $1`
	err := pool.QueryRow(ctx, query, settings.CompanyLegalAddress).Scan(&address)
	return address, err
}

func NewCompanyPhone(t *testing.T, ctx context.Context) *domain.SettingEncrypted {
	now := time.Now()
	return &domain.SettingEncrypted{
		Key:       settings.CompanyPhone,
		Value:     "0612345678",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Company Phone specific helpers - now uses generated ENCX structs
func InsertCompanyPhoneEncrypted(t *testing.T, ctx context.Context, settingEncx *domain.SettingEncryptedEncx, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.encrypted (key, value_encrypted, dek_encrypted, key_version, metadata)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (key) DO UPDATE SET
			value_encrypted = EXCLUDED.value_encrypted,
			dek_encrypted = EXCLUDED.dek_encrypted,
			key_version = EXCLUDED.key_version,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
	`

	_, err := pool.Exec(ctx, query, settings.CompanyPhone, settingEncx.ValueEncrypted, settingEncx.DEKEncrypted, settingEncx.KeyVersion, settingEncx.Metadata)
	require.NoError(t, err)
}

func InsertEncryptedSetting(t *testing.T, ctx context.Context, key string, valueEncrypted []byte, dekEncrypted []byte, keyVersion int, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.encrypted (key, value_encrypted, dek_encrypted, key_version)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key) DO UPDATE SET
			value_encrypted = EXCLUDED.value_encrypted,
			dek_encrypted = EXCLUDED.dek_encrypted,
			key_version = EXCLUDED.key_version,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, key, valueEncrypted, dekEncrypted, keyVersion)
	require.NoError(t, err)
}

func GetEncryptedSettingFromDB(t *testing.T, ctx context.Context, key string, pool *pgxpool.Pool) *domain.SettingEncryptedEncx {
	t.Helper()

	var setting domain.SettingEncryptedEncx

	query := `
	SELECT
	id,
	value_encrypted,
	created_at,
	updated_at,
	dek_encrypted,
	key_version
	FROM settings.encrypted
	WHERE key = $1;`

	err := pool.QueryRow(ctx, query, key).Scan(
		&setting.ID,
		&setting.ValueEncrypted,
		&setting.CreatedAt,
		&setting.UpdatedAt,
		&setting.DEKEncrypted,
		&setting.KeyVersion,
	)
	require.NoError(t, err)
	return &setting
}

// Company Instagram specific helpers
func InsertCompanyInstagram(t *testing.T, ctx context.Context, instagram string, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, settings.CompanyInstagram, instagram)
	require.NoError(t, err)
}

func GetCompanyInstagramFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (string, error) {
	t.Helper()
	var instagram string
	query := `SELECT value FROM settings.plain WHERE key = $1`
	err := pool.QueryRow(ctx, query, settings.CompanyInstagram).Scan(&instagram)
	return instagram, err
}

// OTP Length specific helpers
func InsertOTPLength(t *testing.T, ctx context.Context, length int, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, settings.OTPLength, fmt.Sprintf("%d", length))
	require.NoError(t, err)
}

func GetOTPLengthFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, error) {
	t.Helper()
	var lengthStr string
	query := `SELECT value FROM settings.plain WHERE key = $1`
	err := pool.QueryRow(ctx, query, settings.OTPLength).Scan(&lengthStr)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(lengthStr)
}

// OTP Max Attempts specific helpers
func InsertOTPMaxAttempts(t *testing.T, ctx context.Context, maxAttempts int, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO settings.plain (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`
	_, err := pool.Exec(ctx, query, settings.OTPMaxAttempts, fmt.Sprintf("%d", maxAttempts))
	require.NoError(t, err)
}

func GetOTPMaxAttemptsFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, error) {
	t.Helper()
	var maxAttemptsStr string
	query := `SELECT value FROM settings.plain WHERE key = $1`
	err := pool.QueryRow(ctx, query, settings.OTPMaxAttempts).Scan(&maxAttemptsStr)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(maxAttemptsStr)
}

// Simple test data insert functions for service authentication tests

func InsertTestCompanyName(t *testing.T, ctx context.Context, name string, pool *pgxpool.Pool) {
	InsertCompanyName(t, ctx, name, pool)
}

func InsertTestCompanyEmail(t *testing.T, ctx context.Context, email string, pool *pgxpool.Pool) {
	InsertCompanyEmail(t, ctx, email, pool)
}

// func InsertTestCompanyPhone(t *testing.T, ctx context.Context, phone string, pool *pgxpool.Pool) {
func InsertTestCompanyPhone(t *testing.T, ctx context.Context, phoneSetting *domain.SettingEncryptedEncx, pool *pgxpool.Pool) {
	query := `
		INSERT INTO settings.encrypted (key, value_encrypted, dek_encrypted, key_version, metadata)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (key) DO UPDATE SET
			value_encrypted = EXCLUDED.value_encrypted,
			dek_encrypted = EXCLUDED.dek_encrypted,
			key_version = EXCLUDED.key_version;
	`
	_, err := pool.Exec(
		ctx,
		query,
		phoneSetting.Key,
		phoneSetting.ValueEncrypted,
		phoneSetting.DEKEncrypted,
		phoneSetting.KeyVersion,
		phoneSetting.Metadata,
	)
	require.NoError(t, err)
}

func InsertTestCompanyAddress(t *testing.T, ctx context.Context, address string, pool *pgxpool.Pool) {
	InsertCompanyAddress(t, ctx, address, pool)
}

func InsertTestOTPDuration(t *testing.T, ctx context.Context, duration int, pool *pgxpool.Pool) {
	InsertOTPDuration(t, ctx, duration, pool)
}
