package notification

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	notificationDomain "github.com/Leviosa-care/leviosa/backend/internal/notification/domain"

	"github.com/google/uuid"
)

// ---- Spy email service ----

type spyEmailService struct {
	mu                   sync.Mutex
	bookingConfirmations []notificationDomain.BookingConfirmationRequest
	bookingCancellations []notificationDomain.BookingCancellationRequest
	bookingReminders     []notificationDomain.BookingReminderRequest
	paymentNotifications []notificationDomain.PaymentNotificationRequest
}

func (s *spyEmailService) SendOTPEmail(ctx context.Context, req notificationDomain.OTPEmailRequest) error {
	return nil
}
func (s *spyEmailService) SendWelcomeEmail(ctx context.Context, req notificationDomain.WelcomeEmailRequest) error {
	return nil
}
func (s *spyEmailService) SendVerifyEmailEmail(ctx context.Context, req notificationDomain.VerifyEmailRequest) error {
	return nil
}
func (s *spyEmailService) SendEventNotificationEmail(ctx context.Context, req notificationDomain.EventNotificationRequest) error {
	return nil
}
func (s *spyEmailService) SendPaymentNotificationEmail(ctx context.Context, req notificationDomain.PaymentNotificationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paymentNotifications = append(s.paymentNotifications, req)
	return nil
}
func (s *spyEmailService) SendPaymentFailedEmail(ctx context.Context, req notificationDomain.PaymentNotificationRequest) error {
	return nil
}
func (s *spyEmailService) SendBookingConfirmationEmail(ctx context.Context, req notificationDomain.BookingConfirmationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bookingConfirmations = append(s.bookingConfirmations, req)
	return nil
}
func (s *spyEmailService) SendBookingCancellationEmail(ctx context.Context, req notificationDomain.BookingCancellationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bookingCancellations = append(s.bookingCancellations, req)
	return nil
}
func (s *spyEmailService) SendBookingReminderEmail(ctx context.Context, req notificationDomain.BookingReminderRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bookingReminders = append(s.bookingReminders, req)
	return nil
}

// ---- Stub fetchers ----

type stubUserFetcher struct {
	users map[uuid.UUID]*UserInfo
}

func (f *stubUserFetcher) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	u, ok := f.users[userID]
	if !ok {
		return nil, fmt.Errorf("user %s not found", userID)
	}
	return u, nil
}

type stubRoomFetcher struct {
	rooms map[uuid.UUID]*RoomInfo
}

func (f *stubRoomFetcher) GetRoom(ctx context.Context, roomID uuid.UUID) (*RoomInfo, error) {
	r, ok := f.rooms[roomID]
	if !ok {
		return nil, fmt.Errorf("room %s not found", roomID)
	}
	return r, nil
}

type stubBuildingFetcher struct {
	buildings map[uuid.UUID]*BuildingInfo
}

func (f *stubBuildingFetcher) GetBuilding(ctx context.Context, buildingID uuid.UUID) (*BuildingInfo, error) {
	b, ok := f.buildings[buildingID]
	if !ok {
		return nil, fmt.Errorf("building %s not found", buildingID)
	}
	return b, nil
}

type stubProductFetcher struct {
	products map[string]*catalogDomain.ProductRes
}

func (f *stubProductFetcher) GetProductByID(ctx context.Context, id string) (*catalogDomain.ProductRes, error) {
	p, ok := f.products[id]
	if !ok {
		return nil, fmt.Errorf("product %s not found", id)
	}
	return p, nil
}
func (f *stubProductFetcher) GetAllPublishedProducts(ctx context.Context) ([]*catalogDomain.ProductRes, error) {
	return nil, nil
}
func (f *stubProductFetcher) GetAllProducts(ctx context.Context) ([]*catalogDomain.ProductRes, error) {
	return nil, nil
}

// ---- Helpers ----

func newTestAdapter(
	spy *spyEmailService,
	users map[uuid.UUID]*UserInfo,
	rooms map[uuid.UUID]*RoomInfo,
	buildings map[uuid.UUID]*BuildingInfo,
	products map[string]*catalogDomain.ProductRes,
) *BookingNotificationAdapter {
	return NewBookingNotificationAdapter(
		spy,
		&stubUserFetcher{users: users},
		&stubRoomFetcher{rooms: rooms},
		&stubBuildingFetcher{buildings: buildings},
		&stubProductFetcher{products: products},
	)
}

