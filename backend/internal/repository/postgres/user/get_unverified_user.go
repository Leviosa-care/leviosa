package userRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// GetUnverifiedUser retrieves an unverified user's details by their email hash from the database.
//
// Parameters:
//   - ctx: Context to manage the lifecycle of the operation and handle cancellation.
//   - emailHash: The hashed email of the user to search for.
//
// Returns:
//   - *models.User: A pointer to the populated user model if the user is found.
//   - error: An error if the query fails or no matching user is found.
//   - Returns a "not found" error if no user matches the provided email hash.
//   - Returns a context error if the operation is canceled or the deadline is exceeded.
//   - Returns a database error for any other query-related issues.
func (u *repository) GetUnverifiedUser(ctx context.Context, emailHash string) (*models.User, error) {
	var user models.User
	query := `
        SELECT 
            email_encrypted,
            lastname_encrypted,
            firstname_encrypted,
            gender_encrypted,
            birthdate_encrypted,
            telephone_encrypted,
            postal_code_encrypted,
            city_encrypted,
            address1_encrypted,
            address2_encrypted,
			dek_encrypted
        FROM unverified_users 
        WHERE email_hash = $1;`

	err := u.DB.QueryRowContext(ctx, query, emailHash).Scan(
		&user.EmailEncrypted,
		&user.LastNameEncrypted,
		&user.FirstNameEncrypted,
		&user.GenderEncrypted,
		&user.BirthDateEncrypted,
		&user.TelephoneEncrypted,
		&user.PostalCodeEncrypted,
		&user.CityEncrypted,
		&user.Address1Encrypted,
		&user.Address2Encrypted,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, rp.NewNotFoundErr(err, "unverified user")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	return &user, nil
}
