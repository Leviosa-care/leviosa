package aggregatorHandler

import (
	"net/http"

	mw "github.com/Leviosa-care/core/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Sends an OTP to the provided email address for verification.
	router.HandleFunc("POST /auth/email", mw.EnableCORS(h.CheckEmailSendOTP))

}
