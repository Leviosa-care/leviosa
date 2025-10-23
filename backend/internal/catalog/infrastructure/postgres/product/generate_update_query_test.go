package productRepository

// package productRepository_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateUpdateQuery(t *testing.T) {

	productID := uuid.New().String()

	tests := []struct {
		name          string
		productID     string
		request       *domain.UpdateProductRequest
		expectedQuery string
		expectedArgs  []any
		expectedErr   string // Use string to check for error message
	}{
		{
			name:      "Update all supported fields",
			productID: productID,
			request: &domain.UpdateProductRequest{
				Name:        strPtr("New Name"),
				Description: strPtr("New Description"),
				Status:      statusStrPtr("published"), // Using string literal now
				Metadata:    map[string]any{"key": "value", "num": 123},
			},
			expectedQuery: "UPDATE catalog.products SET name = $1, description = $2, status = $3, metadata = $4 WHERE id = $5;",
			expectedArgs: func() []any {
				metadataJSON, _ := json.Marshal(map[string]any{"key": "value", "num": 123})
				return []any{"New Name", "New Description", "published", metadataJSON, productID} // Using string literal
			}(),
			expectedErr: "",
		},
		{
			name:      "Update only Name",
			productID: productID,
			request: &domain.UpdateProductRequest{
				Name: strPtr("Only Name"),
			},
			expectedQuery: "UPDATE catalog.products SET name = $1 WHERE id = $2;",
			expectedArgs:  []any{"Only Name", productID},
			expectedErr:   "",
		},
		{
			name:      "Update only Description",
			productID: productID,
			request: &domain.UpdateProductRequest{
				Description: strPtr("Only Description"),
			},
			expectedQuery: "UPDATE catalog.products SET description = $1 WHERE id = $2;",
			expectedArgs:  []any{"Only Description", productID},
			expectedErr:   "",
		},
		{
			name:      "Update only Status",
			productID: productID,
			request: &domain.UpdateProductRequest{
				Status: statusStrPtr("draft"), // Using string literal now
			},
			expectedQuery: "UPDATE catalog.products SET status = $1 WHERE id = $2;",
			expectedArgs:  []any{"draft", productID}, // Using string literal
			expectedErr:   "",
		},
		{
			name:      "Update only Metadata",
			productID: productID,
			request: &domain.UpdateProductRequest{
				Metadata: map[string]any{"new_key": "new_value"},
			},
			expectedQuery: "UPDATE catalog.products SET metadata = $1 WHERE id = $2;",
			expectedArgs: func() []any {
				metadataJSON, _ := json.Marshal(map[string]any{"new_key": "new_value"})
				return []any{metadataJSON, productID}
			}(),
			expectedErr: "",
		},
		{
			name:          "No fields provided for update (should return error)",
			productID:     productID,
			request:       &domain.UpdateProductRequest{}, // Empty request
			expectedQuery: "",
			expectedArgs:  nil,
			expectedErr:   "no fields provided for update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily override time.Now() for predictable updated_at
			query, args, err := generateUpdateQuery("catalog", tt.productID, tt.request)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, query)
				assert.Nil(t, args)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedQuery, query)

				// Compare arguments, especially the time.Time and JSON bytes
				require.Len(t, args, len(tt.expectedArgs), "Argument count mismatch")
				for i, expectedArg := range tt.expectedArgs {
					actualArg := args[i]
					switch v := expectedArg.(type) {
					case time.Time:
						assert.WithinDuration(t, v, actualArg.(time.Time), time.Millisecond, "Time argument mismatch at index %d", i)
					case []byte:
						var expectedMap, actualMap map[string]any
						json.Unmarshal(v, &expectedMap)
						json.Unmarshal(actualArg.([]byte), &actualMap)
						assert.Equal(t, expectedMap, actualMap, "JSON metadata argument mismatch at index %d", i)
					default:
						assert.Equal(t, expectedArg, actualArg, "Argument mismatch at index %d", i)
					}
				}
			}
		})
	}
}

// Helper to create a pointer to a string
func strPtr(s string) *string { return &s }

// Helper to create a pointer to a PublishedStatus
func statusStrPtr(s string) *string { return &s }
