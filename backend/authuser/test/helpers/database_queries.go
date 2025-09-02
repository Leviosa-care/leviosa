package helpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// CheckUserExistsByEmailHashSQL checks if a user exists using raw SQL query
func CheckUserExistsByEmailHashSQL(t *testing.T, ctx context.Context, emailHash string, pool *pgxpool.Pool) bool {
	t.Helper()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.users WHERE email_hash = $1)`

	err := pool.QueryRow(ctx, query, emailHash).Scan(&exists)
	require.NoError(t, err, "Failed to check user existence with raw SQL")

	return exists
}

// GetUserByEmailHashSQL retrieves a user using raw SQL query
func GetUserByEmailHashSQL(t *testing.T, ctx context.Context, emailHash string, pool *pgxpool.Pool) *domain.User {
	t.Helper()

	query := `
		SELECT 
			id, state, email_hash, email_encrypted, password_hash,
			picture_encrypted, first_name_encrypted, last_name_encrypted, 
			birth_date_encrypted, gender_encrypted, role_encrypted,
			telephone_hash, telephone_encrypted, postal_code_encrypted,
			city_encrypted, address1_encrypted, address2_encrypted,
			google_id_encrypted, apple_id_encrypted, created_at_encrypted,
			logged_in_at_encrypted, dek_encrypted, key_version
		FROM auth.users 
		WHERE email_hash = $1
	`

	user := &domain.User{}
	var telephoneHash *string // Handle nullable field

	err := pool.QueryRow(ctx, query, emailHash).Scan(
		&user.ID, &user.State, &user.EmailHash, &user.EmailEncrypted,
		&user.PasswordHash, &user.PictureEncrypted, &user.FirstNameEncrypted,
		&user.LastNameEncrypted, &user.BirthDateEncrypted, &user.GenderEncrypted,
		&user.RoleEncrypted, &telephoneHash, &user.TelephoneEncrypted,
		&user.PostalCodeEncrypted, &user.CityEncrypted, &user.Address1Encrypted,
		&user.Address2Encrypted, &user.GoogleIDEncrypted, &user.AppleIDEncrypted,
		&user.CreatedAtEncrypted, &user.LoggedInAtEncrypted, &user.DEKEncrypted,
		&user.KeyVersion,
	)

	if err != nil {
		// Return nil if user not found
		return nil
	}

	// Handle nullable telephone_hash
	if telephoneHash != nil {
		user.TelephoneHash = *telephoneHash
	}

	return user
}

// CountUsersSQL returns the total number of users in the database
func CountUsersSQL(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	var count int
	query := `SELECT COUNT(*) FROM auth.users`

	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err, "Failed to count users with raw SQL")

	return count
}

// GetUserStateSQL gets the user state for a specific email hash
func GetUserStateSQL(t *testing.T, ctx context.Context, emailHash string, pool *pgxpool.Pool) domain.UserState {
	t.Helper()

	var state domain.UserState
	query := `SELECT state FROM auth.users WHERE email_hash = $1`

	err := pool.QueryRow(ctx, query, emailHash).Scan(&state)
	require.NoError(t, err, "Failed to get user state with raw SQL")

	return state
}

// CheckUserHasEncryptedFieldsSQL verifies that encrypted fields are populated
func CheckUserHasEncryptedFieldsSQL(t *testing.T, ctx context.Context, emailHash string, pool *pgxpool.Pool) bool {
	t.Helper()

	query := `
		SELECT 
			email_encrypted IS NOT NULL AND LENGTH(email_encrypted) > 0,
			dek_encrypted IS NOT NULL AND LENGTH(dek_encrypted) > 0,
			key_version > 0
		FROM auth.users 
		WHERE email_hash = $1
	`

	var emailEncrypted, dekEncrypted, keyVersionValid bool

	err := pool.QueryRow(ctx, query, emailHash).Scan(
		&emailEncrypted, &dekEncrypted, &keyVersionValid,
	)

	if err != nil {
		return false
	}

	return emailEncrypted && dekEncrypted && keyVersionValid
}
