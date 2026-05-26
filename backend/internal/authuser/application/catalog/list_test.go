package catalog

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ListPublishedCategories(t *testing.T) {
	t.Run("returns empty slice when cache is cold", func(t *testing.T) {
		cache := NewCatalogCache()
		svc := &Service{cache: cache}

		result, err := svc.ListPublishedCategories(context.Background())
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("returns all seeded categories", func(t *testing.T) {
		cache := NewCatalogCache()
		ctx := context.Background()

		cat1 := &domain.CachedCategory{
			ID:          uuid.New(),
			Name:        "Cat A",
			Description: "Desc A",
			Status:      "published",
		}
		cat2 := &domain.CachedCategory{
			ID:          uuid.New(),
			Name:        "Cat B",
			Description: "Desc B",
			Status:      "published",
		}

		require.NoError(t, cache.UpsertCategory(ctx, cat1))
		require.NoError(t, cache.UpsertCategory(ctx, cat2))

		svc := &Service{cache: cache}
		result, err := svc.ListPublishedCategories(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)

		names := map[string]bool{}
		for _, c := range result {
			names[c.Name] = true
		}
		assert.True(t, names["Cat A"])
		assert.True(t, names["Cat B"])
	})
}

func TestService_ListPublishedProducts(t *testing.T) {
	t.Run("returns empty slice when cache is cold", func(t *testing.T) {
		cache := NewCatalogCache()
		svc := &Service{cache: cache}

		result, err := svc.ListPublishedProducts(context.Background())
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("returns all seeded products", func(t *testing.T) {
		cache := NewCatalogCache()
		ctx := context.Background()

		categoryID := uuid.New()
		prod1 := &domain.CachedProduct{
			ID:          uuid.New(),
			Name:        "Prod A",
			CategoryID:  categoryID,
			Status:      "published",
		}
		prod2 := &domain.CachedProduct{
			ID:          uuid.New(),
			Name:        "Prod B",
			CategoryID:  categoryID,
			Status:      "published",
		}

		require.NoError(t, cache.UpsertProduct(ctx, prod1))
		require.NoError(t, cache.UpsertProduct(ctx, prod2))

		svc := &Service{cache: cache}
		result, err := svc.ListPublishedProducts(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)

		names := map[string]bool{}
		for _, p := range result {
			names[p.Name] = true
		}
		assert.True(t, names["Prod A"])
		assert.True(t, names["Prod B"])
	})
}
