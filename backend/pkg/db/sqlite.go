package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hengadev/leviosa/pkg/envmode"

	"github.com/hengadev/errsx"
)

// constants
const JournalMode = "WAL"
const BusyTimeOut = 5000
const ForeignKeys = "ON"
const MaxOpenConns = 10
const MaxIdleConns = 5
const ConnMaxLifetime = time.Hour

// func SQLite(ctx context.Context, env envmode.Mode, password string) (*sql.DB, error) {
func SQLite(ctx context.Context, env envmode.Mode) (*sql.DB, error) {
	var dsn string
	dsn = fmt.Sprintf(
		"file:%s_%s?_journal_mode=%s&_busy_timeout=%d&_foreign_keys=%s",
		env.String(),
		DefaultDB,
		JournalMode,
		BusyTimeOut,
		ForeignKeys,
	)
	// if env == envmode.Staging || env == envmode.Prod {
	// 	dsn = fmt.Sprintf("%s&_pragma_key=%s&_pragma_cipher_page_size=4096", dsn, password)
	// }
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping %q: %w", dsn, err)
	}
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)
	return db, nil
}

func InitSQLite(db *sql.DB, queries ...string) error {
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// Return the necessary elements to write a query update to target the non zero value.
func WriteSQLiteInsertQuery[S any](s S) (string, []any) {
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
func WriteSQLiteUpdateQuery[T SQLMappable](
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
			tables = append(tables, SQLitePlaceholder(column))
			values = append(values, value.Interface())
		}
	}
	query += strings.Join(tables, ", ")
	query += " WHERE "
	var wherePlaceholder []string
	for key, value := range whereMap {
		minKey := strings.ToLower(key)
		wherePlaceholder = append(wherePlaceholder, SQLitePlaceholder(minKey))
		values = append(values, value)
	}
	query += strings.Join(wherePlaceholder, " AND ") + ";"
	if len(notUpdatedFields) > 0 {
		errs.Set("prohibited fields", strings.Join(notUpdatedFields, ", "))
	}
	return query, values, errs.AsError()
}

func SQLitePlaceholder(name string) string {
	return fmt.Sprintf("%s=?", name)
}

// Helper function to check if a struct field belongs to a list of strings provided.
func isProhibitedField(name string, prohibitedFields []string) error {
	for _, prohibitedField := range prohibitedFields {
		if name == prohibitedField {
			return fmt.Errorf("field %q is prohibited", prohibitedField)
		}
	}
	return nil
}
