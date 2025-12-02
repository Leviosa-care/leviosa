package bookingRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) List(ctx context.Context, filter ports.BookingFilter) ([]*domain.BookingEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			b.id, b.availability_id, b.client_id, b.user_id, b.room_id,
			b.product_id_encrypted, b.slot_start_time_encrypted, b.slot_end_time_encrypted,
			b.client_notes_encrypted, b.partner_notes_encrypted,
			b.total_price_cents, b.currency, b.payment_status, b.payment_intent_id,
			b.status, b.cancelled_at, b.cancellation_reason_encrypted, b.completed_at,
			b.created_at, b.updated_at,
			b.dek_encrypted, b.key_version, b.metadata
		FROM %s.bookings b
	`, r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.ClientID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.client_id = $%d", argIndex))
		args = append(args, *filter.ClientID)
		argIndex++
	}

	if filter.PartnerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.user_id = $%d", argIndex))
		args = append(args, *filter.PartnerID)
		argIndex++
	}

	if filter.RoomID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.room_id = $%d", argIndex))
		args = append(args, *filter.RoomID)
		argIndex++
	}

	if filter.AvailabilityID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.availability_id = $%d", argIndex))
		args = append(args, *filter.AvailabilityID)
		argIndex++
	}

	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, status := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, status)
			argIndex++
		}
		whereConditions = append(whereConditions, fmt.Sprintf("b.status IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filter.PaymentStatus) > 0 {
		placeholders := make([]string, len(filter.PaymentStatus))
		for i, status := range filter.PaymentStatus {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, status)
			argIndex++
		}
		whereConditions = append(whereConditions, fmt.Sprintf("b.payment_status IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.CreatedAfter != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.created_at >= $%d", argIndex))
		args = append(args, *filter.CreatedAfter)
		argIndex++
	}

	if filter.CreatedBefore != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.created_at <= $%d", argIndex))
		args = append(args, *filter.CreatedBefore)
		argIndex++
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add ordering
	orderBy := "b.created_at"
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "created_at", "total_price_cents":
			orderBy = "b." + filter.OrderBy
		}
	}

	orderDirection := "DESC"
	if filter.OrderDirection == "asc" {
		orderDirection = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDirection)

	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("list bookings", err)
	}
	defer rows.Close()

	var bookingsEncx []*domain.BookingEncx
	for rows.Next() {
		bookingEncx := &domain.BookingEncx{}
		err := rows.Scan(
			&bookingEncx.ID,
			&bookingEncx.AvailabilityID,
			&bookingEncx.ClientID,
			&bookingEncx.PartnerID,
			&bookingEncx.RoomID,
			&bookingEncx.ProductIDEncrypted,
			&bookingEncx.SlotStartTimeEncrypted,
			&bookingEncx.SlotEndTimeEncrypted,
			&bookingEncx.ClientNotesEncrypted,
			&bookingEncx.PartnerNotesEncrypted,
			&bookingEncx.TotalPriceCents,
			&bookingEncx.Currency,
			&bookingEncx.PaymentStatus,
			&bookingEncx.PaymentIntentID,
			&bookingEncx.Status,
			&bookingEncx.CancelledAt,
			&bookingEncx.CancellationReasonEncrypted,
			&bookingEncx.CompletedAt,
			&bookingEncx.CreatedAt,
			&bookingEncx.UpdatedAt,
			&bookingEncx.DEKEncrypted,
			&bookingEncx.KeyVersion,
			&bookingEncx.Metadata,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan booking row", err)
		}

		bookingsEncx = append(bookingsEncx, bookingEncx)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate booking rows", err)
	}

	return bookingsEncx, nil
}

