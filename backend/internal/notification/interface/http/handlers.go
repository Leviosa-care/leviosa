package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
)

// Request DTOs
type SendOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type SendWelcomeRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SendVerifyEmailRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SendEventNotificationRequest struct {
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EventName    string `json:"event_name"`
	EventDetails string `json:"event_details"`
}

type SendPaymentNotificationRequest struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Amount      string `json:"amount"`
	Product     string `json:"product"`
	PaymentDate string `json:"payment_date"`
}

type SendOTPSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

type SendGenericSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// Handler implementations following project error handling pattern
func SendOTPEmailHandler(svc ports.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.RespondWithError(w, errs.ErrInvalidInput, http.StatusBadRequest)
			return
		}

		if err := svc.SendOTPEmail(r.Context(), req.Email, req.OTP); err != nil {
			statusCode := classifyError(err)
			httpx.RespondWithError(w, err, statusCode)
			return
		}

		httpx.RespondWithJSON(w, map[string]string{"status": "sent"}, http.StatusOK)
	}
}

func SendWelcomeEmailHandler(svc ports.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendWelcomeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.RespondWithError(w, errs.ErrInvalidInput, http.StatusBadRequest)
			return
		}

		if err := svc.SendWelcomeEmail(r.Context(), req.Email, req.FirstName, req.LastName); err != nil {
			statusCode := classifyError(err)
			httpx.RespondWithError(w, err, statusCode)
			return
		}

		httpx.RespondWithJSON(w, map[string]string{"status": "sent"}, http.StatusOK)
	}
}

func SendVerifyEmailHandler(svc ports.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendVerifyEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.RespondWithError(w, errs.ErrInvalidInput, http.StatusBadRequest)
			return
		}

		if err := svc.SendVerifyEmailEmail(r.Context(), req.Email, req.FirstName, req.LastName); err != nil {
			statusCode := classifyError(err)
			httpx.RespondWithError(w, err, statusCode)
			return
		}

		httpx.RespondWithJSON(w, map[string]string{"status": "sent"}, http.StatusOK)
	}
}

func SendEventNotificationHandler(svc ports.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendEventNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.RespondWithError(w, errs.ErrInvalidInput, http.StatusBadRequest)
			return
		}

		if err := svc.SendEventNotificationEmail(r.Context(), req.Email, req.FirstName, req.LastName, req.EventName, req.EventDetails); err != nil {
			statusCode := classifyError(err)
			httpx.RespondWithError(w, err, statusCode)
			return
		}

		httpx.RespondWithJSON(w, map[string]string{"status": "sent"}, http.StatusOK)
	}
}

func SendPaymentNotificationHandler(svc ports.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendPaymentNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.RespondWithError(w, errs.ErrInvalidInput, http.StatusBadRequest)
			return
		}

		if err := svc.SendPaymentNotificationEmail(r.Context(), req.Email, req.FirstName, req.LastName, req.Amount, req.Product, req.PaymentDate); err != nil {
			statusCode := classifyError(err)
			httpx.RespondWithError(w, err, statusCode)
			return
		}

		httpx.RespondWithJSON(w, map[string]string{"status": "sent"}, http.StatusOK)
	}
}

func SendOTPSMSHandler(svc ports.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendOTPSMSRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.RespondWithError(w, errs.ErrInvalidInput, http.StatusBadRequest)
			return
		}

		if err := svc.SendOTPBySMS(r.Context(), req.PhoneNumber, req.OTP); err != nil {
			statusCode := classifyError(err)
			httpx.RespondWithError(w, err, statusCode)
			return
		}

		httpx.RespondWithJSON(w, map[string]string{"status": "sent"}, http.StatusOK)
	}
}

func SendGenericSMSHandler(svc ports.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendGenericSMSRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.RespondWithError(w, errs.ErrInvalidInput, http.StatusBadRequest)
			return
		}

		if err := svc.SendGenericSMS(r.Context(), req.PhoneNumber, req.Message); err != nil {
			statusCode := classifyError(err)
			httpx.RespondWithError(w, err, statusCode)
			return
		}

		httpx.RespondWithJSON(w, map[string]string{"status": "sent"}, http.StatusOK)
	}
}

// classifyError maps errors to appropriate HTTP status codes
// Following project error handling strategy
func classifyError(err error) int {
	switch {
	case errors.Is(err, errs.ErrInvalidValue), errors.Is(err, errs.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
		return http.StatusServiceUnavailable
	case errors.Is(err, errs.ErrQueryCancelled):
		return http.StatusRequestTimeout
	case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
		return http.StatusServiceUnavailable // Retryable
	default:
		return http.StatusInternalServerError
	}
}
