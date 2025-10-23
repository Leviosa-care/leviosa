package categoryRepository

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// generateUpdateQuery builds the dynamic SQL UPDATE statement and its arguments
// based on the non-nil fields in the UpdateCategoryRequest.
// It returns the query string, a slice of arguments, and an error if validation fails.
func generateUpdateQuery(categoryID string, req *domain.UpdateCategoryRequest) (string, []any, error) {
	// Start building the SET clauses and arguments
	sets := []string{}
	args := []any{}
	argCounter := 1 // For parameterized query placeholders ($1, $2, etc.)

	if req.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, *req.Name)
		argCounter++
	}
	if req.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argCounter))
		args = append(args, *req.Description)
		argCounter++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, *req.Status)
		argCounter++
	}
	// For map[string]any like Metadata, you'd typically store it as JSONB in PostgreSQL.
	// Marshal it to JSON before adding to args.
	if req.Metadata != nil {
		metadataJSON, err := json.Marshal(req.Metadata) // Requires "encoding/json" import
		if err != nil {
			return "", nil, fmt.Errorf("failed to marshal metadata to JSON: %w", err)
		}
		sets = append(sets, fmt.Sprintf("metadata = $%d", argCounter))
		args = append(args, metadataJSON)
		argCounter++
	}

	// If no fields were provided for update (only updated_at would be set)
	if len(sets) == 0 {
		return "", nil, errs.ErrNoFieldsForUpdate
	}

	// Construct the final query
	// The category ID is always the last parameter in the WHERE clause
	query := fmt.Sprintf("UPDATE catalog.categories SET %s WHERE id = $%d;",
		strings.Join(sets, ", "), argCounter)
	args = append(args, categoryID) // Add the category ID to the arguments

	return query, args, nil
}
