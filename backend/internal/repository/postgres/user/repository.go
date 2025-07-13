package userRepository

import (
	"context"
	"database/sql"

	"github.com/hengadev/leviosa/internal/domain/user"

	_ "github.com/jackc/pgx"
)

type repository struct {
	DB     *sql.DB
	schema string
}

func (u *repository) GetDB() *sql.DB {
	return u.DB
}

func New(ctx context.Context, db *sql.DB) (userService.ReadWriter, error) {
	return &repository{DB: db, schema: "users"}, nil
}
