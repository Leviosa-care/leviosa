package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
)

func (h *handler) SetCompanyInstagram(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing set company instagram request",
		"operation", "set_company_instagram",
		"method", r.Method,
		"path", r.URL.Path)

	var request domain.SetCompanyInstagramRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("Handler: Error decoding JSON body: %v", err))
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	response, err := h.svc.SetCompanyInstagram(ctx, &request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "set company instagram")
		return
	}

	logger.InfoContext(ctx, "Handler: Set company instagram completed",
		"operation", "set_company_instagram",
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
