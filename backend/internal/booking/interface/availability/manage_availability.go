package availabilityHandler

// import (
// 	// "encoding/json"
// 	// "errors"
// 	// "fmt"
// 	"net/http"
// 	// "strconv"
// 	"strings"
// 	// "time"
//
// 	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
// 	// "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
// 	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
// 	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
// 	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
//
// 	"github.com/google/uuid"
// )
//
// func (h *handler) GetPartnerAvailabilities(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	logger, err := ctxutil.GetLoggerFromContext(ctx)
// 	if err != nil {
// 		httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Extract partner ID from URL path
// 	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
// 	if len(pathParts) < 3 {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID in path"), http.StatusBadRequest)
// 		return
// 	}
//
// 	partnerID, err := uuid.Parse(pathParts[1])
// 	if err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
// 		return
// 	}
//
// 	// Parse query parameters
// 	filter := ports.AvailabilityFilter{}
// 	if startStr := r.URL.Query().Get("start_time"); startStr != "" {
// 		if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
// 			filter.StartTime = &startTime
// 		}
// 	}
// 	if endStr := r.URL.Query().Get("end_time"); endStr != "" {
// 		if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
// 			filter.EndTime = &endTime
// 		}
// 	}
// 	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
// 		status := domain.AvailabilityStatus(statusStr)
// 		filter.Status = &status
// 	}
// 	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
// 		if limit, err := strconv.Atoi(limitStr); err == nil {
// 			filter.Limit = &limit
// 		}
// 	}
//
// 	// Call service to get partner availabilities
// 	availabilities, err := h.svc.GetPartnerAvailabilities(ctx, partnerID, filter)
// 	if err != nil {
// 		httpx.RespondWithServiceError(w, logger, ctx, err, "get partner availabilities")
// 		return
// 	}
//
// 	// Convert to response DTOs
// 	var responses []domain.AvailabilityResponse
// 	for _, availability := range availabilities {
// 		responses = append(responses, domain.AvailabilityResponse{
// 			ID:              availability.ID,
// 			UserID:          availability.UserID,
// 			RoomID:          availability.RoomID,
// 			StartTime:       availability.StartTime,
// 			EndTime:         availability.EndTime,
// 			MaxCapacity:     availability.MaxCapacity,
// 			CurrentBookings: availability.CurrentBookings,
// 			Status:          availability.Status,
// 			ServiceType:     availability.ServiceType,
// 			PriceCents:      availability.PriceCents,
// 			Notes:           availability.Notes,
// 			RecurrencePattern: availability.RecurrencePattern,
// 			ParentID:        availability.ParentID,
// 			CreatedAt:       availability.CreatedAt,
// 			UpdatedAt:       availability.UpdatedAt,
// 		})
// 	}
//
// 	httpx.RespondWithJSON(w, responses, http.StatusOK)
// }
//
// func (h *handler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	logger, err := ctxutil.GetLoggerFromContext(ctx)
// 	if err != nil {
// 		httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Parse query parameters for filtering
// 	filter := ports.AvailabilityFilter{}
// 	if startStr := r.URL.Query().Get("start_time"); startStr != "" {
// 		if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
// 			filter.StartTime = &startTime
// 		}
// 	}
// 	if endStr := r.URL.Query().Get("end_time"); endStr != "" {
// 		if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
// 			filter.EndTime = &endTime
// 		}
// 	}
// 	if roomIDStr := r.URL.Query().Get("room_id"); roomIDStr != "" {
// 		if roomID, err := uuid.Parse(roomIDStr); err == nil {
// 			filter.RoomID = &roomID
// 		}
// 	}
// 	if partnerIDStr := r.URL.Query().Get("partner_id"); partnerIDStr != "" {
// 		if partnerID, err := uuid.Parse(partnerIDStr); err == nil {
// 			filter.PartnerID = &partnerID
// 		}
// 	}
// 	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
// 		if limit, err := strconv.Atoi(limitStr); err == nil {
// 			filter.Limit = &limit
// 		}
// 	}
//
// 	// Call service to get available slots
// 	availabilities, err := h.svc.GetAvailableSlots(ctx, filter)
// 	if err != nil {
// 		httpx.RespondWithServiceError(w, logger, ctx, err, "get available slots")
// 		return
// 	}
//
// 	// Convert to response DTOs
// 	var responses []domain.AvailabilityResponse
// 	for _, availability := range availabilities {
// 		responses = append(responses, domain.AvailabilityResponse{
// 			ID:              availability.ID,
// 			UserID:          availability.UserID,
// 			RoomID:          availability.RoomID,
// 			StartTime:       availability.StartTime,
// 			EndTime:         availability.EndTime,
// 			MaxCapacity:     availability.MaxCapacity,
// 			CurrentBookings: availability.CurrentBookings,
// 			Status:          availability.Status,
// 			ServiceType:     availability.ServiceType,
// 			PriceCents:      availability.PriceCents,
// 			Notes:           availability.Notes,
// 			RecurrencePattern: availability.RecurrencePattern,
// 			ParentID:        availability.ParentID,
// 			CreatedAt:       availability.CreatedAt,
// 			UpdatedAt:       availability.UpdatedAt,
// 		})
// 	}
//
// 	httpx.RespondWithJSON(w, responses, http.StatusOK)
// }
//
// func (h *handler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
// 	if r.Header.Get("Content-Type") != "application/json" {
// 		httpx.RespondWithError(w, errors.New("unsupported media type: please send 'application/json'"), http.StatusUnsupportedMediaType)
// 		return
// 	}
//
// 	ctx := r.Context()
//
// 	logger, err := ctxutil.GetLoggerFromContext(ctx)
// 	if err != nil {
// 		httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Extract availability ID from URL path
// 	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
// 	if len(pathParts) < 2 {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID in path"), http.StatusBadRequest)
// 		return
// 	}
//
// 	availabilityID, err := uuid.Parse(pathParts[1])
// 	if err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID format"), http.StatusBadRequest)
// 		return
// 	}
//
// 	// Parse request body
// 	var request domain.UpdateAvailabilityRequest
// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()
// 	if err := decoder.Decode(&request); err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
// 		return
// 	}
//
// 	// Call service to update availability
// 	availability, err := h.svc.UpdateAvailability(ctx, availabilityID, request.StartTime, request.EndTime, request.ServiceType, request.PriceCents, request.Notes)
// 	if err != nil {
// 		httpx.RespondWithServiceError(w, logger, ctx, err, "update availability")
// 		return
// 	}
//
// 	// Convert to response DTO
// 	response := domain.AvailabilityResponse{
// 		ID:              availability.ID,
// 		UserID:          availability.UserID,
// 		RoomID:          availability.RoomID,
// 		StartTime:       availability.StartTime,
// 		EndTime:         availability.EndTime,
// 		MaxCapacity:     availability.MaxCapacity,
// 		CurrentBookings: availability.CurrentBookings,
// 		Status:          availability.Status,
// 		ServiceType:     availability.ServiceType,
// 		PriceCents:      availability.PriceCents,
// 		Notes:           availability.Notes,
// 		RecurrencePattern: availability.RecurrencePattern,
// 		ParentID:        availability.ParentID,
// 		CreatedAt:       availability.CreatedAt,
// 		UpdatedAt:       availability.UpdatedAt,
// 	}
//
// 	httpx.RespondWithJSON(w, response, http.StatusOK)
// }
//
// func (h *handler) CancelAvailability(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	logger, err := ctxutil.GetLoggerFromContext(ctx)
// 	if err != nil {
// 		httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Extract availability ID from URL path
// 	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
// 	if len(pathParts) < 3 {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID in path"), http.StatusBadRequest)
// 		return
// 	}
//
// 	availabilityID, err := uuid.Parse(pathParts[1])
// 	if err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID format"), http.StatusBadRequest)
// 		return
// 	}
//
// 	// Call service to cancel availability
// 	err = h.svc.CancelAvailability(ctx, availabilityID)
// 	if err != nil {
// 		httpx.RespondWithServiceError(w, logger, ctx, err, "cancel availability")
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// }
//
// func (h *handler) BlockAvailability(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	logger, err := ctxutil.GetLoggerFromContext(ctx)
// 	if err != nil {
// 		httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Extract availability ID from URL path
// 	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
// 	if len(pathParts) < 3 {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID in path"), http.StatusBadRequest)
// 		return
// 	}
//
// 	availabilityID, err := uuid.Parse(pathParts[1])
// 	if err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid availability ID format"), http.StatusBadRequest)
// 		return
// 	}
//
// 	// Call service to block availability
// 	err = h.svc.BlockAvailability(ctx, availabilityID)
// 	if err != nil {
// 		httpx.RespondWithServiceError(w, logger, ctx, err, "block availability")
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// }
//
// func (h *handler) CheckAvailabilityConflict(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
//
// 	logger, err := ctxutil.GetLoggerFromContext(ctx)
// 	if err != nil {
// 		httpx.RespondWithError(w, err, http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Extract partner ID from URL path
// 	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
// 	if len(pathParts) < 4 {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID in path"), http.StatusBadRequest)
// 		return
// 	}
//
// 	partnerID, err := uuid.Parse(pathParts[1])
// 	if err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid partner ID format"), http.StatusBadRequest)
// 		return
// 	}
//
// 	// Parse query parameters
// 	startTimeStr := r.URL.Query().Get("start_time")
// 	endTimeStr := r.URL.Query().Get("end_time")
// 	excludeIDStr := r.URL.Query().Get("exclude_id")
//
// 	if startTimeStr == "" || endTimeStr == "" {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("start_time and end_time query parameters required"), http.StatusBadRequest)
// 		return
// 	}
//
// 	startTime, err := time.Parse(time.RFC3339, startTimeStr)
// 	if err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid start_time format"), http.StatusBadRequest)
// 		return
// 	}
//
// 	endTime, err := time.Parse(time.RFC3339, endTimeStr)
// 	if err != nil {
// 		httpx.RespondWithError(w, errs.NewInvalidValueErr("invalid end_time format"), http.StatusBadRequest)
// 		return
// 	}
//
// 	var excludeID *uuid.UUID
// 	if excludeIDStr != "" {
// 		if parsed, err := uuid.Parse(excludeIDStr); err == nil {
// 			excludeID = &parsed
// 		}
// 	}
//
// 	// Call service to check conflict
// 	hasConflict, err := h.svc.CheckAvailabilityConflict(ctx, partnerID, startTime, endTime, excludeID)
// 	if err != nil {
// 		httpx.RespondWithServiceError(w, logger, ctx, err, "check availability conflict")
// 		return
// 	}
//
// 	response := struct {
// 		HasConflict bool `json:"has_conflict"`
// 	}{
// 		HasConflict: hasConflict,
// 	}
//
// 	httpx.RespondWithJSON(w, response, http.StatusOK)
// }
