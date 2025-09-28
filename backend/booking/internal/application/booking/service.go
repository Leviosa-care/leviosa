package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/booking/internal/adapters/stripe"
	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

type BookingService struct {
	bookingRepo      ports.BookingRepository
	availabilityRepo ports.AvailabilityRepository
	paymentService   ports.PaymentService
}

// New creates a new instance of the booking service
func New(
	bookingRepo ports.BookingRepository,
	availabilityRepo ports.AvailabilityRepository,
	paymentService ports.PaymentService,
) ports.BookingService {
	return &BookingService{
		bookingRepo:      bookingRepo,
		availabilityRepo: availabilityRepo,
		paymentService:   paymentService,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, availabilityID, clientID uuid.UUID, clientNotes string) (*domain.Booking, error) {
	// Get availability and verify it's bookable
	availability, err := s.availabilityRepo.GetByID(ctx, availabilityID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("availability not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("get availability: %w", err)
	}

	// Check if availability is still available for booking
	if !availability.IsAvailable() {
		return nil, fmt.Errorf("availability is no longer available for booking")
	}

	// Check if there's already a booking for this availability
	existingBooking, err := s.bookingRepo.GetByAvailabilityID(ctx, availabilityID)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("check existing booking: %w", err)
	}

	if existingBooking != nil {
		return nil, fmt.Errorf("availability is already booked")
	}

	// Calculate price (use availability price or default to 0 for now)
	totalPriceCents := 0
	if availability.PriceCents != nil {
		totalPriceCents = *availability.PriceCents
	}

	// Create booking entity
	booking, err := domain.NewBooking(
		availabilityID,
		clientID,
		availability.PartnerID,
		availability.RoomID,
		totalPriceCents,
		"EUR", // Default currency
	)
	if err != nil {
		return nil, fmt.Errorf("create booking entity: %w", err)
	}

	// Set client notes if provided
	if clientNotes != "" {
		booking.SetClientNotes(clientNotes)
	}

	// Create payment intent if booking has a price
	if totalPriceCents > 0 {
		description := fmt.Sprintf("Booking for availability %s", availabilityID.String())
		metadata := map[string]string{
			"booking_id":      booking.ID.String(),
			"availability_id": availabilityID.String(),
			"client_id":       clientID.String(),
			"partner_id":      availability.PartnerID.String(),
		}

		paymentIntentID, _, err := s.paymentService.CreatePaymentIntent(
			ctx,
			totalPriceCents,
			booking.Currency,
			description,
			metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("create payment intent: %w", err)
		}

		// Update booking with payment intent ID
		booking.SetPaymentIntentID(paymentIntentID)
	}

	// Mark availability as booked (atomic operation would be ideal here)
	if err := availability.MarkAsBooked(); err != nil {
		return nil, fmt.Errorf("mark availability as booked: %w", err)
	}

	// Update availability in repository
	if err := s.availabilityRepo.Update(ctx, availability); err != nil {
		return nil, fmt.Errorf("update availability status: %w", err)
	}

	// Persist booking to repository
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		// Rollback availability status on failure
		availability.MarkAsAvailable()
		s.availabilityRepo.Update(ctx, availability)
		return nil, fmt.Errorf("create booking: %w", err)
	}

	return booking, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking: %w", err)
	}

	return booking, nil
}

func (s *BookingService) UpdateBookingNotes(ctx context.Context, id uuid.UUID, clientNotes, partnerNotes string) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for notes update: %w", err)
	}

	// Update notes
	if clientNotes != "" {
		booking.SetClientNotes(clientNotes)
	}
	if partnerNotes != "" {
		booking.SetPartnerNotes(partnerNotes)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update booking notes: %w", err)
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id uuid.UUID, reason string) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for cancellation: %w", err)
	}

	// Check if booking can be cancelled
	if !booking.IsCancellable() {
		return nil, fmt.Errorf("booking cannot be cancelled")
	}

	// Cancel booking
	if err := booking.Cancel(reason); err != nil {
		return nil, fmt.Errorf("cancel booking: %w", err)
	}

	// Mark associated availability as available again
	availability, err := s.availabilityRepo.GetByID(ctx, booking.AvailabilityID)
	if err == nil {
		availability.MarkAsAvailable()
		s.availabilityRepo.Update(ctx, availability) // Best effort, don't fail booking cancellation
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update cancelled booking: %w", err)
	}

	return booking, nil
}

