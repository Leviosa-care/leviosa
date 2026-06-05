package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *PartnerService) GetPublicPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PublicPartnerResponse, error) {
	row, err := s.partnerRepo.GetPublicPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("get public partner by ID: %w", err)
	}

	userDEK, err := s.crypto.DecryptDEKWithVersion(ctx, row.UserDEKEncrypted, row.UserKeyVersion)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user DEK", err)
	}

	firstName, err := decryptStringField(ctx, s.crypto, userDEK, row.FirstNameEncrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt first name: %w", err)
	}

	lastName, err := decryptStringField(ctx, s.crypto, userDEK, row.LastNameEncrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt last name: %w", err)
	}

	picture, _ := decryptStringField(ctx, s.crypto, userDEK, row.PictureEncrypted)

	p := row.PartnerEncx
	tags := p.Tags
	if tags == nil {
		tags = []string{}
	}

	return &domain.PublicPartnerResponse{
		ID:          p.ID,
		FirstName:   firstName,
		LastName:    lastName,
		Picture:     picture,
		Bio:         p.Bio,
		Experience:  p.Experience,
		Occupation:  p.Occupation,
		Quote:       p.Quote,
		Tags:        tags,
		CategoryIDs: p.CategoryIDs,
		CreatedAt:   p.CreatedAt,
	}, nil
}
