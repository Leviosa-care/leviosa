package categoryHandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get category by id",
		"operation", "get_category_by_id",
		"method", r.Method,
		"path", r.URL.Path)

	categoryID := strings.TrimPrefix(r.URL.Path, "/categories/")
	if categoryID == "" || strings.Contains(categoryID, "/") {
		httpx.RespondWithError(w, errors.New("invalid URL path"), http.StatusBadRequest)
		return
	}

	categoryWithImage, err := h.aggr.GetCategoryByIDWithImage(ctx, categoryID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get category by ID")
		return
	}

	logger.InfoContext(ctx, "Handler: get category by id completed",
		"operation", "get_category_by_id",
		"category_id", categoryID,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, categoryWithImage, http.StatusOK)
}
