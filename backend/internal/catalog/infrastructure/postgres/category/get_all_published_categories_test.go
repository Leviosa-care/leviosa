package categoryRepository_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPublishedCategories(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		setup         func(t *testing.T, ctx context.Context)
		expectedCount int
		expectedOrder []string // Expected order of category names
		expectedErr   error
	}{
		{
			name: "Successfully retrieve multiple published categories, ordered by name",
			setup: func(t *testing.T, ctx context.Context) {
				// Insert categories with various statuses
				pub1 := td.NewValidCategory("Beta Published")
				pub1.Status = domain.Published
				pub2 := td.NewValidCategory("Alpha Published")
				pub2.Status = domain.Published
				pub2.Metadata = map[string]any{"color": "green"}

				draft1 := td.NewValidCategory("Draft One")
				draft1.Status = domain.Draft

				archived1 := td.NewValidCategory("Archived One")
				archived1.Status = domain.Archived

				pub3 := td.NewValidCategory("Gamma Published")
				pub3.Status = domain.Published
				pub3.Metadata = map[string]any{"material": "wood"}

				td.InsertCategory(t, ctx, pub1, testPool)
				td.InsertCategory(t, ctx, draft1, testPool)
				td.InsertCategory(t, ctx, archived1, testPool)
				td.InsertCategory(t, ctx, pub2, testPool)
				td.InsertCategory(t, ctx, pub3, testPool)
			},
			expectedCount: 3,
			expectedOrder: []string{"Alpha Published", "Beta Published", "Gamma Published"},
			expectedErr:   nil,
		},
		{
			name: "Successfully retrieve no published categories when only others exist",
			setup: func(t *testing.T, ctx context.Context) {
				// Insert only non-published categories
				draft1 := td.NewValidCategory("Only Draft")
				draft1.Status = domain.Draft
				archived1 := td.NewValidCategory("Only Archived")
				archived1.Status = domain.Archived

				td.InsertCategory(t, ctx, draft1, testPool)
				td.InsertCategory(t, ctx, archived1, testPool)
			},
			expectedCount: 0,
			expectedOrder: []string{},
			expectedErr:   nil,
		},
		{
			name:          "Successfully retrieve no categories when table is empty",
			setup:         func(t *testing.T, ctx context.Context) {}, // Table is clean by default
			expectedCount: 0,
			expectedOrder: []string{},
			expectedErr:   nil,
		},
		// No test for malformed JSON, as discussed, DB handles this on insert.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearCategoriesTable(t, ctx, testPool)

			// Perform test-specific setup
			tt.setup(t, ctx)

			// Call the function under test
			categories, err := repo.GetAllPublishedCategories(ctx)

			if tt.expectedErr != nil {
				assert.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tt.expectedErr, "Expected specific error type")
				assert.Nil(t, categories, "Expected categories slice to be nil on error")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")

				// Ensure consistent return for empty results (assuming it returns an empty slice, not nil)
				if tt.expectedCount == 0 {
					assert.NotNil(t, categories, "Expected categories slice to be non-nil for empty result")
				} else {
					require.NotNil(t, categories, "Expected categories slice to be non-nil on success")
				}
				assert.Len(t, categories, tt.expectedCount, "Mismatched number of categories")

				// Verify order of names
				actualNames := make([]string, len(categories))
				for i, cat := range categories {
					actualNames[i] = cat.Name
				}
				assert.Equal(t, tt.expectedOrder, actualNames, "Categories not ordered as expected")

				// For the "multiple categories" test, perform deeper checks on data mapping
				if tt.name == "Successfully retrieve multiple published categories, ordered by name" {
					// Manually query the database for published categories in expected order for robust comparison
					var dbPublishedCategories []*domain.Category
					dbRows, err := testPool.Query(ctx, `
						SELECT id, name, description, status, image_key, metadata, created_at, updated_at
						FROM catalog.categories
						WHERE status = $1
						ORDER BY name ASC`, domain.Published)
					require.NoError(t, err)
					defer dbRows.Close()

					for dbRows.Next() {
						var (
							cat          domain.Category
							metadataJSON []byte
							imageKey     sql.NullString
						)
						err := dbRows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.Status, &imageKey, &metadataJSON, &cat.CreatedAt, &cat.UpdatedAt)
						require.NoError(t, err)

						if metadataJSON != nil {
							err := json.Unmarshal(metadataJSON, &cat.Metadata)
							require.NoError(t, err)
						} else {
							cat.Metadata = make(map[string]any)
						}
						dbPublishedCategories = append(dbPublishedCategories, &cat)
					}
					require.NoError(t, dbRows.Err())

					// Compare the categories returned by the function with those directly from DB
					require.Len(t, categories, len(dbPublishedCategories), "Length mismatch between retrieved and DB categories")

					for i := range categories {
						retrieved := categories[i]
						expected := dbPublishedCategories[i]

						assert.Equal(t, expected.ID, retrieved.ID, "ID mismatch at index %d", i)
						assert.Equal(t, expected.Name, retrieved.Name, "Name mismatch at index %d", i)
						assert.Equal(t, expected.Description, retrieved.Description, "Description mismatch at index %d", i)
						assert.Equal(t, expected.Status, retrieved.Status, "Status mismatch at index %d", i)

						assert.Equal(t, expected.Metadata, retrieved.Metadata, "Metadata mismatch at index %d", i)

						assert.WithinDuration(t, expected.CreatedAt, retrieved.CreatedAt, time.Millisecond, "CreatedAt mismatch at index %d", i)
						assert.WithinDuration(t, expected.UpdatedAt, retrieved.UpdatedAt, time.Millisecond, "UpdatedAt mismatch at index %d", i)
					}
				}
			}
		})
	}
}
