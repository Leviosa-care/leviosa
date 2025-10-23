package identity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	t.Run("valid roles should marshal to JSON strings", func(t *testing.T) {
		tests := []struct {
			role     Role
			expected string
		}{
			{Visitor, `"visitor"`},
			{Standard, `"standard"`},
			{Administrator, `"administrator"`},
		}

		for _, test := range tests {
			t.Run(test.role.String(), func(t *testing.T) {
				data, err := json.Marshal(test.role)
				require.NoError(t, err)
				assert.Equal(t, test.expected, string(data))
			})
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("valid JSON should unmarshal correctly", func(t *testing.T) {
		tests := []struct {
			json     string
			expected Role
		}{
			{`"visitor"`, Visitor},
			{`"standard"`, Standard},
			{`"administrator"`, Administrator},
		}

		for _, test := range tests {
			t.Run(test.json, func(t *testing.T) {
				var role Role
				err := json.Unmarshal([]byte(test.json), &role)
				require.NoError(t, err)
				assert.Equal(t, test.expected, role)
			})
		}
	})

	t.Run("invalid JSON should return error", func(t *testing.T) {
		invalidJSON := []string{
			`"invalid_role"`,
			`"admin"`,
			`123`,
			`null`,
			`""`,
		}

		for _, invalid := range invalidJSON {
			t.Run(invalid, func(t *testing.T) {
				var role Role
				err := json.Unmarshal([]byte(invalid), &role)
				assert.Error(t, err)
			})
		}
	})
}

func TestValue(t *testing.T) {
	t.Run("should return string representation for database storage", func(t *testing.T) {
		tests := []struct {
			role     Role
			expected driver.Value
		}{
			{Visitor, "visitor"},
			{Administrator, "administrator"},
			{Premium, "premium"},
		}

		for _, test := range tests {
			t.Run(test.role.String(), func(t *testing.T) {
				value, err := test.role.Value()
				require.NoError(t, err)
				assert.Equal(t, test.expected, value)
			})
		}
	})
}

func TestScan(t *testing.T) {
	t.Run("should scan string values correctly", func(t *testing.T) {
		var role Role
		err := role.Scan("administrator")
		require.NoError(t, err)
		assert.Equal(t, Administrator, role)
	})

	t.Run("should scan byte slice values correctly", func(t *testing.T) {
		var role Role
		err := role.Scan([]byte("premium"))
		require.NoError(t, err)
		assert.Equal(t, Premium, role)
	})

	t.Run("should scan int64 values correctly", func(t *testing.T) {
		var role Role
		err := role.Scan(int64(5)) // Administrator
		require.NoError(t, err)
		assert.Equal(t, Administrator, role)
	})

	t.Run("should handle nil values", func(t *testing.T) {
		var role Role
		err := role.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, Visitor, role)
	})

	t.Run("should return error for invalid string values", func(t *testing.T) {
		var role Role
		err := role.Scan("invalid_role")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrInvalidValue))
	})

	t.Run("should return error for invalid int64 values", func(t *testing.T) {
		var role Role
		err := role.Scan(int64(99))
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrInvalidValue))
	})

	t.Run("should return error for unsupported types", func(t *testing.T) {
		var role Role
		err := role.Scan(123.45) // float64
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrInvalidValue))
	})
}

