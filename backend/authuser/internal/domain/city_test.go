package domain

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidateCity(t *testing.T) {
	t.Run("valid cities should pass", func(t *testing.T) {
		validCities := []string{
			"Paris",
			"New York",
			"San Francisco",
			"Los Angeles",
			"Saint-Denis",
			"Aix-en-Provence",
			"México City",
			"São Paulo",
			"北京",                               // Beijing in Chinese
			"Москва",                           // Moscow in Russian
			"القاهرة",                          // Cairo in Arabic
			"مُحَمَّد آباد",                    // Example with Arabic diacritics
			strings.Repeat("A", CityMaxLength), // Exactly 100 characters
		}

		for _, city := range validCities {
			t.Run(city, func(t *testing.T) {
				err := validateCity(city)
				assert.NoError(t, err, "expected %s to be valid", city)
			})
		}
	})

	t.Run("cities with whitespace should pass after trimming", func(t *testing.T) {
		citiesWithSpaces := []string{
			"  Paris  ",
			"\tNew York\t",
			" San Francisco ",
		}

		for _, city := range citiesWithSpaces {
			t.Run(city, func(t *testing.T) {
				err := validateCity(city)
				assert.NoError(t, err, "expected %s to be valid after trimming", city)
			})
		}
	})

	t.Run("empty city should fail", func(t *testing.T) {
		err := validateCity("")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, CityRequiredKey)
		assert.Equal(t, CityRequiredMsg, errMap[CityRequiredKey].Error())
	})

	t.Run("whitespace-only city should fail", func(t *testing.T) {
		err := validateCity("   ")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, CityRequiredKey)
		assert.Equal(t, CityRequiredMsg, errMap[CityRequiredKey].Error())
	})

	t.Run("city too short should fail", func(t *testing.T) {
		shortCities := []string{
			"A",
			"X",
		}

		for _, city := range shortCities {
			t.Run(city, func(t *testing.T) {
				err := validateCity(city)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, CityTooShortKey)
				assert.Equal(t, CityTooShortMsg, errMap[CityTooShortKey].Error())
			})
		}
	})

	t.Run("city too long should fail", func(t *testing.T) {
		longCity := strings.Repeat("A", CityMaxLength+1) // 101 characters
		err := validateCity(longCity)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, CityTooLongKey)
		assert.Equal(t, CityTooLongMsg, errMap[CityTooLongKey].Error())
	})

	t.Run("city with dangerous characters should fail", func(t *testing.T) {
		dangerousCities := []string{
			"Paris<script>",
			"New>York",
			"San;Francisco",
			`Los"Angeles`,
			"Saint'Denis",
			"Mexico&City",
		}

		for _, city := range dangerousCities {
			t.Run(city, func(t *testing.T) {
				err := validateCity(city)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, CityInvalidCharsKey)
				assert.Equal(t, CityInvalidCharsMsg, errMap[CityInvalidCharsKey].Error())
			})
		}
	})

	t.Run("boundary cases", func(t *testing.T) {
		t.Run("minimum valid length", func(t *testing.T) {
			minLengthCity := strings.Repeat("A", CityMinLength) // Exactly 2 characters
			err := validateCity(minLengthCity)
			assert.NoError(t, err)
		})

		t.Run("maximum valid length", func(t *testing.T) {
			maxLengthCity := strings.Repeat("A", CityMaxLength) // Exactly 100 characters
			err := validateCity(maxLengthCity)
			assert.NoError(t, err)
		})

		t.Run("one character too short", func(t *testing.T) {
			tooShortCity := "A" // 1 character
			err := validateCity(tooShortCity)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, CityTooShortKey)
		})

		t.Run("one character too long", func(t *testing.T) {
			tooLongCity := strings.Repeat("A", CityMaxLength+1) // 101 characters
			err := validateCity(tooLongCity)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, CityTooLongKey)
		})
	})

	t.Run("multiple validation errors should be reported", func(t *testing.T) {
		// Too short AND contains dangerous characters
		err := validateCity("<")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))

		// Should have both length and character errors
		assert.Contains(t, errMap, CityTooShortKey)
		assert.Contains(t, errMap, CityInvalidCharsKey)
		assert.Equal(t, CityTooShortMsg, errMap[CityTooShortKey].Error())
		assert.Equal(t, CityInvalidCharsMsg, errMap[CityInvalidCharsKey].Error())
	})

	t.Run("international city names", func(t *testing.T) {
		t.Run("should accept various scripts", func(t *testing.T) {
			internationalCities := []string{
				"北京",       // Chinese
				"東京",       // Japanese
				"서울",       // Korean
				"Москва",   // Cyrillic
				"Αθήνα",    // Greek
				"القاهرة",  // Arabic
				"मुंबई",    // Devanagari
				"กรุงเทพฯ", // Thai
			}

			for _, city := range internationalCities {
				t.Run(city, func(t *testing.T) {
					err := validateCity(city)
					assert.NoError(t, err, "expected international city %s to be valid", city)
				})
			}
		})
	})

	t.Run("common city name patterns", func(t *testing.T) {
		t.Run("hyphenated names should pass", func(t *testing.T) {
			hyphenatedCities := []string{
				"Saint-Denis",
				"Aix-en-Provence",
				"Stratford-upon-Avon",
				"Kingston-upon-Thames",
			}

			for _, city := range hyphenatedCities {
				err := validateCity(city)
				assert.NoError(t, err, "expected hyphenated city %s to be valid", city)
			}
		})

		t.Run("apostrophe names should pass", func(t *testing.T) {
			// Note: These will fail due to our security restrictions
			apostropheCities := []string{
				"L'Aquila",
				"O'Fallon",
			}

			for _, city := range apostropheCities {
				t.Run(city, func(t *testing.T) {
					err := validateCity(city)
					assert.Error(t, err, "apostrophe cities should fail due to security restrictions")

					var errMap errsx.Map
					assert.True(t, errsx.As(err, &errMap))
					assert.Contains(t, errMap, CityInvalidCharsKey)
				})
			}
		})

		t.Run("names with spaces should pass", func(t *testing.T) {
			spacedCities := []string{
				"New York",
				"Los Angeles",
				"San Francisco",
				"Las Vegas",
				"Mexico City",
			}

			for _, city := range spacedCities {
				err := validateCity(city)
				assert.NoError(t, err, "expected spaced city %s to be valid", city)
			}
		})
	})
}

