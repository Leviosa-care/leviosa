package userRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// createNewUser inserts a new user into the specified table ('pending_users', 'users', etc.)
// based on the authentication provider type. The function prepares the appropriate SQL query
// for inserting user data based on the provider (Google, Apple, or Mail) and executes it within
// a transaction if provided.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - tx: The transaction object used for the operation. If nil, the operation is executed outside of a transaction.
//   - user: The user object containing details such as email hash, password hash, personal details, and provider-specific information.
//   - provider: The authentication provider type (Google, Apple, or Mail) which determines how the user data is inserted into the table.
//   - table: The name of the table into which the user data is inserted (e.g., 'users', 'pending_users').
//
// Returns:
//   - error: An error if the user insertion fails, including database-related errors.
//   - If the provider type is unsupported, a validation error is returned.
//   - If the insert operation does not affect any rows, an error indicating no creation is returned.
func (u *Repository) createNewUser(ctx context.Context, tx *sql.Tx, user *models.User, provider models.ProviderType, table string) error {
	var query string
	var args []any

	// TODO: I can add NULL to some field but what is the point of doing so ?
	switch provider {
	case models.Google:
		query = fmt.Sprintf(`
            INSERT INTO %s (
                id,
                email_hash,
                email_encrypted,
                picture_encrypted,
				created_at,
				logged_in_at,
				role,
                lastname_encrypted,
                firstname_encrypted,
                gender_encrypted,
                birthdate_encrypted,
                telephone_encrypted,
				postal_code_encrypted,
				city_encrypted,
				address1_encrypted,
				address2_encrypted,
                google_id_encrypted,
                dek_encrypted
            ) VALUES (?, ?, ?, ?, ?, ?, ? ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, table)
		args = []any{
			user.ID,
			user.EmailHash,
			user.EmailEncrypted,
			user.PictureEncrypted,
			user.CreatedAt,
			user.LoggedInAt,
			user.Role,
			user.LastNameEncrypted,
			user.FirstNameEncrypted,
			user.GenderEncrypted,
			user.BirthDateEncrypted,
			user.TelephoneEncrypted,
			user.PostalCodeEncrypted,
			user.CityEncrypted,
			user.Address1Encrypted,
			user.Address2Encrypted,
			user.GoogleIDEncrypted,
			user.DEKEncrypted,
		}
	case models.Apple:
		query = fmt.Sprintf(`
            INSERT INTO %s (
                id,
                email_hash,
                email_encrypted,
                picture_encrypted,
				created_at,
				logged_in_at,
				role,
                lastname_encrypted,
                firstname_encrypted,
                gender_encrypted,
                birthdate_encrypted,
                telephone_encrypted,
				postal_code_encrypted,
				city_encrypted,
				address1_encrypted,
				address2_encrypted,
                apple_id_encrypted,
                dek_encrypted
            ) VALUES (?, ?, ?, ?, ?, ?, ? ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, table)
		args = []any{
			user.ID,
			user.EmailHash,
			user.EmailEncrypted,
			user.PictureEncrypted,
			user.CreatedAt,
			user.LoggedInAt,
			user.Role,
			user.LastNameEncrypted,
			user.FirstNameEncrypted,
			user.GenderEncrypted,
			user.BirthDateEncrypted,
			user.TelephoneEncrypted,
			user.PostalCodeEncrypted,
			user.CityEncrypted,
			user.Address1Encrypted,
			user.Address2Encrypted,
			user.AppleIDEncrypted,
			user.DEKEncrypted,
		}
	case models.Mail:
		query = fmt.Sprintf(`
            INSERT INTO %s (
                id,
                email_hash,
                email_encrypted,
                password_hash,
				created_at,
				logged_in_at,
                picture_encrypted,
                lastname_encrypted,
                firstname_encrypted,
                gender_encrypted,
                birthdate_encrypted,
                telephone_encrypted,
                postal_code_encrypted,
                city_encrypted,
                address1_encrypted,
                address2_encrypted,
                google_id_encrypted,
                encrypted_apple_id,
				dek_encrypted
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, NULL, ?)`, table)
		args = []any{
			user.ID,
			user.EmailHash,
			user.EmailEncrypted,
			user.PasswordHash,
			user.CreatedAt,
			user.LoggedInAt,
			user.PictureEncrypted,
			user.LastNameEncrypted,
			user.FirstNameEncrypted,
			user.GenderEncrypted,
			user.BirthDateEncrypted,
			user.TelephoneEncrypted,
			user.PostalCodeEncrypted,
			user.CityEncrypted,
			user.Address1Encrypted,
			user.Address2Encrypted,
			user.DEKEncrypted,
		}
	default:
		return rp.NewValidationErr(fmt.Errorf("unsupported provider type: %v", provider), "provider")
	}
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, query, args...)
	} else {
		result, err = u.DB.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return rp.NewNotCreatedErr(fmt.Errorf("failed to create user in %s table: %w", table, err), "pending user")
	}

	// Check if the insert was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotCreatedErr(err, fmt.Sprintf("new user in %s table", table))
	}

	return nil
}
