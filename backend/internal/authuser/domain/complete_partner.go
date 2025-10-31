package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// CompletePartnerRequest combines user profile completion with partner-specific information.
// This is used during partner registration where a user completes their profile
// and adds partner details in a single step.
type CompletePartnerRequest struct {
	// User profile fields (same as CompleteUserRequest)
	Password   string      `json:"password" validate:"required"`
	FirstName  string      `json:"first_name" validate:"required"`
	LastName   string      `json:"last_name" validate:"required"`
	BirthDate  time.Time   `json:"birth_date" validate:"required"`
	Gender     GenderInput `json:"gender" validate:"required"`
	Telephone  string      `json:"telephone" validate:"required"`
	PostalCode string      `json:"postal_code" validate:"required"`
	City       string      `json:"city" validate:"required"`
	Address1   string      `json:"address1" validate:"required"`
	Address2   string      `json:"address2"`

	// Partner-specific fields
	Bio        string `json:"bio,omitempty"`
	Experience string `json:"experience,omitempty"`
	// Certifications []string    `json:"certifications,omitempty"`
	CategoryIDs []uuid.UUID `json:"category_ids,omitempty"`
	ProductIDs  []uuid.UUID `json:"product_ids,omitempty"`
}

// Valid validates the complete partner request including both user and partner fields.
func (r *CompletePartnerRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate user fields using the same validation as CompleteUserRequest
	userRequest := &CompleteUserRequest{
		Password:   r.Password,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		BirthDate:  r.BirthDate,
		Gender:     r.Gender,
		Telephone:  r.Telephone,
		PostalCode: r.PostalCode,
		City:       r.City,
		Address1:   r.Address1,
		Address2:   r.Address2,
	}

	if err := userRequest.Valid(ctx); err != nil {
		// Include user validation errors
		errs.Set("user_validation", err.Error())
	}

	// Validate partner-specific fields
	if r.Bio != "" && len(r.Bio) > 1000 {
		errs.Set("bio", "bio must be 1000 characters or less")
	}

	if r.Experience != "" && len(r.Experience) > 2000 {
		errs.Set("experience", "experience must be 2000 characters or less")
	}

	// Validate certifications array
	// for _, cert := range r.Certifications {
	//	if len(cert) > 200 {
	//		errs.Set("certifications", "each certification must be 200 characters or less")
	//		break
	//	}
	//	if cert == "" {
	//		errs.Set("certifications", "certifications cannot contain empty values")
	//		break
	//	}
	// }

	// Validate UUID arrays (nil UUIDs not allowed)
	for _, categoryID := range r.CategoryIDs {
		if categoryID == uuid.Nil {
			errs.Set("category_ids", "category IDs cannot contain nil values")
			break
		}
	}

	for _, productID := range r.ProductIDs {
		if productID == uuid.Nil {
			errs.Set("product_ids", "product IDs cannot contain nil values")
			break
		}
	}

	return errs.AsError()
}

// ToCompleteUserRequest extracts the user-specific fields for user completion.
func (r *CompletePartnerRequest) ToCompleteUserRequest() *CompleteUserRequest {
	return &CompleteUserRequest{
		Password:   r.Password,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		BirthDate:  r.BirthDate,
		Gender:     r.Gender,
		Telephone:  r.Telephone,
		PostalCode: r.PostalCode,
		City:       r.City,
		Address1:   r.Address1,
		Address2:   r.Address2,
	}
}
