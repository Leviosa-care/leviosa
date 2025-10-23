package image

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

// compile time assertion check if *ServiceImpl implements Service interface
var _ ports.ImageService = (*ImageService)(nil)

// Service defines the public-facing image management business logic.
type ImageService struct {
	repo       ports.ImageRepository
	mediaRepo  ports.ImageMedia
	sharedRepo ports.SharedRepository
}

func New(repo ports.ImageRepository, mediaRepo ports.ImageMedia, sharedRepo ports.SharedRepository) ports.ImageService {
	return &ImageService{
		repo:       repo,
		mediaRepo:  mediaRepo,
		sharedRepo: sharedRepo,
	}
}
