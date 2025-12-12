package http

import (
	"github.com/gorilla/mux"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
)

// RegisterRoutes registers notification HTTP routes
func RegisterRoutes(router *mux.Router, svc ports.NotificationService) {
	// Email routes
	router.HandleFunc("/notifications/email/otp", SendOTPEmailHandler(svc)).Methods("POST")
	router.HandleFunc("/notifications/email/welcome", SendWelcomeEmailHandler(svc)).Methods("POST")
	router.HandleFunc("/notifications/email/verify", SendVerifyEmailHandler(svc)).Methods("POST")
	router.HandleFunc("/notifications/email/event", SendEventNotificationHandler(svc)).Methods("POST")
	router.HandleFunc("/notifications/email/payment", SendPaymentNotificationHandler(svc)).Methods("POST")

	// SMS routes
	router.HandleFunc("/notifications/sms/otp", SendOTPSMSHandler(svc)).Methods("POST")
	router.HandleFunc("/notifications/sms/generic", SendGenericSMSHandler(svc)).Methods("POST")
}
