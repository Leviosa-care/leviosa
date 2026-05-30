package booking

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/hengadev/encx"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Identity crypto mock — passes data through without real encryption.
// This lets us unit-test the scheduler tick without a Vault container.
// ---------------------------------------------------------------------------

type identityCrypto struct {
	pepper []byte
}

var _ encx.CryptoService = (*identityCrypto)(nil)

func (c *identityCrypto) GetPepper() []byte                                        { return c.pepper }
func (c *identityCrypto) GetArgon2Params() *encx.Argon2Params                       { return nil }
func (c *identityCrypto) GetAlias() string                                          { return "test" }
func (c *identityCrypto) GenerateDEK() ([]byte, error)                              { return []byte("test-dek-32-bytes-padding-here!"), nil }
func (c *identityCrypto) EncryptData(_ context.Context, p []byte, _ []byte) ([]byte, error) {
	return p, nil
}
func (c *identityCrypto) DecryptData(_ context.Context, ct []byte, _ []byte) ([]byte, error) {
	return ct, nil
}
func (c *identityCrypto) EncryptDEK(_ context.Context, dek []byte) ([]byte, error) {
	return dek, nil
}
func (c *identityCrypto) DecryptDEKWithVersion(_ context.Context, ct []byte, _ int) ([]byte, error) {
	return ct, nil
}
func (c *identityCrypto) RotateKEK(_ context.Context) error                        { return nil }
func (c *identityCrypto) HashBasic(_ context.Context, v []byte) string              { return string(v) }
func (c *identityCrypto) HashSecure(_ context.Context, v []byte) (string, error)    { return string(v), nil }
func (c *identityCrypto) CompareSecureHashAndValue(_ context.Context, _ any, _ string) (bool, error) {
	return true, nil
}
func (c *identityCrypto) CompareBasicHashAndValue(_ context.Context, _ any, _ string) (bool, error) {
	return true, nil
}
func (c *identityCrypto) EncryptStream(_ context.Context, r io.Reader, w io.Writer, _ []byte) error {
	_, err := io.Copy(w, r)
	return err
}
func (c *identityCrypto) DecryptStream(_ context.Context, r io.Reader, w io.Writer, _ []byte) error {
	_, err := io.Copy(w, r)
	return err
}
func (c *identityCrypto) GetCurrentKEKVersion(_ context.Context, _ string) (int, error) {
	return 1, nil
}
func (c *identityCrypto) GetKMSKeyIDForVersion(_ context.Context, _ string, _ int) (string, error) {
	return "test-key", nil
}

// ---------------------------------------------------------------------------
// Mock repository — only the two methods the scheduler needs.
// We embed the full interface via a wrapper.
// ---------------------------------------------------------------------------

type mockSchedulerRepo struct {
	mu       sync.Mutex
	bookings []*domain.BookingEncx
	markCalled []uuid.UUID
}

func (m *mockSchedulerRepo) FindBookingsDueForReminder(_ context.Context) ([]*domain.BookingEncx, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]*domain.BookingEncx{}, m.bookings...), nil
}

func (m *mockSchedulerRepo) MarkReminderSent(_ context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.markCalled = append(m.markCalled, id)
	return nil
}

// ---------------------------------------------------------------------------
// Mock notification
// ---------------------------------------------------------------------------

type mockReminderNotification struct {
	mu    sync.Mutex
	calls []ports.BookingNotificationData
	err   error // if set, SendBookingReminder returns this error
}

func (m *mockReminderNotification) SendBookingConfirmation(_ context.Context, _ ports.BookingNotificationData) error {
	return nil
}
func (m *mockReminderNotification) SendBookingCancellation(_ context.Context, _ ports.BookingNotificationData) error {
	return nil
}
func (m *mockReminderNotification) SendBookingReminder(_ context.Context, data ports.BookingNotificationData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, data)
	return m.err
}
func (m *mockReminderNotification) SendPaymentConfirmation(_ context.Context, _ ports.BookingNotificationData) error {
	return nil
}
func (m *mockReminderNotification) SendPaymentFailed(_ context.Context, _ ports.BookingNotificationData) error {
	return nil
}

