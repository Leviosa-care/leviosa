package booking

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

// reminderRepo is the subset of BookingRepository that the scheduler needs.
type reminderRepo interface {
	FindBookingsDueForReminder(ctx context.Context) ([]*domain.BookingEncx, error)
	MarkReminderSent(ctx context.Context, bookingID uuid.UUID) error
}

// ReminderScheduler is an in-process background worker that periodically
// queries upcoming bookings and sends reminder notifications.
type ReminderScheduler struct {
	repo         reminderRepo
	notification ports.BookingNotificationService
	crypto       encx.CryptoService

	interval       time.Duration // how often to tick (default: 15 minutes)
	reminderWindow time.Duration // how far ahead to look (default: 24 hours)

	mu     sync.Mutex
	stopCh chan struct{}
}

// NewReminderScheduler creates a new scheduler. Accepts the full
// ports.BookingRepository (which satisfies the reminderRepo interface).
// Use WithInterval and WithReminderWindow options to override defaults.
func NewReminderScheduler(
	bookingRepo ports.BookingRepository,
	notification ports.BookingNotificationService,
	crypto encx.CryptoService,
	opts ...SchedulerOption,
) *ReminderScheduler {
	s := &ReminderScheduler{
		repo:           bookingRepo,
		notification:   notification,
		crypto:         crypto,
		interval:       15 * time.Minute,
		reminderWindow: 24 * time.Hour,
		stopCh:         make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// NewReminderSchedulerForTest creates a scheduler with a narrower repository
// interface suitable for unit testing without the full BookingRepository.
// Accepts the same SchedulerOption set as NewReminderScheduler so tests can
// override interval and reminder window.
func NewReminderSchedulerForTest(
	repo reminderRepo,
	notification ports.BookingNotificationService,
	crypto encx.CryptoService,
	opts ...SchedulerOption,
) *ReminderScheduler {
	s := &ReminderScheduler{
		repo:           repo,
		notification:   notification,
		crypto:         crypto,
		interval:       15 * time.Minute,
		reminderWindow: 24 * time.Hour,
		stopCh:         make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// SchedulerOption configures optional ReminderScheduler settings.
type SchedulerOption func(*ReminderScheduler)

// WithInterval sets the tick interval. Must be > 0.
func WithInterval(d time.Duration) SchedulerOption {
	return func(s *ReminderScheduler) {
		if d > 0 {
			s.interval = d
		}
	}
}

// WithReminderWindow sets how far ahead the scheduler looks for upcoming bookings.
// Must be > 0.
func WithReminderWindow(d time.Duration) SchedulerOption {
	return func(s *ReminderScheduler) {
		if d > 0 {
			s.reminderWindow = d
		}
	}
}

// Start begins the scheduler loop. It blocks until the context is cancelled
// or Stop is called. Start is safe to call from a goroutine.
//
// The scheduler is non-blocking: if a tick takes longer than the interval,
// the next scheduled tick is skipped rather than queued. This is achieved by
// resetting the timer after each tick completes, ensuring a full interval
// elapses before the next tick regardless of how long the previous one took.
func (s *ReminderScheduler) Start(ctx context.Context) {
	slog.InfoContext(ctx, "booking reminder scheduler started",
		"interval", s.interval,
		"reminder_window", s.reminderWindow,
	)

	timer := time.NewTimer(0) // fire immediately for first tick
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "booking reminder scheduler stopped")
			return
		case <-s.stopCh:
			slog.InfoContext(ctx, "booking reminder scheduler stopped via Stop")
			return
		case <-timer.C:
			s.tick(ctx)
			// Reset with a full interval. If the tick took longer than the
			// interval, the next tick is effectively skipped — the timer fires
			// one full interval from now rather than immediately.
			timer.Reset(s.interval)
		}
	}
}

// Stop signals the scheduler to stop. It is safe to call from any goroutine.
func (s *ReminderScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-s.stopCh:
		// already closed
	default:
		close(s.stopCh)
	}
}

// TickOnce runs a single scheduler tick. Exposed for integration testing.
func (s *ReminderScheduler) TickOnce(ctx context.Context) {
	s.tick(ctx)
}

// tick runs one scheduler cycle: find eligible bookings, send reminders, mark sent.
func (s *ReminderScheduler) tick(ctx context.Context) {
	if s.repo == nil {
		return
	}
	now := time.Now()
	windowEnd := now.Add(s.reminderWindow)

	bookingsEncx, err := s.repo.FindBookingsDueForReminder(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "reminder scheduler: failed to query bookings", "error", err)
		return
	}

	for _, bx := range bookingsEncx {
		booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bx)
		if err != nil {
			slog.ErrorContext(ctx, "reminder scheduler: failed to decrypt booking",
				"booking_id", bx.ID,
				"error", err,
			)
			continue
		}

		// Past bookings: mark as reminded without sending so they don't
		// reappear on every tick (slot_start_time is encrypted so the DB
		// cannot filter them out; we do it here after decryption).
		if booking.SlotStartTime.Before(now) {
			if err := s.repo.MarkReminderSent(ctx, booking.ID); err != nil {
				slog.ErrorContext(ctx, "reminder scheduler: failed to mark past booking as reminded",
					"booking_id", booking.ID,
					"error", err,
				)
			}
			continue
		}

		// Bookings beyond the reminder window: skip without marking.
		if booking.SlotStartTime.After(windowEnd) {
			continue
		}

		// Send reminder (fire-and-forget from the adapter's perspective).
		data := s.buildNotificationData(booking)
		if err := s.notification.SendBookingReminder(ctx, data); err != nil {
			slog.ErrorContext(ctx, "reminder scheduler: failed to send reminder",
				"booking_id", booking.ID,
				"error", err,
			)
		}

		// Mark as reminded regardless of notification outcome to prevent retry storms.
		if err := s.repo.MarkReminderSent(ctx, booking.ID); err != nil {
			slog.ErrorContext(ctx, "reminder scheduler: failed to mark reminder sent",
				"booking_id", booking.ID,
				"error", err,
			)
		}
	}

	elapsed := time.Since(now)
	slog.DebugContext(ctx, "reminder scheduler tick completed",
		"candidates", len(bookingsEncx),
		"elapsed", elapsed,
	)
}

// buildNotificationData creates BookingNotificationData from a decrypted booking.
func (s *ReminderScheduler) buildNotificationData(booking *domain.Booking) ports.BookingNotificationData {
	data := ports.BookingNotificationData{
		BookingID:       booking.ID,
		PartnerID:       booking.PartnerID,
		RoomID:          booking.RoomID,
		ProductID:       booking.ProductID,
		SlotStartTime:   booking.SlotStartTime,
		SlotEndTime:     booking.SlotEndTime,
		TotalPriceCents: booking.TotalPriceCents,
		Currency:        booking.Currency,
	}

	if booking.ClientID != nil {
		data.ClientID = *booking.ClientID
	}

	if booking.IsGuestBooking() {
		data.IsGuestBooking = true
		data.GuestEmail = booking.GuestEmail
		data.GuestPhone = booking.GuestPhone
		data.ClientName = booking.GuestDisplayName()
		data.ClientEmail = booking.GuestEmail
		data.ClientPhone = booking.GuestPhone
	}

	data.Token = booking.Token

	return data
}
