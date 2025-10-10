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

	// Encrypt user data using the new generated function
	userEncx, err := domain.ProcessUserEncx(ctx, s.crypto, user)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("user during partner creation", err)
	}

	// Create user in database
	if err := s.userRepo.CreateUser(ctx, userEncx); err != nil {
		return nil, err
	}

	// Create partner entity
	partner := request.ToPartner(user.ID)

	// Encrypt partner data using the new generated function
	partnerEncx, err := domain.ProcessPartnerEncx(ctx, s.crypto, partner)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("partner during creation", err)
	}

	// Create partner in database
	if err := s.partnerRepo.CreatePartner(ctx, partnerEncx); err != nil {
		return nil, err
	}

	// Add specializations
	for _, specializationID := range request.SpecializationIDs {
		// Verify specialization exists and is active
		specEncx, err := s.specializationRepo.GetSpecializationByID(ctx, specializationID)
		if err != nil {
			if errors.Is(err, errs.ErrRepositoryNotFound) {
				return nil, errs.ErrInvalidInput
			}
			return nil, err
		}

		// Decrypt specialization data using the new generated function
		spec, err := domain.DecryptSpecializationEncx(ctx, s.crypto, specEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("specialization during partner creation", err)
		}

		if !spec.IsActive {
			return nil, errs.ErrInvalidInput
		}

		// Add association
		if err := s.partnerRepo.AddPartnerSpecialization(ctx, partnerEncx.ID, specializationID); err != nil {
			return nil, err
		}
	}

	// Get complete partner with user and specializations for response
	completePartnerEncx, err := s.partnerRepo.GetPartnerWithUser(ctx, partnerEncx.ID)
	if err != nil {
		return nil, err
	}

	// Get specializations
	specializationsEncx, err := s.partnerRepo.GetPartnerSpecializations(ctx, partnerEncx.ID)
	if err != nil {
		return nil, err
	}

	// Decrypt all data for response using the new generated functions
	completePartnerUser, err := domain.DecryptUserEncx(ctx, s.crypto, completePartnerEncx.User)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user in complete partner response", err)
	}

	completePartner, err := domain.DecryptPartnerEncx(ctx, s.crypto, completePartnerEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("complete partner response", err)
	}

	// Update the user reference in complete partner
	completePartner.User = completePartnerUser

	specResponses := make([]domain.SpecializationResponse, 0, len(specializationsEncx))
	for _, specEncx := range specializationsEncx {
		spec, err := domain.DecryptSpecializationEncx(ctx, s.crypto, specEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("specialization in partner response", err)
		}
		specResponses = append(specResponses, *spec.ToResponse())
	}

	completePartner.Specializations = specResponses
	return completePartner.ToCompleteResponse(), nil
}