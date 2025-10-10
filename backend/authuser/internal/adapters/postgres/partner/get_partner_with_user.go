package partnerRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetPartnerWithUser(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			p.id, p.user_id, p.bio_encrypted, p.experience_encrypted, p.certifications_encrypted,
			p.is_verified, p.verified_at_encrypted, p.verified_by_user_id,
			p.dek_encrypted, p.key_version, p.created_at, p.updated_at,
			u.id, u.state, u.email_hash, u.email_encrypted, u.password_hash_secure,
			u.picture_encrypted, u.first_name_encrypted, u.last_name_encrypted,
			u.birth_date_encrypted, u.gender_encrypted, u.role_encrypted,
			u.telephone_hash, u.telephone_encrypted, u.postal_code_encrypted,
			u.city_encrypted, u.address1_encrypted, u.address2_encrypted, u.stripe_customer_id_encrypted,
			u.google_id_encrypted, u.apple_id_encrypted, u.created_at_encrypted,
			u.logged_in_at_encrypted, u.dek_encrypted, u.key_version
		FROM %s.partners p
		INNER JOIN %s.users u ON p.user_id = u.id
		WHERE p.id = $1
	`, r.schema, r.schema)

	partner := &domain.PartnerEncx{}
	user := &domain.UserEncx{}

	err := r.pool.QueryRow(ctx, query, partnerID).Scan(
		&partner.ID,
		&partner.UserID,
		&partner.BioEncrypted,
		&partner.ExperienceEncrypted,
		&partner.CertificationsEncrypted,
		&partner.IsVerified,
		&partner.VerifiedAtEncrypted,
		&partner.VerifiedByUserID,
		&partner.DEKEncrypted,
		&partner.KeyVersion,
		&partner.CreatedAt,
		&partner.UpdatedAt,
		&user.ID,
		&user.State,
		&user.EmailHash,
		&user.EmailEncrypted,
		&user.PasswordHashSecure,
		&user.PictureEncrypted,
		&user.FirstNameEncrypted,
		&user.LastNameEncrypted,
		&user.BirthDateEncrypted,
		&user.GenderEncrypted,
		&user.RoleEncrypted,
		&user.TelephoneHash,
		&user.TelephoneEncrypted,
		&user.PostalCodeEncrypted,
		&user.CityEncrypted,
		&user.Address1Encrypted,
		&user.Address2Encrypted,
		&user.GoogleIDEncrypted,
		&user.AppleIDEncrypted,
		&user.StripeCustomerIDEncrypted,
		&user.CreatedAtEncrypted,
		&user.LoggedInAtEncrypted,
		&user.DEKEncrypted,
		&user.KeyVersion,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get partner with user", err)
	}

	partner.User = user
	return partner, nil
}