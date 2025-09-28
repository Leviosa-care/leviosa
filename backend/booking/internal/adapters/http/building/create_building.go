package buildingHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
)

func (h *handler) CreateBuilding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Log incoming request
	logger.InfoContext(ctx, "Handler: Processing create building request",
		"operation", "create_building",
		"method", r.Method,
		"path", r.URL.Path,
		"user_agent", r.Header.Get("User-Agent"))

	// Parse request body
	var request domain.CreateBuildingRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_building",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to create building
	building, err := h.svc.CreateBuilding(ctx, request.Name, request.Address, request.City, request.PostalCode, request.Country)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		var statusCode int

		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid request validation"
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrUniqueViolation):
			logLevel = "warn"
			errorContext = "building name already exists"
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "database connection failure"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "database resource exhaustion"
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			logLevel = "error"
			errorContext = "transaction conflict"
			statusCode = http.StatusServiceUnavailable
		default:
			logLevel = "error"
			errorContext = "unexpected error"
			statusCode = http.StatusInternalServerError
		}

		if logLevel == "error" {
			logger.ErrorContext(ctx, "Handler: Create building failed",
				"error", err,
				"operation", "create_building",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Create building failed",
				"error", err,
				"operation", "create_building",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Set additional fields if provided
	if request.Description != "" {
		building.SetDescription(request.Description)
	}
	if request.Phone != "" || request.Email != "" {
		building.SetContactInfo(request.Phone, request.Email)
	}

	// Update building with additional fields if any were set
	if request.Description != "" || request.Phone != "" || request.Email != "" {
		building, err = h.svc.UpdateBuilding(ctx, building.ID, building.Name, building.Address, building.City, building.PostalCode, building.Country, building.Description)
		if err != nil {
			logger.ErrorContext(ctx, "Handler: Failed to update building with additional fields",
				"error", err,
				"building_id", building.ID,
				"operation", "create_building")
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		if request.Phone != "" || request.Email != "" {
			building, err = h.svc.UpdateBuildingContactInfo(ctx, building.ID, request.Phone, request.Email)
			if err != nil {
				logger.ErrorContext(ctx, "Handler: Failed to update building contact info",
					"error", err,
					"building_id", building.ID,
					"operation", "create_building")
				httpx.RespondWithError(w, err, http.StatusInternalServerError)
				return
			}
		}
	}

	// Convert to response DTO
	response := domain.BuildingResponse{
		ID:          building.ID,
		Name:        building.Name,
		Address:     building.Address,
		City:        building.City,
		PostalCode:  building.PostalCode,
		Country:     building.Country,
		Description: building.Description,
		Phone:       building.Phone,
		Email:       building.Email,
		IsActive:    building.IsActive,
		CreatedAt:   building.CreatedAt,
		UpdatedAt:   building.UpdatedAt,
	}

	logger.InfoContext(ctx, "Handler: Building created successfully",
		"building_id", building.ID,
		"building_name", building.Name,
		"operation", "create_building")

	httpx.RespondWithJSON(w, response, http.StatusCreated)
}