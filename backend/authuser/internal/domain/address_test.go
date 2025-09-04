package domain

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidateAddress1(t *testing.T) {
	t.Run("valid address1 should pass", func(t *testing.T) {
		validAddresses := []string{
			"123 Main St",
			"456 Elm Street Apartment 2B",
			"789 Oak Avenue, Unit 45",
			"1000 Broadway Suite 500",
			"12345 Very Long Street Name With Many Words",
			strings.Repeat("A", AddressMaxLength), // Exactly 200 characters
		}

		for _, address := range validAddresses {
			t.Run(address, func(t *testing.T) {
				err := validateAddress1(address)
				assert.NoError(t, err, "expected %s to be valid", address)
			})
		}
	})

	t.Run("invalid address1 should fail", func(t *testing.T) {
		t.Run("empty address1", func(t *testing.T) {
			err := validateAddress1("")
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, AddressRequiredKey)
			assert.Equal(t, "address1 is required", errMap[AddressRequiredKey].Error())
		})

		t.Run("too short address1", func(t *testing.T) {
			err := validateAddress1("123")
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, AddressTooShortKey)
			assert.Equal(t, "address1 must be at least 5 characters long", errMap[AddressTooShortKey].Error())
		})

		t.Run("too long address1", func(t *testing.T) {
			longAddress := strings.Repeat("A", AddressMaxLength+1) // 201 characters
			err := validateAddress1(longAddress)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, AddressTooLongKey)
			assert.Equal(t, "address1 must be no more than 200 characters long", errMap[AddressTooLongKey].Error())
		})
	})
}

func TestValidateAddress2(t *testing.T) {
	t.Run("empty address2 should pass", func(t *testing.T) {
		err := validateAddress2("")
		assert.NoError(t, err, "address2 is optional, empty should be valid")
	})

	t.Run("valid address2 should pass", func(t *testing.T) {
		validAddresses := []string{
			"Apartment 2B",
			"Unit 45",
			"Suite 500",
			"Second Floor",
			strings.Repeat("A", AddressMaxLength), // Exactly 200 characters
		}

		for _, address := range validAddresses {
			t.Run(address, func(t *testing.T) {
				err := validateAddress2(address)
				assert.NoError(t, err, "expected %s to be valid", address)
			})
		}
	})

	t.Run("invalid address2 should fail", func(t *testing.T) {
		t.Run("too short address2", func(t *testing.T) {
			err := validateAddress2("Apt")
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, AddressTooShortKey)
			assert.Equal(t, "address2 must be at least 5 characters long", errMap[AddressTooShortKey].Error())
		})

		t.Run("too long address2", func(t *testing.T) {
			longAddress := strings.Repeat("A", AddressMaxLength+1) // 201 characters
			err := validateAddress2(longAddress)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, AddressTooLongKey)
			assert.Equal(t, "address2 must be no more than 200 characters long", errMap[AddressTooLongKey].Error())
		})
	})
}

