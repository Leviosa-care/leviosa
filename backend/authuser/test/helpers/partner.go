package helpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
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

// InsertSpecialization directly inserts a specialization into the database for testing
func InsertSpecialization(t *testing.T, ctx context.Context, spec *domain.Specialization, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		INSERT INTO auth.specializations (
			id, name, name_encrypted, name_hash,
			display_name, display_name_encrypted, display_name_hash,
			description, description_encrypted, description_hash,
			is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		)`

	_, err := pool.Exec(ctx, query,
		spec.ID,
		spec.Name, spec.NameEncrypted, spec.NameHash,
		spec.DisplayName, spec.DisplayNameEncrypted, spec.DisplayNameHash,
		spec.Description, spec.DescriptionEncrypted, spec.DescriptionHash,
		spec.IsActive,
	)
	require.NoError(t, err, "Failed to insert test specialization")
}

// GetSpecializationFromDB retrieves a specialization directly from the database
func GetSpecializationFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) (*domain.Specialization, error) {
	t.Helper()

	var spec domain.Specialization
	query := `
		SELECT id, name, name_encrypted, name_hash,
			   display_name, display_name_encrypted, display_name_hash,
			   description, description_encrypted, description_hash,
			   is_active, created_at, updated_at
		FROM auth.specializations
		WHERE id = $1`

	err := pool.QueryRow(ctx, query, id).Scan(
		&spec.ID,
		&spec.Name, &spec.NameEncrypted, &spec.NameHash,
		&spec.DisplayName, &spec.DisplayNameEncrypted, &spec.DisplayNameHash,
		&spec.Description, &spec.DescriptionEncrypted, &spec.DescriptionHash,
		&spec.IsActive,
		&spec.CreatedAt, &spec.UpdatedAt,
	)

	return &spec, err
}

// InsertPartner directly inserts a partner into the database for testing
func InsertPartner(t *testing.T, ctx context.Context, partner *domain.Partner, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		INSERT INTO auth.partners (
			id, user_id, bio, bio_encrypted, bio_hash,
			experience, experience_encrypted, experience_hash,
			certifications, certifications_encrypted, certifications_hash,
			is_verified, verified_at, verified_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, NOW(), NOW()
		)`

	_, err := pool.Exec(ctx, query,
		partner.ID, partner.UserID,
		partner.Bio, partner.BioEncrypted, partner.BioHash,
		partner.Experience, partner.ExperienceEncrypted, partner.ExperienceHash,
		partner.Certifications, partner.CertificationsEncrypted, partner.CertificationsHash,
		partner.IsVerified, partner.VerifiedAt, partner.VerifiedBy,
	)
	require.NoError(t, err, "Failed to insert test partner")
}

// GetPartnerFromDB retrieves a partner directly from the database
func GetPartnerFromDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) (*domain.Partner, error) {
	t.Helper()

	var partner domain.Partner
	query := `
		SELECT id, user_id, bio, bio_encrypted, bio_hash,
			   experience, experience_encrypted, experience_hash,
			   certifications, certifications_encrypted, certifications_hash,
			   is_verified, verified_at, verified_by, created_at, updated_at
		FROM auth.partners
		WHERE id = $1`

	err := pool.QueryRow(ctx, query, id).Scan(
		&partner.ID, &partner.UserID,
		&partner.Bio, &partner.BioEncrypted, &partner.BioHash,
		&partner.Experience, &partner.ExperienceEncrypted, &partner.ExperienceHash,
		&partner.Certifications, &partner.CertificationsEncrypted, &partner.CertificationsHash,
		&partner.IsVerified, &partner.VerifiedAt, &partner.VerifiedBy,
		&partner.CreatedAt, &partner.UpdatedAt,
	)

	return &partner, err
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