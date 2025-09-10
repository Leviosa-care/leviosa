package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/contracts/identity"
)

// ClearUsersTable truncates the users table for clean test state
func ClearUsersTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE auth.users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// NewTestUser creates a User domain object with basic test data (plaintext fields only)
func NewTestUser(email, firstName, lastName string) *domain.User {
	return &domain.User{
		ID:         uuid.New(),
		State:      domain.Unverified,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		Password:   "qPDAR0.4Z8{vpCO]",
		Telephone:  "0612345678",
		Role:       identity.Standard.String(),
		CreatedAt:  time.Now(),
		LoggedInAt: time.Now(),
	}
}

// NewTestUserWithEncryption creates a User with all encrypted/hashed fields populated
func NewTestUserWithEncryption(email, firstName, lastName string, crypto encx.CryptoService) (*domain.User, error) {
	user := NewTestUser(email, firstName, lastName)

	// Use crypto service to process the struct and populate encrypted/hashed fields
	err := crypto.ProcessStruct(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("process user struct for encryption: %w", err)
	}

	return user, nil
}

// InsertUser performs atomic insertion of a User into the database
func InsertUser(t *testing.T, ctx context.Context, user *domain.User, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		INSERT INTO auth.users (
			id, state, email_hash, email_encrypted, password_hash,
			first_name_encrypted, last_name_encrypted, telephone_hash, telephone_encrypted,
			picture_encrypted, birth_date_encrypted, gender_encrypted,
			address1_encrypted, address2_encrypted, city_encrypted, postal_code_encrypted,
			role_encrypted, created_at_encrypted, logged_in_at_encrypted, stripe_customer_id_encrypted,
			google_id_encrypted, apple_id_encrypted, dek_encrypted, key_version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`
	_, err := pool.Exec(ctx, query,
		user.ID, user.State, user.EmailHash, user.EmailEncrypted, user.PasswordHash,
		user.FirstNameEncrypted, user.LastNameEncrypted, user.TelephoneHash, user.TelephoneEncrypted,
		user.PictureEncrypted, user.BirthDateEncrypted, user.GenderEncrypted,
		user.Address1Encrypted, user.Address2Encrypted, user.CityEncrypted, user.PostalCodeEncrypted,
		user.RoleEncrypted, user.CreatedAtEncrypted, user.LoggedInAtEncrypted, user.StripeCustomerIDEncrypted,
		user.GoogleIDEncrypted, user.AppleIDEncrypted, user.DEKEncrypted, user.KeyVersion)
	require.NoError(t, err)
}

// InsertTestUser convenience function that creates and inserts a test user
func InsertTestUser(t *testing.T, ctx context.Context, email, firstName, lastName string, pool *pgxpool.Pool, crypto encx.CryptoService) {
	t.Helper()
	user, err := NewTestUserWithEncryption(email, firstName, lastName, crypto)
	require.NoError(t, err, "Failed to create encrypted test user")
	InsertUser(t, ctx, user, pool)
}

// GetUserByIDDecrypted retrieves a user by user ID for test verification
func GetUserByIDDecrypted(t *testing.T, ctx context.Context, userID string, pool *pgxpool.Pool, crypto encx.CryptoService) (*domain.User, error) {
	t.Helper()
	var user domain.User
	query := `
		SELECT 
			id, state, email_hash, email_encrypted, password_hash,
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
		&user.ID, &user.State, &user.EmailHash, &user.EmailEncrypted,
		&user.PasswordHash, &user.PictureEncrypted, &user.FirstNameEncrypted,
		&user.LastNameEncrypted, &user.BirthDateEncrypted, &user.GenderEncrypted,
		&user.RoleEncrypted, &user.TelephoneHash, &user.TelephoneEncrypted,
		&user.PostalCodeEncrypted, &user.CityEncrypted, &user.Address1Encrypted,
		&user.Address2Encrypted, &user.StripeCustomerIDEncrypted, &user.GoogleIDEncrypted, &user.AppleIDEncrypted,
		&user.CreatedAtEncrypted, &user.LoggedInAtEncrypted, &user.DEKEncrypted,
		&user.KeyVersion,
	)

	if err != nil {
		return nil, err
	}

	// Use crypto service to decrypt the struct fields
	err = crypto.DecryptStruct(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("decrypt user struct: %w", err)
	}

	return &user, nil
}

