package notification

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	notificationDomain "github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
	notificationPorts "github.com/Leviosa-care/leviosa/backend/internal/notification/ports"

	"github.com/google/uuid"
)

// detach returns a context that is never cancelled but carries the same values.
// Used for fire-and-forget goroutines so that HTTP request cancellation does not
// abort in-flight notification enrichment or SMTP delivery.
func detach(ctx context.Context) context.Context {
	return context.WithoutCancel(ctx)
}

// BookingNotificationAdapter bridges the booking service's notification port
// to the notification module's email service. It enriches BookingNotificationData
// with user/product/room details and delegates email dispatch.
//
// All notification methods are fire-and-forget: errors are logged but never
// propagated to the caller, ensuring booking transactions are never rolled
// back due to a notification failure.
type BookingNotificationAdapter struct {
	emailService    notificationPorts.EmailService
	smsService      notificationPorts.SMSService
	frontendOrigin  string
	userFetcher     UserFetcher
	roomFetcher     RoomFetcher
	buildingFetcher BuildingFetcher
	productFetcher  catalogPorts.PublicProductService
}

// UserFetcher retrieves user details by ID.
type UserFetcher interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
}

// UserInfo holds the subset of user fields needed for notifications.
type UserInfo struct {
	Email     string
	FirstName string
	LastName  string
	Phone     string
}

// RoomFetcher retrieves room details by ID.
type RoomFetcher interface {
	GetRoom(ctx context.Context, roomID uuid.UUID) (*RoomInfo, error)
}

// RoomInfo holds the subset of room fields needed for notifications.
type RoomInfo struct {
	Name       string
	BuildingID uuid.UUID
}

// BuildingFetcher retrieves building details by ID.
type BuildingFetcher interface {
	GetBuilding(ctx context.Context, buildingID uuid.UUID) (*BuildingInfo, error)
}

// BuildingInfo holds the subset of building fields needed for notifications.
type BuildingInfo struct {
	Name    string
	Address string
}

// NewBookingNotificationAdapter creates a new adapter.
func NewBookingNotificationAdapter(
	emailService notificationPorts.EmailService,
	smsService notificationPorts.SMSService,
	frontendOrigin string,
	userFetcher UserFetcher,
	roomFetcher RoomFetcher,
	buildingFetcher BuildingFetcher,
	productFetcher catalogPorts.PublicProductService,
) *BookingNotificationAdapter {
	return &BookingNotificationAdapter{
		emailService:    emailService,
		smsService:      smsService,
		frontendOrigin:  frontendOrigin,
		userFetcher:     userFetcher,
		roomFetcher:     roomFetcher,
		buildingFetcher: buildingFetcher,
		productFetcher:  productFetcher,
	}
}

// SendBookingConfirmation sends a booking confirmation email (fire-and-forget).
func (a *BookingNotificationAdapter) SendBookingConfirmation(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	go a.sendAsync(detach(ctx), "booking confirmation", data, a.sendBookingConfirmationEmail)
	return nil
}

// SendBookingCancellation sends a booking cancellation email (fire-and-forget).
func (a *BookingNotificationAdapter) SendBookingCancellation(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	go a.sendAsync(detach(ctx), "booking cancellation", data, a.sendBookingCancellationEmail)
	return nil
}

// SendBookingReminder sends a booking reminder email (fire-and-forget).
func (a *BookingNotificationAdapter) SendBookingReminder(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	go a.sendAsync(detach(ctx), "booking reminder", data, a.sendBookingReminderEmail)
	return nil
}

// SendPaymentConfirmation sends a payment confirmation email (fire-and-forget).
func (a *BookingNotificationAdapter) SendPaymentConfirmation(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	go a.sendAsync(detach(ctx), "payment confirmation", data, a.sendPaymentConfirmationEmail)
	return nil
}

// SendPaymentFailed sends a payment failed email (fire-and-forget).
func (a *BookingNotificationAdapter) SendPaymentFailed(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	go a.sendAsync(detach(ctx), "payment failed", data, a.sendPaymentFailedEmail)
	return nil
}

