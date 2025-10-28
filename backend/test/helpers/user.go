package helpers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
)

// ClearUsersTable truncates the users table for clean test state
func ClearUsersTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE auth.users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// generateStrongPassword creates a cryptographically secure password for testing
// that won't be flagged by pwned password validation
func generateStrongPassword(t *testing.T) string {
	t.Helper()
	bytes := make([]byte, 16) // 32 character hex string
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	return fmt.Sprintf("TestPass_%s_2024!", hex.EncodeToString(bytes))
}

// NewTestUser creates a User domain object with basic test data (plaintext fields only)
func NewTestUser(t *testing.T, email, firstName, lastName string) *domain.User {
	t.Helper()
	return &domain.User{
		ID:        uuid.New(),
		State:     domain.Unverified,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		// Password:   "qPDAR0.4Z8{vpCO]",
		Password:   generateStrongPassword(t),
		Telephone:  "0612345678",
		Role:       identity.Standard.String(),
		CreatedAt:  time.Now(),
		LoggedInAt: time.Now(),
	}
}

// NewTestUserEncx creates a UserEncx domain object with basic test data (plaintext fields only)
func NewTestUserEncx(t *testing.T) *domain.UserEncx {
	t.Helper()
	return &domain.UserEncx{
		ID:                  uuid.New(),
		State:               domain.Unverified,
		EmailEncrypted:      []byte("email_encrypted"),
		EmailHash:           "email_hash",
		PasswordHashSecure:  "password_hash_secure",
		PictureEncrypted:    []byte("picture_encrypted"),
		CreatedAtEncrypted:  []byte("created_at_encrypted"),
		LoggedInAtEncrypted: []byte("logged_in_at_encrypted"),
		RoleEncrypted:       []byte("role_encrypted"),
		BirthDateEncrypted:  []byte("birthday_encrypted"),
		LastNameEncrypted:   []byte("lastname_encrypted"),
		FirstNameEncrypted:  []byte("firstname_encrypted"),
		GenderEncrypted:     []byte("gender_encrypted"),
		TelephoneHash:       "telephone_hash",
		TelephoneEncrypted:  []byte("telephone_encrypted"),
		PostalCodeEncrypted: []byte("postalcode_encrypted"),
		CityEncrypted:       []byte("city_encrypted"),
		Address1Encrypted:   []byte("address1_encrypted"),
		Address2Encrypted:   []byte("address2_encrypted"),
		GoogleIDEncrypted:   []byte("google_id_encrypted"),
		AppleIDEncrypted:    []byte("apple_id_encrypted"),
		DEKEncrypted:        []byte("dek_encrypted"),
		KeyVersion:          1,
		Metadata:            encx.EncryptionMetadata{},
	}
}

// GetUserEncxByID retrieves a UserEncx by user ID for test verification using the new Encx approach
func GetUserEnxByID(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) (*domain.UserEncx, error) {
	t.Helper()

	var userEncx domain.UserEncx
	query := `
		SELECT
			id, state, email_hash, email_encrypted, password_hash_secure,
			picture_encrypted, first_name_encrypted, last_name_encrypted,
			birth_date_encrypted, gender_encrypted, role_encrypted,
			telephone_hash, telephone_encrypted, postal_code_encrypted,
			city_encrypted, address1_encrypted, address2_encrypted, stripe_customer_id_encrypted,
			google_id_encrypted, apple_id_encrypted, created_at_encrypted,
			logged_in_at_encrypted, dek_encrypted, key_version
		FROM auth.users
		WHERE id = $1
	`

	err := pool.QueryRow(ctx, query, userID).Scan(
		&userEncx.ID, &userEncx.State, &userEncx.EmailHash, &userEncx.EmailEncrypted,
		&userEncx.PasswordHashSecure, &userEncx.PictureEncrypted, &userEncx.FirstNameEncrypted,
		&userEncx.LastNameEncrypted, &userEncx.BirthDateEncrypted, &userEncx.GenderEncrypted,
		&userEncx.RoleEncrypted, &userEncx.TelephoneHash, &userEncx.TelephoneEncrypted,
		&userEncx.PostalCodeEncrypted, &userEncx.CityEncrypted, &userEncx.Address1Encrypted,
		&userEncx.Address2Encrypted, &userEncx.StripeCustomerIDEncrypted, &userEncx.GoogleIDEncrypted, &userEncx.AppleIDEncrypted,
		&userEncx.CreatedAtEncrypted, &userEncx.LoggedInAtEncrypted, &userEncx.DEKEncrypted,
		&userEncx.KeyVersion,
	)

	return &userEncx, err
}

