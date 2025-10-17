package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearPartnersTable removes all test data from partners and related tables
func ClearPartnersTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	// Clear partner_specializations first (foreign key dependency)
	_, err := pool.Exec(ctx, "DELETE FROM auth.partner_specializations")
	require.NoError(t, err, "Failed to clear partner_specializations table")

	// Clear partners table
	_, err = pool.Exec(ctx, "DELETE FROM auth.partners")
	require.NoError(t, err, "Failed to clear partners table")

	// Clear users table (partners depend on users)
	_, err = pool.Exec(ctx, "DELETE FROM auth.users")
	require.NoError(t, err, "Failed to clear users table")
}

// ClearSpecializationsTable removes all test data from specializations table
func ClearSpecializationsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	// Clear partner_specializations first (foreign key dependency)
	_, err := pool.Exec(ctx, "DELETE FROM auth.partner_specializations")
	require.NoError(t, err, "Failed to clear partner_specializations table")

	// Clear specializations table
	_, err = pool.Exec(ctx, "DELETE FROM auth.specializations")
	require.NoError(t, err, "Failed to clear specializations table")
}

// ClearPartnerTestData clears all partner-related test data
func ClearPartnerTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	// Clear all partner-related tables in correct order
	ClearPartnersTable(t, ctx, pool)
	ClearSpecializationsTable(t, ctx, pool)
}

// NewTestSpecialization creates a test specialization with given parameters
func NewTestSpecialization(name, displayName, description string) *domain.Specialization {
	return &domain.Specialization{
		ID:          uuid.New(),
		Name:        name,
		DisplayName: displayName,
		Description: description,
		IsActive:    true,
	}
}

// InsertSpecialization directly inserts a specialization into the database for testing using the new Encx approach
func InsertSpecialization(t *testing.T, ctx context.Context, spec *domain.Specialization, pool *pgxpool.Pool, crypto encx.CryptoService) {
	t.Helper()

	// Process the specialization to get encrypted data
	specEncx, err := domain.ProcessSpecializationEncx(ctx, crypto, spec)
	require.NoError(t, err, "Failed to process specialization for encryption")

	query := `
		INSERT INTO auth.specializations (
			id, name_encrypted, description_encrypted, display_name_encrypted,
			is_active, created_at, updated_at, dek_encrypted, key_version, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	_, err = pool.Exec(ctx, query,
		specEncx.ID,
		specEncx.NameEncrypted, specEncx.DescriptionEncrypted, specEncx.DisplayNameEncrypted,
		specEncx.IsActive, specEncx.CreatedAt, specEncx.UpdatedAt,
		specEncx.DEKEncrypted, specEncx.KeyVersion, specEncx.Metadata,
	)
	require.NoError(t, err, "Failed to insert test specialization")
}

// GetSpecializationFromDB retrieves a specialization directly from the database using the new Encx approach
func GetSpecializationFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id uuid.UUID, crypto encx.CryptoService) (*domain.Specialization, error) {
	t.Helper()

	var specEncx domain.SpecializationEncx
	query := `
		SELECT id, name_encrypted, description_encrypted, display_name_encrypted,
			   is_active, created_at, updated_at, dek_encrypted, key_version, metadata
		FROM auth.specializations
		WHERE id = $1`

	err := pool.QueryRow(ctx, query, id).Scan(
		&specEncx.ID,
		&specEncx.NameEncrypted, &specEncx.DescriptionEncrypted, &specEncx.DisplayNameEncrypted,
		&specEncx.IsActive, &specEncx.CreatedAt, &specEncx.UpdatedAt,
		&specEncx.DEKEncrypted, &specEncx.KeyVersion, &specEncx.Metadata,
	)
	if err != nil {
		return nil, err
	}

	// Decrypt the specialization using the new generated function
	spec, err := domain.DecryptSpecializationEncx(ctx, crypto, &specEncx)
	if err != nil {
		return nil, fmt.Errorf("decrypt specialization: %w", err)
	}

	return spec, nil
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

// AddPartnerSpecializationToDB directly adds a partner-specialization association to the database
func AddPartnerSpecializationToDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, partnerID, specializationID uuid.UUID) {
	t.Helper()

	query := `
		INSERT INTO auth.partner_specializations (partner_id, specialization_id, created_at)
		VALUES ($1, $2, NOW())`

	_, err := pool.Exec(ctx, query, partnerID, specializationID)
	require.NoError(t, err, "Failed to add partner specialization association")
}

// GetPartnerSpecializationCount returns the number of specializations associated with a partner
func GetPartnerSpecializationCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, partnerID uuid.UUID) int {
	t.Helper()

	var count int
	query := `SELECT COUNT(*) FROM auth.partner_specializations WHERE partner_id = $1`

	err := pool.QueryRow(ctx, query, partnerID).Scan(&count)
	require.NoError(t, err, "Failed to get partner specialization count")

	return count
}

// CheckPartnerSpecializationExists checks if a partner-specialization association exists
func CheckPartnerSpecializationExists(t *testing.T, ctx context.Context, pool *pgxpool.Pool, partnerID, specializationID uuid.UUID) bool {
	t.Helper()

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM auth.partner_specializations
			WHERE partner_id = $1 AND specialization_id = $2
		)`

	err := pool.QueryRow(ctx, query, partnerID, specializationID).Scan(&exists)
	require.NoError(t, err, "Failed to check partner specialization existence")

	return exists
}

