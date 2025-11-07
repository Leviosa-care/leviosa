package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

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

// NewTestPartner creates a Partner domain object with basic test data (plaintext fields only)
func NewTestPartner(t *testing.T, userID uuid.UUID) *domain.Partner {
	t.Helper()

	categoryIDs := []uuid.UUID{uuid.New(), uuid.New()}
	productIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	return &domain.Partner{
		ID:         uuid.New(),
		UserID:     userID,
		Bio:        "Test partner bio with relevant experience and qualifications",
		Experience: "5+ years of professional experience in relevant field",
		// Certifications:           []string{"Certification 1", "Certification 2", "Advanced Certification"},
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
		CategoryIDs:              categoryIDs,
		ProductIDs:               productIDs,
		StripeConnectedAccountID: "acct_test123456789",
		StripeAccountStatus:      domain.StripeAccountStatusPending,
		StripeOnboardingComplete: false,
	}
}

// NewTestPartnerEncx creates a PartnerEncx with a test user that exists in the database
func NewTestPartnerEncx(t *testing.T) *domain.PartnerEncx {
	t.Helper()
	now := time.Now()

	// First create a test user that the partner can be associated with
	userEncx := NewTestUserEncx(t)

	return &domain.PartnerEncx{
		ID:         uuid.New(),
		UserID:     userEncx.ID,
		Bio:        "bio",
		Experience: "experience",
		// CertificationsEncrypted:           []byte("certifications encrypted"),
		CategoryIDs:                       nil,
		ProductIDs:                        nil,
		StripeConnectedAccountIDEncrypted: []byte("stripe_connected_account_id_encrypted"),
		StripeAccountStatus:               domain.StripeAccountStatusPending,
		StripeOnboardingComplete:          false,
		CreatedAt:                         now,
		UpdatedAt:                         now,
		DEKEncrypted:                      []byte("dek encrypted"),
		KeyVersion:                        1,
		Metadata:                          encx.EncryptionMetadata{},
	}
}

// NewTestPartnerEncxWithUserID creates a PartnerEncx with a specific user ID
func NewTestPartnerEncxWithUserID(t *testing.T, userID uuid.UUID) *domain.PartnerEncx {
	t.Helper()
	now := time.Now()

	return &domain.PartnerEncx{
		ID:         uuid.New(),
		UserID:     userID,
		Bio:        "bio",
		Experience: "experience",
		// CertificationsEncrypted:           []byte("certifications encrypted"),
		CategoryIDs:                       nil,
		ProductIDs:                        nil,
		StripeConnectedAccountIDEncrypted: []byte("stripe_connected_account_id_encrypted"),
		StripeAccountStatus:               domain.StripeAccountStatusPending,
		StripeOnboardingComplete:          false,
		CreatedAt:                         now,
		UpdatedAt:                         now,
		DEKEncrypted:                      []byte("dek encrypted"),
		KeyVersion:                        1,
		Metadata:                          encx.EncryptionMetadata{},
	}
}

// CreateTestUserForPartner creates a test user and inserts it into the database, returning the user ID
func CreateTestUserForPartner(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	// Create a test user
	userEncx := NewTestUserEncx(t)

	// Insert the user into the database
	err := InsertUserEncx(t, ctx, userEncx, pool)
	require.NoError(t, err, "Failed to create test user for partner")

	return userEncx.ID
}

// CreateTestUserForPartnerWithUniqueEmail creates a test user with a unique email and inserts it into the database
func CreateTestUserForPartnerWithUniqueEmail(t *testing.T, ctx context.Context, pool *pgxpool.Pool, emailSuffix string) uuid.UUID {
	t.Helper()

	// Create a test user with unique email
	userEncx := NewTestUserEncx(t)
	// Make email unique by adding suffix
	uniqueEmail := fmt.Sprintf("testuser%s@example.com", emailSuffix)
	userEncx.EmailHash = uniqueEmail
	userEncx.EmailEncrypted = []byte(uniqueEmail)

	// Insert the user into the database
	err := InsertUserEncx(t, ctx, userEncx, pool)
	require.NoError(t, err, "Failed to create test user for partner")

	return userEncx.ID
}

