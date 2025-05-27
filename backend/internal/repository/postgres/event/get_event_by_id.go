package eventRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (e *repository) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	event := &models.Event{}
	query := fmt.Sprintf(`
        SELECT
            title_encrypted,
            description_encrypted,
            city_encrypted,
            postal_code_encrypted,
            address1_encrypted,
            address2_encrypted,
            placecount,
            freeplace,
            begin_at_encrypted,
            end_at_encrypted
        FROM %s
        WHERE id = $1;`, pg.QualifiedTable(e.schema, "events"))

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
