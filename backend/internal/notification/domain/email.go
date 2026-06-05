package domain

// OTP Email
type OTPEmailRequest struct {
	ToEmail   string `json:"to_email"`
	OTP       string `json:"otp"`
	FromEmail string `json:"from_email"`
	LogoURL   string `json:"logo_url"`
}

// Welcome Email
type WelcomeEmailRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	FromEmail   string `json:"from_email"`
	CompanyName string `json:"company_name"`
	LogoURL     string `json:"logo_url"`
}

// Verify Email
type VerifyEmailRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	FromEmail   string `json:"from_email"`
	CompanyName string `json:"company_name"`
	LogoURL     string `json:"logo_url"`
}

// Event Notification
type EventNotificationRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	Event       string `json:"event"`
	Details     string `json:"details"`
	FromEmail   string `json:"from_email"`
	CompanyName string `json:"company_name"`
	LogoURL     string `json:"logo_url"`
}

// Payment Notification
type PaymentNotificationRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	Amount      string `json:"amount"`
	Product     string `json:"product"`
	PaymentDate string `json:"payment_date"`
	FromEmail   string `json:"from_email"`
	CompanyName string `json:"company_name"`
	LogoURL     string `json:"logo_url"`
	RetryURL    string `json:"retry_url,omitempty"` // non-empty for payment_failed emails
}

// Booking Confirmation Email
type BookingConfirmationRequest struct {
	ToEmail         string `json:"to_email"`
	ToFirstName     string `json:"to_first_name"`
	ToLastName      string `json:"to_last_name"`
	BookingID       string `json:"booking_id"`
	ProductName     string `json:"product_name"`
	RoomName        string `json:"room_name"`
	Building        string `json:"building"`
	Address         string `json:"address"`
	Date            string `json:"date"`
	Time            string `json:"time"`
	PartnerName     string `json:"partner_name"`
	Amount          string `json:"amount"`
	Year            int    `json:"year"`
	BookingTokenURL string `json:"booking_token_url,omitempty"` // non-empty only for guest bookings
}

// Booking Cancellation Email
type BookingCancellationRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	BookingID   string `json:"booking_id"`
	ProductName string `json:"product_name"`
	RoomName    string `json:"room_name"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Reason      string `json:"reason"`
	Year        int    `json:"year"`
}

// Booking Reminder Email
type BookingReminderRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	BookingID   string `json:"booking_id"`
	ProductName string `json:"product_name"`
	RoomName    string `json:"room_name"`
	Building    string `json:"building"`
	Address     string `json:"address"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	PartnerName string `json:"partner_name"`
	Year        int    `json:"year"`
}

// EmailRequest is the generic email structure used by SMTP client
type EmailRequest struct {
	To         string
	Subject    string
	Template   string
	Data       interface{}
	CarbonCopy map[string]string // email -> name
	Images     map[string]string // path -> rename
}
