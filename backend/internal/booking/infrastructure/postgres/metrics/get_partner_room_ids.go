package metricsRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

// GetPartnerRoomIDs retrieves all room IDs a partner has access to
func (r *Repository) GetPartnerRoomIDs(ctx context.Context, partnerID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT DISTINCT room_id
		FROM booking.room_allocations
		WHERE user_id_hash = $1
			AND is_active = true
		ORDER BY room_id
	`

	// Hash the partner ID for querying (GDPR-compliant hashing)
	partnerIDBytes, err := encx.SerializeValue(partnerID)
	if err != nil {
		return nil, fmt.Errorf("serialize partner ID for hashing: %w", err)
	}
	partnerIDHash := r.crypto.HashBasic(ctx, partnerIDBytes)

	rows, err := r.pool.Query(ctx, query, partnerIDHash)
	if err != nil {
		return nil, errs.ClassifyPgError("query partner room IDs", err)
	}
	defer rows.Close()

	roomIDs := []uuid.UUID{}

	for rows.Next() {
		var roomID uuid.UUID
		if err := rows.Scan(&roomID); err != nil {
			return nil, errs.ClassifyPgError("scan room ID", err)
		}
		roomIDs = append(roomIDs, roomID)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room IDs", err)
	}

	return roomIDs, nil
}
