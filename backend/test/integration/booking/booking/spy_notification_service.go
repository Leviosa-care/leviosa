package booking

import (
	"context"
	"sync"

	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

// SpyNotificationService records all notification calls for test assertions.
// It implements bookingPorts.BookingNotificationService.
type SpyNotificationService struct {
	mu                    sync.Mutex
	BookingConfirmations  []bookingPorts.BookingNotificationData
	PaymentConfirmations  []bookingPorts.BookingNotificationData
	PaymentFailed         []bookingPorts.BookingNotificationData
	BookingCancellations  []bookingPorts.BookingNotificationData
	BookingReminders      []bookingPorts.BookingNotificationData
}

func NewSpyNotificationService() *SpyNotificationService {
	return &SpyNotificationService{}
}

func (s *SpyNotificationService) SendBookingConfirmation(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BookingConfirmations = append(s.BookingConfirmations, data)
	return nil
}

func (s *SpyNotificationService) SendPaymentConfirmation(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.PaymentConfirmations = append(s.PaymentConfirmations, data)
	return nil
}

func (s *SpyNotificationService) SendPaymentFailed(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.PaymentFailed = append(s.PaymentFailed, data)
	return nil
}

func (s *SpyNotificationService) SendBookingCancellation(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BookingCancellations = append(s.BookingCancellations, data)
	return nil
}

func (s *SpyNotificationService) SendBookingReminder(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BookingReminders = append(s.BookingReminders, data)
	return nil
}
