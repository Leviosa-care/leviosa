package userRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// FindAccountByID retrieves a user's account details from the 'users' table based on the provided user ID.
// The function performs a database query to fetch the user's data and maps it to a User model.
// If the user is not found or an error occurs during the operation, appropriate errors are returned.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - id: The unique identifier of the user to find in the database.
//
// Returns:
//   - *models.User: A pointer to the User model containing the retrieved account details.
//   - error: An error if the query fails or the user is not found.
//   - If the query returns no rows, a "not found" error is returned.
//   - If context-related errors occur (e.g., deadline exceeded, canceled), a context error is returned.
//   - For any other query failures, a database error is returned.
func (u *repository) FindAccountByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `
        SELECT 
            email_encrypted,
            picture_encrypted,
            role,
            lastname_encrypted,
            firstname_encrypted,
            gender_encrypted,
            birthdate_encrypted,
            telephone_encrypted,
            postal_code_encrypted,
            city_encrypted,
            address1_encrypted,
            encrypted_address2,
			dek_encrypted
        FROM users
        WHERE id = $1;`
	if err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&user.EmailEncrypted,
		&user.PictureEncrypted,
		&user.Role,
		&user.LastNameEncrypted,
		&user.FirstNameEncrypted,
		&user.GenderEncrypted,
		&user.BirthDateEncrypted,
		&user.TelephoneEncrypted,
		&user.PostalCodeEncrypted,
		&user.CityEncrypted,
		&user.Address1Encrypted,
		&user.Address2Encrypted,
		&user.DEKEncrypted,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &models.User{}, rp.NewNotFoundErr(err, "unverified user")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	return &user, nil
}
