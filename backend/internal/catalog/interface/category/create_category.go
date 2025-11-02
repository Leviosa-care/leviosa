package categoryHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: that is admin only so use the context to check that
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

	logger.InfoContext(ctx, "Handler: Processing create category",
		"operation", "create_category",
		"method", r.Method,
		"path", r.URL.Path)

	var category domain.CreateCategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&category); err != nil {
		logger.ErrorContext(ctx, "Handler: create category failed",
			"operation", "create_category",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		// TODO: Use a specific handler error or a generic bad request message.
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call Service Layer with io.Reader and *multipart.FileHeader
	categoryID, err := h.svc.CreateCategory(ctx, &category)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			// This covers validation errors from errs.Category.Valid()
			// and invalid metadata JSON from the repository (via errs.NewInvalidValueErr).
			httpx.RespondWithError(w, err, http.StatusBadRequest) // Send back the wrapped errs error message
		case errors.Is(err, errs.ErrAlreadyExists):
			// Catches the unique constraint violation mhttpxed from the repository.
			httpx.RespondWithError(w, err, http.StatusConflict) // Send back the wrapped errs error message
		case errors.Is(err, errs.ErrDomainNotFound):
			// If a dependency of category creation (e.g., a parent entity) wasn't found.
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrDomainNotCreated):
			// This would indicate a semantic problem preventing creation.
			// Re-evaluate if this specific error is still necessary with `ErrUniqueViolation` handled.
			// If it means "semantically invalid input despite basic validation", then 422 is good.
			httpx.RespondWithError(w, errors.New("failed to create category due to an unprocessable entity"), http.StatusUnprocessableEntity)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			// Catch all other unexpected internal errors.
			// Log the full underlying error for operations/debugging.
			logger.ErrorContext(ctx, "Handler: create category failed",
				"operation", "create_category",
				"error_context", "internal server error during category creation",
				"status_code", http.StatusInternalServerError,
				"error", err)
			// Return a generic message to the client to avoid leaking internal details.
			httpx.RespondWithError(w, errors.New("internal server error occurred"), http.StatusInternalServerError)
		default:
			// Catch any unhandled errors that didn't match specific cases.
			logger.ErrorContext(ctx, "Handler: create category failed",
				"operation", "create_category",
				"error_context", "unexpected error from service during category creation",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: create category completed",
		"operation", "create_category",
		"category_id", categoryID,
		"status_code", http.StatusCreated)

	httpx.RespondWithJSON(
		w,
		struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}{
			ID:      categoryID,
			Message: "Category created successfully!",
		},
		http.StatusCreated,
	)
}
