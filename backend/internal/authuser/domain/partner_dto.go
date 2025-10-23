package domain

import (
	"context"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type PartnerResponse struct {
	ID               uuid.UUID                `json:"id"`
	UserID           uuid.UUID                `json:"user_id"`
	Bio              string                   `json:"bio"`
	Experience       string                   `json:"experience"`
	Certifications   []string                 `json:"certifications"`
	CategoryIDs      []uuid.UUID              `json:"category_ids,omitempty"`
	ProductIDs       []uuid.UUID              `json:"product_ids,omitempty"`
	IsVerified       bool                     `json:"is_verified"`
	VerifiedAt       *time.Time               `json:"verified_at,omitempty"`
	VerifiedByUserID *uuid.UUID               `json:"verified_by_user_id,omitempty"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
	User             *UserResponse            `json:"user,omitempty"`
	Specializations  []SpecializationResponse `json:"specializations,omitempty"`
}

type CompletePartnerResponse struct {
	ID               uuid.UUID                `json:"id"`
	UserID           uuid.UUID                `json:"user_id"`
	Bio              string                   `json:"bio"`
	Experience       string                   `json:"experience"`
	Certifications   []string                 `json:"certifications"`
	CategoryIDs      []uuid.UUID              `json:"category_ids,omitempty"`
	ProductIDs       []uuid.UUID              `json:"product_ids,omitempty"`
	IsVerified       bool                     `json:"is_verified"`
	VerifiedAt       *time.Time               `json:"verified_at,omitempty"`
	VerifiedByUserID *uuid.UUID               `json:"verified_by_user_id,omitempty"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
	User             *UserResponse            `json:"user"`
	Specializations  []SpecializationResponse `json:"specializations"`
}

// NOTE: This DTO is deprecated. Partner registration now uses CompletePartnerRequest via /auth/complete/partner
// Keeping this for potential future admin-initiated partner creation, but it's currently unused.
type CreatePartnerRequest struct {
	UserID         string        `json:"user_id"`          // ID of existing user to create partner profile for
	Bio            string        `json:"bio,omitempty"`
	Experience     string        `json:"experience,omitempty"`
	Certifications []string      `json:"certifications,omitempty"`
	CategoryIDs    []uuid.UUID   `json:"category_ids,omitempty"`
	ProductIDs     []uuid.UUID   `json:"product_ids,omitempty"`
}

func (r *CreatePartnerRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate user ID
	if strings.TrimSpace(r.UserID) == "" {
		errs.Set("user_id", "user ID is required")
	} else if _, err := uuid.Parse(r.UserID); err != nil {
		errs.Set("user_id", "invalid user ID format")
	}

	// Validate partner-specific fields
	bio := strings.TrimSpace(r.Bio)
	if len(bio) > 1000 {
		errs.Set("bio", "bio must be 1000 characters or less")
	}

	experience := strings.TrimSpace(r.Experience)
	if len(experience) > 2000 {
		errs.Set("experience", "experience must be 2000 characters or less")
	}

	if len(r.Certifications) > 20 {
		errs.Set("certifications", "maximum 20 certifications allowed")
	}

	for _, cert := range r.Certifications {
		cert = strings.TrimSpace(cert)
		if cert == "" {
			errs.Set("certifications", "certification cannot be empty")
		} else if len(cert) > 200 {
			errs.Set("certifications", "each certification must be 200 characters or less")
		}
	}

	return errs.AsError()
}

// NOTE: the old way with the old partner API
// type CreatePartnerRequest struct {
// 	// User fields (required for partner creation)
// 	Email      string `json:"email"`
// 	Password   string `json:"password"`
// 	FirstName  string `json:"first_name"`
// 	LastName   string `json:"last_name"`
// 	Telephone  string `json:"telephone"`
// 	BirthDate  string `json:"birthdate"` // ISO format
// 	Gender     string `json:"gender"`
// 	PostalCode string `json:"postal_code"`
// 	City       string `json:"city"`
// 	Address1   string `json:"address1"`
// 	Address2   string `json:"address2,omitempty"`
//
// 	// Partner-specific fields
// 	Bio               string      `json:"bio,omitempty"`
// 	Experience        string      `json:"experience,omitempty"`
// 	Certifications    []string    `json:"certifications,omitempty"`
// 	SpecializationIDs []uuid.UUID `json:"specialization_ids,omitempty"`
// }
// func (r *CreatePartnerRequest) Valid(ctx context.Context) error {
// 	var errs errsx.Map
//
// 	// Validate user fields
// 	if strings.TrimSpace(r.Email) == "" {
// 		errs.Set("email", "email is required")
// 	}
//
// 	if strings.TrimSpace(r.Password) == "" {
// 		errs.Set("password", "password is required")
// 	} else if err := ValidatePassword(r.Password); err != nil {
// 		errs.Set("password", err)
// 	}
//
// 	if strings.TrimSpace(r.FirstName) == "" {
// 		errs.Set("first_name", "first name is required")
// 	}
//
// 	if strings.TrimSpace(r.LastName) == "" {
// 		errs.Set("last_name", "last name is required")
// 	}
//
// 	if strings.TrimSpace(r.Telephone) == "" {
// 		errs.Set("telephone", "telephone is required")
// 	}
//
// 	if strings.TrimSpace(r.BirthDate) == "" {
// 		errs.Set("birthdate", "birthdate is required")
// 	} else {
// 		if _, err := time.Parse("2006-01-02", r.BirthDate); err != nil {
// 			errs.Set("birthdate", "birthdate must be in YYYY-MM-DD format")
// 		}
// 	}
//
// 	if strings.TrimSpace(r.Gender) == "" {
// 		errs.Set("gender", "gender is required")
// 	}
//
// 	if strings.TrimSpace(r.PostalCode) == "" {
// 		errs.Set("postal_code", "postal code is required")
// 	}
//
// 	if strings.TrimSpace(r.City) == "" {
// 		errs.Set("city", "city is required")
// 	}
//
// 	if strings.TrimSpace(r.Address1) == "" {
// 		errs.Set("address1", "address is required")
// 	}
//
// 	// Validate partner-specific fields
// 	bio := strings.TrimSpace(r.Bio)
// 	if len(bio) > 1000 {
// 		errs.Set("bio", "bio must be 1000 characters or less")
// 	}
//
// 	experience := strings.TrimSpace(r.Experience)
// 	if len(experience) > 2000 {
// 		errs.Set("experience", "experience must be 2000 characters or less")
// 	}
//
// 	if len(r.Certifications) > 20 {
// 		errs.Set("certifications", "maximum 20 certifications allowed")
// 	}
//
// 	for _, cert := range r.Certifications {
// 		cert = strings.TrimSpace(cert)
// 		if cert == "" {
// 			errs.Set("certifications", "certification cannot be empty")
// 		} else if len(cert) > 200 {
// 			errs.Set("certifications", "each certification must be 200 characters or less")
// 		}
// 	}
//
// 	if len(r.SpecializationIDs) == 0 {
// 		errs.Set("specialization_ids", "at least one specialization is required")
// 	}
//
// 	for _, id := range r.SpecializationIDs {
// 		if id == uuid.Nil {
// 			errs.Set("specialization_ids", "invalid specialization ID")
// 		}
// 	}
//
// 	return errs.AsError()
// }

// TODO: that thing should be deleted
func (r *CreatePartnerRequest) ToUser() (*User, error) {
	birthDate, err := time.Parse("2006-01-02", r.BirthDate)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:         uuid.New(),
		State:      Unverified,
		Email:      strings.TrimSpace(r.Email),
		Password:   r.Password,
		FirstName:  strings.TrimSpace(r.FirstName),
		LastName:   strings.TrimSpace(r.LastName),
		Telephone:  strings.TrimSpace(r.Telephone),
		BirthDate:  birthDate,
		Gender:     strings.TrimSpace(r.Gender),
		PostalCode: strings.TrimSpace(r.PostalCode),
		City:       strings.TrimSpace(r.City),
		Address1:   strings.TrimSpace(r.Address1),
		Address2:   strings.TrimSpace(r.Address2),
		Role:       identity.PartnerStr,
	}, nil
}

