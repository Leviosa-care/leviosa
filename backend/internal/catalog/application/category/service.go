package category

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

// check if *ServiceImpl implements Service interface
var _ ports.CategoryService = (*CategoryService)(nil)

type CategoryService struct {
	repo       ports.CategoryRepository
	sharedRepo ports.SharedRepository
}

func New(repo ports.CategoryRepository, sharedRepo ports.SharedRepository) ports.CategoryService {
	return &CategoryService{
		repo:       repo,
		sharedRepo: sharedRepo,
	}
}
