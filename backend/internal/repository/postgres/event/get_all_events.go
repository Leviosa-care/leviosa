package eventRepository

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (e *repository) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	query := `
        SELECT 
            id,
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
        FROM events;`
	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(
			&event.ID,
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
			return nil, rp.NewDatabaseErr(err)
		}
		events = append(events, &event)
	}
	if err = rows.Err(); err != nil {
		return nil, rp.NewDatabaseErr(err)
	}
	return events, nil
}
