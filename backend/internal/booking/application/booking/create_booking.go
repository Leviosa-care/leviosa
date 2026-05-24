package booking

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	bookingContracts "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/booking"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// CreateBooking creates a new booking for a specific product slot within an availability.
//
// Flow:
//  1. Validate and fetch availability
//  2. Validate and fetch product
//  3. Calculate slot end time based on product duration
//  4. Validate slot is within availability bounds
//  5. Validate slot alignment to 10-minute base
//  6. Check for overlapping bookings
//  7. Get product price from Stripe
//  8. Create payment intent if price > 0
//  9. Create and persist booking entity
//
// Parameters:
//   - availabilityID: The availability block to book within
//   - clientID: The user making the booking
//   - productID: The product/service being booked
//   - slotStartTime: Desired start time (must be aligned to 10-minute boundaries)
//   - clientNotes: Optional notes from client
//
// Returns:
//   - Created booking with encrypted fields
//   - Error if validation fails or slot is unavailable
func (s *BookingService) CreateBooking(
	ctx context.Context,
	availabilityID uuid.UUID,
	clientID *uuid.UUID,
	productID uuid.UUID,
	slotStartTime time.Time,
	clientNotes string,
	guestFirstName, guestLastName, guestEmail, guestPhone string,
) (*domain.Booking, error) {
	// 1. Fetch and decrypt availability
	availabilityEncx, err := s.availabilityRepo.GetByID(ctx, availabilityID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("availability not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("get availability: %w", err)
	}

	availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availabilityEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("availability", err)
	}

	// Verify availability is still open for bookings
	if !availability.IsAvailable() {
		return nil, fmt.Errorf("availability is no longer available for booking")
	}

	// 2. Fetch product from catalog service
	product, err := s.productService.GetProductByID(ctx, productID.String())
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("product not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("get product: %w", err)
	}

	// 3. Calculate slot end time based on product duration
	slotEndTime := slotStartTime.Add(time.Duration(product.Duration) * time.Minute)

	// 4. Validate slot is within availability time bounds
	if slotStartTime.Before(availability.StartTime) {
		return nil, errs.NewInvalidValueErr("slot_start_time: slot start time is before availability start time")
	}

	if slotEndTime.After(availability.EndTime) {
		return nil, errs.NewInvalidValueErr("slot_end_time: slot end time extends beyond availability end time")
	}

	// 5. Validate slot start time is aligned to 10-minute base
	if !isAlignedToBaseSlot(slotStartTime) {
		return nil, errs.NewInvalidValueErr(
			fmt.Sprintf("slot_start_time: must be aligned to %d-minute boundaries (e.g., :00, :10, :20)",
				bookingContracts.BaseTimeSlotMinutes))
	}

	// 6. Check for overlapping bookings
	existingBookingsEncx, err := s.bookingRepo.GetBookingsByAvailability(ctx, availabilityID)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("check existing bookings: %w", err)
	}

	// Decrypt existing bookings for overlap detection
	existingBookings := make([]*domain.Booking, 0, len(existingBookingsEncx))
	for _, bookingEncx := range existingBookingsEncx {
		existingBooking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt existing booking %s: %w", bookingEncx.ID, err)
		}
		existingBookings = append(existingBookings, existingBooking)
	}

	// Use slot calculator's overlap detection
	if hasOverlap(slotStartTime, slotEndTime, existingBookings) {
		return nil, errs.NewConflictErr(fmt.Errorf("time slot is already booked"))
	}

	// 7. Get product price from pricing service
	// Default currency for Europe-based business
	const defaultCurrency = "EUR"

	price, err := s.priceService.GetActiveOneTimePriceByProductID(ctx, productID.String(), defaultCurrency)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("no active price found for product: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("get product price: %w", err)
	}

	totalPriceCents := price.Amount

	// 8. Create booking entity with slot information
	booking, err := domain.NewBooking(
		availabilityID,
		clientID,
		availability.UserID, // partner ID
		availability.RoomID,
		totalPriceCents,
		price.Currency,
		guestFirstName,
		guestLastName,
		guestEmail,
		guestPhone,
	)
	if err != nil {
		return nil, fmt.Errorf("create booking entity: %w", err)
	}

	// Set slot information (new fields)
	booking.ProductID = productID
	booking.SlotStartTime = slotStartTime
	booking.SlotEndTime = slotEndTime

	// Set client notes if provided
	if clientNotes != "" {
		booking.SetClientNotes(clientNotes)
	}

	// 9. Create payment intent if booking has a price
	if totalPriceCents > 0 {
		description := fmt.Sprintf("Booking for %s - %s to %s",
			product.Name,
			slotStartTime.Format("Jan 02, 15:04"),
			slotEndTime.Format("15:04"))

		metadata := map[string]string{
			"booking_id":      booking.ID.String(),
			"availability_id": availabilityID.String(),
			"partner_id":      availability.UserID.String(),
			"product_id":      productID.String(),
			"slot_start":      slotStartTime.Format(time.RFC3339),
			"slot_end":        slotEndTime.Format(time.RFC3339),
		}
		if clientID != nil {
			metadata["client_id"] = clientID.String()
		}
		if booking.IsGuestBooking() {
			metadata["guest_name"]  = booking.GuestDisplayName()
			metadata["guest_email"] = guestEmail
			metadata["guest_phone"] = guestPhone
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

		booking.SetPaymentIntentID(paymentIntentID)
	}

	// 10. Encrypt and persist booking to repository
	bookingEncx, err := domain.ProcessBookingEncx(ctx, s.crypto, booking)
	if err != nil {
		return nil, fmt.Errorf("encrypt booking: %w", err)
	}

	if err := s.bookingRepo.Create(ctx, bookingEncx); err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	// 11. Send booking confirmation notification (best effort - don't fail booking on notification error)
	if s.notificationService != nil {
		notificationData := s.buildNotificationData(booking, product.Name)
		if err := s.notificationService.SendBookingConfirmation(ctx, notificationData); err != nil {
			slog.WarnContext(ctx, "failed to send booking confirmation notification",
				"booking_id", booking.ID,
				"is_guest", booking.IsGuestBooking(),
				"err", err)
		}
	}

	return booking, nil
}

// isAlignedToBaseSlot checks if a time is aligned to the base time slot boundary.
// For a 10-minute base slot, valid times are :00, :10, :20, :30, :40, :50 with 0 seconds.
func isAlignedToBaseSlot(t time.Time) bool {
	return t.Minute()%bookingContracts.BaseTimeSlotMinutes == 0 && t.Second() == 0
}
