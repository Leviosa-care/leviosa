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

	if req.Name == nil && req.Description == nil && req.Status == nil {
		httpx.RespondWithError(w, errors.New("no updatable fields provided in request body"), http.StatusBadRequest)
		return
	}

	req.ID = categoryID

	if err := h.svc.UpdateCategory(ctx, &req); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "update category")
		return
	}

	logger.InfoContext(ctx, "Handler: modify category completed",
		"operation", "modify_category",
		"category_id", categoryID,
		"status_code", http.StatusOK)

	w.WriteHeader(http.StatusNoContent)
}
