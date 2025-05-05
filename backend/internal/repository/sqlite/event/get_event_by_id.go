package eventRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (e *EventRepository) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	event := &models.Event{}
	query := `
        SELECT 
            encrypted_title,
            encrypted_description,
            encrypted_city,
            encrypted_postal_code,
            encrypted_address1,
            encrypted_address2,
            placecount,
            freeplace,
            encrypted_begin_at,
            encrypted_end_at
        FROM events 
        WHERE id = ?;`

	if err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&event.Title,
		&event.Description,
		&event.City,
		&event.PostalCode,
		&event.Address1,
		&event.Address2,
		&event.PlaceCount,
		&event.FreePlace,
		&event.BeginAt,
		&event.EndAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, rp.NewNotFoundErr(err, "user")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	event.ID = id
	return event, nil
}
