package allocationRepository

import (
	"context"

	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) (ports.RoomAllocationRepository, error) {
	if pool == nil {
		return nil, errs.ErrDatabase
	}

	return &Repository{
		pool:   pool,
		schema: "booking",
	}, nil
}