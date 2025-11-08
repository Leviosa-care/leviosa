package roomRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/hengadev/encx"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	pool   *pgxpool.Pool
	schema string
	crypto *encx.CryptoService
}

func New(ctx context.Context, pool *pgxpool.Pool, crypto *encx.CryptoService) (ports.RoomRepository, error) {
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