// sendAsync runs the notification dispatch in a goroutine (fire-and-forget).
func (a *BookingNotificationAdapter) sendAsync(
	ctx context.Context,
	label string,
	data bookingPorts.BookingNotificationData,
	sendFn func(context.Context, bookingPorts.BookingNotificationData) error,
) {
	if err := sendFn(ctx, data); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to send %s notification", label),
			"booking_id", data.BookingID,
			"client_id", data.ClientID,
			"error", err,
		)
	}
}

// enrichData fills in missing notification data fields by querying the relevant services.
func (a *BookingNotificationAdapter) enrichData(ctx context.Context, data *bookingPorts.BookingNotificationData) {
	// Fetch client details (skip for guest bookings — already populated)
	if !data.IsGuestBooking && (data.ClientEmail == "" || data.ClientName == "") {
		user, err := a.userFetcher.GetUserByID(ctx, data.ClientID)
		if err != nil {
			slog.WarnContext(ctx, "failed to fetch client details for notification",
				"client_id", data.ClientID,
				"error", err,
			)
		} else {
			data.ClientEmail = user.Email
			data.ClientName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
			data.ClientPhone = user.Phone
		}
	}

	// Fetch partner details
	if data.PartnerEmail == "" || data.PartnerName == "" {
		user, err := a.userFetcher.GetUserByID(ctx, data.PartnerID)
		if err != nil {
			slog.WarnContext(ctx, "failed to fetch partner details for notification",
				"partner_id", data.PartnerID,
				"error", err,
			)
		} else {
			data.PartnerName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
			data.PartnerEmail = user.Email
		}
	}

	// Fetch product name
	if data.ProductName == "" && data.ProductID != uuid.Nil {
		product, err := a.productFetcher.GetProductByID(ctx, data.ProductID.String())
		if err != nil {
			slog.WarnContext(ctx, "failed to fetch product details for notification",
				"product_id", data.ProductID,
				"error", err,
			)
		} else {
			data.ProductName = product.Name
		}
	}

	// Fetch room name and building details
	if data.RoomName == "" && data.RoomID != uuid.Nil {
		room, err := a.roomFetcher.GetRoom(ctx, data.RoomID)
		if err != nil {
			slog.WarnContext(ctx, "failed to fetch room details for notification",
				"room_id", data.RoomID,
				"error", err,
			)
		} else {
			data.RoomName = room.Name

			if data.BuildingName == "" {
				building, err := a.buildingFetcher.GetBuilding(ctx, room.BuildingID)
				if err != nil {
					slog.WarnContext(ctx, "failed to fetch building details for notification",
						"building_id", room.BuildingID,
						"error", err,
					)
				} else {
					data.BuildingName = building.Name
					data.Address = building.Address
				}
			}
		}
	}
}

func (a *BookingNotificationAdapter) sendBookingConfirmationEmail(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	a.enrichData(ctx, &data)

	// Send email confirmation if address is available.
	if data.ClientEmail == "" {
		slog.WarnContext(ctx, "no client email available for booking confirmation",
			"booking_id", data.BookingID,
		)
	} else {
		formattedDate := data.SlotStartTime.Format("Monday, 02 January 2006")
		formattedTime := fmt.Sprintf("%s – %s",
			data.SlotStartTime.Format("15:04"),
			data.SlotEndTime.Format("15:04"),
		)

		req := notificationDomain.BookingConfirmationRequest{
			ToEmail:         data.ClientEmail,
			ToFirstName:     firstName(data.ClientName),
			ToLastName:      lastName(data.ClientName),
			BookingID:       data.BookingID.String(),
			ProductName:     data.ProductName,
			RoomName:        data.RoomName,
			Building:        data.BuildingName,
			Address:         data.Address,
			Date:            formattedDate,
			Time:            formattedTime,
			PartnerName:     data.PartnerName,
			Amount:          formatAmount(data.TotalPriceCents, data.Currency),
			Year:            time.Now().Year(),
			BookingTokenURL: a.buildTokenURL(data.Token),
		}

		if err := a.emailService.SendBookingConfirmationEmail(ctx, req); err != nil {
			slog.ErrorContext(ctx, "failed to send booking confirmation email",
				"booking_id", data.BookingID,
				"error", err,
			)
		}
	}

	// Send confirmation SMS if phone is available.
	a.sendBookingConfirmationSMS(ctx, data)

	return nil
}

