package domain

import (
	"bytes"
	"html/template"
	"testing"
)

func TestBookingConfirmationTemplateRenders(t *testing.T) {
	tmpl := parseTemplateOrFail(t, "booking_confirmation")

	data := struct {
		FirstName       string
		LastName        string
		BookingID       string
		ProductName     string
		RoomName        string
		Building        string
		Address         string
		Date            string
		Time            string
		PartnerName     string
		Amount          string
		Year            int
		BookingTokenURL string
	}{
		FirstName:   "Alice",
		LastName:    "Smith",
		BookingID:   "550e8400-e29b-41d4-a716-446655440000",
		ProductName: "Massage Therapy",
		RoomName:    "Room A",
		Building:    "Main Building",
		Address:     "123 Rue de Paris, 75001 Paris",
		Date:        "Monday, 15 June 2026",
		Time:        "10:00 – 11:00",
		PartnerName: "Bob Jones",
		Amount:      "€50.00",
		Year:        2026,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("failed to execute booking_confirmation template: %v", err)
	}

	assertContains(t, buf.String(), "Massage Therapy")
	assertContains(t, buf.String(), "Bob Jones")
	assertContains(t, buf.String(), "Room A")
	assertContains(t, buf.String(), "Main Building")
	assertContains(t, buf.String(), "123 Rue de Paris")
	assertContains(t, buf.String(), "Monday, 15 June 2026")
	assertContains(t, buf.String(), "10:00 – 11:00")
	assertContains(t, buf.String(), "€50.00")
	assertContains(t, buf.String(), "550e8400-e29b-41d4-a716-446655440000")
	assertContains(t, buf.String(), "2026")
	assertContains(t, buf.String(), "Politique d'annulation")

	// No token URL for registered users — guest link section must not appear.
	if bytes.Contains(buf.Bytes(), []byte("bookings?token=")) {
		t.Error("token link section should not appear when BookingTokenURL is empty")
	}
}

func TestBookingConfirmationTemplate_GuestTokenLink(t *testing.T) {
	tmpl := parseTemplateOrFail(t, "booking_confirmation")

	data := struct {
		FirstName       string
		LastName        string
		BookingID       string
		ProductName     string
		RoomName        string
		Building        string
		Address         string
		Date            string
		Time            string
		PartnerName     string
		Amount          string
		Year            int
		BookingTokenURL string
	}{
		FirstName:       "Jane",
		LastName:        "Doe",
		BookingID:       "550e8400-e29b-41d4-a716-446655440000",
		ProductName:     "Massage Therapy",
		RoomName:        "Room A",
		Building:        "Main Building",
		Address:         "123 Rue de Paris",
		Date:            "Monday, 15 June 2026",
		Time:            "10:00 – 11:00",
		PartnerName:     "Bob Jones",
		Amount:          "€50.00",
		Year:            2026,
		BookingTokenURL: "https://app.leviosa.com/bookings?token=abc123",
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("failed to execute booking_confirmation template with token: %v", err)
	}

	assertContains(t, buf.String(), "https://app.leviosa.com/bookings?token=abc123")
	assertContains(t, buf.String(), "Voir ma réservation")
}

func TestBookingCancellationTemplateRenders(t *testing.T) {
	tmpl := parseTemplateOrFail(t, "booking_cancellation")

	data := struct {
		FirstName   string
		LastName    string
		BookingID   string
		ProductName string
		RoomName    string
		Date        string
		Time        string
		Reason      string
		Year        int
	}{
		FirstName:   "Alice",
		LastName:    "Smith",
		BookingID:   "550e8400-e29b-41d4-a716-446655440000",
		ProductName: "Massage Therapy",
		RoomName:    "Room A",
		Date:        "Monday, 15 June 2026",
		Time:        "10:00 – 11:00",
		Reason:      "Client needs to reschedule",
		Year:        2026,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("failed to execute booking_cancellation template: %v", err)
	}

	assertContains(t, buf.String(), "550e8400-e29b-41d4-a716-446655440000")
	assertContains(t, buf.String(), "Massage Therapy")
	assertContains(t, buf.String(), "Room A")
	assertContains(t, buf.String(), "Monday, 15 June 2026")
	assertContains(t, buf.String(), "10:00 – 11:00")
	assertContains(t, buf.String(), "Client needs to reschedule")
	assertContains(t, buf.String(), "2026")
}

func TestBookingReminderTemplateRenders(t *testing.T) {
	tmpl := parseTemplateOrFail(t, "booking_reminder")

	data := struct {
		FirstName   string
		LastName    string
		BookingID   string
		ProductName string
		RoomName    string
		Building    string
		Address     string
		Date        string
		Time        string
		PartnerName string
		Year        int
	}{
		FirstName:   "Alice",
		LastName:    "Smith",
		BookingID:   "550e8400-e29b-41d4-a716-446655440000",
		ProductName: "Massage Therapy",
		RoomName:    "Room A",
		Building:    "Main Building",
		Address:     "123 Rue de Paris, 75001 Paris",
		Date:        "Monday, 15 June 2026",
		Time:        "10:00 – 11:00",
		PartnerName: "Bob Jones",
		Year:        2026,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("failed to execute booking_reminder template: %v", err)
	}

	assertContains(t, buf.String(), "Massage Therapy")
	assertContains(t, buf.String(), "Bob Jones")
	assertContains(t, buf.String(), "Room A")
	assertContains(t, buf.String(), "Main Building")
	assertContains(t, buf.String(), "123 Rue de Paris")
	assertContains(t, buf.String(), "Monday, 15 June 2026")
	assertContains(t, buf.String(), "10:00 – 11:00")
	assertContains(t, buf.String(), "550e8400-e29b-41d4-a716-446655440000")
	assertContains(t, buf.String(), "2026")
}

func TestPaymentFailedTemplateRenders(t *testing.T) {
	tmpl := parseTemplateOrFail(t, "payment_failed")

	data := struct {
		FirstName   string
		LastName    string
		Amount      string
		Product     string
		PaymentDate string
		CompanyName string
		Year        int
		LogoURL     string
	}{
		FirstName:   "Alice",
		LastName:    "Smith",
		Amount:      "€50.00",
		Product:     "Massage Therapy",
		PaymentDate: "15 June 2026",
		CompanyName: "Leviosa",
		Year:        2026,
		LogoURL:     "https://example.com/logo.png",
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("failed to execute payment_failed template: %v", err)
	}

	assertContains(t, buf.String(), "€50.00")
	assertContains(t, buf.String(), "Massage Therapy")
	assertContains(t, buf.String(), "15 June 2026")
	assertContains(t, buf.String(), "Réessayer le paiement")
	assertContains(t, buf.String(), "2026")
}

func parseTemplateOrFail(t *testing.T, name string) *template.Template {
	t.Helper()
	tmpl, err := template.ParseFS(EmailTemplates, "templates/"+name+".html")
	if err != nil {
		t.Fatalf("failed to parse template %s: %v", name, err)
	}
	return tmpl
}

func assertContains(t *testing.T, body, substr string) {
	t.Helper()
	if !bytes.Contains([]byte(body), []byte(substr)) {
		t.Errorf("template output missing expected substring %q", substr)
	}
}
