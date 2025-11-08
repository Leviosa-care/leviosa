package roomRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Update(ctx context.Context, room *domain.Room) error {
	// Encrypt sensitive fields
	if err := r.crypto.EncryptStruct(ctx, room); err != nil {
		return fmt.Errorf("encrypt room data: %w", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s.rooms SET
			building_id = $2,
			name_encrypted = $3,
			description_encrypted = $4,
			room_number_encrypted = $5,
			capacity = $6,
			equipment_encrypted = $7,
			hourly_rate_cents = $8,
			is_active = $9,
			updated_at = $10
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		room.ID,
		room.BuildingID,
		room.NameEncrypted,
		room.DescriptionEncrypted,
		room.RoomNumberEncrypted,
		room.Capacity,
		room.EquipmentEncrypted,
		room.HourlyRateCents,
		room.IsActive,
		room.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("update room", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