func (s *BookingService) CompleteBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for completion: %w", err)
	}

	// Complete booking
	if err := booking.Complete(); err != nil {
		return nil, fmt.Errorf("complete booking: %w", err)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update completed booking: %w", err)
	}

	return booking, nil
}

func (s *BookingService) MarkNoShow(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for no-show: %w", err)
	}

	// Mark as no-show
	if err := booking.MarkNoShow(); err != nil {
		return nil, fmt.Errorf("mark booking as no-show: %w", err)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update no-show booking: %w", err)
	}

	return booking, nil
}

func (s *BookingService) ProcessPayment(ctx context.Context, id uuid.UUID, paymentIntentID string) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for payment: %w", err)
	}

	// Verify payment intent status with Stripe
	paymentInfo, err := s.paymentService.RetrievePaymentIntent(ctx, paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("retrieve payment intent: %w", err)
	}

	// Set payment intent ID if not already set
	if booking.PaymentIntentID == nil || *booking.PaymentIntentID != paymentIntentID {
		booking.SetPaymentIntentID(paymentIntentID)
	}

	// Update payment status based on Stripe payment intent status
	switch paymentInfo.Status {
	case stripe.PaymentIntentStatusSucceeded:
		if err := booking.MarkPaymentPaid(); err != nil {
			return nil, fmt.Errorf("mark payment as paid: %w", err)
		}
	case stripe.PaymentIntentStatusRequiresPaymentMethod,
		 stripe.PaymentIntentStatusRequiresConfirmation,
		 stripe.PaymentIntentStatusRequiresAction,
		 stripe.PaymentIntentStatusProcessing:
		// Payment is still pending, no action needed
	case stripe.PaymentIntentStatusCanceled:
		if err := booking.MarkPaymentFailed(); err != nil {
			return nil, fmt.Errorf("mark payment as failed: %w", err)
		}
	case stripe.PaymentIntentStatusPaymentFailed:
		if err := booking.MarkPaymentFailed(); err != nil {
			return nil, fmt.Errorf("mark payment as failed: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown payment status: %s", paymentInfo.Status)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update booking payment: %w", err)
	}

	return booking, nil
}

func (s *BookingService) RefundBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for refund: %w", err)
	}

	// Validate booking can be refunded
	if booking.PaymentStatus != domain.PaymentStatusPaid {
		return nil, fmt.Errorf("booking payment must be paid to process refund")
	}

	if booking.PaymentIntentID == nil {
		return nil, fmt.Errorf("booking has no payment intent ID")
	}

	// Process Stripe refund
	refundID, err := s.paymentService.RefundPayment(
		ctx,
		*booking.PaymentIntentID,
		0, // 0 = full refund
		stripe.RefundReasonRequestedByCustomer,
	)
	if err != nil {
		return nil, fmt.Errorf("process stripe refund: %w", err)
	}

	// Mark booking as refunded
	if err := booking.RefundPayment(); err != nil {
		return nil, fmt.Errorf("refund booking payment: %w", err)
	}

	// Store refund ID in booking notes for tracking
	notes := fmt.Sprintf("Refund processed: %s", refundID)
	booking.SetPartnerNotes(notes)

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update refunded booking: %w", err)
	}

	return booking, nil
}

func (s *BookingService) GetClientBookings(ctx context.Context, clientID uuid.UUID, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookings, err := s.bookingRepo.GetByClientID(ctx, clientID, filter)
	if err != nil {
		return nil, fmt.Errorf("get client bookings: %w", err)
	}

	return bookings, nil
}

func (s *BookingService) GetPartnerBookings(ctx context.Context, partnerID uuid.UUID, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookings, err := s.bookingRepo.GetByPartnerID(ctx, partnerID, filter)
	if err != nil {
		return nil, fmt.Errorf("get partner bookings: %w", err)
	}

	return bookings, nil
}

func (s *BookingService) GetUpcomingBookings(ctx context.Context, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookings, err := s.bookingRepo.GetUpcoming(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get upcoming bookings: %w", err)
	}

	return bookings, nil
}