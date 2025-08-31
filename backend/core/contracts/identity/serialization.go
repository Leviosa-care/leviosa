package identity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/core/errs"
)

// MarshalJSON implements json.Marshaler interface
func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON implements json.Unmarshaler interface
func (r *Role) UnmarshalJSON(data []byte) error {
	var roleStr string
	if err := json.Unmarshal(data, &roleStr); err != nil {
		return err
	}

	role, err := ParseRole(roleStr)
	if err != nil {
		return err
	}

	*r = role
	return nil
}

// Value implements driver.Valuer interface for database storage
func (r Role) Value() (driver.Value, error) {
	return r.String(), nil
}

// Scan implements sql.Scanner interface for database scanning
func (r *Role) Scan(value interface{}) error {
	if value == nil {
		*r = Visitor
		return nil
	}

	switch v := value.(type) {
	case string:
		role, err := ParseRole(v)
		if err != nil {
			return err
		}
		*r = role
		return nil
	case []byte:
		role, err := ParseRole(string(v))
		if err != nil {
			return err
		}
		*r = role
		return nil
	case int64:
		role := Role(v)
		if !role.IsValid() {
			return errs.NewInvalidValueErr(fmt.Sprintf("invalid role value: %d", v))
		}
		*r = role
		return nil
	default:
		return errs.NewInvalidValueErr(fmt.Sprintf("cannot scan %T into Role", value))
	}
}

