package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPartnersByCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract category IDs from query parameters
	categoryIDs := r.URL.Query()["category_id"]
	if len(categoryIDs) == 0 {
		logger.WarnContext(ctx, "Handler: Missing category IDs in request",
			"operation", "get_partners_by_categories",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("at least one category_id query parameter is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partners by categories request",
		"operation", "get_partners_by_categories",
		"method", r.Method,
		"path", r.URL.Path,
		"category_ids", categoryIDs,
		"category_count", len(categoryIDs),
		"user_agent", r.Header.Get("User-Agent"))

	partners, err := h.svc.GetAllPartnersByCategories(ctx, categoryIDs)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get partners by categories")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partners by categories completed",
		"operation", "get_partners_by_categories",
		"method", r.Method,
		"path", r.URL.Path,
		"category_ids", categoryIDs,
		"category_count", len(categoryIDs),
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
