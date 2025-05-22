package userRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// GetUserByEmail retrieves a user's details from the database using their email hash.
//
// Parameters:
//   - ctx: Context for managing the lifecycle of the operation (e.g., handling timeouts and cancellations).
//   - emailHash: The hashed email of the user to search for.
//
// Returns:
//   - *models.User: A pointer to a populated User model if the user is found.
//   - error: An error if the query fails or no matching user is found.
//   - Returns a "not found" error if no user matches the provided email hash.
//   - Returns a context error if the operation is canceled or times out.
//   - Returns a database error for other query-related issues.
func (u *repository) GetUserByEmail(ctx context.Context, emailHash string) (*models.User, error) {
	var user models.User
	query := `
        SELECT 
            email_encrypted,
            picture_encrypted,
            created_at,
            logged_in,
            role,
            birthdate_encrypted,
            lastname_encrypted,
            firstname_encrypted,
            gender_encrypted,
            telephone_encrypted,
            postal_code_encrypted,
            city_encrypted,
            address1_encrypted,
            address2_encrypted,
            google_id_encrypted,
            apple_id_encrypted,
			dek_encrypted
        FROM users 
        WHERE email_hash = $1;`

	err := u.DB.QueryRowContext(ctx, query, emailHash).Scan(
		&user.EmailEncrypted,
		&user.PictureEncrypted,
		&user.CreatedAt,
		&user.LoggedInAt,
		&user.Role,
		&user.BirthDateEncrypted,
		&user.LastNameEncrypted,
		&user.FirstNameEncrypted,
		&user.GenderEncrypted,
		&user.TelephoneEncrypted,
		&user.PostalCodeEncrypted,
		&user.CityEncrypted,
		&user.Address1Encrypted,
		&user.Address2Encrypted,
		&user.GoogleIDEncrypted,
		&user.AppleIDEncrypted,
		&user.DEKEncrypted,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, rp.NewNotFoundErr(err, "user")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	return &user, nil
}
