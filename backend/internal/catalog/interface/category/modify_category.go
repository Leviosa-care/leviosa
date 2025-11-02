package categoryHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) ModifyCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: this an admin only request

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing modify category",
		"operation", "modify_category",
		"method", r.Method,
		"path", r.URL.Path)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[0] != "" || parts[1] != "admin" || parts[2] != "categories" {
		httpx.RespondWithError(w, errors.New("invalid URL path format. Expected /admin/categories/{id}"), http.StatusBadRequest)
		return
	}
	categoryID := parts[3] // The ID should be the last part
	if categoryID == "" {
		httpx.RespondWithError(w, errors.New("category ID is missing from the URL"), http.StatusBadRequest)
		return
	}

	var req domain.UpdateCategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&req); err != nil {
		logger.ErrorContext(ctx, "Handler: modify category failed",
			"operation", "modify_category",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	if req.Name == nil && req.Description == nil && req.Status == nil && req.Metadata == nil {
		httpx.RespondWithError(w, errors.New("no updatable fields provided in request body"), http.StatusBadRequest)
		return
	}

	req.ID = categoryID

	if err := h.svc.UpdateCategory(ctx, &req); err != nil {
		// TODO: better error handling to return the proper status
		// 204 (no content) : since there is no body in response, if there was return 200
		// 400 (bad request) : mal formed input, invalid format
		// 404 (not found) : no category with given ID
		// 409 (conflit) : unique constraint violated
		// 500 (Internal Server error) : server error, something broke
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)

		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)

		case errors.Is(err, errs.ErrConflict):
			// This covers unique constraint violations (e.g., duplicate name) or other errs-level conflicts.
			httpx.RespondWithError(w, err, http.StatusConflict)

		case errors.Is(err, errs.ErrExternalService):
			// If you had an external service call in update (e.g., image re-upload)
			httpx.RespondWithError(w, errors.New("external service error during category update"), http.StatusServiceUnavailable)

		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, "Handler: modify category failed",
				"operation", "modify_category",
				"error_context", "internal server error during category update",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)

		default:
			logger.ErrorContext(ctx, "Handler: modify category failed",
				"operation", "modify_category",
				"error_context", "unexpected error from service during category update",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: modify category completed",
		"operation", "modify_category",
		"category_id", categoryID,
		"status_code", http.StatusOK)

	w.WriteHeader(http.StatusNoContent)
}
