package domain

// OTP Email
type OTPEmailRequest struct {
	ToEmail   string `json:"to_email"`
	OTP       string `json:"otp"`
	FromEmail string `json:"from_email"`
}

// Welcome Email
type WelcomeEmailRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	FromEmail   string `json:"from_email"`
	CompanyName string `json:"company_name"`
}

// Verify Email
type VerifyEmailRequest struct {
	ToEmail     string `json:"to_email"`
	ToFirstName string `json:"to_first_name"`
	ToLastName  string `json:"to_last_name"`
	FromEmail   string `json:"from_email"`
	CompanyName string `json:"company_name"`
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
