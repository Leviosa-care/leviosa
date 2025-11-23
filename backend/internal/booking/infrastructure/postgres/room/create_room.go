package roomRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, room *domain.RoomEncx) error {
	// Encrypt sensitive fields

	query := fmt.Sprintf(`
		INSERT INTO %s.rooms (
			id, building_id, name_encrypted, name_hash, description_encrypted,
			room_number_encrypted, room_number_hash, capacity, equipment_encrypted,
			operating_start_time, operating_end_time,
			is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		room.ID,
		room.BuildingID,
		room.NameEncrypted,
		room.NameHash,
		room.DescriptionEncrypted,
		room.RoomNumberEncrypted,
		room.RoomNumberHash,
		room.Capacity,
		room.EquipmentEncrypted,
		room.OperatingStartTime,
		room.OperatingEndTime,
		room.IsActive,
		room.CreatedAt,
		room.UpdatedAt,
		room.DEKEncrypted,
		room.KeyVersion,
		room.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create room", err)
	}

	return nil
}
