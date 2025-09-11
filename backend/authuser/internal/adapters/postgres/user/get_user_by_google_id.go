package userRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) GetUserByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, state, email_hash, email_encrypted, password_hash,
			picture_encrypted, first_name_encrypted, last_name_encrypted, 
			birth_date_encrypted, gender_encrypted, role_encrypted,
			telephone_hash, telephone_encrypted, postal_code_encrypted,
			city_encrypted, address1_encrypted, address2_encrypted, stripe_customer_id_encrypted,
			google_id_encrypted, apple_id_encrypted, created_at_encrypted,
			logged_in_at_encrypted, dek_encrypted, key_version
		FROM %s.users 
		WHERE google_id_encrypted = $1
	`, r.schema)

	user := &domain.User{}

	// Only nullable string fields need special handling
	var telephoneHash sql.NullString

	err := r.pool.QueryRow(ctx, query, googleID).Scan(
		&user.ID, &user.State, &user.EmailHash, &user.EmailEncrypted,
		&user.PasswordHash, &user.PictureEncrypted, &user.FirstNameEncrypted,
		&user.LastNameEncrypted, &user.BirthDateEncrypted, &user.GenderEncrypted,
		&user.RoleEncrypted, &telephoneHash, &user.TelephoneEncrypted,
		&user.PostalCodeEncrypted, &user.CityEncrypted, &user.Address1Encrypted,
		&user.Address2Encrypted, &user.StripeCustomerIDEncrypted, &user.GoogleIDEncrypted, &user.AppleIDEncrypted,
		&user.CreatedAtEncrypted, &user.LoggedInAtEncrypted, &user.DEKEncrypted,
		&user.KeyVersion,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get user by Google ID", err)
	}

	// Handle nullable string fields
	if telephoneHash.Valid {
		user.TelephoneHash = telephoneHash.String
	}

	return user, nil
}