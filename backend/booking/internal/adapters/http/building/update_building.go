package buildingHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	"github.com/google/uuid"
)

func (h *handler) UpdateBuilding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract building ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for update building",
			"path", r.URL.Path,
			"operation", "update_building")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID in path"), http.StatusBadRequest)
		return
	}

	buildingIDStr := pathParts[1]
	buildingID, err := uuid.Parse(buildingIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid building ID format",
			"building_id", buildingIDStr,
			"error", err,
			"operation", "update_building")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing update building request",
		"building_id", buildingID,
		"operation", "update_building",
		"method", r.Method,
		"path", r.URL.Path)

	// Parse request body
	var request domain.UpdateBuildingRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"building_id", buildingID,
			"operation", "update_building")
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to update building
	building, err := h.svc.UpdateBuilding(ctx, buildingID, request.Name, request.Address, request.City, request.PostalCode, request.Country, request.Description)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
			errorContext = "building not found"
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
			errorContext = "invalid request validation"
		case errors.Is(err, errs.ErrUniqueViolation):
			statusCode = http.StatusConflict
			errorContext = "building name already exists"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database resource exhaustion"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			statusCode = http.StatusServiceUnavailable
			errorContext = "transaction conflict"
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Update building failed",
				"error", err,
				"building_id", buildingID,
				"operation", "update_building",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Update building failed",
				"error", err,
				"building_id", buildingID,
				"operation", "update_building",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
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

	logger.InfoContext(ctx, "Handler: Building updated successfully",
		"building_id", buildingID,
		"building_name", building.Name,
		"operation", "update_building")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) UpdateBuildingContactInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract building ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for update building contact",
			"path", r.URL.Path,
			"operation", "update_building_contact")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID in path"), http.StatusBadRequest)
		return
	}

	buildingIDStr := pathParts[1]
	buildingID, err := uuid.Parse(buildingIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid building ID format",
			"building_id", buildingIDStr,
			"error", err,
			"operation", "update_building_contact")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing update building contact info request",
		"building_id", buildingID,
		"operation", "update_building_contact",
		"method", r.Method,
		"path", r.URL.Path)

	// Parse request body
	var request domain.UpdateBuildingContactRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"building_id", buildingID,
			"operation", "update_building_contact")
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	// Call service to update building contact info
	building, err := h.svc.UpdateBuildingContactInfo(ctx, buildingID, request.Phone, request.Email)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
			errorContext = "building not found"
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
			errorContext = "invalid contact info validation"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
			errorContext = "database resource exhaustion"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			statusCode = http.StatusServiceUnavailable
			errorContext = "transaction conflict"
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Update building contact info failed",
				"error", err,
				"building_id", buildingID,
				"operation", "update_building_contact",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Update building contact info failed",
				"error", err,
				"building_id", buildingID,
				"operation", "update_building_contact",
				"context", errorContext,
				"status_code", statusCode)
		}

		httpx.RespondWithError(w, err, statusCode)
		return
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

	logger.InfoContext(ctx, "Handler: Building contact info updated successfully",
		"building_id", buildingID,
		"operation", "update_building_contact")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}