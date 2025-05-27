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

func (e *repository) InsertEvent(
	ctx context.Context,
	tx *sql.Tx,
	event *models.Event,
) error {
	query := fmt.Sprintf(`INSERT INTO %s (
                id,
                encrypted_title,
                encrypted_description,
				encrypted_postal_code,
				encrypted_city, 
				encrypted_address1, 
				encrypted_address2,
                placecount,
                freeplace,
                encrypted_begin_at,
                encrypted_end_at,
                encrypted_price_id,
                day,
                month,
                year
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`, pg.QualifiedTable(e.schema, "events"))

	result, err := tx.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.PostalCode,
		event.City,
		event.Address1,
		event.Address2,
		event.PlaceCount,
		event.FreePlace,
		event.BeginAt,
		event.EndAt,
		event.PriceID,
		event.Day,
		event.Month,
		event.Year,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), "event")
	}
	return nil
}
