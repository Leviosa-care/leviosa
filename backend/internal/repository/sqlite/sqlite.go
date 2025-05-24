package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	testdb "github.com/hengadev/leviosa/pkg/testdatabase"
)

type sqliteRepository interface {
	GetDB() *sql.DB
}

type repoConstructor[T sqliteRepository] func(context.Context, *sql.DB) T

func recoverDB() {
	if r := recover(); r != nil {
		fmt.Println("In memory SQLite database setup failed.")
	}
}

type teardownFunc func()

func SetupRepository[T sqliteRepository](t testing.TB, ctx context.Context, version int64, constructor repoConstructor[T]) (T, teardownFunc) {
	t.Helper()
	var repo T
	db, err := testdb.NewSQLite(ctx)
	if err != nil {
		t.Errorf("database connection: %s", err)
	}
	// NOTE: I think it is better to do it that way instead
	// if err := testdb.Setup(ctx, repo.GetDB(), version); err != nil {
	if err := testdb.MigrateSQLite(ctx, db, version); err != nil {
		t.Errorf("database setup: %s", err)
	}
	repo = constructor(ctx, db)
	teardown := func() {
		defer recoverDB()
		err := db.Close()
		if err != nil {
			t.Errorf("database teardown: %s", err)
		}
	}
	return repo, teardown
}
