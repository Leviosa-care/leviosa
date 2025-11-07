package helpers

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// CorruptPartnerDEK overwrites a partner's DEK with random invalid bytes
// to simulate decryption failures in tests
func CorruptPartnerDEK(t *testing.T, ctx context.Context, partnerID uuid.UUID, pool *pgxpool.Pool) {
	t.Helper()

	// Generate random garbage bytes
	corruptedDEK := make([]byte, 64)
	_, err := rand.Read(corruptedDEK)
	require.NoError(t, err, "Failed to generate corrupted DEK")

	query := `UPDATE auth.partners SET dek_encrypted = $1 WHERE id = $2`
	_, err = pool.Exec(ctx, query, corruptedDEK, partnerID)
	require.NoError(t, err, "Failed to corrupt partner DEK")
}

// SetInvalidKeyVersion sets a partner's key_version to an invalid value
// to simulate decryption failures due to key version mismatch
func SetInvalidKeyVersion(t *testing.T, ctx context.Context, partnerID uuid.UUID, pool *pgxpool.Pool, version int) {
	t.Helper()

	query := `UPDATE auth.partners SET key_version = $1 WHERE id = $2`
	_, err := pool.Exec(ctx, query, version, partnerID)
	require.NoError(t, err, "Failed to set invalid key version")
}

// NullifyPartnerDEKFields sets a partner's encryption fields to NULL
// to simulate missing encryption metadata
func NullifyPartnerDEKFields(t *testing.T, ctx context.Context, partnerID uuid.UUID, pool *pgxpool.Pool) {
	t.Helper()

	query := `UPDATE auth.partners SET dek_encrypted = NULL, key_version = NULL WHERE id = $2`
	_, err := pool.Exec(ctx, query, partnerID)
	require.NoError(t, err, "Failed to nullify partner DEK fields")
}
