package promotionCodeRepository

import (
	"fmt"
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *PromotionCodeRepository) DeactivatePromotionCode(ctx context.Context, promotionCodeID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.promotion_codes 
		SET active = false,
		    updated_at = $2
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, promotionCodeID, time.Now())
	if err != nil {
		return errs.ClassifyPgError("deactivate promotion code", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "promotion code")
	}

	return nil
}