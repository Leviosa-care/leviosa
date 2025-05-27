package eventRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

// TODO:
// - move that function to registrationRepository
// - move part of this function to the service that uses that function

func (e *repository) GetEventForUser(ctx context.Context, userID string) (*models.EventUser, error) {
	tx, err := e.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, rp.NewDatabaseErr(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback()

	now := time.Now()
	day, month, year := now.Day(), int(now.Month()), now.Year()
	statements := []struct {
		condition string
		field     string
	}{
		{condition: fmt.Sprintf("(month < %d AND year = %d) OR (year < %d) OR (day < %d AND month = %d AND year = %d) LIMIT 3", month, year, year, day, month, year), field: "past"},
		{condition: fmt.Sprintf("(year > %d) OR (month > %d AND year = %d) OR (day > %d AND month = %d AND year = %d) LIMIT 3", year, month, year, day, month, year), field: "next"},
		{condition: fmt.Sprintf("(day > %d AND month = %d AND year = %d) OR (month = %d + 1 AND year = %d) OR (month = 1 AND year = %d + 1) LIMIT 1", day, month, year, year, month, year), field: "incoming"},
	}

	var res models.EventUser
	for _, statement := range statements {
		query := fmt.Sprintf("SELECT * FROM %s WHERE %s;", pg.QualifiedTable(e.schema, "events"), statement.condition)
		rows, err := tx.QueryContext(ctx, query, userID)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
				return nil, rp.NewContextErr(err)
			default:
				return nil, rp.NewDatabaseErr(err)
			}
		}
		defer rows.Close()
		for rows.Next() {
			var priceID string
			var beginAt string
			event := &models.Event{}
			if err := rows.Scan(
				&event.ID,
				&event.PlaceCount,
				&beginAt,
				&priceID,
				&event.Day,
				&event.Month,
				&event.Year,
			); err != nil {
				return &res, rp.NewDatabaseErr(err)
			}
			event.BeginAt, err = parseBeginAt(beginAt, event.Day, event.Month, event.Year)
			if err != nil {
				return &res, fmt.Errorf("%s: %w", "error parsing time", err)
			}
			var usedCount int
			query := fmt.Sprintf("SELECT COUNT(userid) FROM %s_%s;", pg.QualifiedTable(e.schema, "event"), event.ID)
			if err := tx.QueryRowContext(ctx, query).Scan(&usedCount); err != nil {
				switch {
				case errors.Is(err, sql.ErrNoRows):
					return &res, rp.NewNotFoundErr(err, "user")
				case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
					return &res, rp.NewContextErr(err)
				default:
					return &res, rp.NewDatabaseErr(err)
				}
			}
			event.FreePlace = event.PlaceCount - usedCount
			switch statement.field {
			case "past":
				res.PastEvents = append(res.PastEvents, event)
			case "next":
				res.NextEvents = append(res.NextEvents, event)
			case "incoming":
				res.IncomingEvents = append(res.IncomingEvents, event)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return &res, rp.NewDatabaseErr(fmt.Errorf("commit transaction: %w", err))
	}

	return &res, nil
}