var _ ports.BookingNotificationService = (*mockReminderNotification)(nil)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func makeEncryptedBooking(t *testing.T, crypto encx.CryptoService, slotStart time.Time) *domain.BookingEncx {
	t.Helper()
	clientID := uuid.New()
	b := &domain.Booking{
		ID:              uuid.New(),
		AvailabilityID:  uuid.New(),
		ClientID:        &clientID,
		PartnerID:       uuid.New(),
		RoomID:          uuid.New(),
		ProductID:       uuid.New(),
		SlotStartTime:   slotStart,
		SlotEndTime:     slotStart.Add(time.Hour),
		TotalPriceCents: 5000,
		Currency:        "EUR",
		PaymentStatus:   domain.PaymentStatusPaid,
		Status:          domain.BookingStatusConfirmed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	encx, err := domain.ProcessBookingEncx(context.Background(), crypto, b)
	require.NoError(t, err)
	return encx
}

func makeGuestEncryptedBooking(t *testing.T, crypto encx.CryptoService, slotStart time.Time) *domain.BookingEncx {
	t.Helper()
	b, err := domain.NewBooking(
		uuid.New(), nil, uuid.New(), uuid.New(),
		5000, "EUR",
		"Jean", "Dupont", "jean@example.com", "+33612345678",
	)
	require.NoError(t, err)
	b.SlotStartTime = slotStart
	b.SlotEndTime = slotStart.Add(time.Hour)
	b.ProductID = uuid.New()
	b.Status = domain.BookingStatusConfirmed
	b.PaymentStatus = domain.PaymentStatusPaid

	encx, err := domain.ProcessBookingEncx(context.Background(), crypto, b)
	require.NoError(t, err)
	return encx
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestReminderSchedulerTick_SendsRemindersForEligibleBookings(t *testing.T) {
	ctx := context.Background()
	crypto := &identityCrypto{pepper: make([]byte, 32)}

	now := time.Now()
	inWindow := now.Add(12 * time.Hour)     // within 24h window
	outsideWindow := now.Add(48 * time.Hour) // outside 24h window
	inPast := now.Add(-1 * time.Hour)        // past booking

	encxInWindow := makeEncryptedBooking(t, crypto, inWindow)
	encxOutside := makeEncryptedBooking(t, crypto, outsideWindow)
	encxPast := makeEncryptedBooking(t, crypto, inPast)

	repo := &mockSchedulerRepo{
		bookings: []*domain.BookingEncx{encxInWindow, encxOutside, encxPast},
	}
	notif := &mockReminderNotification{}

	scheduler := NewReminderSchedulerForTest(repo, notif, crypto, 24*time.Hour)
	scheduler.tick(ctx)

	notif.mu.Lock()
	calls := notif.calls
	notif.mu.Unlock()

	repo.mu.Lock()
	marks := repo.markCalled
	repo.mu.Unlock()

	// Only the in-window booking should trigger a reminder notification
	assert.Len(t, calls, 1)
	assert.Equal(t, encxInWindow.ID, calls[0].BookingID)

	// All three bookings pass FindBookingsDueForReminder, but only the in-window
	// one gets a notification. However, MarkReminderSent is called only for
	// bookings that were actually sent reminders (in-window).
	assert.Len(t, marks, 1)
	assert.Equal(t, encxInWindow.ID, marks[0])
}

func TestReminderSchedulerTick_MarkSentEvenOnNotificationError(t *testing.T) {
	ctx := context.Background()
	crypto := &identityCrypto{pepper: make([]byte, 32)}

	now := time.Now()
	inWindow := now.Add(12 * time.Hour)

	encxBooking := makeEncryptedBooking(t, crypto, inWindow)

	repo := &mockSchedulerRepo{
		bookings: []*domain.BookingEncx{encxBooking},
	}
	notif := &mockReminderNotification{
		err: assert.AnError,
	}

	scheduler := NewReminderSchedulerForTest(repo, notif, crypto, 24*time.Hour)
	scheduler.tick(ctx)

	notif.mu.Lock()
	calls := notif.calls
	notif.mu.Unlock()

	repo.mu.Lock()
	marks := repo.markCalled
	repo.mu.Unlock()

	// Notification was attempted
	assert.Len(t, calls, 1)

	// MarkReminderSent must still be called even though notification failed
	assert.Len(t, marks, 1)
	assert.Equal(t, encxBooking.ID, marks[0])
}

func TestReminderSchedulerTick_MultipleBookingsInWindow(t *testing.T) {
	ctx := context.Background()
	crypto := &identityCrypto{pepper: make([]byte, 32)}

	now := time.Now()
	b1Time := now.Add(6 * time.Hour)
	b2Time := now.Add(18 * time.Hour)

	encx1 := makeEncryptedBooking(t, crypto, b1Time)
	encx2 := makeEncryptedBooking(t, crypto, b2Time)

	repo := &mockSchedulerRepo{
		bookings: []*domain.BookingEncx{encx1, encx2},
	}
	notif := &mockReminderNotification{}

	scheduler := NewReminderSchedulerForTest(repo, notif, crypto, 24*time.Hour)
	scheduler.tick(ctx)

	notif.mu.Lock()
	calls := notif.calls
	notif.mu.Unlock()

	repo.mu.Lock()
	marks := repo.markCalled
	repo.mu.Unlock()

	assert.Len(t, calls, 2)
	assert.Len(t, marks, 2)
}

func TestReminderSchedulerTick_GuestBooking(t *testing.T) {
	ctx := context.Background()
	crypto := &identityCrypto{pepper: make([]byte, 32)}

	now := time.Now()
	inWindow := now.Add(12 * time.Hour)

	encxGuest := makeGuestEncryptedBooking(t, crypto, inWindow)

	repo := &mockSchedulerRepo{
		bookings: []*domain.BookingEncx{encxGuest},
	}
	notif := &mockReminderNotification{}

	scheduler := NewReminderSchedulerForTest(repo, notif, crypto, 24*time.Hour)
	scheduler.tick(ctx)

	notif.mu.Lock()
	calls := notif.calls
	notif.mu.Unlock()

	require.Len(t, calls, 1)
	assert.True(t, calls[0].IsGuestBooking)
	assert.Equal(t, "jean@example.com", calls[0].GuestEmail)
}

func TestReminderSchedulerBuildNotificationData(t *testing.T) {
	clientID := uuid.New()
	now := time.Now()

	booking := &domain.Booking{
		ID:              uuid.New(),
		ClientID:        &clientID,
		PartnerID:       uuid.New(),
		RoomID:          uuid.New(),
		ProductID:       uuid.New(),
		SlotStartTime:   now.Add(12 * time.Hour),
		SlotEndTime:     now.Add(13 * time.Hour),
		TotalPriceCents: 5000,
		Currency:        "EUR",
	}

	scheduler := &ReminderScheduler{}
	data := scheduler.buildNotificationData(booking)

	assert.Equal(t, booking.ID, data.BookingID)
	assert.Equal(t, clientID, data.ClientID)
	assert.Equal(t, booking.PartnerID, data.PartnerID)
	assert.False(t, data.IsGuestBooking)
}

func TestReminderSchedulerBuildNotificationData_Guest(t *testing.T) {
	now := time.Now()

	booking := &domain.Booking{
		ID:              uuid.New(),
		ClientID:        nil,
		PartnerID:       uuid.New(),
		RoomID:          uuid.New(),
		ProductID:       uuid.New(),
		SlotStartTime:   now.Add(12 * time.Hour),
		SlotEndTime:     now.Add(13 * time.Hour),
		TotalPriceCents: 5000,
		Currency:        "EUR",
		GuestFirstName:  "Jean",
		GuestLastName:   "Dupont",
		GuestEmail:      "jean@example.com",
		GuestPhone:      "+33612345678",
	}

	scheduler := &ReminderScheduler{}
	data := scheduler.buildNotificationData(booking)

	assert.True(t, data.IsGuestBooking)
	assert.Equal(t, "jean@example.com", data.ClientEmail)
	assert.Equal(t, "+33612345678", data.ClientPhone)
	assert.Equal(t, "Jean Dupont", data.ClientName)
}

func TestReminderSchedulerOptions(t *testing.T) {
	s := NewReminderScheduler(nil, nil, nil,
		WithInterval(5*time.Minute),
		WithReminderWindow(48*time.Hour),
	)
	assert.Equal(t, 5*time.Minute, s.interval)
	assert.Equal(t, 48*time.Hour, s.reminderWindow)
}

func TestReminderSchedulerOptions_ZeroIgnored(t *testing.T) {
	s := NewReminderScheduler(nil, nil, nil,
		WithInterval(0),
		WithReminderWindow(0),
	)
	assert.Equal(t, 15*time.Minute, s.interval)
	assert.Equal(t, 24*time.Hour, s.reminderWindow)
}

func TestReminderSchedulerStop(t *testing.T) {
	s := NewReminderScheduler(nil, nil, nil)

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		s.Start(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not stop")
	}
}
