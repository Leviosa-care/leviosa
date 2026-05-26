package messagingHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/domain"
	"github.com/google/uuid"
)

func (h *handler) ListThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	threads, err := h.svc.ListThreads(ctx, userID)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	httpx.RespondWithJSON(w, threads, http.StatusOK)
}

func (h *handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	role, err := getRoleFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	var request domain.CreateThreadRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidInputErr(err), http.StatusBadRequest)
		return
	}

	thread, err := h.svc.CreateThread(ctx, userID, request.ParticipantID, role)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, domain.ErrCannotInitiateThread), errors.Is(err, domain.ErrNoBookingRelationship):
			statusCode = http.StatusForbidden
		case errors.Is(err, domain.ErrThreadAlreadyExists):
			statusCode = http.StatusConflict
		case errors.Is(err, domain.ErrInvalidParticipantID):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	logger.InfoContext(ctx, "Handler: Thread created", "thread_id", thread.ID, "operation", "create_thread")

	httpx.RespondWithJSON(w, map[string]uuid.UUID{"id": thread.ID}, http.StatusCreated)
}

func (h *handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	threadIDStr := r.PathValue("id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidInputErr(err), http.StatusBadRequest)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		var parsed int
		if _, err := fmt.Sscanf(l, "%d", &parsed); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	response, err := h.svc.GetMessages(ctx, threadID, userID, limit, cursor)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, domain.ErrNotThreadParticipant):
			statusCode = http.StatusForbidden
		case errors.Is(err, domain.ErrThreadNotFound):
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpx.RespondWithError(w, errors.New("unsupported media type"), http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	threadIDStr := r.PathValue("id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidInputErr(err), http.StatusBadRequest)
		return
	}

	var request domain.SendMessageRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		httpx.RespondWithError(w, errs.NewInvalidInputErr(err), http.StatusBadRequest)
		return
	}

	message, err := h.svc.SendMessage(ctx, threadID, userID, request.Body)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, domain.ErrNotThreadParticipant):
			statusCode = http.StatusForbidden
		case errors.Is(err, domain.ErrEmptyMessageBody):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	httpx.RespondWithJSON(w, message, http.StatusCreated)
}

func (h *handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	threadIDStr := r.PathValue("id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		httpx.RespondWithError(w, errs.NewInvalidInputErr(err), http.StatusBadRequest)
		return
	}

	if err := h.svc.MarkAsRead(ctx, threadID, userID); err != nil {
		var statusCode int
		switch {
		case errors.Is(err, domain.ErrNotThreadParticipant):
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusInternalServerError
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	httpx.RespondWithJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (h *handler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	count, err := h.svc.GetUnreadCount(ctx, userID)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	httpx.RespondWithJSON(w, domain.UnreadCountResponse{UnreadCount: count}, http.StatusOK)
}
