package http

import (
	"net/http"

	"github.com/Leviosa-care/notification/internal/adapters/http/mail"
	"github.com/Leviosa-care/notification/internal/adapters/http/sms"
)

func RegisterRoutes(
	router *http.ServeMux,
	mailHandler mailHandler.Handler,
	smsHandler smsHandler.Handler,
) {
	smsHandler.RegisterRoutes(router)
	mailHandler.RegisterRoutes(router)
}
