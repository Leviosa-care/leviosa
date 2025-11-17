package partnerHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

func (h *handler) GetAllPartnersByProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract product IDs from query parameters
	productIDs := r.URL.Query()["product_id"]
	if len(productIDs) == 0 {
		logger.WarnContext(ctx, "Handler: Missing product IDs in request",
			"operation", "get_partners_by_products",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr("at least one product_id query parameter is required"), http.StatusBadRequest)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing get partners by products request",
		"operation", "get_partners_by_products",
		"method", r.Method,
		"path", r.URL.Path,
		"product_ids", productIDs,
		"product_count", len(productIDs),
		"user_agent", r.Header.Get("User-Agent"))

	partners, err := h.svc.GetAllPartnersByProducts(ctx, productIDs)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get partners by products")
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get partners by products completed",
		"operation", "get_partners_by_products",
		"method", r.Method,
		"path", r.URL.Path,
		"product_ids", productIDs,
		"product_count", len(productIDs),
		"status_code", http.StatusOK,
		"partner_count", len(partners))

	httpx.RespondWithJSON(w, partners, http.StatusOK)
}
