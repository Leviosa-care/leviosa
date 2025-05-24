package registerRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/register"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (r *repository) AddRegistration(ctx context.Context, reg *registerService.Registration, day, year int, month string) error {
	// TODO: change that function because it no longer works
	tablename := getTablename(day, year, month)
	query := fmt.Sprintf("INSERT INTO %s (userid, eventid, beginat) values ($1,$2,$3);", tablename)
	result, err := r.DB.ExecContext(ctx, query, reg.UserID, reg.ProductID, reg.StartTime)
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
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), "registration")
	}
	return nil
}