func makeBookingData() bookingPorts.BookingNotificationData {
	return bookingPorts.BookingNotificationData{
		BookingID:       uuid.New(),
		ClientID:        uuid.New(),
		PartnerID:       uuid.New(),
		RoomID:          uuid.New(),
		ProductID:       uuid.New(),
		SlotStartTime:   time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
		SlotEndTime:     time.Date(2026, 6, 15, 11, 0, 0, 0, time.UTC),
		TotalPriceCents: 5000,
		Currency:        "EUR",
	}
}

// ---- Tests ----

func TestSendBookingConfirmation_MapsFieldsCorrectly(t *testing.T) {
	data := makeBookingData()

	users := map[uuid.UUID]*UserInfo{
		data.ClientID:  {Email: "client@example.com", FirstName: "Alice", LastName: "Smith", Phone: "+33612345678"},
		data.PartnerID: {Email: "partner@example.com", FirstName: "Bob", LastName: "Jones"},
	}
	rooms := map[uuid.UUID]*RoomInfo{
		data.RoomID: {Name: "Room A", BuildingID: uuid.New()},
	}
	buildingID := rooms[data.RoomID].BuildingID
	buildings := map[uuid.UUID]*BuildingInfo{
		buildingID: {Name: "Main Building", Address: "123 Rue de Paris, 75001 Paris"},
	}
	products := map[string]*catalogDomain.ProductRes{
		data.ProductID.String(): {Name: "Massage Therapy"},
	}

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, users, rooms, buildings, products)

	// Call the internal method directly to verify field mapping without goroutine races.
	err := adapter.sendBookingConfirmationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.bookingConfirmations) != 1 {
		t.Fatalf("expected 1 booking confirmation, got %d", len(spy.bookingConfirmations))
	}

	got := spy.bookingConfirmations[0]

	assertEqual(t, "ToEmail", "client@example.com", got.ToEmail)
	assertEqual(t, "ToFirstName", "Alice", got.ToFirstName)
	assertEqual(t, "ToLastName", "Smith", got.ToLastName)
	assertEqual(t, "ProductName", "Massage Therapy", got.ProductName)
	assertEqual(t, "RoomName", "Room A", got.RoomName)
	assertEqual(t, "Building", "Main Building", got.Building)
	assertEqual(t, "Address", "123 Rue de Paris, 75001 Paris", got.Address)
	assertEqual(t, "Amount", "€50.00", got.Amount)
	assertEqual(t, "PartnerName", "Bob Jones", got.PartnerName)
	assertEqual(t, "Date", "Monday, 15 June 2026", got.Date)
	assertEqual(t, "Time", "10:00 – 11:00", got.Time)
	assertEqual(t, "BookingID", data.BookingID.String(), got.BookingID)
}

func TestSendPaymentConfirmation_MapsFieldsCorrectly(t *testing.T) {
	data := makeBookingData()
	data.ProductName = "Pre-populated Product" // already known

	users := map[uuid.UUID]*UserInfo{
		data.ClientID:  {Email: "client@example.com", FirstName: "Carol", LastName: "White"},
		data.PartnerID: {Email: "partner@example.com", FirstName: "Dave", LastName: "Brown"},
	}

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, users, nil, nil, nil)

	err := adapter.sendPaymentConfirmationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.paymentNotifications) != 1 {
		t.Fatalf("expected 1 payment notification, got %d", len(spy.paymentNotifications))
	}

	got := spy.paymentNotifications[0]

	assertEqual(t, "ToEmail", "client@example.com", got.ToEmail)
	assertEqual(t, "ToFirstName", "Carol", got.ToFirstName)
	assertEqual(t, "ToLastName", "White", got.ToLastName)
	assertEqual(t, "Amount", "€50.00", got.Amount)
	assertEqual(t, "Product", "Pre-populated Product", got.Product)
}

func TestSendBookingConfirmation_GuestBooking(t *testing.T) {
	data := makeBookingData()
	data.IsGuestBooking = true
	data.GuestEmail = "guest@example.com"
	data.ClientEmail = "guest@example.com"
	data.ClientName = "Jane Doe"
	data.ClientID = uuid.Nil

	// Only partner needs to be fetched
	users := map[uuid.UUID]*UserInfo{
		data.PartnerID: {Email: "partner@example.com", FirstName: "Bob", LastName: "Jones"},
	}

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, users, nil, nil, nil)

	err := adapter.sendBookingConfirmationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.bookingConfirmations) != 1 {
		t.Fatalf("expected 1 booking confirmation, got %d", len(spy.bookingConfirmations))
	}

	got := spy.bookingConfirmations[0]
	assertEqual(t, "ToEmail", "guest@example.com", got.ToEmail)
	assertEqual(t, "ToFirstName", "Jane", got.ToFirstName)
	assertEqual(t, "ToLastName", "Doe", got.ToLastName)
}

