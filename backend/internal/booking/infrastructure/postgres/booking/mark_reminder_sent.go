package bookingRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// MarkReminderSent sets reminded_at to now() for a given booking.
func (r *Repository) MarkReminderSent(ctx context.Context, bookingID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.bookings SET reminded_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, bookingID)
	if err != nil {
		return errs.ClassifyPgError("mark reminder sent", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
