package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/hengadev/encx"
)

func (s *PartnerService) GetPublicPartners(ctx context.Context) ([]*domain.PublicPartnerResponse, error) {
	rows, err := s.partnerRepo.GetPublicPartners(ctx)
	if err != nil {
		return nil, fmt.Errorf("get public partners: %w", err)
	}

	partners := make([]*domain.PublicPartnerResponse, 0, len(rows))
	for _, row := range rows {
		userDEK, err := s.crypto.DecryptDEKWithVersion(ctx, row.UserDEKEncrypted, row.UserKeyVersion)
		if err != nil {
			return nil, fmt.Errorf("decrypt user DEK: %w", err)
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

		partners = append(partners, &domain.PublicPartnerResponse{
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
		})
	}

	return partners, nil
}

// decryptStringField decrypts an encrypted byte slice using the given DEK
// and deserializes it as a string. Returns empty string if data is nil/empty.
func decryptStringField(ctx context.Context, crypto encx.CryptoService, dek []byte, encrypted []byte) (string, error) {
	if len(encrypted) == 0 {
		return "", nil
	}

	decrypted, err := crypto.DecryptData(ctx, encrypted, dek)
	if err != nil {
		return "", err
	}

	var result string
	if err := encx.DeserializeValue(decrypted, &result); err != nil {
		return "", err
	}

	return result, nil
}
