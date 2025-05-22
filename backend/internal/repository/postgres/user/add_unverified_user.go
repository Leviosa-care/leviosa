package userRepository

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// AddUnverifiedUser inserts a new unverified user into the 'unverified_users' table.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - user: The user object containing details to be stored in the 'unverified_users' table.
//     This includes email hash, password hash, personal details, and encrypted birthdate.
//
// Returns:
//   - error: An error if the insertion fails, including database or context-related errors. Returns nil if successful.
//   - If no rows are affected by the insertion, a "not created" error is returned.
func (u *repository) AddUnverifiedUser(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO unverified_users (
            email_hash,
            email_encrypted,
            password_hash,
            lastname_encrypted,
            firstname_encrypted,
            gender_encrypted,
            birthdate_encrypted,
            telephone_encrypted,
            telephone_hash,
            created_at,
            postal_code_encrypted,
            city_encrypted,
            address1_encrypted,
            address2_encrypted,
			dek_encrypted
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`
	result, err := u.DB.ExecContext(
		ctx,
		query,
		user.EmailHash,
		user.EmailEncrypted,
		user.PasswordHash,
		user.LastNameEncrypted,
		user.FirstNameEncrypted,
		user.GenderEncrypted,
		user.BirthDateEncrypted,
		user.TelephoneEncrypted,
		user.TelephoneHash,
		user.CreatedAt,
		user.PostalCodeEncrypted,
		user.CityEncrypted,
		user.Address1Encrypted,
		user.Address2Encrypted,
		user.DEKEncrypted,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), "unverified user with provided emailHash")
	}
	return nil
}
