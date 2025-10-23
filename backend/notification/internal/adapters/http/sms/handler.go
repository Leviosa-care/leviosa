package smsHandler

import (
	"net/http"

	"github.com/Leviosa-care/notification/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc ports.SMSService
}