func (r *CreatePartnerRequest) ToPartner(userID uuid.UUID) *Partner {
	// Clean certifications
	cleanCertifications := make([]string, 0, len(r.Certifications))
	for _, cert := range r.Certifications {
		if clean := strings.TrimSpace(cert); clean != "" {
			cleanCertifications = append(cleanCertifications, clean)
		}
	}

	return &Partner{
		ID:             uuid.New(),
		UserID:         userID,
		Bio:            strings.TrimSpace(r.Bio),
		Experience:     strings.TrimSpace(r.Experience),
		Certifications: cleanCertifications,
		CategoryIDs:    r.CategoryIDs,
		ProductIDs:     r.ProductIDs,
		IsVerified:     false, // Partners start unverified
	}
}

type UpdatePartnerRequest struct {
	Bio            *string   `json:"bio,omitempty"`
	Experience     *string   `json:"experience,omitempty"`
	Certifications *[]string `json:"certifications,omitempty"`
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
	if r.Certifications != nil {
		if len(*r.Certifications) > 20 {
			errs.Set("certifications", "maximum 20 certifications allowed")
		}

		for _, cert := range *r.Certifications {
			cert = strings.TrimSpace(cert)
			if cert == "" {
				errs.Set("certifications", "certification cannot be empty")
			} else if len(cert) > 200 {
				errs.Set("certifications", "each certification must be 200 characters or less")
			}
		}
	}

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

type GetPartnersResponse struct {
	Partners []CompletePartnerResponse `json:"partners"`
	Total    int                       `json:"total"`
}

type GetPartnerSpecializationsResponse struct {
	Specializations []SpecializationResponse `json:"specializations"`
}
