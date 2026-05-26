package sse

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Leviosa-care/leviosa/backend/internal/messaging/domain"
	"github.com/google/uuid"
)

// Event represents a single SSE event sent to clients.
type Event struct {
	ID   string
	Data any
}

// Subscriber is a channel that receives SSE events for a thread.
type Subscriber = chan Event

// Broker manages SSE subscriptions per thread. It is safe for concurrent use.
type Broker struct {
	mu          sync.Mutex
	subscribers map[uuid.UUID]map[Subscriber]struct{}
	logger      *slog.Logger
}

// NewBroker creates a new Broker.
func NewBroker(logger *slog.Logger) *Broker {
	return &Broker{
		subscribers: make(map[uuid.UUID]map[Subscriber]struct{}),
		logger:      logger,
	}
}

// Subscribe registers a subscriber for the given thread and returns the channel.
func (b *Broker) Subscribe(threadID uuid.UUID) Subscriber {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 8)
	if _, ok := b.subscribers[threadID]; !ok {
		b.subscribers[threadID] = make(map[Subscriber]struct{})
	}
	b.subscribers[threadID][ch] = struct{}{}
	return ch
}

// Unsubscribe removes a subscriber for the given thread.
func (b *Broker) Unsubscribe(threadID uuid.UUID, ch Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subs, ok := b.subscribers[threadID]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(b.subscribers, threadID)
		}
	}
}

// Publish sends a message event to all subscribers of the given thread.
func (b *Broker) Publish(threadID uuid.UUID, msg *domain.MessageResponse) {
	b.mu.Lock()
	subs, ok := b.subscribers[threadID]
	// Copy subscriber keys so we don't hold the lock while sending.
	targets := make([]Subscriber, 0, len(subs))
	for sub := range subs {
		targets = append(targets, sub)
	}
	b.mu.Unlock()

	if !ok || len(targets) == 0 {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		b.logger.Error("sse broker: failed to marshal message", "error", err)
		return
	}

	ev := Event{
		ID:   msg.ID.String(),
		Data: string(data),
	}

	for _, ch := range targets {
		select {
		case ch <- ev:
		default:
			// Drop event if subscriber channel is full; client will
			// miss this message but SSE reconnect will catch up via REST.
			b.logger.Warn("sse broker: subscriber channel full, dropping event",
				"thread_id", threadID, "message_id", msg.ID)
		}
	}
}

// SubscriberCount returns the number of active subscribers for a thread (for diagnostics).
func (b *Broker) SubscriberCount(threadID uuid.UUID) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.subscribers[threadID])
}

// FormatSSE formats an SSE event as per the spec: "event: <type>\ndata: <payload>\nid: <id>\n\n"
func FormatSSE(eventType string, ev Event) ([]byte, error) {
	dataStr, ok := ev.Data.(string)
	if !ok {
		dataBytes, err := json.Marshal(ev.Data)
		if err != nil {
			return nil, fmt.Errorf("marshal sse data: %w", err)
		}
		dataStr = string(dataBytes)
	}

	return []byte(fmt.Sprintf("event: %s\ndata: %s\nid: %s\n\n", eventType, dataStr, ev.ID)), nil
}