// sendBookingConfirmationSMS sends a booking confirmation SMS with a token URL.
// Bookings with a null token (created before issue 002) send the SMS without
// a token URL rather than failing.
func (a *BookingNotificationAdapter) sendBookingConfirmationSMS(ctx context.Context, data bookingPorts.BookingNotificationData) {
	phone := data.ClientPhone
	if phone == "" {
		return
	}

	formattedDate := data.SlotStartTime.Format("02/01/2006")
	formattedTime := fmt.Sprintf("%s-%s",
		data.SlotStartTime.Format("15:04"),
		data.SlotEndTime.Format("15:04"),
	)

	message := fmt.Sprintf("Confirmation: %s le %s de %s. ",
		data.ProductName,
		formattedDate,
		formattedTime,
	)

	// Append token URL if available (legacy bookings may have an empty token)
	if data.Token != "" {
		message += fmt.Sprintf("Voir votre réservation: %s/bookings?token=%s",
			a.frontendOrigin, data.Token)
	}

	if a.smsService == nil {
		return
	}

	if err := a.smsService.SendSMS(ctx, notificationDomain.GenericSMSRequest{
		PhoneNumber: phone,
		Message:     message,
	}); err != nil {
		slog.ErrorContext(ctx, "failed to send booking confirmation SMS",
			"booking_id", data.BookingID,
			"error", err,
		)
	}
}

func (a *BookingNotificationAdapter) sendBookingCancellationEmail(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	a.enrichData(ctx, &data)

	formattedDate := data.SlotStartTime.Format("Monday, 02 January 2006")
	formattedTime := fmt.Sprintf("%s – %s",
		data.SlotStartTime.Format("15:04"),
		data.SlotEndTime.Format("15:04"),
	)

	// Send cancellation email to the client
	if data.ClientEmail == "" {
		slog.WarnContext(ctx, "no client email available for booking cancellation",
			"booking_id", data.BookingID,
			"client_id", data.ClientID,
		)
	} else {
		clientReq := notificationDomain.BookingCancellationRequest{
			ToEmail:     data.ClientEmail,
			ToFirstName: firstName(data.ClientName),
			ToLastName:  lastName(data.ClientName),
			BookingID:   data.BookingID.String(),
			ProductName: data.ProductName,
			RoomName:    data.RoomName,
			Date:        formattedDate,
			Time:        formattedTime,
			Reason:      data.CancellationReason,
			Year:        time.Now().Year(),
		}
		if err := a.emailService.SendBookingCancellationEmail(ctx, clientReq); err != nil {
			slog.ErrorContext(ctx, "failed to send cancellation email to client",
				"booking_id", data.BookingID,
				"client_email", data.ClientEmail,
				"error", err,
			)
		}
	}

	// Send cancellation email to the partner
	if data.PartnerEmail == "" {
		slog.WarnContext(ctx, "no partner email available for booking cancellation",
			"booking_id", data.BookingID,
			"partner_id", data.PartnerID,
		)
	} else {
		partnerReq := notificationDomain.BookingCancellationRequest{
			ToEmail:     data.PartnerEmail,
			ToFirstName: firstName(data.PartnerName),
			ToLastName:  lastName(data.PartnerName),
			BookingID:   data.BookingID.String(),
			ProductName: data.ProductName,
			RoomName:    data.RoomName,
			Date:        formattedDate,
			Time:        formattedTime,
			Reason:      data.CancellationReason,
			Year:        time.Now().Year(),
		}
		if err := a.emailService.SendBookingCancellationEmail(ctx, partnerReq); err != nil {
			slog.ErrorContext(ctx, "failed to send cancellation email to partner",
				"booking_id", data.BookingID,
				"partner_email", data.PartnerEmail,
				"error", err,
			)
		}
	}

	return nil
}