// GetUserByEmailHash retrieves a user by email hash for test verification
func GetUserByEmailHash(t *testing.T, ctx context.Context, emailHash string, pool *pgxpool.Pool, crypto encx.CryptoService) (*domain.User, error) {
	t.Helper()
	var user domain.User
	query := `
		SELECT id, state, email_encrypted, password_hash,
		       first_name_encrypted, last_name_encrypted, telephone_hash, telephone_encrypted,
		       picture_encrypted, birth_date_encrypted, gender_encrypted,
		       address1_encrypted, city_encrypted, postal_code_encrypted,
		       role_encrypted, created_at_encrypted, logged_in_at_encrypted,
		       dek_encrypted, key_version
		FROM auth.users WHERE email_hash = $1
	`
	err := pool.QueryRow(ctx, query, emailHash).Scan(
		&user.ID, &user.State, &user.EmailEncrypted, &user.PasswordHash,
		&user.FirstNameEncrypted, &user.LastNameEncrypted, &user.TelephoneHash, &user.TelephoneEncrypted,
		&user.PictureEncrypted, &user.BirthDateEncrypted, &user.GenderEncrypted,
		&user.Address1Encrypted, &user.CityEncrypted, &user.PostalCodeEncrypted,
		&user.RoleEncrypted, &user.CreatedAtEncrypted, &user.LoggedInAtEncrypted,
		&user.DEKEncrypted, &user.KeyVersion,
	)
	if err != nil {
		return nil, err
	}

	// Use crypto service to decrypt the struct fields
	err = crypto.DecryptStruct(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("decrypt user struct: %w", err)
	}

	return &user, nil
}

// CountUsers returns the total number of users in the auth.users table
func CountUsers(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()
	var count int
	query := `SELECT COUNT(*) FROM auth.users`
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err)
	return count
}

// CountPendingUsers returns the number of users with pending state
func CountPendingUsers(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()
	var count int
	query := `SELECT COUNT(*) FROM auth.users WHERE state = $1`
	err := pool.QueryRow(ctx, query, domain.Pending).Scan(&count)
	require.NoError(t, err)
	return count
}

// CheckUserExistsByID checks if a user exists by user ID
func CheckUserExistsByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (bool, error) {
	t.Helper()
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.users WHERE id = $1)`
	err := pool.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// InsertUserWithEncryption convenience function that creates, encrypts and inserts a user
func InsertUserWithEncryption(t *testing.T, ctx context.Context, user *domain.User, pool *pgxpool.Pool, crypto encx.CryptoService) {
	t.Helper()

	// Process encryption on the user
	err := crypto.ProcessStruct(ctx, user)
	require.NoError(t, err, "Failed to encrypt user")

	InsertUser(t, ctx, user, pool)
}

// GetUserByIDFromDB retrieves a user by UUID for test verification
func GetUserByIDFromDB(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool, crypto encx.CryptoService) *domain.User {
	t.Helper()

	user, err := GetUserByIDDecrypted(t, ctx, userID.String(), pool, crypto)
	require.NoError(t, err, "Failed to get user by ID from database")

	return user
}

// NewTestUserWithPassword creates a User domain object with a specific password
func NewTestUserWithPassword(email, firstName, lastName, password string) *domain.User {
	return &domain.User{
		ID:         uuid.New(),
		State:      domain.Active,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		Password:   password,
		Telephone:  "0612345678",
		Role:       identity.Standard.String(),
		CreatedAt:  time.Now(),
		LoggedInAt: time.Now(),
	}
}

// NewTestUserWithPasswordAndEncryption creates a User with specific password and all encrypted/hashed fields populated
func NewTestUserWithPasswordAndEncryption(email, firstName, lastName, password string, crypto encx.CryptoService) (*domain.User, error) {
	user := NewTestUserWithPassword(email, firstName, lastName, password)

	// Use crypto service to process the struct and populate encrypted/hashed fields
	err := crypto.ProcessStruct(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("process user struct for encryption: %w", err)
	}

	return user, nil
}

// InsertTestUserWithPassword convenience function that creates and inserts a test user with specific password
func InsertTestUserWithPassword(t *testing.T, ctx context.Context, email, firstName, lastName, password string, pool *pgxpool.Pool, crypto encx.CryptoService) uuid.UUID {
	t.Helper()
	user, err := NewTestUserWithPasswordAndEncryption(email, firstName, lastName, password, crypto)
	require.NoError(t, err, "Failed to create encrypted test user with password")
	InsertUser(t, ctx, user, pool)
	return user.ID
}

// GetUserByID returns a user by ID, returns nil if not found (for test assertions)
func GetUserByID(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) *domain.User {
	t.Helper()

	var user domain.User
	query := `
		SELECT
			id, state, email_hash, email_encrypted, password_hash,
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
		&user.ID, &user.State, &user.EmailHash, &user.EmailEncrypted,
		&user.PasswordHash, &user.PictureEncrypted, &user.FirstNameEncrypted,
		&user.LastNameEncrypted, &user.BirthDateEncrypted, &user.GenderEncrypted,
		&user.RoleEncrypted, &user.TelephoneHash, &user.TelephoneEncrypted,
		&user.PostalCodeEncrypted, &user.CityEncrypted, &user.Address1Encrypted,
		&user.Address2Encrypted, &user.StripeCustomerIDEncrypted, &user.GoogleIDEncrypted, &user.AppleIDEncrypted,
		&user.CreatedAtEncrypted, &user.LoggedInAtEncrypted, &user.DEKEncrypted,
		&user.KeyVersion,
	)

	if err != nil {
		// Return nil if user not found (for test assertions)
		return nil
	}

	return &user
}
