package booking

import (
	"time"

	bookingContracts "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/booking"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
)

// AvailableSlot represents a bookable time slot
type AvailableSlot struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// CalculateAvailableSlots generates all available slots for a product within an availability.
// It considers existing bookings to exclude overlapping time slots.
//
// Parameters:
//   - availability: The availability block to calculate slots within
//   - productDuration: Duration of the product in minutes
//   - existingBookings: All bookings for this availability (to detect overlaps)
//
// Returns:
//   - Slice of AvailableSlot aligned to 10-minute boundaries
func CalculateAvailableSlots(
	availability *domain.Availability,
	productDuration int, // minutes
	existingBookings []*domain.Booking,
) []AvailableSlot {
	slots := []AvailableSlot{}
	duration := time.Duration(productDuration) * time.Minute

	// Align start to base slot (round down to nearest 10-minute boundary)
	current := alignToBaseSlot(availability.StartTime)

	// If alignment pushed us before availability start, move forward one base slot
	if current.Before(availability.StartTime) {
		current = current.Add(time.Duration(bookingContracts.BaseTimeSlotMinutes) * time.Minute)
	}

	// Generate slots starting from aligned time
	for {
		slotEnd := current.Add(duration)

		// Stop if slot would extend beyond availability end time
		if slotEnd.After(availability.EndTime) {
			break
		}

		// Only add slot if it doesn't overlap with existing bookings
		if !hasOverlap(current, slotEnd, existingBookings) {
			slots = append(slots, AvailableSlot{
				StartTime: current,
				EndTime:   slotEnd,
				Duration:  duration,
			})
		}

		// Move to next base slot boundary (10 minutes)
		current = current.Add(time.Duration(bookingContracts.BaseTimeSlotMinutes) * time.Minute)
	}

	return slots
}

// alignToBaseSlot rounds a time down to the nearest base time slot boundary.
// For a 10-minute base slot, times are aligned to :00, :10, :20, :30, :40, :50.
//
// Example:
//   - 10:03 → 10:00
//   - 10:27 → 10:20
//   - 11:55 → 11:50
func alignToBaseSlot(t time.Time) time.Time {
	minutes := t.Minute()
	aligned := (minutes / bookingContracts.BaseTimeSlotMinutes) * bookingContracts.BaseTimeSlotMinutes
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), aligned, 0, 0, t.Location())
}

// hasOverlap checks if a proposed time slot overlaps with any existing bookings.
// Cancelled bookings are ignored.
//
// Overlap detection uses the interval intersection formula:
//   - Two intervals [A_start, A_end) and [B_start, B_end) overlap if:
//     A_start < B_end AND B_start < A_end
//
// Parameters:
//   - slotStart: Proposed slot start time
//   - slotEnd: Proposed slot end time
//   - bookings: Existing bookings to check against
//
// Returns:
//   - true if slot overlaps with any non-cancelled booking
//   - false if slot is available (no overlaps)
func hasOverlap(slotStart, slotEnd time.Time, bookings []*domain.Booking) bool {
	for _, booking := range bookings {
		// Skip cancelled bookings - they don't block slots
		if booking.Status == domain.BookingStatusCancelled {
			continue
		}

		// Overlap check: (start1 < end2) AND (start2 < end1)
		if slotStart.Before(booking.SlotEndTime) && booking.SlotStartTime.Before(slotEnd) {
			return true
		}
	}
	return false
}