func (a *BookingNotificationAdapter) sendBookingReminderEmail(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	a.enrichData(ctx, &data)

	if data.ClientEmail == "" {
		slog.WarnContext(ctx, "no client email available for booking reminder",
			"booking_id", data.BookingID,
		)
		return nil
	}

	formattedDate := data.SlotStartTime.Format("Monday, 02 January 2006")
	formattedTime := fmt.Sprintf("%s – %s",
		data.SlotStartTime.Format("15:04"),
		data.SlotEndTime.Format("15:04"),
	)

	req := notificationDomain.BookingReminderRequest{
		ToEmail:     data.ClientEmail,
		ToFirstName: firstName(data.ClientName),
		ToLastName:  lastName(data.ClientName),
		BookingID:   data.BookingID.String(),
		ProductName: data.ProductName,
		RoomName:    data.RoomName,
		Building:    data.BuildingName,
		Address:     data.Address,
		Date:        formattedDate,
		Time:        formattedTime,
		PartnerName: data.PartnerName,
		Year:        time.Now().Year(),
	}

	return a.emailService.SendBookingReminderEmail(ctx, req)
}

func (a *BookingNotificationAdapter) sendPaymentConfirmationEmail(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	a.enrichData(ctx, &data)

	if data.ClientEmail == "" {
		return fmt.Errorf("no client email available for payment confirmation")
	}

	req := notificationDomain.PaymentNotificationRequest{
		ToEmail:     data.ClientEmail,
		ToFirstName: firstName(data.ClientName),
		ToLastName:  lastName(data.ClientName),
		Amount:      formatAmount(data.TotalPriceCents, data.Currency),
		Product:     data.ProductName,
		PaymentDate: time.Now().Format("02 January 2006"),
	}

	return a.emailService.SendPaymentNotificationEmail(ctx, req)
}

func (a *BookingNotificationAdapter) sendPaymentFailedEmail(ctx context.Context, data bookingPorts.BookingNotificationData) error {
	a.enrichData(ctx, &data)

	if data.ClientEmail == "" {
		return fmt.Errorf("no client email available for payment failed notification")
	}

	req := notificationDomain.PaymentNotificationRequest{
		ToEmail:     data.ClientEmail,
		ToFirstName: firstName(data.ClientName),
		ToLastName:  lastName(data.ClientName),
		Amount:      formatAmount(data.TotalPriceCents, data.Currency),
		Product:     data.ProductName,
		PaymentDate: time.Now().Format("02 January 2006"),
		RetryURL:    a.frontendOrigin + "/bookings",
	}

	return a.emailService.SendPaymentFailedEmail(ctx, req)
}

// buildTokenURL constructs the guest booking URL from a token.
// Returns an empty string for registered-user bookings (no token).
func (a *BookingNotificationAdapter) buildTokenURL(token string) string {
	if token == "" {
		return ""
	}
	return fmt.Sprintf("%s/bookings?token=%s", a.frontendOrigin, token)
}

// formatAmount converts cents to a human-readable currency string.
func formatAmount(cents int, currency string) string {
	amount := float64(cents) / 100.0
	switch currency {
	case "EUR":
		return fmt.Sprintf("€%.2f", amount)
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	case "GBP":
		return fmt.Sprintf("£%.2f", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, currency)
	}
}

// firstName extracts the first name from a full name string.
func firstName(fullName string) string {
	for i, c := range fullName {
		if c == ' ' {
			return fullName[:i]
		}
	}
	return fullName
}

// lastName extracts the last name from a full name string.
func lastName(fullName string) string {
	for i, c := range fullName {
		if c == ' ' {
			return fullName[i+1:]
		}
	}
	return ""
}
