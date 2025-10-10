package userRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) CreateUser(ctx context.Context, user *domain.UserEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.users (
			id, state, email_hash, email_encrypted, password_hash_secure,
			picture_encrypted, first_name_encrypted, last_name_encrypted,
			birth_date_encrypted, gender_encrypted, role_encrypted,
			telephone_hash, telephone_encrypted, postal_code_encrypted,
			city_encrypted, address1_encrypted, address2_encrypted, stripe_customer_id_encrypted,
			google_id_encrypted, apple_id_encrypted, created_at_encrypted,
			logged_in_at_encrypted, dek_encrypted, key_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.State, user.EmailHash, user.EmailEncrypted,
		user.PasswordHashSecure, user.PictureEncrypted, user.FirstNameEncrypted,
		user.LastNameEncrypted, user.BirthDateEncrypted, user.GenderEncrypted,
		user.RoleEncrypted, user.TelephoneHash, user.TelephoneEncrypted,
		user.PostalCodeEncrypted, user.CityEncrypted, user.Address1Encrypted,
		user.Address2Encrypted, user.StripeCustomerIDEncrypted,
		user.GoogleIDEncrypted, user.AppleIDEncrypted, user.CreatedAtEncrypted, user.LoggedInAtEncrypted,
		user.DEKEncrypted, user.KeyVersion,
	)
	if err != nil {
		return errs.ClassifyPgError("create user", err)
	}

	return nil
}
