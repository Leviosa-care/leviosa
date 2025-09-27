package bookingRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	query := fmt.Sprintf(`
		SELECT
			id, availability_id, client_id, partner_id, room_id,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted,
			created_at, updated_at
		FROM %s.bookings
		WHERE id = $1
	`, r.schema)

	booking := &domain.Booking{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&booking.ID,
		&booking.AvailabilityID,
		&booking.ClientID,
		&booking.PartnerID,
		&booking.RoomID,
		&booking.ClientNotesEncrypted,
		&booking.PartnerNotesEncrypted,
		&booking.TotalPriceCents,
		&booking.Currency,
		&booking.PaymentStatus,
		&booking.PaymentIntentID,
		&booking.Status,
		&booking.CancelledAt,
		&booking.CancellationReasonEncrypted,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get booking by id", err)
	}

	// Decrypt sensitive fields
	if err := r.crypto.DecryptStruct(ctx, booking); err != nil {
		return nil, fmt.Errorf("decrypt booking data: %w", err)
	}

	return booking, nil
}

func (r *Repository) Update(ctx context.Context, booking *domain.Booking) error {
	// Encrypt sensitive fields
	if err := r.crypto.EncryptStruct(ctx, booking); err != nil {
		return fmt.Errorf("encrypt booking data: %w", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s.bookings SET
			availability_id = $2,
			client_id = $3,
			partner_id = $4,
			room_id = $5,
			client_notes_encrypted = $6,
			partner_notes_encrypted = $7,
			total_price_cents = $8,
			currency = $9,
			payment_status = $10,
			payment_intent_id = $11,
			status = $12,
			cancelled_at = $13,
			cancellation_reason_encrypted = $14,
			updated_at = $15
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		booking.ID,
		booking.AvailabilityID,
		booking.ClientID,
		booking.PartnerID,
		booking.RoomID,
		booking.ClientNotesEncrypted,
		booking.PartnerNotesEncrypted,
		booking.TotalPriceCents,
		booking.Currency,
		booking.PaymentStatus,
		booking.PaymentIntentID,
		booking.Status,
		booking.CancelledAt,
		booking.CancellationReasonEncrypted,
		booking.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("update booking", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	// Hard delete for GDPR compliance
	query := fmt.Sprintf(`
		DELETE FROM %s.bookings
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete booking", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}

func (r *Repository) GetByAvailabilityID(ctx context.Context, availabilityID uuid.UUID) (*domain.Booking, error) {
	query := fmt.Sprintf(`
		SELECT
			id, availability_id, client_id, partner_id, room_id,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted,
			created_at, updated_at
		FROM %s.bookings
		WHERE availability_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, r.schema)

	booking := &domain.Booking{}
	err := r.pool.QueryRow(ctx, query, availabilityID).Scan(
		&booking.ID,
		&booking.AvailabilityID,
		&booking.ClientID,
		&booking.PartnerID,
		&booking.RoomID,
		&booking.ClientNotesEncrypted,
		&booking.PartnerNotesEncrypted,
		&booking.TotalPriceCents,
		&booking.Currency,
		&booking.PaymentStatus,
		&booking.PaymentIntentID,
		&booking.Status,
		&booking.CancelledAt,
		&booking.CancellationReasonEncrypted,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get booking by availability id", err)
	}

	// Decrypt sensitive fields
	if err := r.crypto.DecryptStruct(ctx, booking); err != nil {
		return nil, fmt.Errorf("decrypt booking data: %w", err)
	}

	return booking, nil
}

func (r *Repository) GetByPaymentIntentID(ctx context.Context, paymentIntentID string) (*domain.Booking, error) {
	query := fmt.Sprintf(`
		SELECT
			id, availability_id, client_id, partner_id, room_id,
			client_notes_encrypted, partner_notes_encrypted,
			total_price_cents, currency, payment_status, payment_intent_id,
			status, cancelled_at, cancellation_reason_encrypted,
			created_at, updated_at
		FROM %s.bookings
		WHERE payment_intent_id = $1
	`, r.schema)

	booking := &domain.Booking{}
	err := r.pool.QueryRow(ctx, query, paymentIntentID).Scan(
		&booking.ID,
		&booking.AvailabilityID,
		&booking.ClientID,
		&booking.PartnerID,
		&booking.RoomID,
		&booking.ClientNotesEncrypted,
		&booking.PartnerNotesEncrypted,
		&booking.TotalPriceCents,
		&booking.Currency,
		&booking.PaymentStatus,
		&booking.PaymentIntentID,
		&booking.Status,
		&booking.CancelledAt,
		&booking.CancellationReasonEncrypted,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get booking by payment intent id", err)
	}

	// Decrypt sensitive fields
	if err := r.crypto.DecryptStruct(ctx, booking); err != nil {
		return nil, fmt.Errorf("decrypt booking data: %w", err)
	}

	return booking, nil
}

func (r *Repository) List(ctx context.Context, filter ports.BookingFilter) ([]*domain.Booking, error) {
	query := fmt.Sprintf(`
		SELECT
			b.id, b.availability_id, b.client_id, b.partner_id, b.room_id,
			b.client_notes_encrypted, b.partner_notes_encrypted,
			b.total_price_cents, b.currency, b.payment_status, b.payment_intent_id,
			b.status, b.cancelled_at, b.cancellation_reason_encrypted,
			b.created_at, b.updated_at
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
		whereConditions = append(whereConditions, fmt.Sprintf("b.partner_id = $%d", argIndex))
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

	var bookings []*domain.Booking
	for rows.Next() {
		booking := &domain.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.AvailabilityID,
			&booking.ClientID,
			&booking.PartnerID,
			&booking.RoomID,
			&booking.ClientNotesEncrypted,
			&booking.PartnerNotesEncrypted,
			&booking.TotalPriceCents,
			&booking.Currency,
			&booking.PaymentStatus,
			&booking.PaymentIntentID,
			&booking.Status,
			&booking.CancelledAt,
			&booking.CancellationReasonEncrypted,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan booking row", err)
		}

		// Decrypt sensitive fields
		if err := r.crypto.DecryptStruct(ctx, booking); err != nil {
			return nil, fmt.Errorf("decrypt booking data: %w", err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate booking rows", err)
	}

	return bookings, nil
}

func (r *Repository) GetByClientID(ctx context.Context, clientID uuid.UUID, filter ports.BookingFilter) ([]*domain.Booking, error) {
	filter.ClientID = &clientID
	return r.List(ctx, filter)
}

func (r *Repository) GetByPartnerID(ctx context.Context, partnerID uuid.UUID, filter ports.BookingFilter) ([]*domain.Booking, error) {
	filter.PartnerID = &partnerID
	return r.List(ctx, filter)
}

func (r *Repository) GetUpcoming(ctx context.Context, filter ports.BookingFilter) ([]*domain.Booking, error) {
	// Force confirmed status and join with availabilities to filter by future start times
	filter.Status = []domain.BookingStatus{domain.BookingStatusConfirmed}

	// This would require a join with availabilities table for more complex filtering
	// For now, return all confirmed bookings and let the service layer filter by time
	return r.List(ctx, filter)
}