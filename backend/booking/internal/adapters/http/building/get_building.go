package buildingHandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/google/uuid"
)

func (h *handler) GetBuilding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract building ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		logger.WarnContext(ctx, "Handler: Invalid URL path for get building",
			"path", r.URL.Path,
			"operation", "get_building")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID in path"), http.StatusBadRequest)
		return
	}

	buildingIDStr := pathParts[1]
	buildingID, err := uuid.Parse(buildingIDStr)
	if err != nil {
		logger.WarnContext(ctx, "Handler: Invalid building ID format",
			"building_id", buildingIDStr,
			"error", err,
			"operation", "get_building")
		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid building ID format"), http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get building request",
		"building_id", buildingID,
		"operation", "get_building",
		"method", r.Method,
		"path", r.URL.Path)

	// Call service to get building
	building, err := h.svc.GetBuilding(ctx, buildingID)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			statusCode = http.StatusNotFound
			errorContext = "building not found"
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
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		if statusCode >= 500 {
			logger.ErrorContext(ctx, "Handler: Get building failed",
				"error", err,
				"building_id", buildingID,
				"operation", "get_building",
				"context", errorContext,
				"status_code", statusCode)
		} else {
			logger.WarnContext(ctx, "Handler: Get building failed",
				"error", err,
				"building_id", buildingID,
				"operation", "get_building",
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

	logger.InfoContext(ctx, "Handler: Building retrieved successfully",
		"building_id", buildingID,
		"building_name", building.Name,
		"operation", "get_building")

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) GetAllBuildings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get all buildings request",
		"operation", "get_all_buildings",
		"method", r.Method,
		"path", r.URL.Path)

	// Parse query parameters for filtering
	activeOnly := r.URL.Query().Get("active_only") == "true"

	// Call service to get all buildings
	buildings, err := h.svc.GetAllBuildings(ctx, activeOnly)
	if err != nil {
		var statusCode int
		var errorContext string

		switch {
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
		default:
			statusCode = http.StatusInternalServerError
			errorContext = "unexpected error"
		}

		logger.ErrorContext(ctx, "Handler: Get all buildings failed",
			"error", err,
			"operation", "get_all_buildings",
			"context", errorContext,
			"status_code", statusCode)

		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Convert to response DTOs
	var responses []domain.BuildingResponse
	for _, building := range buildings {
		responses = append(responses, domain.BuildingResponse{
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
		})
	}

	logger.InfoContext(ctx, "Handler: Buildings retrieved successfully",
		"building_count", len(buildings),
		"active_only", activeOnly,
		"operation", "get_all_buildings")

	httpx.RespondWithJSON(w, responses, http.StatusOK)
}