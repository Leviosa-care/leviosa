package postgres

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewRoomScheduleRepository(pool *pgxpool.Pool) *RoomScheduleRepository {
	return &RoomScheduleRepository{pool: pool}
}

// GetRoomHoursForDate returns room hours for a specific date
// Prioritizes specific dates over recurring patterns based on priority
func (r *RoomScheduleRepository) GetRoomHoursForDate(
	ctx context.Context,
	roomID uuid.UUID,
	date time.Time,
) (*domain.RoomAvailabilitySchedule, error) {
	dayOfWeek := int(date.Weekday())

	query := `
		SELECT id, room_id, day_of_week, specific_date, open_time, close_time,
		       priority, is_active, created_at, updated_at
		FROM booking.room_availability_schedules
		WHERE room_id = $1
		  AND is_active = TRUE
		  AND (
		      specific_date = $2::DATE OR
		      day_of_week = $3
		  )
		ORDER BY priority DESC, specific_date DESC NULLS LAST
		LIMIT 1
	`

	var schedule domain.RoomAvailabilitySchedule
	var dayOfWeekPtr *int
	var specificDatePtr *time.Time

	err := r.pool.QueryRow(ctx, query, roomID, date, dayOfWeek).Scan(
		&schedule.ID,
		&schedule.RoomID,
		&dayOfWeekPtr,
		&specificDatePtr,
		&schedule.OpenTime,
		&schedule.CloseTime,
		&schedule.Priority,
		&schedule.IsActive,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		return nil, errs.ClassifyPgError("get room hours for date", err)
	}

	schedule.DayOfWeek = dayOfWeekPtr
	schedule.SpecificDate = specificDatePtr

	return &schedule, nil
}

// Create adds a new room availability schedule
func (r *RoomScheduleRepository) Create(ctx context.Context, schedule *domain.RoomAvailabilitySchedule) error {
	query := `
		INSERT INTO booking.room_availability_schedules (
			id, room_id, day_of_week, specific_date, open_time, close_time,
			priority, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()
	schedule.ID = uuid.New()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now

	_, err := r.pool.Exec(ctx, query,
		schedule.ID,
		schedule.RoomID,
		schedule.DayOfWeek,
		schedule.SpecificDate,
		schedule.OpenTime,
		schedule.CloseTime,
		schedule.Priority,
		schedule.IsActive,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)

	if err != nil {
		return errs.ClassifyPgError("create room schedule", err)
	}

	return nil
}

// GetByID retrieves a schedule by its ID
func (r *RoomScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomAvailabilitySchedule, error) {
	query := `
		SELECT id, room_id, day_of_week, specific_date, open_time, close_time,
		       priority, is_active, created_at, updated_at
		FROM booking.room_availability_schedules
		WHERE id = $1 AND is_active = TRUE
	`

	var schedule domain.RoomAvailabilitySchedule
	var dayOfWeekPtr *int
	var specificDatePtr *time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&schedule.ID,
		&schedule.RoomID,
		&dayOfWeekPtr,
		&specificDatePtr,
		&schedule.OpenTime,
		&schedule.CloseTime,
		&schedule.Priority,
		&schedule.IsActive,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		return nil, errs.ClassifyPgError("get room schedule by ID", err)
	}

	schedule.DayOfWeek = dayOfWeekPtr
	schedule.SpecificDate = specificDatePtr

	return &schedule, nil
}

// Update modifies an existing schedule
func (r *RoomScheduleRepository) Update(ctx context.Context, schedule *domain.RoomAvailabilitySchedule) error {
	query := `
		UPDATE booking.room_availability_schedules
		SET day_of_week = $1,
		    specific_date = $2,
		    open_time = $3,
		    close_time = $4,
		    priority = $5,
		    is_active = $6,
		    updated_at = $7
		WHERE id = $8
	`

	schedule.UpdatedAt = time.Now()

	result, err := r.pool.Exec(ctx, query,
		schedule.DayOfWeek,
		schedule.SpecificDate,
		schedule.OpenTime,
		schedule.CloseTime,
		schedule.Priority,
		schedule.IsActive,
		schedule.UpdatedAt,
		schedule.ID,
	)

	if err != nil {
		return errs.ClassifyPgError("update room schedule", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrRepositoryNotUpdated
	}

	return nil
}

// Delete soft-deletes a schedule by setting is_active to false
func (r *RoomScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE booking.room_availability_schedules
		SET is_active = FALSE,
		    updated_at = $1
		WHERE id = $2
	`

	result, err := r.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return errs.ClassifyPgError("delete room schedule", err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrRepositoryNotDeleted
	}

	return nil
}

// ListByRoom returns all active schedules for a room
func (r *RoomScheduleRepository) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]*domain.RoomAvailabilitySchedule, error) {
	query := `
		SELECT id, room_id, day_of_week, specific_date, open_time, close_time,
		       priority, is_active, created_at, updated_at
		FROM booking.room_availability_schedules
		WHERE room_id = $1 AND is_active = TRUE
		ORDER BY priority DESC, specific_date DESC NULLS LAST
	`

	rows, err := r.pool.Query(ctx, query, roomID)
	if err != nil {
		return nil, errs.ClassifyPgError("list room schedules", err)
	}
	defer rows.Close()

	var schedules []*domain.RoomAvailabilitySchedule

	for rows.Next() {
		var schedule domain.RoomAvailabilitySchedule
		var dayOfWeekPtr *int
		var specificDatePtr *time.Time

		err := rows.Scan(
			&schedule.ID,
			&schedule.RoomID,
			&dayOfWeekPtr,
			&specificDatePtr,
			&schedule.OpenTime,
			&schedule.CloseTime,
			&schedule.Priority,
			&schedule.IsActive,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)

		if err != nil {
			return nil, errs.ClassifyPgError("scan room schedule row", err)
		}

		schedule.DayOfWeek = dayOfWeekPtr
		schedule.SpecificDate = specificDatePtr

		schedules = append(schedules, &schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room schedule rows", err)
	}

	return schedules, nil
}
