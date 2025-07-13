package eventRepository

import (
	"context"
	"database/sql"

	"github.com/hengadev/leviosa/internal/domain/event"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	DB     *sql.DB
	schema string
}

func (e *repository) GetDB() *sql.DB {
	return e.DB
}

func New(ctx context.Context, db *sql.DB) (eventService.ReadWriter, error) {

	return &repository{DB: db, schema: "events"}, nil
}
