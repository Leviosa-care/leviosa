package booking

const (
	// BaseTimeSlotMinutes defines the minimum time increment for booking slots.
	// All booking slots must start on boundaries that are multiples of this value.
	// Example: With 10-minute slots, valid start times are :00, :10, :20, :30, :40, :50
	BaseTimeSlotMinutes = 10
)