func TestSendBookingConfirmation_MissingClientEmail_ReturnsError(t *testing.T) {
	data := makeBookingData()
	// No users, no guest email
	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, nil, nil, nil, nil)

	err := adapter.sendBookingConfirmationEmail(context.Background(), data)
	if err == nil {
		t.Fatal("expected error when client email is missing")
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		cents    int
		currency string
		want     string
	}{
		{5000, "EUR", "€50.00"},
		{1000, "USD", "$10.00"},
		{2500, "GBP", "£25.00"},
		{999, "JPY", "9.99 JPY"},
	}
	for _, tt := range tests {
		got := formatAmount(tt.cents, tt.currency)
		if got != tt.want {
			t.Errorf("formatAmount(%d, %q) = %q, want %q", tt.cents, tt.currency, got, tt.want)
		}
	}
}

func TestFirstNameLastName(t *testing.T) {
	assertEqual(t, "firstName", "Alice", firstName("Alice Smith"))
	assertEqual(t, "firstName", "Alice", firstName("Alice"))
	assertEqual(t, "lastName", "Smith", lastName("Alice Smith"))
	assertEqual(t, "lastName", "", lastName("Alice"))
}

func TestSendBookingCancellation_SendsToClientAndPartner(t *testing.T) {
	data := makeBookingData()
	data.CancellationReason = "Client needs to reschedule"
	cancelledAt := time.Date(2026, 6, 14, 9, 0, 0, 0, time.UTC)
	data.CancelledAt = &cancelledAt

	users := map[uuid.UUID]*UserInfo{
		data.ClientID:  {Email: "client@example.com", FirstName: "Alice", LastName: "Smith"},
		data.PartnerID: {Email: "partner@example.com", FirstName: "Bob", LastName: "Jones"},
	}
	rooms := map[uuid.UUID]*RoomInfo{
		data.RoomID: {Name: "Room A", BuildingID: uuid.New()},
	}
	buildingID := rooms[data.RoomID].BuildingID
	buildings := map[uuid.UUID]*BuildingInfo{
		buildingID: {Name: "Main Building", Address: "123 Rue de Paris"},
	}
	products := map[string]*catalogDomain.ProductRes{
		data.ProductID.String(): {Name: "Massage Therapy"},
	}

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, users, rooms, buildings, products)

	err := adapter.sendBookingCancellationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.bookingCancellations) != 2 {
		t.Fatalf("expected 2 cancellation emails (client + partner), got %d", len(spy.bookingCancellations))
	}

	clientEmail := spy.bookingCancellations[0]
	assertEqual(t, "client ToEmail", "client@example.com", clientEmail.ToEmail)
	assertEqual(t, "client ToFirstName", "Alice", clientEmail.ToFirstName)
	assertEqual(t, "client ToLastName", "Smith", clientEmail.ToLastName)
	assertEqual(t, "client BookingID", data.BookingID.String(), clientEmail.BookingID)
	assertEqual(t, "client ProductName", "Massage Therapy", clientEmail.ProductName)
	assertEqual(t, "client RoomName", "Room A", clientEmail.RoomName)
	assertEqual(t, "client Date", "Monday, 15 June 2026", clientEmail.Date)
	assertEqual(t, "client Time", "10:00 – 11:00", clientEmail.Time)
	assertEqual(t, "client Reason", "Client needs to reschedule", clientEmail.Reason)

	partnerEmail := spy.bookingCancellations[1]
	assertEqual(t, "partner ToEmail", "partner@example.com", partnerEmail.ToEmail)
	assertEqual(t, "partner ToFirstName", "Bob", partnerEmail.ToFirstName)
	assertEqual(t, "partner ToLastName", "Jones", partnerEmail.ToLastName)
	assertEqual(t, "partner BookingID", data.BookingID.String(), partnerEmail.BookingID)
	assertEqual(t, "partner ProductName", "Massage Therapy", partnerEmail.ProductName)
	assertEqual(t, "partner Reason", "Client needs to reschedule", partnerEmail.Reason)
}