func TestValidateAddress(t *testing.T) {
	t.Run("valid addresses should pass", func(t *testing.T) {
		validAddresses := []string{
			"123 Main St",
			"456 Elm Street",
			"789 Oak Avenue, Unit 45",
			"1000 Broadway Suite 500",
			"12345 Very Long Street Name",
			"Rue de la Paix 123",
			"Bahnhofstraße 456",
			"улица Пушкина 789",
		}

		for _, address := range validAddresses {
			t.Run(address, func(t *testing.T) {
				err := validateAddress(address, "testField")
				assert.NoError(t, err, "expected %s to be valid", address)
			})
		}
	})

	t.Run("addresses with whitespace should pass after trimming", func(t *testing.T) {
		addressesWithSpaces := []string{
			"  123 Main St  ",
			"\t456 Elm Street\t",
			" 789 Oak Avenue ",
		}

		for _, address := range addressesWithSpaces {
			t.Run(address, func(t *testing.T) {
				err := validateAddress(address, "testField")
				assert.NoError(t, err, "expected %s to be valid after trimming", address)
			})
		}
	})

	t.Run("empty address should fail", func(t *testing.T) {
		err := validateAddress("", "testField")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, AddressRequiredKey)
		assert.Equal(t, "testField is required", errMap[AddressRequiredKey].Error())
	})

	t.Run("whitespace-only address should fail", func(t *testing.T) {
		err := validateAddress("   ", "testField")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, AddressRequiredKey)
		assert.Equal(t, "testField is required", errMap[AddressRequiredKey].Error())
	})

	t.Run("address too short should fail", func(t *testing.T) {
		shortAddresses := []string{
			"123",
			"A",
			"12B",
			"Main",
		}

		for _, address := range shortAddresses {
			t.Run(address, func(t *testing.T) {
				err := validateAddress(address, "testField")
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, AddressTooShortKey)
				assert.Equal(t, "testField must be at least 5 characters long", errMap[AddressTooShortKey].Error())
			})
		}
	})

	t.Run("address too long should fail", func(t *testing.T) {
		longAddress := strings.Repeat("A", AddressMaxLength+1) // 201 characters
		err := validateAddress(longAddress, "testField")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, AddressTooLongKey)
		assert.Equal(t, "testField must be no more than 200 characters long", errMap[AddressTooLongKey].Error())
	})

	t.Run("address with dangerous characters should fail", func(t *testing.T) {
		dangerousAddresses := []string{
			"123 Main<script>",
			"456 Elm>alert",
			"789 Oak;DROP",
			`123 Main"Street`,
			"456 Elm'Avenue",
			"789 Oak&Boulevard",
		}

		for _, address := range dangerousAddresses {
			t.Run(address, func(t *testing.T) {
				err := validateAddress(address, "testField")
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, AddressInvalidCharsKey)
				assert.Equal(t, "testField contains invalid characters", errMap[AddressInvalidCharsKey].Error())
			})
		}
	})

	t.Run("boundary cases", func(t *testing.T) {
		t.Run("minimum valid length", func(t *testing.T) {
			minLengthAddress := strings.Repeat("A", AddressMinLength) // Exactly 5 characters
			err := validateAddress(minLengthAddress, "testField")
			assert.NoError(t, err)
		})

		t.Run("maximum valid length", func(t *testing.T) {
			maxLengthAddress := strings.Repeat("A", AddressMaxLength) // Exactly 200 characters
			err := validateAddress(maxLengthAddress, "testField")
			assert.NoError(t, err)
		})

		t.Run("one character too short", func(t *testing.T) {
			tooShortAddress := strings.Repeat("A", AddressMinLength-1) // 4 characters
			err := validateAddress(tooShortAddress, "testField")
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, AddressTooShortKey)
		})

		t.Run("one character too long", func(t *testing.T) {
			tooLongAddress := strings.Repeat("A", AddressMaxLength+1) // 201 characters
			err := validateAddress(tooLongAddress, "testField")
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, AddressTooLongKey)
		})
	})

	t.Run("multiple validation errors should be reported", func(t *testing.T) {
		// Too short AND contains dangerous characters
		err := validateAddress("<>&", "testField")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))

		// Should have both length and character errors
		assert.Contains(t, errMap, AddressTooShortKey)
		assert.Contains(t, errMap, AddressInvalidCharsKey)
		assert.Equal(t, "testField must be at least 5 characters long", errMap[AddressTooShortKey].Error())
		assert.Equal(t, "testField contains invalid characters", errMap[AddressInvalidCharsKey].Error())
	})

	t.Run("different field names should be reflected in error messages", func(t *testing.T) {
		testCases := []struct {
			fieldName string
			expected  string
		}{
			{"address1", "address1 is required"},
			{"address2", "address2 is required"},
			{"homeAddress", "homeAddress is required"},
		}

		for _, tc := range testCases {
			t.Run(tc.fieldName, func(t *testing.T) {
				err := validateAddress("", tc.fieldName)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, AddressRequiredKey)
				assert.Equal(t, tc.expected, errMap[AddressRequiredKey].Error())
			})
		}
	})
}

