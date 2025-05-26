package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/hengadev/leviosa/pkg/config"
	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/hengadev/errsx"
)

func Postgres(ctx context.Context, env envmode.Mode, conf *config.PostgresSecrets) (*sql.DB, error) {
	var dsn string
	dsn = fmt.Sprintf(
		"postgres://%s:%s@localhost:%d/%s_%s?sslmode=disable",
		conf.User,
		conf.Password,
		conf.Port,
		env.String(),
		DefaultDB,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping %q: %w", dsn, err)
	}
	return db, nil
}

func InitPostgres(db *sql.DB, queries ...string) error {
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: update the next functions since the placeholders are not the same

// Return the necessary elements to write a query update to target the non zero value.
func WritePostgresInsertQuery[S any](s S) (string, []any) {
	var tables []string
	var values []any
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	vf := reflect.VisibleFields(t)
	tableName := fmt.Sprintf("%ss", strings.ToLower(t.Name()))
	query := fmt.Sprintf("UPDATE %s set ", tableName)
	for _, f := range vf {
		value := v.FieldByName(f.Name)
		if !value.IsZero() && value.CanInterface() {
			tables = append(tables, fmt.Sprintf("%s=?", strings.ToLower(f.Name)))
			values = append(values, value.Interface())
		}
	}
	query += strings.Join(tables, ", ")
	query += " WHERE id=?;"
	return query, values
}

// Return the necessary elements to write a query update to target the non zero value.
func WritePostgresUpdateQuery[T SQLMappable](
	object T,
	whereMap map[string]any,
) (string, []any, error) {
	var errs errsx.Map
	var tables []string
	var values []any
	var notUpdatedFields []string
	v := reflect.ValueOf(object)
	t := reflect.TypeOf(object)
	vf := reflect.VisibleFields(t)
	tableName := fmt.Sprintf("%ss", strings.ToLower(t.Name()))
	query := fmt.Sprintf("UPDATE %s set ", tableName)
	for _, f := range vf {
		value := v.FieldByName(f.Name)
		if !value.IsZero() && value.CanInterface() {
			if err := isProhibitedField(f.Name, object.GetProhibitedFields()); err != nil {
				notUpdatedFields = append(notUpdatedFields, f.Name)
				continue
			}
			column := object.GetSQLColumnMapping()[f.Name]
			tables = append(tables, PostgresPlaceholder(column))
			values = append(values, value.Interface())
		}
	}
	query += strings.Join(tables, ", ")
	query += " WHERE "
	var wherePlaceholder []string
	for key, value := range whereMap {
		minKey := strings.ToLower(key)
		wherePlaceholder = append(wherePlaceholder, PostgresPlaceholder(minKey))
		values = append(values, value)
	}
	query += strings.Join(wherePlaceholder, " AND ") + ";"
	if len(notUpdatedFields) > 0 {
		errs.Set("prohibited fields", strings.Join(notUpdatedFields, ", "))
	}
	return query, values, errs.AsError()
}

func PostgresPlaceholder(name string) string {
	return fmt.Sprintf("%s=?", name)
}
