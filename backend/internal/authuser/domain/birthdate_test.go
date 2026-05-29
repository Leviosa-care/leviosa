package domain

import (
	"testing"
	"time"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidateBirthDate(t *testing.T) {
	now := time.Now()

	t.Run("valid birth dates should pass", func(t *testing.T) {
		validBirthDates := []time.Time{
			now.AddDate(-MinAgeYears, 0, -1), // 18 years and 1 day old
			now.AddDate(-20, 0, 0),           // 20 years old
			now.AddDate(-30, -6, -15),        // 30 years, 6 months, 15 days old
			now.AddDate(-50, 0, 0),           // 50 years old
			now.AddDate(-MaxAgeYears, 0, 1),  // 119 years, 11 months, 29 days old (just under 120)
		}

		for i, birthdate := range validBirthDates {
			t.Run(birthdate.Format("2006-01-02"), func(t *testing.T) {
				err := validateBirthDate(birthdate)
				assert.NoError(t, err, "expected valid birth date #%d to pass: %s", i+1, birthdate.Format("2006-01-02"))
			})
		}
	})

	t.Run("future birth date should fail", func(t *testing.T) {
		futureDates := []time.Time{
			now.Add(time.Hour),    // 1 hour from now
			now.AddDate(0, 0, 1),  // Tomorrow
			now.AddDate(0, 1, 0),  // 1 month from now
			now.AddDate(1, 0, 0),  // 1 year from now
			now.AddDate(10, 0, 0), // 10 years from now
		}

		for _, birthdate := range futureDates {
			t.Run(birthdate.Format("2006-01-02"), func(t *testing.T) {
				err := validateBirthDate(birthdate)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, BirthDateFutureKey)
				assert.Equal(t, BirthDateFutureMsg, errMap[BirthDateFutureKey].Error())
			})
		}
	})

	t.Run("too young birth date should fail", func(t *testing.T) {
		tooYoungDates := []time.Time{
			now,                             // Born today
			now.AddDate(0, 0, -1),           // Born yesterday
			now.AddDate(-1, 0, 0),           // 1 year old
			now.AddDate(-17, 0, 0),          // 17 years old
			now.AddDate(-MinAgeYears, 0, 1), // 17 years, 11 months, 29 days old (just under 18)
		}

		for _, birthdate := range tooYoungDates {
			t.Run(birthdate.Format("2006-01-02"), func(t *testing.T) {
				err := validateBirthDate(birthdate)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, BirthDateTooYoungKey)
				assert.Equal(t, BirthDateTooYoungMsg, errMap[BirthDateTooYoungKey].Error())
			})
		}
	})

	t.Run("too old birth date should fail", func(t *testing.T) {
		tooOldDates := []time.Time{
			now.AddDate(-MaxAgeYears, 0, -1),  // 120 years and 1 day old
			now.AddDate(-MaxAgeYears-1, 0, 0), // 121 years old
			now.AddDate(-150, 0, 0),           // 150 years old
			now.AddDate(-200, 0, 0),           // 200 years old
		}

		for _, birthdate := range tooOldDates {
			t.Run(birthdate.Format("2006-01-02"), func(t *testing.T) {
				err := validateBirthDate(birthdate)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, BirthDateTooOldKey)
				assert.Equal(t, BirthDateTooOldMsg, errMap[BirthDateTooOldKey].Error())
			})
		}
	})

	t.Run("boundary cases", func(t *testing.T) {
		t.Run("exactly minimum age should pass", func(t *testing.T) {
			// Exactly 18 years old
			exactMinAge := now.AddDate(-MinAgeYears, 0, 0)
			err := validateBirthDate(exactMinAge)
			assert.NoError(t, err, "expected exactly %d years old to be valid", MinAgeYears)
		})

		t.Run("just under minimum age should fail", func(t *testing.T) {
			// 18 years minus 1 day
			justUnderMinAge := now.AddDate(-MinAgeYears, 0, 1)
			err := validateBirthDate(justUnderMinAge)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, BirthDateTooYoungKey)
		})

		t.Run("just within maximum age should pass", func(t *testing.T) {
			// 1 day under 120 years — testing "exactly 120" is inherently flaky because
			// the validator calls time.Now() independently, making the boundary shift by
			// nanoseconds and causing a spurious "too old" error on exact matches.
			justWithinMaxAge := now.AddDate(-MaxAgeYears, 0, 1)
			err := validateBirthDate(justWithinMaxAge)
			assert.NoError(t, err, "expected just within %d years old to be valid", MaxAgeYears)
		})

		t.Run("just over maximum age should fail", func(t *testing.T) {
			// 120 years plus 1 day
			justOverMaxAge := now.AddDate(-MaxAgeYears, 0, -1)
			err := validateBirthDate(justOverMaxAge)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, BirthDateTooOldKey)
		})

		t.Run("exactly current time should fail", func(t *testing.T) {
			// Born at exactly this moment
			err := validateBirthDate(now)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, BirthDateTooYoungKey) // Should fail as too young, not future
		})

		t.Run("one second in future should fail as future", func(t *testing.T) {
			// Born one second from now
			oneSecondFuture := now.Add(time.Second)
			err := validateBirthDate(oneSecondFuture)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, BirthDateFutureKey)
		})
	})

	t.Run("edge cases with leap years", func(t *testing.T) {
		// Test leap year boundary conditions
		t.Run("leap year birth date", func(t *testing.T) {
			// Born on leap day, exactly 20 years ago (assuming 2004 was a leap year)
			leapYearBirth := time.Date(now.Year()-20, 2, 29, 0, 0, 0, 0, time.UTC)
			// Only test if the leap day actually exists for that year
			if leapYearBirth.Day() == 29 {
				err := validateBirthDate(leapYearBirth)
				assert.NoError(t, err, "leap year birth date should be valid")
			}
		})
	})

	t.Run("different time zones should work", func(t *testing.T) {
		// Test with different time zones
		utc := time.UTC
		est := time.FixedZone("EST", -5*60*60)
		pst := time.FixedZone("PST", -8*60*60)

		validAge := 25

		birthDates := []time.Time{
			time.Date(now.Year()-validAge, 6, 15, 10, 30, 0, 0, utc),
			time.Date(now.Year()-validAge, 6, 15, 10, 30, 0, 0, est),
			time.Date(now.Year()-validAge, 6, 15, 10, 30, 0, 0, pst),
		}

		for i, birthdate := range birthDates {
			t.Run(birthdate.Location().String(), func(t *testing.T) {
				err := validateBirthDate(birthdate)
				assert.NoError(t, err, "expected birth date #%d with timezone %s to be valid", i+1, birthdate.Location())
			})
		}
	})

	t.Run("historical dates", func(t *testing.T) {
		t.Run("valid historical dates", func(t *testing.T) {
			historicalDates := []time.Time{
				time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC),   // Post-WWII
				time.Date(1960, 7, 15, 0, 0, 0, 0, time.UTC),  // 1960s
				time.Date(1975, 12, 25, 0, 0, 0, 0, time.UTC), // 1970s
				time.Date(1990, 3, 10, 0, 0, 0, 0, time.UTC),  // 1990s
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),   // Y2K
			}

			for _, birthdate := range historicalDates {
				t.Run(birthdate.Format("2006-01-02"), func(t *testing.T) {
					err := validateBirthDate(birthdate)
					// Only check if the person wouldn't be too old
					age := now.Year() - birthdate.Year()
					if age <= MaxAgeYears {
						assert.NoError(t, err, "expected historical date %s to be valid", birthdate.Format("2006-01-02"))
					} else {
						assert.Error(t, err, "expected historical date %s to be too old", birthdate.Format("2006-01-02"))
					}
				})
			}
		})
	})
}