// InsertPartnerEncx directly inserts a partner into the database for testing using the new Encx approach
func InsertPartnerEncx(t *testing.T, ctx context.Context, partnerEncx *domain.PartnerEncx, pool *pgxpool.Pool) error {
	t.Helper()

	query := `
		INSERT INTO auth.partners (
			id, user_id, bio, experience,
			category_ids, product_ids,
			stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			dek_encrypted, key_version, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`

	_, err := pool.Exec(ctx, query,
		partnerEncx.ID,
		partnerEncx.UserID,
		partnerEncx.Bio, partnerEncx.Experience,
		partnerEncx.CategoryIDs, partnerEncx.ProductIDs,
		partnerEncx.StripeConnectedAccountIDEncrypted, partnerEncx.StripeAccountStatus, partnerEncx.StripeOnboardingComplete,
		partnerEncx.DEKEncrypted, partnerEncx.KeyVersion,
		partnerEncx.CreatedAt, partnerEncx.UpdatedAt,
	)
	return err
}

// GetPartnerEncxByUserID retrieves a partner by user ID directly from the database using the new Encx approach
func GetPartnerEncxByUserID(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) (*domain.PartnerEncx, error) {
	t.Helper()

	var partnerEncx domain.PartnerEncx
	query := `
		SELECT id, bio, experience,
			   category_ids, product_ids,
			   stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			   dek_encrypted, key_version, created_at, updated_at
		FROM auth.partners
		WHERE user_id = $1`

	err := pool.QueryRow(ctx, query, userID).Scan(
		&partnerEncx.ID,
		&partnerEncx.Bio, &partnerEncx.Experience,
		&partnerEncx.CategoryIDs, &partnerEncx.ProductIDs,
		&partnerEncx.StripeConnectedAccountIDEncrypted, &partnerEncx.StripeAccountStatus, &partnerEncx.StripeOnboardingComplete,
		&partnerEncx.DEKEncrypted, &partnerEncx.KeyVersion,
		&partnerEncx.CreatedAt, &partnerEncx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	partnerEncx.UserID = userID

	return &partnerEncx, nil
}

// GetPartnerEncxByID retrieves a partner by user ID directly from the database using the new Encx approach
func GetPartnerEncxByID(t *testing.T, ctx context.Context, ID uuid.UUID, pool *pgxpool.Pool) (*domain.PartnerEncx, error) {
	t.Helper()

	var partnerEncx domain.PartnerEncx
	query := `
		SELECT id, user_id, bio, experience,
			   category_ids, product_ids,
			   stripe_connected_account_id_encrypted, stripe_account_status, stripe_onboarding_complete,
			   dek_encrypted, key_version, created_at, updated_at
		FROM auth.partners
		WHERE id = $1`

	err := pool.QueryRow(ctx, query, ID).Scan(
		&partnerEncx.ID, &partnerEncx.UserID,
		&partnerEncx.Bio, &partnerEncx.Experience,
		&partnerEncx.CategoryIDs, &partnerEncx.ProductIDs,
		&partnerEncx.StripeConnectedAccountIDEncrypted, &partnerEncx.StripeAccountStatus, &partnerEncx.StripeOnboardingComplete,
		&partnerEncx.DEKEncrypted, &partnerEncx.KeyVersion,
		&partnerEncx.CreatedAt, &partnerEncx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	//
	// partnerEncx.ID = ID

	return &partnerEncx, nil
}

// CountPartners returns the total number of partners in the auth.partners table
func CountPartners(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, error) {
	t.Helper()

	var count int
	query := `SELECT COUNT(*) FROM auth.partners`
	err := pool.QueryRow(ctx, query).Scan(&count)

	return count, err
}

// CheckPartnerExistsByUserID checks if a partner exists for a given user ID
func CheckPartnerExistsByUserID(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) (bool, error) {
	t.Helper()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.partners WHERE user_id = $1)`
	err := pool.QueryRow(ctx, query, userID).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}

// DeletePartnerEncx directly deletes a partner from the database for testing using the new Encx approach
func DeletePartnerEncx(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) error {
	t.Helper()

	query := `DELETE FROM auth.partners WHERE user_id = $1`

	_, err := pool.Exec(ctx, query, userID)
	return err
}

// DeleteUserEncx directly deletes a user from the database for testing (helper for cascade delete tests)
func DeleteUserEncx(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) error {
	t.Helper()

	query := `DELETE FROM auth.users WHERE id = $1`

	_, err := pool.Exec(ctx, query, userID)
	return err
}
