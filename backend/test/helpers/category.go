package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// cleanCategoriesTable truncates the categories table after each test
func ClearCategoriesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE catalog.categories RESTART IDENTITY CASCADE;")
	require.NoError(t, err)
}

func NewValidCategory(name string) *domain.Category {
	return &domain.Category{
		ID:          uuid.New(),
		Name:        name,
		Description: fmt.Sprintf("Description for %s", name),
		Status:      domain.Draft, // Using a specific status for valid cases
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:   time.Now().UTC().Truncate(time.Microsecond),
	}
}

func InsertCategory(t *testing.T, ctx context.Context, category *domain.Category, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO catalog.categories (
			id,
			name,
			description,
			status,
			created_at,
			metadata
		) VALUES ($1, $2, $3, $4, $5, $6)`

	// Marshal metadata to JSONB for insertion
	metadataJSON, err := json.Marshal(category.Metadata)
	require.NoError(t, err, "Failed to marshal metadata for test insertion")

	_, err = pool.Exec(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.Status,
		category.CreatedAt,
		metadataJSON,
	)
	require.NoError(t, err, fmt.Sprintf("Failed to pre-insert category '%s'", category.Name))
}

// GetCategoryByID fetches a category from the database by its ID.
// This is a test helper function used to verify the state of the database.
func GetCategoryByID(t *testing.T, ctx context.Context, categoryID uuid.UUID, pool *pgxpool.Pool) (*domain.Category, error) {
	t.Helper()

	query := `
		SELECT
			id,
			name,
			description,
			status,
			created_at,
			updated_at,
			metadata
		FROM catalog.categories
		WHERE id = $1
	`

	var cat domain.Category
	var metadataJSON []byte

	err := pool.QueryRow(ctx, query, categoryID).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Description,
		&cat.Status,
		&cat.CreatedAt,
		&cat.UpdatedAt,
		&metadataJSON,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("category with ID %s not found", categoryID)
		}
		return nil, fmt.Errorf("failed to get category from database: %w", err)
	}

	if metadataJSON != nil {
		err = json.Unmarshal(metadataJSON, &cat.Metadata)
		require.NoError(t, err, "Failed to unmarshal metadata from database")
	}

	return &cat, nil
}
