package productHandler

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

func (h *handler) CreateProductWithPrice(w http.ResponseWriter, r *http.Request) {
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

	logger.InfoContext(ctx, "Handler: Processing create product with price",
		"operation", "create_product_with_price",
		"method", r.Method,
		"path", r.URL.Path)

	var request domain.CreateProductWithPriceRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
	if err := decoder.Decode(&request); err != nil {
		logger.ErrorContext(ctx, "Handler: create product with price failed",
			"operation", "create_product_with_price",
			"error_context", "invalid JSON request body",
			"status_code", http.StatusBadRequest,
			"error", err)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// productID, priceID, err := h.productPriceService.CreateProductWithPrice(ctx, &request)
	productID, priceID, err := h.aggr.CreateProductWithPrice(ctx, &request)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrExternalService):
			httpx.RespondWithError(w, errors.New("failed to create price in external payment system"), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrAlreadyExists):
			httpx.RespondWithError(w, err, http.StatusConflict)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrDomainNotCreated):
			httpx.RespondWithError(w, errors.New("failed to create product due to an unprocessable entity"), http.StatusUnprocessableEntity)
		case errors.Is(err, errs.ErrRateLimit):
			httpx.RespondWithError(w, errors.New("rate limit exceeded"), http.StatusTooManyRequests)
		case errors.Is(err, errs.ErrExpiredToken), errors.Is(err, errs.ErrAccountLocked):
			httpx.RespondWithError(w, errors.New("authentication or authorization issue"), http.StatusUnauthorized)
		case errors.Is(err, errs.ErrParsing), errors.Is(err, errs.ErrFormat):
			httpx.RespondWithError(w, fmt.Errorf("input format error: %w", err), http.StatusBadRequest)
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			logger.ErrorContext(ctx, "Handler: create product with price failed",
				"operation", "create_product_with_price",
				"error_context", "internal server error during product creation",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("internal server error occurred"), http.StatusInternalServerError)
		default:
			logger.ErrorContext(ctx, "Handler: create product with price failed",
				"operation", "create_product_with_price",
				"error_context", "unexpected error from service during product with price creation",
				"status_code", http.StatusInternalServerError,
				"error", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	logger.InfoContext(ctx, "Handler: create product with price completed",
		"operation", "create_product_with_price",
		"product_id", productID,
		"status_code", http.StatusCreated)

	httpx.RespondWithJSON(
		w,
		struct {
			ProductID string `json:"product_id"`
			PriceID   string `json:"price_id"`
			Message   string `json:"message"`
		}{
			ProductID: productID,
			PriceID:   priceID,
			Message:   "Product created successfully!",
		},
		http.StatusCreated,
	)
}

// NOTE: the old way
// func (h *handler) CreateProductWithPrice(w http.ResponseWriter, r *http.Request) {
// 	if r.Header.Get("Content-Type") != "application/json" {
// 		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
// 		return
// 	}
//
// 	ctx := r.Context()
//
// 	var request domain.CreateProductWithPriceRequest
// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields() // Prevent clients from sending unexpected fields
// 	if err := decoder.Decode(&request); err != nil {
// 		log.Printf("Handler: Error decoding JSON body: %v", err)
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
// 		return
// 	}
//
// 	productID, err := h.productService.CreateProduct(ctx, &request.Product)
// 	if err != nil {
// 		log.Printf("Handler: Service CreateProduct failed: %v", err)
// 		switch {
// 		case errors.Is(err, errs.ErrInvalidValue):
// 			httpx.RespondWithError(w, err, http.StatusBadRequest)
// 		case errors.Is(err, errs.ErrAlreadyExists):
// 			httpx.RespondWithError(w, err, http.StatusConflict)
// 		case errors.Is(err, errs.ErrDomainNotFound):
// 			httpx.RespondWithError(w, err, http.StatusNotFound)
// 		case errors.Is(err, errs.ErrDomainNotCreated):
// 			httpx.RespondWithError(w, errors.New("failed to create product due to an unprocessable entity"), http.StatusUnprocessableEntity)
// 		case errors.Is(err, errs.ErrRateLimit):
// 			httpx.RespondWithError(w, errors.New("rate limit exceeded"), http.StatusTooManyRequests)
// 		case errors.Is(err, errs.ErrExpiredToken), errors.Is(err, errs.ErrAccountLocked):
// 			httpx.RespondWithError(w, errors.New("authentication or authorization issue"), http.StatusUnauthorized)
// 		case errors.Is(err, errs.ErrParsing), errors.Is(err, errs.ErrFormat):
// 			httpx.RespondWithError(w, fmt.Errorf("input format error: %w", err), http.StatusBadRequest)
// 		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
// 			log.Printf("Handler: Internal server error during product creation: %v", err)
// 			httpx.RespondWithError(w, errors.New("internal server error occurred"), http.StatusInternalServerError)
// 		default:
// 			log.Printf("Handler: Unhandled error from service during product creation: %v", err)
// 			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
// 		}
// 		return
// 	}
//
// 	priceID, err := h.priceService.CreatePrice(ctx, productID, &request.Price)
// 	if err != nil {
// 		log.Printf("Handler: Service CreatePrice failed: %v", err)
// 		switch {
// 		case errors.Is(err, errs.ErrInvalidValue):
// 			httpx.RespondWithError(w, err, http.StatusBadRequest)
// 		case errors.Is(err, errs.ErrDomainNotFound):
// 			// This case handles the scenario where the product created just before is not found.
// 			// This is a serious data consistency issue and should be treated as an internal error.
// 			httpx.RespondWithError(w, errors.New("internal server error: parent product not found"), http.StatusInternalServerError)
// 		case errors.Is(err, errs.ErrExternalService):
// 			httpx.RespondWithError(w, errors.New("failed to create price in external payment system"), http.StatusServiceUnavailable)
// 		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
// 			log.Printf("Handler: Internal server error during price creation: %v", err)
// 			httpx.RespondWithError(w, errors.New("an internal server error occurred during price creation"), http.StatusInternalServerError)
// 		default:
// 			log.Printf("Handler: Unhandled error from price service during creation: %v", err)
// 			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
// 		}
// 		return
// 	}
//
// 	httpx.RespondWithJSON(
// 		w,
// 		struct {
// 			ProductID string `json:"product_id"`
// 			PriceID   string `json:"price_id"`
// 			Message   string `json:"message"`
// 		}{
// 			ProductID: productID,
// 			PriceID:   priceID,
// 			Message:   "Product created successfully!",
// 		},
// 		http.StatusCreated,
// 	)
// 	log.Printf("Handler: Product with ID %s and price with ID %s created successfully.", productID, priceID)
// }
//
