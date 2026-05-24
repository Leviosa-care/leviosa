package domain

import "errors"

// Building errors
var (
	ErrInvalidBuildingName    = errors.New("building name cannot be empty")
	ErrInvalidBuildingAddress = errors.New("building address cannot be empty")
	ErrInvalidBuildingCity    = errors.New("building city cannot be empty")
	ErrInvalidBuildingCountry = errors.New("building country cannot be empty")
	ErrInvalidBuildingID      = errors.New("invalid building ID")
)

// Room errors
var (
	ErrInvalidRoomName     = errors.New("room name cannot be empty")
	ErrInvalidRoomCapacity = errors.New("room capacity must be greater than 0")
	ErrInvalidRoomID       = errors.New("invalid room ID")
)

// Room allocation errors
var (
	ErrInvalidPartnerID                   = errors.New("invalid partner ID")
	ErrInvalidAllocationStartDate         = errors.New("allocation start date cannot be empty for dedicated allocations")
	ErrInvalidAllocationEndDate           = errors.New("allocation end date must be after start date")
	ErrCannotUpdateSharedAllocationPeriod = errors.New("cannot update time period for shared allocations")
)

// Availability errors
var (
	ErrInvalidTimeSlot                = errors.New("end time must be after start time")
	ErrCannotCreatePastAvailability   = errors.New("cannot create availability in the past")
	ErrInvalidAvailabilityCapacity    = errors.New("availability capacity must be greater than 0")
	ErrCannotUpdateToPastTime         = errors.New("cannot update availability to past time")
	ErrCannotUpdateBookedAvailability = errors.New("cannot update booked availability")
	ErrAvailabilityNotAvailable       = errors.New("availability is not available for booking")
	ErrInvalidRecurrenceType          = errors.New("invalid recurrence type")
	ErrInvalidRecurrenceInterval      = errors.New("recurrence interval must be greater than 0")
	ErrMissingWeeklyDays              = errors.New("weekly recurrence requires days of week")
	ErrInvalidWeekDay                 = errors.New("week day must be between 0 and 6")
	ErrInvalidAvailabilityID          = errors.New("invalid availability ID")
)

// Booking errors
var (
	ErrInvalidClientID              = errors.New("invalid client ID")
	ErrAmbiguousBookingIdentity     = errors.New("booking must have either a client ID or guest fields, not both")
	ErrInvalidBookingPrice          = errors.New("booking price cannot be negative")
	ErrCannotMarkRefundedAsPaid     = errors.New("cannot mark refunded payment as paid")
	ErrCannotRefundUnpaidBooking    = errors.New("cannot refund unpaid booking")
	ErrCannotCancelCompletedBooking = errors.New("cannot cancel completed booking")
	ErrBookingAlreadyCancelled      = errors.New("booking is already cancelled")
	ErrCannotCompleteBooking        = errors.New("cannot complete cancelled booking")
	ErrBookingAlreadyCompleted      = errors.New("booking is already completed")
	ErrCannotMarkCancelledAsNoShow  = errors.New("cannot mark cancelled booking as no-show")
	ErrCannotMarkCompletedAsNoShow  = errors.New("cannot mark completed booking as no-show")
)

