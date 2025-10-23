package domain

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestGenderInput_ValidateGender(t *testing.T) {
	t.Run("valid predefined genders should pass", func(t *testing.T) {
		validGenders := []struct {
			gender       Gender
			customGender string
		}{
			{GenderMan, ""},
			{GenderWoman, ""},
			{GenderNonBinary, ""},
			{GenderPreferNotToSay, ""},
		}

		for _, tc := range validGenders {
			t.Run(string(tc.gender), func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       tc.gender,
					CustomGender: tc.customGender,
				}
				err := genderInput.ValidateGender()
				assert.NoError(t, err, "expected predefined gender %s to be valid", tc.gender)
			})
		}
	})

	t.Run("valid custom genders should pass", func(t *testing.T) {
		validCustomGenders := []string{
			"Non-binary",
			"Agender",
			"Genderfluid",
			"Two-spirit",
			"Demigender",
			"Pangender",
			"Bigender",
			"Gender questioning",
			"Other",
			"Prefer to self-describe",
			strings.Repeat("A", CustomGenderMaxLength), // Exactly 50 characters
		}

		for _, customGender := range validCustomGenders {
			t.Run(customGender, func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       GenderCustom,
					CustomGender: customGender,
				}
				err := genderInput.ValidateGender()
				assert.NoError(t, err, "expected custom gender %s to be valid", customGender)
			})
		}
	})

	t.Run("custom gender with whitespace should pass after trimming", func(t *testing.T) {
		customGendersWithSpaces := []string{
			"  Non-binary  ",
			"\tAgender\t",
			" Genderfluid ",
		}

		for _, customGender := range customGendersWithSpaces {
			t.Run(customGender, func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       GenderCustom,
					CustomGender: customGender,
				}
				err := genderInput.ValidateGender()
				assert.NoError(t, err, "expected custom gender %s to be valid after trimming", customGender)
			})
		}
	})

	t.Run("invalid gender value should fail", func(t *testing.T) {
		invalidGenders := []Gender{
			"invalid",
			"unknown",
			"male",
			"female",
			"other",
		}

		for _, gender := range invalidGenders {
			t.Run(string(gender), func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       gender,
					CustomGender: "",
				}
				err := genderInput.ValidateGender()
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, GenderInvalidKey)
				assert.Contains(t, errMap[GenderInvalidKey].Error(), GenderInvalidMsg)
				assert.Contains(t, errMap[GenderInvalidKey].Error(), string(gender))
			})
		}
	})

	t.Run("custom gender without custom value should fail", func(t *testing.T) {
		testCases := []struct {
			name         string
			customGender string
		}{
			{"empty string", ""},
			{"whitespace only", "   "},
			{"tabs only", "\t\t"},
			{"mixed whitespace", " \t \n "},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       GenderCustom,
					CustomGender: tc.customGender,
				}
				err := genderInput.ValidateGender()
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, CustomGenderRequiredKey)
				assert.Equal(t, CustomGenderRequiredMsg, errMap[CustomGenderRequiredKey].Error())
			})
		}
	})

	t.Run("custom gender too long should fail", func(t *testing.T) {
		longCustomGender := strings.Repeat("A", CustomGenderMaxLength+1) // 51 characters
		genderInput := &GenderInput{
			Gender:       GenderCustom,
			CustomGender: longCustomGender,
		}
		err := genderInput.ValidateGender()
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, CustomGenderTooLongKey)
		assert.Equal(t, CustomGenderTooLongMsg, errMap[CustomGenderTooLongKey].Error())
	})

	t.Run("custom gender with dangerous characters should fail", func(t *testing.T) {
		dangerousCustomGenders := []string{
			"Non-binary<script>",
			"Agender>alert",
			"Gender;DROP",
			`Two"spirit`,
			"Demi'gender",
			"Pan&gender",
		}

		for _, customGender := range dangerousCustomGenders {
			t.Run(customGender, func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       GenderCustom,
					CustomGender: customGender,
				}
				err := genderInput.ValidateGender()
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, CustomGenderCharsKey)
				assert.Equal(t, CustomGenderCharsMsg, errMap[CustomGenderCharsKey].Error())
			})
		}
	})

	t.Run("boundary cases", func(t *testing.T) {
		t.Run("maximum valid custom gender length", func(t *testing.T) {
			maxLengthCustomGender := strings.Repeat("A", CustomGenderMaxLength) // Exactly 50 characters
			genderInput := &GenderInput{
				Gender:       GenderCustom,
				CustomGender: maxLengthCustomGender,
			}
			err := genderInput.ValidateGender()
			assert.NoError(t, err)
		})

		t.Run("one character too long custom gender", func(t *testing.T) {
			tooLongCustomGender := strings.Repeat("A", CustomGenderMaxLength+1) // 51 characters
			genderInput := &GenderInput{
				Gender:       GenderCustom,
				CustomGender: tooLongCustomGender,
			}
			err := genderInput.ValidateGender()
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, CustomGenderTooLongKey)
		})
	})

	t.Run("multiple validation errors should be reported", func(t *testing.T) {
		// Too long AND contains dangerous characters
		invalidCustomGender := strings.Repeat("A", CustomGenderMaxLength+1) + "<script>"
		genderInput := &GenderInput{
			Gender:       GenderCustom,
			CustomGender: invalidCustomGender,
		}
		err := genderInput.ValidateGender()
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))

		// Should have both length and character errors
		assert.Contains(t, errMap, CustomGenderTooLongKey)
		assert.Contains(t, errMap, CustomGenderCharsKey)
		assert.Equal(t, CustomGenderTooLongMsg, errMap[CustomGenderTooLongKey].Error())
		assert.Equal(t, CustomGenderCharsMsg, errMap[CustomGenderCharsKey].Error())
	})

	t.Run("predefined genders should ignore custom gender value", func(t *testing.T) {
		// Even if custom gender is provided, predefined genders should pass
		testCases := []struct {
			gender       Gender
			customGender string
		}{
			{GenderMan, "invalid<script>"},
			{GenderWoman, strings.Repeat("A", 100)}, // Way too long
			{GenderNonBinary, ""},
			{GenderPreferNotToSay, "whatever"},
		}

		for _, tc := range testCases {
			t.Run(string(tc.gender), func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       tc.gender,
					CustomGender: tc.customGender,
				}
				err := genderInput.ValidateGender()
				assert.NoError(t, err, "predefined gender %s should pass regardless of custom gender value", tc.gender)
			})
		}
	})

	t.Run("gender constants should have expected values", func(t *testing.T) {
		assert.Equal(t, "man", string(GenderMan))
		assert.Equal(t, "woman", string(GenderWoman))
		assert.Equal(t, "non_binary", string(GenderNonBinary))
		assert.Equal(t, "prefer_not_to_say", string(GenderPreferNotToSay))
		assert.Equal(t, "custom", string(GenderCustom))
	})

	t.Run("common custom gender values", func(t *testing.T) {
		commonCustomGenders := []string{
			"Agender",
			"Androgyne",
			"Bigender",
			"Demiboy",
			"Demigirl",
			"Genderfluid",
			"Genderqueer",
			"Non-binary",
			"Pangender",
			"Two-spirit",
			"Third gender",
			"Gender variant",
			"Gender non-conforming",
		}

		for _, customGender := range commonCustomGenders {
			t.Run(customGender, func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       GenderCustom,
					CustomGender: customGender,
				}
				err := genderInput.ValidateGender()
				assert.NoError(t, err, "common custom gender %s should be valid", customGender)
			})
		}
	})
}