func TestSendBookingCancellation_GuestBooking(t *testing.T) {
	data := makeBookingData()
	data.IsGuestBooking = true
	data.GuestEmail = "guest@example.com"
	data.ClientEmail = "guest@example.com"
	data.ClientName = "Jane Doe"
	data.ClientID = uuid.Nil
	data.CancellationReason = "Schedule conflict"

	users := map[uuid.UUID]*UserInfo{
		data.PartnerID: {Email: "partner@example.com", FirstName: "Bob", LastName: "Jones"},
	}

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, users, nil, nil, nil)

	err := adapter.sendBookingCancellationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.bookingCancellations) != 2 {
		t.Fatalf("expected 2 cancellation emails (guest client + partner), got %d", len(spy.bookingCancellations))
	}

	assertEqual(t, "guest client ToEmail", "guest@example.com", spy.bookingCancellations[0].ToEmail)
	assertEqual(t, "guest client ToFirstName", "Jane", spy.bookingCancellations[0].ToFirstName)
	assertEqual(t, "partner ToEmail", "partner@example.com", spy.bookingCancellations[1].ToEmail)
}

func TestSendBookingCancellation_MissingClientEmail_StillSendsToPartner(t *testing.T) {
	data := makeBookingData()
	data.CancellationReason = "No reason"

	users := map[uuid.UUID]*UserInfo{
		data.PartnerID: {Email: "partner@example.com", FirstName: "Bob", LastName: "Jones"},
	}

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, users, nil, nil, nil)

	err := adapter.sendBookingCancellationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Partner should still receive the cancellation email even if client email is missing
	if len(spy.bookingCancellations) != 1 {
		t.Fatalf("expected 1 cancellation email (partner only), got %d", len(spy.bookingCancellations))
	}
	assertEqual(t, "partner ToEmail", "partner@example.com", spy.bookingCancellations[0].ToEmail)
}

func TestSendBookingCancellation_MissingPartnerEmail_StillSendsToClient(t *testing.T) {
	data := makeBookingData()
	data.CancellationReason = "No reason"

	users := map[uuid.UUID]*UserInfo{
		data.ClientID: {Email: "client@example.com", FirstName: "Alice", LastName: "Smith"},
	}

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, users, nil, nil, nil)

	err := adapter.sendBookingCancellationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Client should still receive the cancellation email even if partner email is missing
	if len(spy.bookingCancellations) != 1 {
		t.Fatalf("expected 1 cancellation email (client only), got %d", len(spy.bookingCancellations))
	}
	assertEqual(t, "client ToEmail", "client@example.com", spy.bookingCancellations[0].ToEmail)
}

func TestSendBookingCancellation_BothEmailsMissing_LogsAndReturnsNil(t *testing.T) {
	data := makeBookingData()
	data.CancellationReason = "No reason"

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, nil, nil, nil, nil)

	err := adapter.sendBookingCancellationEmail(context.Background(), data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spy.bookingCancellations) != 0 {
		t.Fatalf("expected 0 cancellation emails, got %d", len(spy.bookingCancellations))
	}
}

func TestSendBookingCancellation_FireAndForget(t *testing.T) {
	data := makeBookingData()
	data.ClientEmail = "client@example.com"
	data.PartnerEmail = "partner@example.com"
	data.ClientName = "Test Client"
	data.PartnerName = "Test Partner"
	data.CancellationReason = "Test"

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, nil, nil, nil, nil)

	// SendBookingCancellation returns nil immediately (fire-and-forget)
	err := adapter.SendBookingCancellation(context.Background(), data)
	if err != nil {
		t.Fatalf("expected nil error from fire-and-forget, got: %v", err)
	}

	// Wait for goroutine to execute
	time.Sleep(100 * time.Millisecond)
}

func TestSendBookingConfirmation_FireAndForget(t *testing.T) {
	data := makeBookingData()
	data.ClientEmail = "client@example.com"
	data.ClientName = "Test Client"

	spy := &spyEmailService{}
	adapter := newTestAdapter(spy, nil, nil, nil, nil)

	// SendBookingConfirmation returns nil immediately (fire-and-forget)
	err := adapter.SendBookingConfirmation(context.Background(), data)
	if err != nil {
		t.Fatalf("expected nil error from fire-and-forget, got: %v", err)
	}

	// Wait for goroutine to execute
	time.Sleep(100 * time.Millisecond)

	// The goroutine will fail to send (no client in fetcher) but that's OK,
	// the important thing is SendBookingConfirmation itself returned nil.
}

func assertEqual(t *testing.T, field, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}
