package buildingHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) UpdateBuilding(w http.ResponseWriter, r *http.Request) {
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

	request.ID = buildingID

	// Call service to update building
	building, err := h.svc.UpdateBuilding(ctx, &request)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			statusCode = http.StatusNotFound
			errorContext = "building not found"
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
			errorContext = "invalid request validation"
		case errors.Is(err, errs.ErrConflict):
			statusCode = http.StatusConflict
			errorContext = "building conflict"
		case errors.Is(err, errs.ErrNotEncrypted), errors.Is(err, errs.ErrNotDecrypted):
			statusCode = http.StatusInternalServerError
			errorContext = "encryption error"
		case errors.Is(err, errs.ErrQueryFailed), errors.Is(err, errs.ErrUnexpectedError):
			statusCode = http.StatusInternalServerError
			errorContext = "database operation failed"
		case errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
			errorContext = "request timeout"
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
	}

	logger.InfoContext(ctx, "Handler: Building updated successfully",
		"building_id", buildingID,
		"building_name", building.Name,
		"operation", "update_building")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
