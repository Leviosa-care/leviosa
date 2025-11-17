package categoryHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAdminAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get admin all categories",
		"operation", "get_admin_all_categories",
		"method", r.Method,
		"path", r.URL.Path)

	categoryWithImages, err := h.aggr.GetAdminAllCategoriesWithImages(ctx)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get admin all categories")
		return
	}

	count := len(categoryWithImages)
	logger.InfoContext(ctx, "Handler: get admin all categories completed",
		"operation", "get_admin_all_categories",
		"count", count,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, categoryWithImages, http.StatusOK)
}
