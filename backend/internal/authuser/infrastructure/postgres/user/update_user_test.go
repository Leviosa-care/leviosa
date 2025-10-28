package userRepository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateUser TEST_PATH=internal/authuser/infrastructure/postgres/user/update_user_test.go

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "updateuser@example.com"

		// Create and insert initial userEncx
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailHash = email
		userEncx.EmailEncrypted = []byte(email)

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Update user data
		firstnameEncrypted := []byte("Jane")
		lastnameEncrypted := []byte("Smith")
		passwordHashSecure := "newpassword123"
		genderEncrypted := []byte("woman")
		telephoneEncrypted := []byte("+33987654321")
		postalCodeEncrypted := []byte("75001")
		cityEncrypted := []byte("Paris")
		address1Encrypted := []byte("123 Updated Street")
		address2Encrypted := []byte("Apt 456")
		birthdayEncrypted := []byte("new_birthday_encrypted")

		userEncx.FirstNameEncrypted = firstnameEncrypted
		userEncx.LastNameEncrypted = lastnameEncrypted
		userEncx.State = domain.Pending
		userEncx.PasswordHashSecure = passwordHashSecure
		userEncx.GenderEncrypted = genderEncrypted
		userEncx.TelephoneEncrypted = telephoneEncrypted
		userEncx.PostalCodeEncrypted = postalCodeEncrypted
		userEncx.CityEncrypted = cityEncrypted
		userEncx.Address1Encrypted = address1Encrypted
		userEncx.Address2Encrypted = address2Encrypted
		userEncx.BirthDateEncrypted = birthdayEncrypted

		// Act
		err = repo.UpdateUser(ctx, userEncx)

		// Assert
		assert.NoError(t, err)

		// Verify the update by retrieving the user
		// updatedUser, err := repo.GetUserByID(ctx, userEncx.ID)
		updatedUser, err := td.GetUserEnxByID(t, ctx, userEncx.ID, testPool)
		require.NoError(t, err)
		require.NotNil(t, updatedUser)

		// Verify data was updated
		assert.Equal(t, userEncx.ID, updatedUser.ID)
		assert.Equal(t, domain.Pending, updatedUser.State)
		assert.Equal(t, firstnameEncrypted, updatedUser.FirstNameEncrypted)
		assert.Equal(t, lastnameEncrypted, updatedUser.LastNameEncrypted)
		assert.Equal(t, genderEncrypted, updatedUser.GenderEncrypted)
		assert.Equal(t, telephoneEncrypted, updatedUser.TelephoneEncrypted)
		assert.Equal(t, postalCodeEncrypted, updatedUser.PostalCodeEncrypted)
		assert.Equal(t, cityEncrypted, updatedUser.CityEncrypted)
		assert.Equal(t, address1Encrypted, updatedUser.Address1Encrypted)
		assert.Equal(t, address2Encrypted, updatedUser.Address2Encrypted)
		assert.Equal(t, birthdayEncrypted, updatedUser.BirthDateEncrypted)
	})

	t.Run("should return not found error when updating non-existent user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user := td.NewTestUserEncx(t)

		// Act
		err := repo.UpdateUser(ctx, user)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should successfully update user state from unverified to pending", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		// email := "statechange@example.com"

		// Create userEncx with unverified state
		userEncx := td.NewTestUserEncx(t)
		userEncx.State = domain.Unverified

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Update user state to pending
		userEncx.State = domain.Pending

		// Act
		err = repo.UpdateUser(ctx, userEncx)

		// Assert
		assert.NoError(t, err)

		// Verify the state change
		updatedUser, err := td.GetUserEnxByID(t, ctx, userEncx.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, domain.Pending, updatedUser.State)
	})

	t.Run("should successfully update user with empty optional fields", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create userEncx with some fields populated
		userEncx := td.NewTestUserEncx(t)

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Update user to remove optional fields
		userEncx.Address2Encrypted = []byte("")
		userEncx.TelephoneEncrypted = []byte("")

		// Act
		err = repo.UpdateUser(ctx, userEncx)

		// Assert
		assert.NoError(t, err)

		// Verify the update
		updatedUser, err := td.GetUserEnxByID(t, ctx, userEncx.ID, testPool)
		require.NoError(t, err)

		assert.Empty(t, updatedUser.Address2Encrypted)
		assert.Empty(t, updatedUser.TelephoneEncrypted)
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})

	t.Run("should successfully update all user fields in complete user flow", func(t *testing.T) {
		// Arrange - simulates the complete user registration flow
		td.ClearUsersTable(t, ctx, testPool)
		email := "complete@example.com"

		// Start with minimal unverified userEncx
		userEncx := &domain.UserEncx{
			ID:                 uuid.New(),
			State:              domain.Unverified,
			EmailEncrypted:     []byte(email),
			CreatedAtEncrypted: []byte("created_at_encrypted"),
			DEKEncrypted:       []byte("dek_encrypted"),
			KeyVersion:         1,
		}

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Update with complete user information

		firstName := []byte("Complete")
		lastName := []byte("Registration")
		password := "securepassword123"
		gender := []byte("non_binary")
		telephone := []byte("+33456789012")
		postalCode := []byte("69000")
		city := []byte("Lyon")
		address1 := []byte("456 Complete Avenue")
		address2 := []byte("Building C")
		birthDate := []byte("new_birthday")

		userEncx.State = domain.Pending
		userEncx.FirstNameEncrypted = firstName
		userEncx.LastNameEncrypted = lastName
		userEncx.PasswordHashSecure = password
		userEncx.GenderEncrypted = gender
		userEncx.TelephoneEncrypted = telephone
		userEncx.PostalCodeEncrypted = postalCode
		userEncx.CityEncrypted = city
		userEncx.Address1Encrypted = address1
		userEncx.Address2Encrypted = address2
		userEncx.BirthDateEncrypted = birthDate

		// Act
		err = repo.UpdateUser(ctx, userEncx)

		// Assert
		assert.NoError(t, err)

		// Verify all fields were updated correctly
		completeUser, err := repo.GetUserByID(ctx, userEncx.ID)
		require.NoError(t, err)

		assert.Equal(t, domain.Pending, completeUser.State)
		assert.Equal(t, firstName, completeUser.FirstNameEncrypted)
		assert.Equal(t, lastName, completeUser.LastNameEncrypted)
		assert.Equal(t, gender, completeUser.GenderEncrypted)
		assert.Equal(t, telephone, completeUser.TelephoneEncrypted)
		assert.Equal(t, postalCode, completeUser.PostalCodeEncrypted)
		assert.Equal(t, city, completeUser.CityEncrypted)
		assert.Equal(t, address1, completeUser.Address1Encrypted)
		assert.Equal(t, address2, completeUser.Address2Encrypted)
		assert.NotEmpty(t, completeUser.PasswordHashSecure) // Password should be hashed
	})
}
