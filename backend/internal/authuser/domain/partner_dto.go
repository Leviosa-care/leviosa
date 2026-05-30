package domain

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type PartnerResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Bio        string    `json:"bio"`
	Experience string    `json:"experience"`
	Occupation string    `json:"occupation"`
	Quote      string    `json:"quote"`
	Tags       []string  `json:"tags"`
	// Certifications []string    `json:"certifications"`
	CategoryIDs             []uuid.UUID       `json:"category_ids,omitempty"`
	ProductIDs              []uuid.UUID       `json:"product_ids,omitempty"`
	StripeAccountStatus     StripeAccountStatus `json:"stripe_account_status"`
	StripeOnboardingComplete bool              `json:"stripe_onboarding_complete"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
}
type UpdatePartnerRequest struct {
	Bio        *string  `json:"bio,omitempty"`
	Experience *string  `json:"experience,omitempty"`
	Occupation *string  `json:"occupation,omitempty"`
	Quote      *string  `json:"quote,omitempty"`
	Tags       *[]string `json:"tags,omitempty"`
	// Certifications *[]string `json:"certifications,omitempty"`
}

func (r *UpdatePartnerRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Bio validation if provided
	if r.Bio != nil {
		bio := strings.TrimSpace(*r.Bio)
		if len(bio) > 1000 {
			errs.Set("bio", "bio must be 1000 characters or less")
		}
	}

	// Experience validation if provided
	if r.Experience != nil {
		experience := strings.TrimSpace(*r.Experience)
		if len(experience) > 2000 {
			errs.Set("experience", "experience must be 2000 characters or less")
		}
	}

	// Occupation validation if provided
	if r.Occupation != nil {
		occupation := strings.TrimSpace(*r.Occupation)
		if len(occupation) > 200 {
			errs.Set("occupation", "occupation must be 200 characters or less")
		}
	}

	// Quote validation if provided
	if r.Quote != nil {
		quote := strings.TrimSpace(*r.Quote)
		if len(quote) > 300 {
			errs.Set("quote", "quote must be 300 characters or less")
		}
	}

	// Tags validation if provided
	if r.Tags != nil {
		if len(*r.Tags) > 20 {
			errs.Set("tags", "maximum 20 tags allowed")
		}
		for _, tag := range *r.Tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				errs.Set("tags", "tag cannot be empty")
			} else if len(tag) > 100 {
				errs.Set("tags", "each tag must be 100 characters or less")
			}
		}
	}

	// Certifications validation if provided
	// if r.Certifications != nil {
	//	if len(*r.Certifications) > 20 {
	//		errs.Set("certifications", "maximum 20 certifications allowed")
	//	}
	//
	//	for _, cert := range *r.Certifications {
	//		cert = strings.TrimSpace(cert)
	//		if cert == "" {
	//			errs.Set("certifications", "certification cannot be empty")
	//		} else if len(cert) > 200 {
	//			errs.Set("certifications", "each certification must be 200 characters or less")
	//		}
	//	}
	// }

	return errs.AsError()
}

// PublicPartnerResponse is the DTO returned by the unauthenticated GET /partners endpoint.
// It embeds the partner's public fields alongside the linked user's name and picture.
type PublicPartnerResponse struct {
	ID          uuid.UUID  `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Picture     string     `json:"picture,omitempty"`
	Bio         string     `json:"bio"`
	Experience  string     `json:"experience"`
	Occupation  string     `json:"occupation"`
	Quote       string     `json:"quote"`
	Tags        []string   `json:"tags"`
	CategoryIDs []uuid.UUID `json:"category_ids,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// PublicPartnerRow holds the raw row from the public partners query,
// which JOINs partners with users. The user fields are encrypted
// and will be decrypted by the application layer.
type PublicPartnerRow struct {
	PartnerEncx              *PartnerEncx
	FirstNameEncrypted       []byte
	LastNameEncrypted        []byte
	PictureEncrypted         []byte
	UserDEKEncrypted         []byte
	UserKeyVersion           int
}

type VerifyPartnerRequest struct {
	PartnerID uuid.UUID `json:"partner_id"`
}

func (r *VerifyPartnerRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if r.PartnerID == uuid.Nil {
		errs.Set("partner_id", "partner ID is required")
	}

	return errs.AsError()
}
