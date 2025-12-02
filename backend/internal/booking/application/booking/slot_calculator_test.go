package booking

import (
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAvailableSlots(t *testing.T) {
	location := time.UTC

	t.Run("basic slot generation - 30min product in 2hr availability", func(t *testing.T) {
		// Availability: 10:00-12:00 (2 hours)
		// Product: 30 minutes
		// Expected: 10 slots (10:00, 10:10, 10:20, ..., 11:30)
		// Note: Slots are every 10 minutes, last slot at 11:30-12:00
		availability := &domain.Availability{
			StartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			EndTime:   time.Date(2025, 12, 1, 12, 0, 0, 0, location),
		}

		slots := CalculateAvailableSlots(availability, 30, nil)

		assert.Len(t, slots, 10, "Should generate 10 slots for 30-min product in 2-hour window")

		// Verify first slot
		assert.Equal(t, time.Date(2025, 12, 1, 10, 0, 0, 0, location), slots[0].StartTime)
		assert.Equal(t, time.Date(2025, 12, 1, 10, 30, 0, 0, location), slots[0].EndTime)

		// Verify last slot
		assert.Equal(t, time.Date(2025, 12, 1, 11, 30, 0, 0, location), slots[9].StartTime)
		assert.Equal(t, time.Date(2025, 12, 1, 12, 0, 0, 0, location), slots[9].EndTime)

		// Verify all slots are 10 minutes apart
		for i := 1; i < len(slots); i++ {
			diff := slots[i].StartTime.Sub(slots[i-1].StartTime)
			assert.Equal(t, 10*time.Minute, diff, "Slots should be 10 minutes apart")
		}
	})

	t.Run("alignment - availability starts at non-aligned time", func(t *testing.T) {
		// Availability: 10:03-11:27 (misaligned)
		// Product: 30 minutes
		// Expected: First slot at 10:10 (aligned), last possible at 10:50
		availability := &domain.Availability{
			StartTime: time.Date(2025, 12, 1, 10, 3, 0, 0, location),
			EndTime:   time.Date(2025, 12, 1, 11, 27, 0, 0, location),
		}

		slots := CalculateAvailableSlots(availability, 30, nil)

		assert.Len(t, slots, 5, "Should generate 5 aligned slots")

		// First slot should be at 10:10 (first aligned boundary after 10:03)
		assert.Equal(t, time.Date(2025, 12, 1, 10, 10, 0, 0, location), slots[0].StartTime)

		// Last slot should be at 10:50 (last slot that fits before 11:27)
		assert.Equal(t, time.Date(2025, 12, 1, 10, 50, 0, 0, location), slots[4].StartTime)
		assert.Equal(t, time.Date(2025, 12, 1, 11, 20, 0, 0, location), slots[4].EndTime)
	})

	t.Run("overlap exclusion - existing booking blocks slots", func(t *testing.T) {
		// Availability: 10:00-12:00 (10 slots for 30-min product)
		// Product: 30 minutes
		// Existing booking: 10:20-10:50
		// Blocked slots: Any slot that overlaps with 10:20-10:50
		//   - 10:00-10:30 overlaps (ends at 10:30 which is after 10:20)
		//   - 10:10-10:40 overlaps
		//   - 10:20-10:50 overlaps (exact match)
		//   - 10:30-11:00 overlaps (starts at 10:30 which is before 10:50)
		//   - 10:40-11:10 overlaps (starts at 10:40 which is before 10:50)
		// Available slots: 10:50, 11:00, 11:10, 11:20, 11:30 (5 slots)
		availability := &domain.Availability{
			StartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			EndTime:   time.Date(2025, 12, 1, 12, 0, 0, 0, location),
		}

		existingBooking := &domain.Booking{
			ID:            uuid.New(),
			Status:        domain.BookingStatusConfirmed,
			SlotStartTime: time.Date(2025, 12, 1, 10, 20, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 10, 50, 0, 0, location),
		}

		slots := CalculateAvailableSlots(availability, 30, []*domain.Booking{existingBooking})

		// Should have 10 total slots - 5 blocked = 5 available
		assert.Len(t, slots, 5, "Should exclude overlapping slots")

		// Verify first available slot starts after the booking ends
		assert.Equal(t, time.Date(2025, 12, 1, 10, 50, 0, 0, location), slots[0].StartTime)
		assert.Equal(t, time.Date(2025, 12, 1, 11, 20, 0, 0, location), slots[0].EndTime)

		// Verify no slot starts before the booking ends
		for _, slot := range slots {
			assert.False(t, slot.StartTime.Before(existingBooking.SlotEndTime),
				"All available slots should start at or after booking end time")
		}
	})

	t.Run("cancelled booking does not block slots", func(t *testing.T) {
		// Availability: 10:00-11:00 (1 hour)
		// Product: 30 minutes
		// Expected: 4 slots (10:00, 10:10, 10:20, 10:30)
		// Cancelled booking: 10:20-10:50 (should NOT block any slots)
		availability := &domain.Availability{
			StartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			EndTime:   time.Date(2025, 12, 1, 11, 0, 0, 0, location),
		}

		cancelledBooking := &domain.Booking{
			ID:            uuid.New(),
			Status:        domain.BookingStatusCancelled,
			SlotStartTime: time.Date(2025, 12, 1, 10, 20, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 10, 50, 0, 0, location),
		}

		slots := CalculateAvailableSlots(availability, 30, []*domain.Booking{cancelledBooking})

		// Should have 4 slots (cancelled booking doesn't block)
		assert.Len(t, slots, 4, "Cancelled bookings should not block slots")

		// Verify slot at 10:20 is available (would be blocked if booking was active)
		hasSlotAt1020 := false
		for _, slot := range slots {
			if slot.StartTime.Equal(time.Date(2025, 12, 1, 10, 20, 0, 0, location)) {
				hasSlotAt1020 = true
				break
			}
		}
		assert.True(t, hasSlotAt1020, "Slot at 10:20 should be available despite cancelled booking")
	})

	t.Run("no slots available - product too long for availability", func(t *testing.T) {
		// Availability: 10:00-10:25 (25 minutes)
		// Product: 30 minutes
		// Expected: No slots fit
		availability := &domain.Availability{
			StartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			EndTime:   time.Date(2025, 12, 1, 10, 25, 0, 0, location),
		}

		slots := CalculateAvailableSlots(availability, 30, nil)

		assert.Len(t, slots, 0, "Should return empty when product doesn't fit")
	})

	t.Run("slot exactly fits at end", func(t *testing.T) {
		// Availability: 10:00-10:30 (exactly 30 minutes)
		// Product: 30 minutes
		// Expected: 1 slot at 10:00-10:30
		availability := &domain.Availability{
			StartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			EndTime:   time.Date(2025, 12, 1, 10, 30, 0, 0, location),
		}

		slots := CalculateAvailableSlots(availability, 30, nil)

		assert.Len(t, slots, 1, "Should generate exactly 1 slot that fits perfectly")
		assert.Equal(t, availability.StartTime, slots[0].StartTime)
		assert.Equal(t, availability.EndTime, slots[0].EndTime)
	})

	t.Run("60-minute product in 2-hour availability", func(t *testing.T) {
		// Availability: 09:00-11:00 (2 hours)
		// Product: 60 minutes
		// Expected: 7 slots (09:00, 09:10, 09:20, ..., 10:00)
		availability := &domain.Availability{
			StartTime: time.Date(2025, 12, 1, 9, 0, 0, 0, location),
			EndTime:   time.Date(2025, 12, 1, 11, 0, 0, 0, location),
		}

		slots := CalculateAvailableSlots(availability, 60, nil)

		assert.Len(t, slots, 7, "Should generate 7 slots for 60-min product")

		// First slot
		assert.Equal(t, time.Date(2025, 12, 1, 9, 0, 0, 0, location), slots[0].StartTime)
		assert.Equal(t, time.Date(2025, 12, 1, 10, 0, 0, 0, location), slots[0].EndTime)

		// Last slot
		assert.Equal(t, time.Date(2025, 12, 1, 10, 0, 0, 0, location), slots[6].StartTime)
		assert.Equal(t, time.Date(2025, 12, 1, 11, 0, 0, 0, location), slots[6].EndTime)
	})
}

func TestAlignToBaseSlot(t *testing.T) {
	location := time.UTC

	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "already aligned - :00",
			input:    time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			expected: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
		},
		{
			name:     "already aligned - :10",
			input:    time.Date(2025, 12, 1, 10, 10, 0, 0, location),
			expected: time.Date(2025, 12, 1, 10, 10, 0, 0, location),
		},
		{
			name:     "round down - :03 to :00",
			input:    time.Date(2025, 12, 1, 10, 3, 0, 0, location),
			expected: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
		},
		{
			name:     "round down - :27 to :20",
			input:    time.Date(2025, 12, 1, 10, 27, 45, 0, location),
			expected: time.Date(2025, 12, 1, 10, 20, 0, 0, location),
		},
		{
			name:     "round down - :55 to :50",
			input:    time.Date(2025, 12, 1, 11, 55, 0, 0, location),
			expected: time.Date(2025, 12, 1, 11, 50, 0, 0, location),
		},
		{
			name:     "seconds and nanoseconds removed",
			input:    time.Date(2025, 12, 1, 14, 33, 59, 999999999, location),
			expected: time.Date(2025, 12, 1, 14, 30, 0, 0, location),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := alignToBaseSlot(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasOverlap(t *testing.T) {
	location := time.UTC

	t.Run("no overlap - slot before booking", func(t *testing.T) {
		slotStart := time.Date(2025, 12, 1, 10, 0, 0, 0, location)
		slotEnd := time.Date(2025, 12, 1, 10, 30, 0, 0, location)

		booking := &domain.Booking{
			Status:        domain.BookingStatusConfirmed,
			SlotStartTime: time.Date(2025, 12, 1, 11, 0, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 11, 30, 0, 0, location),
		}

		overlap := hasOverlap(slotStart, slotEnd, []*domain.Booking{booking})
		assert.False(t, overlap, "No overlap when slot ends before booking starts")
	})

	t.Run("no overlap - slot after booking", func(t *testing.T) {
		slotStart := time.Date(2025, 12, 1, 11, 0, 0, 0, location)
		slotEnd := time.Date(2025, 12, 1, 11, 30, 0, 0, location)

		booking := &domain.Booking{
			Status:        domain.BookingStatusConfirmed,
			SlotStartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 10, 30, 0, 0, location),
		}

		overlap := hasOverlap(slotStart, slotEnd, []*domain.Booking{booking})
		assert.False(t, overlap, "No overlap when slot starts after booking ends")
	})

	t.Run("overlap - slot starts during booking", func(t *testing.T) {
		slotStart := time.Date(2025, 12, 1, 10, 20, 0, 0, location)
		slotEnd := time.Date(2025, 12, 1, 10, 50, 0, 0, location)

		booking := &domain.Booking{
			Status:        domain.BookingStatusConfirmed,
			SlotStartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 10, 30, 0, 0, location),
		}

		overlap := hasOverlap(slotStart, slotEnd, []*domain.Booking{booking})
		assert.True(t, overlap, "Overlap when slot starts during booking")
	})

	t.Run("overlap - slot contains booking", func(t *testing.T) {
		slotStart := time.Date(2025, 12, 1, 10, 0, 0, 0, location)
		slotEnd := time.Date(2025, 12, 1, 11, 0, 0, 0, location)

		booking := &domain.Booking{
			Status:        domain.BookingStatusConfirmed,
			SlotStartTime: time.Date(2025, 12, 1, 10, 20, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 10, 40, 0, 0, location),
		}

		overlap := hasOverlap(slotStart, slotEnd, []*domain.Booking{booking})
		assert.True(t, overlap, "Overlap when slot contains booking")
	})

	t.Run("overlap - exact match", func(t *testing.T) {
		slotStart := time.Date(2025, 12, 1, 10, 0, 0, 0, location)
		slotEnd := time.Date(2025, 12, 1, 10, 30, 0, 0, location)

		booking := &domain.Booking{
			Status:        domain.BookingStatusConfirmed,
			SlotStartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 10, 30, 0, 0, location),
		}

		overlap := hasOverlap(slotStart, slotEnd, []*domain.Booking{booking})
		assert.True(t, overlap, "Overlap when slot exactly matches booking")
	})

	t.Run("no overlap - cancelled booking ignored", func(t *testing.T) {
		slotStart := time.Date(2025, 12, 1, 10, 0, 0, 0, location)
		slotEnd := time.Date(2025, 12, 1, 10, 30, 0, 0, location)

		booking := &domain.Booking{
			Status:        domain.BookingStatusCancelled,
			SlotStartTime: time.Date(2025, 12, 1, 10, 0, 0, 0, location),
			SlotEndTime:   time.Date(2025, 12, 1, 10, 30, 0, 0, location),
		}

		overlap := hasOverlap(slotStart, slotEnd, []*domain.Booking{booking})
		assert.False(t, overlap, "Cancelled bookings should be ignored")
	})

	t.Run("multiple bookings - checks all", func(t *testing.T) {
		slotStart := time.Date(2025, 12, 1, 10, 30, 0, 0, location)
		slotEnd := time.Date(2025, 12, 1, 11, 0, 0, 0, location)

		bookings := []*domain.Booking{
			{
				Status:        domain.BookingStatusConfirmed,
				SlotStartTime: time.Date(2025, 12, 1, 9, 0, 0, 0, location),
				SlotEndTime:   time.Date(2025, 12, 1, 9, 30, 0, 0, location),
			},
			{
				Status:        domain.BookingStatusConfirmed,
				SlotStartTime: time.Date(2025, 12, 1, 10, 40, 0, 0, location),
				SlotEndTime:   time.Date(2025, 12, 1, 11, 10, 0, 0, location),
			},
		}

		overlap := hasOverlap(slotStart, slotEnd, bookings)
		assert.True(t, overlap, "Should detect overlap with second booking")
	})
}