func TestGender_String(t *testing.T) {
	testCases := []struct {
		gender   Gender
		expected string
	}{
		{GenderMan, "man"},
		{GenderWoman, "woman"},
		{GenderNonBinary, "non_binary"},
		{GenderPreferNotToSay, "prefer_not_to_say"},
		{GenderCustom, "custom"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.gender), func(t *testing.T) {
			result := tc.gender.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateCustomGender(t *testing.T) {
	t.Run("valid custom genders should pass", func(t *testing.T) {
		validGenders := []string{
			"Non-binary",
			"Agender",
			"Two-spirit",
			"A", // Minimum meaningful length
			strings.Repeat("A", CustomGenderMaxLength), // Maximum length
		}

		for _, customGender := range validGenders {
			t.Run(customGender, func(t *testing.T) {
				genderInput := &GenderInput{
					Gender:       GenderCustom,
					CustomGender: customGender,
				}
				err := validateCustomGender(genderInput)
				assert.NoError(t, err, "expected custom gender %s to be valid", customGender)
			})
		}
	})

	t.Run("empty custom gender should fail", func(t *testing.T) {
		genderInput := &GenderInput{
			Gender:       GenderCustom,
			CustomGender: "",
		}
		err := validateCustomGender(genderInput)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, CustomGenderRequiredKey)
	})

	t.Run("too long custom gender should fail", func(t *testing.T) {
		genderInput := &GenderInput{
			Gender:       GenderCustom,
			CustomGender: strings.Repeat("A", CustomGenderMaxLength+1),
		}
		err := validateCustomGender(genderInput)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, CustomGenderTooLongKey)
	})

	t.Run("dangerous characters should fail", func(t *testing.T) {
		genderInput := &GenderInput{
			Gender:       GenderCustom,
			CustomGender: "Test<script>",
		}
		err := validateCustomGender(genderInput)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, CustomGenderCharsKey)
	})
}

