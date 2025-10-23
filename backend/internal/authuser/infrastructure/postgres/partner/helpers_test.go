package partnerRepository_test

import (
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func NewTestPartnerEncx() *domain.PartnerEncx {
	now := time.Now()
	return &domain.PartnerEncx{
		ID:                      uuid.New(),
		UserID:                  uuid.New(),
		IsVerified:              true,
		CreatedAt:               now,
		UpdatedAt:               now,
		BioEncrypted:            []byte("bio encrypted"),
		CertificationsEncrypted: []byte("certifications encrypted"),
		ExperienceEncrypted:     []byte("experience encrypted"),
		DEKEncrypted:            []byte("dek encrypted"),
		Metadata:                encx.EncryptionMetadata{},
	}
}
