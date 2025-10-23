package productRepository

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// generateUpdateQuery builds the dynamic SQL UPDATE statement and its arguments
// based on the non-nil fields in the UpdateProductRequest.
// It returns the query string, a slice of arguments, and an error if validation fails.
func generateUpdateQuery(schema, productID string, req *domain.UpdateProductRequest) (string, []any, error) {
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
	if req.CategoryID != nil {
		// CategoryID is a *string in UpdateProductRequest, but UUID in DB
		catUUID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return "", nil, fmt.Errorf("invalid category ID format: %w", err)
		}
		sets = append(sets, fmt.Sprintf("category_id = $%d", argCounter))
		args = append(args, catUUID)
		argCounter++
	}
	if req.Duration != nil {
		sets = append(sets, fmt.Sprintf("duration = $%d", argCounter))
		args = append(args, *req.Duration)
		argCounter++
	}
	if req.Availability != nil {
		sets = append(sets, fmt.Sprintf("availability = $%d", argCounter))
		args = append(args, *req.Availability)
		argCounter++
	}
	if req.BufferTime != nil {
		sets = append(sets, fmt.Sprintf("buffer_time = $%d", argCounter))
		args = append(args, *req.BufferTime)
		argCounter++
	}
	if req.CancellationHours != nil {
		sets = append(sets, fmt.Sprintf("cancellation_hours = $%d", argCounter))
		args = append(args, *req.CancellationHours)
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
	// The product ID is always the last parameter in the WHERE clause
	query := fmt.Sprintf("UPDATE %s.products SET %s WHERE id = $%d;",
		schema, strings.Join(sets, ", "), argCounter)
	args = append(args, productID) // Add the product ID to the arguments

	return query, args, nil
}
