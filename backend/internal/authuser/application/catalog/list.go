package catalog

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
)

// TODO: do the implementation for this

func (s *Service) ListPublishedCategories(ctx context.Context) ([]domain.CachedCategory, error) {
	return []domain.CachedCategory{}, nil
}

func (s *Service) ListPublishedProducts(ctx context.Context) ([]domain.CachedProduct, error) {
	return []domain.CachedProduct{}, nil
}
