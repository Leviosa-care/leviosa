package eventRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (e *repository) GetPriceID(ctx context.Context, eventID string) (string, error) {
	var priceID string
	query := fmt.Sprintf(`
        SELECT 
            price_id_encrypted
        FROM %s 
        WHERE id = $1;`, pg.QualifiedTable(e.schema, "events"))
	err := e.DB.QueryRowContext(ctx, query, eventID).Scan(&priceID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", rp.NewNotFoundErr(err, "price ID")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return "", rp.NewContextErr(err)
		default:
			return "", rp.NewDatabaseErr(err)
		}
	}
	return priceID, nil
}
