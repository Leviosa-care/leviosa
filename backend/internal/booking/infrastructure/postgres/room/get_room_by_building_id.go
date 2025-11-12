package roomRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (r *Repository) GetByBuildingID(ctx context.Context, buildingID uuid.UUID, activeOnly bool) ([]*domain.RoomEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, building_id, name_encrypted, name_hash, description_encrypted,
			room_number_encrypted, room_number_hash, capacity, equipment_encrypted,
			hourly_rate_cents, is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.rooms
		WHERE building_id = $1
	`, r.schema)

	args := []interface{}{buildingID}
	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY name_encrypted ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("get rooms by building id", err)
	}
	defer rows.Close()

	var roomsEncx []*domain.RoomEncx
	for rows.Next() {
		roomEncx := &domain.RoomEncx{}
		err := rows.Scan(
			&roomEncx.ID,
			&roomEncx.BuildingID,
			&roomEncx.NameEncrypted,
			&roomEncx.NameHash,
			&roomEncx.DescriptionEncrypted,
			&roomEncx.RoomNumberEncrypted,
			&roomEncx.RoomNumberHash,
			&roomEncx.Capacity,
			&roomEncx.EquipmentEncrypted,
			&roomEncx.HourlyRateCents,
			&roomEncx.IsActive,
			&roomEncx.CreatedAt,
			&roomEncx.UpdatedAt,
			&roomEncx.DEKEncrypted,
			&roomEncx.KeyVersion,
			&roomEncx.Metadata,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan room row", err)
		}

		roomsEncx = append(roomsEncx, roomEncx)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room rows", err)
	}

	return roomsEncx, nil
}
