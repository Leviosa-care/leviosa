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
	// Certifications []string    `json:"certifications"`
	CategoryIDs             []uuid.UUID       `json:"category_ids,omitempty"`
	ProductIDs              []uuid.UUID       `json:"product_ids,omitempty"`
	StripeAccountStatus     StripeAccountStatus `json:"stripe_account_status"`
	StripeOnboardingComplete bool              `json:"stripe_onboarding_complete"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
}
type UpdatePartnerRequest struct {
	Bio        *string `json:"bio,omitempty"`
	Experience *string `json:"experience,omitempty"`
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
