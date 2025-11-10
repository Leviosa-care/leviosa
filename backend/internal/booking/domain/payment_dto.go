package domain

type ProcessPaymentRequest struct {
	PaymentIntentID string `json:"payment_intent_id" validate:"required"`
}

// Payment Intent DTOs
type CreatePaymentIntentResponse struct {
	PaymentIntentID string `json:"payment_intent_id"`
	ClientSecret    string `json:"client_secret"`
	Amount          int    `json:"amount"`
	Currency        string `json:"currency"`
}

type PaymentIntentStatusResponse struct {
	PaymentIntentID string                      `json:"payment_intent_id"`
	Status          string                      `json:"status"`
	Amount          int                         `json:"amount"`
	Currency        string                      `json:"currency"`
	LastError       *PaymentIntentErrorResponse `json:"last_error,omitempty"`
}

type PaymentIntentErrorResponse struct {
	Code        string `json:"code"`
	DeclineCode string `json:"decline_code,omitempty"`
	Message     string `json:"message"`
	Type        string `json:"type"`
}

type RefundResponse struct {
	RefundID string `json:"refund_id"`
	Amount   int    `json:"amount"`
	Status   string `json:"status"`
}
