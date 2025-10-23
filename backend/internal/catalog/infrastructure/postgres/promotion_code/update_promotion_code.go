package promotionCodeRepository

import (
	"fmt"
	"context"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *PromotionCodeRepository) UpdatePromotionCode(ctx context.Context, promotionCodeID uuid.UUID, req *domain.UpdatePromotionCodeRequest) error {
	var setParts []string
	var args []interface{}
	argIndex := 1

	// Build dynamic SET clause based on provided fields
	if req.Active != nil {
		setParts = append(setParts, fmt.Sprintf("active = $%d", argIndex))
		args = append(args, *req.Active)
		argIndex++
	}

	if req.Metadata != nil {
		setParts = append(setParts, fmt.Sprintf("metadata = $%d", argIndex))
		args = append(args, req.Metadata)
		argIndex++
	}

	// Always update updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	if len(setParts) == 1 { // Only updated_at was added
		return nil // Nothing to update
	}

	// Add promotion code ID as last parameter
	args = append(args, promotionCodeID)

	query := fmt.Sprintf(`
		UPDATE %s.promotion_codes 
		SET %s 
		WHERE id = $%d
	`, r.schema, strings.Join(setParts, ", "), argIndex)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errs.ClassifyPgError("update promotion code", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "promotion code")
	}

	return nil
}