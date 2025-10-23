package userRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) UpdateUser(ctx context.Context, user *domain.UserEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.users SET
			state = $2,
			email_hash = $3,
			email_encrypted = $4,
			password_hash_secure = $5,
			picture_encrypted = $6,
			first_name_encrypted = $7,
			last_name_encrypted = $8,
			birth_date_encrypted = $9,
			gender_encrypted = $10,
			role_encrypted = $11,
			telephone_hash = $12,
			telephone_encrypted = $13,
			postal_code_encrypted = $14,
			city_encrypted = $15,
			address1_encrypted = $16,
			address2_encrypted = $17,
			stripe_customer_id_encrypted = $18,
			google_id_encrypted = $19,
			apple_id_encrypted = $20,
			created_at_encrypted = $21,
			logged_in_at_encrypted = $22,
			dek_encrypted = $23,
			key_version = $24
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		user.ID, user.State, user.EmailHash, user.EmailEncrypted,
		user.PasswordHashSecure, user.PictureEncrypted, user.FirstNameEncrypted,
		user.LastNameEncrypted, user.BirthDateEncrypted, user.GenderEncrypted,
		user.RoleEncrypted, user.TelephoneHash, user.TelephoneEncrypted,
		user.PostalCodeEncrypted, user.CityEncrypted, user.Address1Encrypted,
		user.Address2Encrypted, user.StripeCustomerIDEncrypted, user.GoogleIDEncrypted, user.AppleIDEncrypted,
		user.CreatedAtEncrypted, user.LoggedInAtEncrypted, user.DEKEncrypted,
		user.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("update user", err)
	}

	// Check if any row was actually updated
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
