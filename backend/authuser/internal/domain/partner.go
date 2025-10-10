package domain

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type Partner struct {
	ID               uuid.UUID       `json:"id"`
	UserID           uuid.UUID       `json:"user_id"`
	Bio              string          `json:"bio" encx:"encrypt"`
	Experience       string          `json:"experience" encx:"encrypt"`
	Certifications   []string        `json:"certifications" encx:"encrypt"`
	IsVerified       bool            `json:"is_verified"`
	VerifiedAt       *time.Time      `json:"verified_at" encx:"encrypt"`
	VerifiedByUserID *uuid.UUID      `json:"verified_by_user_id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`

	// Embedded user data (populated when joining with users table)
	User           *User                      `json:"user,omitempty"`
	Specializations []SpecializationResponse `json:"specializations,omitempty"`
}

func (p *Partner) Valid(ctx context.Context) error {
	var errs errsx.Map

	// UserID is required
	if p.UserID == uuid.Nil {
		errs.Set("user_id", "user ID is required")
	}

	// Bio length validation
	bio := strings.TrimSpace(p.Bio)
	if len(bio) > 1000 {
		errs.Set("bio", "bio must be 1000 characters or less")
	}

	// Experience length validation
	experience := strings.TrimSpace(p.Experience)
	if len(experience) > 2000 {
		errs.Set("experience", "experience must be 2000 characters or less")
	}

	// Certifications validation
	if len(p.Certifications) > 20 {
		errs.Set("certifications", "maximum 20 certifications allowed")
	}

	for _, cert := range p.Certifications {
		cert = strings.TrimSpace(cert)
		if cert == "" {
			errs.Set("certifications", "certification cannot be empty")
		} else if len(cert) > 200 {
			errs.Set("certifications", "each certification must be 200 characters or less")
		}
	}

	return errs.AsError()
}

func (p *Partner) ToResponse() *PartnerResponse {
	resp := &PartnerResponse{
		ID:               p.ID,
		UserID:           p.UserID,
		Bio:              p.Bio,
		Experience:       p.Experience,
		Certifications:   p.Certifications,
		IsVerified:       p.IsVerified,
		VerifiedAt:       p.VerifiedAt,
		VerifiedByUserID: p.VerifiedByUserID,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
		Specializations:  p.Specializations,
	}

	// Include user data if available
	if p.User != nil {
		resp.User = p.User.ToResponse()
	}

	return resp
}

// ToCompleteResponse returns a response with both partner and user data
func (p *Partner) ToCompleteResponse() *CompletePartnerResponse {
	resp := &CompletePartnerResponse{
		ID:               p.ID,
		UserID:           p.UserID,
		Bio:              p.Bio,
		Experience:       p.Experience,
		Certifications:   p.Certifications,
		IsVerified:       p.IsVerified,
		VerifiedAt:       p.VerifiedAt,
		VerifiedByUserID: p.VerifiedByUserID,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
		Specializations:  p.Specializations,
	}

	// Include user data if available
	if p.User != nil {
		resp.User = p.User.ToResponse()
	}

	return resp
}

// MarshalCertifications converts certifications slice to JSON for database storage
func (p *Partner) MarshalCertifications() ([]byte, error) {
	if len(p.Certifications) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(p.Certifications)
}

// UnmarshalCertifications converts JSON from database to certifications slice
func (p *Partner) UnmarshalCertifications(data []byte) error {
	if len(data) == 0 {
		p.Certifications = []string{}
		return nil
	}
	return json.Unmarshal(data, &p.Certifications)
}