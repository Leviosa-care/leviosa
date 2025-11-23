package availabilityHelpers

import (
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func NewTestAvailability(t *testing.T) *domain.Availability {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)    // Tomorrow
	endTime := startTime.Add(2 * time.Hour) // 2 hours later
	priceCents := 15000                     // $150.00

	return &domain.Availability{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: uuid.New(),

		StartTime: startTime,
		EndTime:   endTime,

		ServiceType: "Consultation",
		PriceCents:  &priceCents,
		MaxCapacity: 1,

		Notes:       "Regular consultation slot",
		IsRecurring: false,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewTestAvailabilityWithParams creates a test availability with custom parameters
func NewTestAvailabilityWithParams(t *testing.T, userID, roomID uuid.UUID, startTime, endTime time.Time, serviceType string, priceCents *int, maxCapacity int, status domain.AvailabilityStatus) *domain.Availability {
	now := time.Now()
	return &domain.Availability{
		ID:     uuid.New(),
		UserID: userID,
		RoomID: roomID,

		StartTime: startTime,
		EndTime:   endTime,

		ServiceType: serviceType,
		PriceCents:  priceCents,
		MaxCapacity: maxCapacity,

		Notes:       "Test availability",
		IsRecurring: false,
		Status:      status,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewTestRecurringAvailability creates a test recurring availability
func NewTestRecurringAvailability(t *testing.T) *domain.Availability {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	until := now.Add(30 * 24 * time.Hour) // 30 days from now
	priceCents := 10000                   // $100.00

	return &domain.Availability{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: uuid.New(),

		StartTime: startTime,
		EndTime:   endTime,

		ServiceType: "Therapy Session",
		PriceCents:  &priceCents,
		MaxCapacity: 1,

		Notes:       "Weekly therapy session",
		IsRecurring: true,
		RecurrencePattern: &domain.RecurrencePattern{
			Type:       "weekly",
			Interval:   1,
			Until:      &until,
			DaysOfWeek: []int{1, 3, 5}, // Monday, Wednesday, Friday
		},
		Status: domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewTestAvailabilityForuser creates a test availability for a specific user
func NewTestAvailabilityForuser(t *testing.T, userID uuid.UUID) *domain.Availability {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	priceCents := 8000 // $80.00

	return &domain.Availability{
		ID:     uuid.New(),
		UserID: userID,
		RoomID: uuid.New(),

		StartTime: startTime,
		EndTime:   endTime,

		ServiceType: "Check-up",
		PriceCents:  &priceCents,
		MaxCapacity: 1,

		Notes:       "Regular health check-up",
		IsRecurring: false,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewTestAvailabilityForRoom creates a test availability for a specific room
func NewTestAvailabilityForRoom(t *testing.T, roomID uuid.UUID) *domain.Availability {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	priceCents := 20000 // $200.00

	return &domain.Availability{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: roomID,

		StartTime: startTime,
		EndTime:   endTime,

		ServiceType: "Specialized Treatment",
		PriceCents:  &priceCents,
		MaxCapacity: 2,

		Notes:       "Specialized treatment session",
		IsRecurring: false,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewBookedTestAvailability creates a test availability that is already booked
func NewBookedTestAvailability(t *testing.T) *domain.Availability {
	avail := NewTestAvailability(t)
	avail.Status = domain.AvailabilityStatusBooked
	avail.UpdatedAt = time.Now()
	return avail
}

// NewPastTestAvailability creates a test availability in the past
func NewPastTestAvailability(t *testing.T) *domain.Availability {
	now := time.Now()
	startTime := now.Add(-24 * time.Hour) // Yesterday
	endTime := startTime.Add(1 * time.Hour)

	return &domain.Availability{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: uuid.New(),

		StartTime: startTime,
		EndTime:   endTime,

		ServiceType: "Past Session",
		PriceCents:  nil, // Free
		MaxCapacity: 1,

		Notes:       "A session that already happened",
		IsRecurring: false,
		Status:      domain.AvailabilityStatusCancelled,

		CreatedAt: startTime.Add(-2 * time.Hour),
		UpdatedAt: endTime,
	}
}

func NewTestAvailabilityEncx(t *testing.T) *domain.AvailabilityEncx {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	priceCents := 15000

	return &domain.AvailabilityEncx{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: uuid.New(),

		StartTime: startTime,
		EndTime:   endTime,

		PriceCents:  &priceCents,
		MaxCapacity: 1,
		IsRecurring: false,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,

		ServiceTypeEncrypted:       []byte("encrypted_consultation"),
		NotesEncrypted:             []byte("encrypted_regular_consultation_slot"),
		RecurrencePattern: nil,

		DEKEncrypted: []byte("mock_dek_data"),
		KeyVersion:   1,
		Metadata:     encx.EncryptionMetadata{},
	}
}

// NewTestAvailabilityEncxWithPartner creates an encrypted test availability with specific user ID
func NewTestAvailabilityEncxWithPartner(t *testing.T, userID uuid.UUID) *domain.AvailabilityEncx {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	priceCents := 8000

	return &domain.AvailabilityEncx{
		ID:     uuid.New(),
		UserID: userID,
		RoomID: uuid.New(),

		StartTime: startTime,
		EndTime:   endTime,

		PriceCents:  &priceCents,
		MaxCapacity: 1,
		IsRecurring: false,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,

		ServiceTypeEncrypted:       []byte("encrypted_check_up"),
		NotesEncrypted:             []byte("encrypted_regular_health_check_up"),
		RecurrencePattern: nil,

		DEKEncrypted: []byte("mock_dek_data"),
		KeyVersion:   1,
		Metadata:     encx.EncryptionMetadata{},
	}
}

// NewTestAvailabilityEncxWithRoom creates an encrypted test availability with specific room ID
func NewTestAvailabilityEncxWithRoom(t *testing.T, roomID uuid.UUID) *domain.AvailabilityEncx {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	priceCents := 20000

	return &domain.AvailabilityEncx{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: roomID,

		StartTime: startTime,
		EndTime:   endTime,

		PriceCents:  &priceCents,
		MaxCapacity: 2,
		IsRecurring: false,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,

		ServiceTypeEncrypted:       []byte("encrypted_specialized_treatment"),
		NotesEncrypted:             []byte("encrypted_specialized_treatment_session"),
		RecurrencePattern: nil,

		DEKEncrypted: []byte("mock_dek_data"),
		KeyVersion:   1,
		Metadata:     encx.EncryptionMetadata{},
	}
}

// NewTestRecurringAvailabilityEncx creates an encrypted test recurring availability
func NewTestRecurringAvailabilityEncx(t *testing.T) *domain.AvailabilityEncx {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	priceCents := 10000

	return &domain.AvailabilityEncx{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: uuid.New(),

		StartTime: startTime,
		EndTime:   endTime,

		PriceCents:  &priceCents,
		MaxCapacity: 1,
		IsRecurring: true,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,

		ServiceTypeEncrypted:       []byte("encrypted_therapy_session"),
		NotesEncrypted:             []byte("encrypted_weekly_therapy_session"),
		RecurrencePattern: &domain.RecurrencePattern{
			Type:       "weekly",
			Interval:   1,
			DaysOfWeek: []int{1, 3, 5},
		},

		DEKEncrypted: []byte("mock_dek_data"),
		KeyVersion:   1,
		Metadata:     encx.EncryptionMetadata{},
	}
}

// NewTestAvailabilityEncxWithPartnerAndRoom creates an encrypted test availability with specific user and room IDs
func NewTestAvailabilityEncxWithPartnerAndRoom(t *testing.T, userID, roomID uuid.UUID) *domain.AvailabilityEncx {
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	priceCents := 15000 // $150.00

	return &domain.AvailabilityEncx{
		ID:     uuid.New(),
		UserID: userID,
		RoomID: roomID,

		StartTime: startTime,
		EndTime:   endTime,

		PriceCents:  &priceCents,
		MaxCapacity: 1,
		IsRecurring: false,
		Status:      domain.AvailabilityStatusAvailable,

		CreatedAt: now,
		UpdatedAt: now,

		ServiceTypeEncrypted:       []byte("encrypted_consultation"),
		NotesEncrypted:             []byte("encrypted_regular_consultation_slot"),
		RecurrencePattern: nil,

		DEKEncrypted: []byte("mock_dek_data"),
		KeyVersion:   1,
		Metadata:     encx.EncryptionMetadata{},
	}
}
