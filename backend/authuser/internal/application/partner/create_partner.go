package partner

import (
	"context"
	"errors"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
)

func (s *PartnerService) CreatePartner(ctx context.Context, request *domain.CreatePartnerRequest) (*domain.CompletePartnerResponse, error) {
	// Validate request
	if err := request.Valid(ctx); err != nil {
		return nil, errs.ErrInvalidInput
	}

	// Check if user with email already exists
	exists, err := s.userRepo.ExistsByEmailHash(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.ErrUniqueViolation
	}

	// Create user entity
	user, err := request.ToUser()
	if err != nil {
		return nil, errs.ErrInvalidInput
	}

	// Encrypt user data
	if err := s.crypto.Encrypt(ctx, user); err != nil {
		return nil, errs.ErrInvalidValue
	}

	// Create user in database
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Create partner entity
	partner := request.ToPartner(user.ID)

	// Encrypt partner data
	if err := s.crypto.Encrypt(ctx, partner); err != nil {
		return nil, errs.ErrInvalidValue
	}

	// Create partner in database
	if err := s.partnerRepo.CreatePartner(ctx, partner); err != nil {
		return nil, err
	}

	// Add specializations
	for _, specializationID := range request.SpecializationIDs {
		// Verify specialization exists and is active
		spec, err := s.specializationRepo.GetSpecializationByID(ctx, specializationID)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				return nil, errs.ErrInvalidInput
			}
			return nil, err
		}
		if !spec.IsActive {
			return nil, errs.ErrInvalidInput
		}

		// Add association
		if err := s.partnerRepo.AddPartnerSpecialization(ctx, partner.ID, specializationID); err != nil {
			return nil, err
		}
	}

	// Get complete partner with user and specializations for response
	completePartner, err := s.partnerRepo.GetPartnerWithUser(ctx, partner.ID)
	if err != nil {
		return nil, err
	}

	// Get specializations
	specializations, err := s.partnerRepo.GetPartnerSpecializations(ctx, partner.ID)
	if err != nil {
		return nil, err
	}

	// Decrypt all data for response
	if err := s.crypto.Decrypt(ctx, completePartner.User); err != nil {
		return nil, err
	}
	if err := s.crypto.Decrypt(ctx, completePartner); err != nil {
		return nil, err
	}

	specResponses := make([]domain.SpecializationResponse, 0, len(specializations))
	for _, spec := range specializations {
		if err := s.crypto.Decrypt(ctx, spec); err != nil {
			return nil, err
		}
		specResponses = append(specResponses, *spec.ToResponse())
	}

	completePartner.Specializations = specResponses
	return completePartner.ToCompleteResponse(), nil
}