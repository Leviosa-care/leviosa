package bookingRepository

import (
	"context"

	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
	"github.com/hengadev/encx"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	pool   *pgxpool.Pool
	schema string
	crypto *encx.Encx
}

func New(ctx context.Context, pool *pgxpool.Pool, crypto *encx.Encx) (ports.BookingRepository, error) {
	if pool == nil {
		return nil, errs.ErrDatabase
	}
	if crypto == nil {
		return nil, errs.ErrInvalidValue
	}

	return &Repository{
		pool:   pool,
		schema: "booking",
		crypto: crypto,
	}, nil
}