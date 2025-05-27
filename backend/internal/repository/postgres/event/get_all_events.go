package eventRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (e *repository) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	query := fmt.Sprintf(`
        SELECT 
            id,
            title,_encrypted
            description,_encrypted
			city,_encrypted
			postal_code,_encrypted
			address1,_encrypted
			address2,_encrypted
            placecount,
            freeplace,
            begin_at,_encrypted
            end_at_encrypted
        FROM %s;`, pg.QualifiedTable(e.schema, "events"))
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
