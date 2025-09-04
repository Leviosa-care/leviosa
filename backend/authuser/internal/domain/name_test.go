package domain

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	t.Run("valid names should pass", func(t *testing.T) {
		validNames := []string{
			"Jo",
			"John",
			"Jean-Pierre",
			"María José",
			"李小明",
			"François",
			"O'Connor",
			"van der Berg",
			strings.Repeat("A", NameMaxLength), // Exactly 50 characters
		}

		for _, name := range validNames {
			t.Run(name, func(t *testing.T) {
				err := validateName(name, "firstName")
				assert.NoError(t, err, "expected %s to be valid", name)
			})
		}
	})

	t.Run("valid names with whitespace should pass", func(t *testing.T) {
		namesWithSpaces := []string{
			"  John  ",
			"\tMaria\t",
			" Jean-Pierre ",
		}

		for _, name := range namesWithSpaces {
			t.Run(name, func(t *testing.T) {
				err := validateName(name, "firstName")
				assert.NoError(t, err, "expected %s to be valid after trimming", name)
			})
		}
	})

	t.Run("empty name should fail", func(t *testing.T) {
		err := validateName("", "firstName")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, NameRequiredKey)
		assert.Equal(t, "firstName is required", errMap[NameRequiredKey].Error())
	})

	t.Run("whitespace-only name should fail", func(t *testing.T) {
		err := validateName("   ", "firstName")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, NameRequiredKey)
		assert.Equal(t, "firstName is required", errMap[NameRequiredKey].Error())
	})

	t.Run("name too short should fail", func(t *testing.T) {
		shortNames := []string{
			"A",
			"X",
		}

		for _, name := range shortNames {
			t.Run(name, func(t *testing.T) {
				err := validateName(name, "lastName")
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, NameTooShortKey)
				assert.Equal(t, "lastName must be at least 2 characters long", errMap[NameTooShortKey].Error())
			})
		}
	})

	t.Run("name too long should fail", func(t *testing.T) {
		longName := strings.Repeat("A", NameMaxLength+1) // 51 characters
		err := validateName(longName, "firstName")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, NameTooLongKey)
		assert.Equal(t, "firstName must be no more than 50 characters long", errMap[NameTooLongKey].Error())
	})

	t.Run("name with dangerous characters should fail", func(t *testing.T) {
		dangerousNames := []string{
			"John<script>",
			"Maria>alert",
			"Jean;DROP",
			`John"Doe`,
			"Mary'Smith",
			"Tom&Jerry",
		}

		for _, name := range dangerousNames {
			t.Run(name, func(t *testing.T) {
				err := validateName(name, "firstName")
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, NameInvalidCharsKey)
				assert.Equal(t, "firstName contains invalid characters", errMap[NameInvalidCharsKey].Error())
			})
		}
	})

	t.Run("boundary cases", func(t *testing.T) {
		t.Run("minimum valid length", func(t *testing.T) {
			err := validateName("AB", "firstName")
			assert.NoError(t, err)
		})

		t.Run("maximum valid length", func(t *testing.T) {
			maxLengthName := strings.Repeat("A", NameMaxLength)
			err := validateName(maxLengthName, "firstName")
			assert.NoError(t, err)
		})

		t.Run("one character too short", func(t *testing.T) {
			err := validateName("A", "firstName")
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, NameTooShortKey)
		})

		t.Run("one character too long", func(t *testing.T) {
			tooLongName := strings.Repeat("A", NameMaxLength+1)
			err := validateName(tooLongName, "firstName")
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, NameTooLongKey)
		})
	})

	t.Run("multiple validation errors should be reported", func(t *testing.T) {
		// Too short AND contains dangerous characters
		err := validateName("<", "firstName")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))

		// Should have both length and character errors
		assert.Contains(t, errMap, NameTooShortKey)
		assert.Contains(t, errMap, NameInvalidCharsKey)
		assert.Equal(t, "firstName must be at least 2 characters long", errMap[NameTooShortKey].Error())
		assert.Equal(t, "firstName contains invalid characters", errMap[NameInvalidCharsKey].Error())
	})

	t.Run("different field names should be reflected in error messages", func(t *testing.T) {
		testCases := []struct {
			fieldName string
			expected  string
		}{
			{"firstName", "firstName is required"},
			{"lastName", "lastName is required"},
			{"middleName", "middleName is required"},
		}

		for _, tc := range testCases {
			t.Run(tc.fieldName, func(t *testing.T) {
				err := validateName("", tc.fieldName)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, NameRequiredKey)
				assert.Equal(t, tc.expected, errMap[NameRequiredKey].Error())
			})
		}
	})
}

