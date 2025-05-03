package sqliteutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hengadev/errsx"
)

type SQLMappable interface {
	GetSQLColumnMapping() map[string]string
	GetProhibitedFields() []string
}

// Return the necessary elements to write a query update to target the non zero value.
func WriteUpdateQuery[T SQLMappable](
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
			tables = append(tables, placeholder(column))
			values = append(values, value.Interface())
		}
	}
	query += strings.Join(tables, ", ")
	query += " WHERE "
	var wherePlaceholder []string
	for key, value := range whereMap {
		minKey := strings.ToLower(key)
		wherePlaceholder = append(wherePlaceholder, placeholder(minKey))
		values = append(values, value)
	}
	query += strings.Join(wherePlaceholder, " AND ") + ";"
	if len(notUpdatedFields) > 0 {
		errs.Set("prohibited fields", strings.Join(notUpdatedFields, ", "))
	}
	return query, values, errs.AsError()
}

func placeholder(name string) string {
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
