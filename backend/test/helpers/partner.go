package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearPartnersTable removes all test data from partners table
func ClearPartnersTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	// Clear partners table
	_, err := pool.Exec(ctx, "DELETE FROM auth.partners")
	require.NoError(t, err, "Failed to clear partners table")

	// Clear users table (partners depend on users)
	_, err = pool.Exec(ctx, "DELETE FROM auth.users")
	require.NoError(t, err, "Failed to clear users table")
}

// ClearPartnerTestData clears all partner-related test data
func ClearPartnerTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	// Clear all partner-related tables in correct order
	ClearPartnersTable(t, ctx, pool)
}

// InsertPartner directly inserts a partner into the database for testing using the new Encx approach
func InsertPartner(t *testing.T, ctx context.Context, partner *domain.Partner, pool *pgxpool.Pool, crypto encx.CryptoService) {
	t.Helper()

	// Process the partner to get encrypted data
	partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
	require.NoError(t, err, "Failed to process partner for encryption")

	query := `
		INSERT INTO auth.partners (
			id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			verified_at_encrypted, is_verified, verified_by, created_at, updated_at,
			dek_encrypted, key_version, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`

	_, err = pool.Exec(ctx, query,
		partnerEncx.ID, partnerEncx.UserID,
		partnerEncx.BioEncrypted, partnerEncx.ExperienceEncrypted, partnerEncx.CertificationsEncrypted,
		partnerEncx.VerifiedAtEncrypted, partnerEncx.IsVerified, partnerEncx.VerifiedByUserID,
		partnerEncx.CreatedAt, partnerEncx.UpdatedAt,
		partnerEncx.DEKEncrypted, partnerEncx.KeyVersion, partnerEncx.Metadata,
	)
	require.NoError(t, err, "Failed to insert test partner")
}

// GetPartnerFromDB retrieves a partner directly from the database using the new Encx approach
func GetPartnerFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id uuid.UUID, crypto encx.CryptoService) (*domain.Partner, error) {
	t.Helper()

	var partnerEncx domain.PartnerEncx
	query := `
		SELECT id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
			   verified_at_encrypted, is_verified, verified_by, created_at, updated_at,
			   dek_encrypted, key_version, metadata
		FROM auth.partners
		WHERE id = $1`

	err := pool.QueryRow(ctx, query, id).Scan(
		&partnerEncx.ID, &partnerEncx.UserID,
		&partnerEncx.BioEncrypted, &partnerEncx.ExperienceEncrypted, &partnerEncx.CertificationsEncrypted,
		&partnerEncx.VerifiedAtEncrypted, &partnerEncx.IsVerified, &partnerEncx.VerifiedByUserID,
		&partnerEncx.CreatedAt, &partnerEncx.UpdatedAt,
		&partnerEncx.DEKEncrypted, &partnerEncx.KeyVersion, &partnerEncx.Metadata,
	)
	if err != nil {
		return nil, err
	}

	// Decrypt the partner using the new generated function
	partner, err := domain.DecryptPartnerEncx(ctx, crypto, &partnerEncx)
	if err != nil {
		return nil, fmt.Errorf("decrypt partner: %w", err)
	}

	return partner, nil
}

// GetPartnerByUserID retrieves a partner by user ID directly from the database using the new Encx approach
func GetPartnerByUserID(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) (*domain.PartnerEncx, error) {
	t.Helper()

	var partnerEncx domain.PartnerEncx
	query := `
		SELECT id, user_id, bio_encrypted, experience_encrypted, certifications_encrypted,
		       category_ids, product_ids, is_verified, verified_at_encrypted, verified_by_user_id,
		       dek_encrypted, key_version, created_at, updated_at
		FROM auth.partners
		WHERE user_id = $1`

	err := pool.QueryRow(ctx, query, userID).Scan(
		&partnerEncx.ID, &partnerEncx.UserID,
		&partnerEncx.BioEncrypted, &partnerEncx.ExperienceEncrypted, &partnerEncx.CertificationsEncrypted,
		&partnerEncx.CategoryIDs, &partnerEncx.ProductIDs,
		&partnerEncx.IsVerified, &partnerEncx.VerifiedAtEncrypted, &partnerEncx.VerifiedByUserID,
		&partnerEncx.DEKEncrypted, &partnerEncx.KeyVersion,
		&partnerEncx.CreatedAt, &partnerEncx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &partnerEncx, nil
}
