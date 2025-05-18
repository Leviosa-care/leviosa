package eventRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (e *repository) AddEvent(ctx context.Context, event *models.Event) (string, error) {
	tx, err := e.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", rp.NewDatabaseErr(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback()
	if err := e.InsertEvent(ctx, tx, event); err != nil {
		return "", err
	}
	if err := e.LoopQuery(ctx, tx, "product_id", "event_products", event.Products, event.ID); err != nil {
		return "", err
	}
	if err := e.LoopQuery(ctx, tx, "offer_id", "event_offers", event.Offers, event.ID); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return event.ID, rp.NewDatabaseErr(fmt.Errorf("commit transaction: %w", err))
	}
	return event.ID, nil
}
