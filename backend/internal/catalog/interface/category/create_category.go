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
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call Service Layer with io.Reader and *multipart.FileHeader
	categoryID, err := h.svc.CreateCategory(ctx, &category)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "create category")
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