// GetUserEncxByEmailHash retrieves a UserEncx instance by email hash for test verification using the new Encx approach
func GetUserEncxByEmailHash(t *testing.T, ctx context.Context, emailHash string, pool *pgxpool.Pool, crypto encx.CryptoService) (*domain.UserEncx, error) {
	t.Helper()

	var userEncx domain.UserEncx

	query := `
		SELECT id, state, email_encrypted, password_hash_secure,
		       first_name_encrypted, last_name_encrypted, telephone_hash, telephone_encrypted,
		       picture_encrypted, birth_date_encrypted, gender_encrypted,
		       address1_encrypted, city_encrypted, postal_code_encrypted,
		       role_encrypted, created_at_encrypted, logged_in_at_encrypted,
		       dek_encrypted, key_version
		FROM auth.users WHERE email_hash = $1
	`
	err := pool.QueryRow(ctx, query, emailHash).Scan(
		&userEncx.ID, &userEncx.State, &userEncx.EmailEncrypted, &userEncx.PasswordHashSecure,
		&userEncx.FirstNameEncrypted, &userEncx.LastNameEncrypted, &userEncx.TelephoneHash, &userEncx.TelephoneEncrypted,
		&userEncx.PictureEncrypted, &userEncx.BirthDateEncrypted, &userEncx.GenderEncrypted,
		&userEncx.Address1Encrypted, &userEncx.CityEncrypted, &userEncx.PostalCodeEncrypted,
		&userEncx.RoleEncrypted, &userEncx.CreatedAtEncrypted, &userEncx.LoggedInAtEncrypted,
		&userEncx.DEKEncrypted, &userEncx.KeyVersion,
	)

	return &userEncx, err
}

// CountUsers returns the total number of users in the auth.users table
func CountUsers(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, error) {
	t.Helper()

	var count int

	query := `SELECT COUNT(*) FROM auth.users`
	err := pool.QueryRow(ctx, query).Scan(&count)

	return count, err
}

// CountPendingUsers returns the number of users with pending state
func CountPendingUsers(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, error) {
	t.Helper()

	var count int

	query := `SELECT COUNT(*) FROM auth.users WHERE state = $1`
	err := pool.QueryRow(ctx, query, domain.Pending).Scan(&count)

	return count, err
}

// CheckUserEncxExistsByID checks if a UserEncx instance exists by user ID
func CheckUserEncxExistsByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (bool, error) {
	t.Helper()

	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM auth.users WHERE id = $1)`
	err := pool.QueryRow(ctx, query, userID).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}

// InsertUserEncx convenience function that inserts a UserEncx instance
func InsertUserEncx(t *testing.T, ctx context.Context, userEncx *domain.UserEncx, pool *pgxpool.Pool) error {
	t.Helper()

	query := `
		INSERT INTO auth.users (
			id, state, email_hash, email_encrypted, password_hash_secure,
			first_name_encrypted, last_name_encrypted, telephone_hash, telephone_encrypted,
			picture_encrypted, birth_date_encrypted, gender_encrypted,
			address1_encrypted, address2_encrypted, city_encrypted, postal_code_encrypted,
			role_encrypted, created_at_encrypted, logged_in_at_encrypted, stripe_customer_id_encrypted,
			google_id_encrypted, apple_id_encrypted, dek_encrypted, key_version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`
	_, err := pool.Exec(ctx, query,
		userEncx.ID, userEncx.State, userEncx.EmailHash, userEncx.EmailEncrypted, userEncx.PasswordHashSecure,
		userEncx.FirstNameEncrypted, userEncx.LastNameEncrypted, userEncx.TelephoneHash, userEncx.TelephoneEncrypted,
		userEncx.PictureEncrypted, userEncx.BirthDateEncrypted, userEncx.GenderEncrypted,
		userEncx.Address1Encrypted, userEncx.Address2Encrypted, userEncx.CityEncrypted, userEncx.PostalCodeEncrypted,
		userEncx.RoleEncrypted, userEncx.CreatedAtEncrypted, userEncx.LoggedInAtEncrypted, userEncx.StripeCustomerIDEncrypted,
		userEncx.GoogleIDEncrypted, userEncx.AppleIDEncrypted, userEncx.DEKEncrypted, userEncx.KeyVersion)

	return err
}

// // NewTestUserWithPassword creates a User domain object with a specific password
// func NewTestUserWithPassword(t *testing.T, email, firstName, lastName, password string) *domain.User {
// 	t.Helper()
// 	return &domain.User{
// 		ID:         uuid.New(),
// 		State:      domain.Active,
// 		Email:      email,
// 		FirstName:  firstName,
// 		LastName:   lastName,
// 		Password:   password,
// 		Telephone:  "0612345678",
// 		Role:       identity.Standard.String(),
// 		CreatedAt:  time.Now(),
// 		LoggedInAt: time.Now(),
// 	}
// }
