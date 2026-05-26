package messagingHandler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/infrastructure/sse"
	"github.com/google/uuid"
)

// StreamThreadEvents opens an SSE connection that streams new messages for the
// given thread. The connection stays open until the client disconnects or the
// server shuts down.
func (h *handler) StreamThreadEvents(w http.ResponseWriter, r *http.Request) {
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

	// Verify the user is a participant of this thread by trying to fetch messages.
	// We only need to check access, so fetch 1 message.
	_, err = h.svc.GetMessages(ctx, threadID, userID, 1, "")
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

	// Set SSE headers. Do NOT wrap with EnableCORS middleware — set headers
	// directly because this is a streaming response.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

	w.WriteHeader(http.StatusOK)

	// Flush headers
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Send initial comment to keep the connection alive (some proxies
	// close idle connections quickly).
	if _, err := fmt.Fprint(w, ": connected\n\n"); err != nil {
		return
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Subscribe to events for this thread.
	sub := h.broker.Subscribe(threadID)
	defer h.broker.Unsubscribe(threadID, sub)

	logger, _ := ctxutil.GetLoggerFromContext(ctx)
	if logger != nil {
		logger.InfoContext(ctx, "SSE: client subscribed",
			"thread_id", threadID, "user_id", userID)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if _, err := fmt.Fprint(w, ": ping\n\n"); err != nil {
				return
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

		case ev, ok := <-sub:
			if !ok {
				return
			}
			payload, err := sse.FormatSSE("message", ev)
			if err != nil {
				if logger != nil {
					logger.Error("SSE: format event error", "error", err)
				}
				continue
			}
			if _, err := w.Write(payload); err != nil {
				return
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}
