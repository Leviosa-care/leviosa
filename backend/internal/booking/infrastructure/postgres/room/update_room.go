package roomRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Update(ctx context.Context, room *domain.RoomEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.rooms SET
			building_id = $2,
			name_encrypted = $3,
			name_hash = $4,
			description_encrypted = $5,
			room_number_encrypted = $6,
			room_number_hash = $7,
			capacity = $8,
			equipment_encrypted = $9,
			operating_start_time = $10,
			operating_end_time = $11,
			is_active = $12,
			updated_at = $13,
			dek_encrypted = $14,
			key_version = $15,
			metadata = $16
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
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
		room.UpdatedAt,
		room.DEKEncrypted,
		room.KeyVersion,
		room.Metadata,
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
