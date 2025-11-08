package roomRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, room *domain.Room) error {
	// Encrypt sensitive fields
	if err := r.crypto.EncryptStruct(ctx, room); err != nil {
		return fmt.Errorf("encrypt room data: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.rooms (
			id, building_id, name_encrypted, description_encrypted,
			room_number_encrypted, capacity, equipment_encrypted,
			hourly_rate_cents, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		room.ID,
		room.BuildingID,
		room.NameEncrypted,
		room.DescriptionEncrypted,
		room.RoomNumberEncrypted,
		room.Capacity,
		room.EquipmentEncrypted,
		room.HourlyRateCents,
		room.IsActive,
		room.CreatedAt,
		room.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create room", err)
	}

	return nil
}
