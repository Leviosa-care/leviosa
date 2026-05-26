package catalog

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
)

func (s *Service) ListPublishedCategories(ctx context.Context) ([]domain.CachedCategory, error) {
	return s.cache.ListCategories(), nil
}

func (s *Service) ListPublishedProducts(ctx context.Context) ([]domain.CachedProduct, error) {
	return s.cache.ListProducts(), nil
}
